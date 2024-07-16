package app

import (
	"context"
	"database/sql"
	"io"
	stdhttp "net/http"
	"net/url"
	"os"
	"path"
	"time"

	//nolint:depguard,staticcheck // uses this internally to do the logging
	"github.com/hashicorp/go-multierror"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/http"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	_ "github.com/mattn/go-sqlite3" //nolint:blank-imports // database driver
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/panics"
	"github.com/kava-forge/eve-alts/pkg/repository"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type Dependencies struct {
	logfile  io.WriteCloser
	logger   logging.Logger
	db       database.Connection
	staticDB database.Connection

	telemetry    *telemetry.Telemeter
	stats        *telemetry.Stats
	statsHandler stdhttp.Handler

	httpClient     http.Client
	esiClient      esi.Client
	callbackServer *esi.CallbackServer

	appRepo    *repository.AppSqliteRepository
	staticRepo *repository.StaticSqliteRepository
}

var _ dependencies = (*Dependencies)(nil)

func CreateDependencies(ctx context.Context, appName, buildVersion string, conf Config) (deps *Dependencies, err error) {
	deps = &Dependencies{}

	var logger logging.Logger

	deps.logfile = &lumberjack.Logger{
		Filename:   path.Join(conf.Logging.Directory, "app.log"),
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		LocalTime:  true,
		Compress:   false,
	}

	switch conf.Logging.Format {
	case "json":
		logger = logging.NewJSONFileLogger(deps.logfile)
	case "json-stdout":
		logger = logging.NewJSONLogger()
	case "logfmt":
		logger = logging.NewLogfmtFileLogger(deps.logfile)
	case "logfmt-stdout":
		logger = logging.NewLogfmtLogger()
	default:
		// TODO: revisit?
		logger = logging.NewJSONFileLogger(deps.logfile)
	}

	logger = logging.WithLevel(logger, conf.Logging.Level)
	logger = logging.With(logger, "timestamp", logging.DefaultTimestampUTC, "caller", logging.DefaultCaller)
	deps.logger = logger

	level.Info(logger).Message("used configuration file", keys.Path, conf.ConfigFile)

	logging.PatchStdLib(deps.logger)

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get hostname")
	}

	if deps.telemetry, deps.statsHandler, err = telemetry.NewTelemeter(ctx, deps, appName, buildVersion, host, telemetry.Options{
		PrometheusNamespace: conf.Telemetry.PrometheusNamespace,
		JaegerURL:           conf.Telemetry.JaegerHostPort,
		TraceProbability:    conf.Telemetry.TraceProbability,
	}); err != nil {
		return nil, errors.Wrap(err, "could not create Telemeter")
	}
	deps.telemetry.MakeDefault()

	if deps.stats, err = telemetry.NewStats(appName, deps.telemetry); err != nil {
		return nil, errors.Wrap(err, "could not create Stats")
	}

	deps.httpClient = http.NewTelemeterClient(deps.Logger(), deps.Telemetry())

	fh, err := os.Create(conf.Database.StaticLocation)
	if err != nil {
		return nil, errors.Wrap(err, "could not open static data file")
	}
	// deferutil.CheckDeferLog(deps.logger, fh.Close)
	if err := esi.WriteDatabase(fh); err != nil {
		return nil, errors.Wrap(err, "could not write static data to disk")
	}
	_ = fh.Close()

	staticdb, err := sql.Open("sqlite3", conf.Database.StaticLocation)
	if err != nil {
		return nil, errors.Wrap(err, "could not open static database", keys.Path, conf.Database.StaticLocation)
	}

	deps.staticDB = &database.WrappedConnection{DB: staticdb}
	deps.staticRepo = repository.NewStaticData(deps)

	appdb, err := sql.Open("sqlite3", conf.Database.Location)
	if err != nil {
		return nil, errors.Wrap(err, "could not open app database", keys.Path, conf.Database.Location)
	}

	deps.db = &database.WrappedConnection{DB: appdb}
	deps.appRepo = repository.NewAppData(deps)

	callbackURL := &url.URL{
		Scheme: conf.Serving.CallbackScheme,
		Host:   conf.Serving.HostPort,
		Path:   conf.Serving.CallbackPath,
	}
	deps.esiClient, err = esi.NewClient(deps, callbackURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "could not create esi client")
	}

	deps.callbackServer = esi.NewCallbackServer(deps.Logger(), conf.Serving.HostPort, conf.Serving.CallbackPath)

	return deps, nil
}

func (d *Dependencies) DB() database.Connection                { return d.db }
func (d *Dependencies) StaticDB() database.Connection          { return d.staticDB }
func (d *Dependencies) Logger() logging.Logger                 { return d.logger }
func (d *Dependencies) HTTPClient() http.Client                { return d.httpClient }
func (d *Dependencies) ESIClient() esi.Client                  { return d.esiClient }
func (d *Dependencies) ESICallbackServer() *esi.CallbackServer { return d.callbackServer }

func (d *Dependencies) Telemetry() *telemetry.Telemeter { return d.telemetry }
func (d *Dependencies) Stats() *telemetry.Stats         { return d.stats }
func (d *Dependencies) StatsHandler() stdhttp.Handler   { return d.statsHandler }

func (d *Dependencies) AppRepo() repository.AppData       { return d.appRepo }
func (d *Dependencies) StaticRepo() repository.StaticData { return d.staticRepo }

func (d *Dependencies) Close(ctx context.Context, timeout time.Duration) error {
	defer level.Debug(d.Logger()).Message("done deps close")
	done := make(chan struct{})

	if timeout == 0 {
		return d.close(ctx, done)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		if err := d.close(ctx, done); err != nil {
			level.Error(d.Logger()).Err("error closing dependencies", err)
		}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Dependencies) close(ctx context.Context, done chan struct{}) error {
	defer panics.Handler(d.Logger())
	defer close(done)

	var errs error

	if err := d.db.Close(ctx); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := d.staticDB.Close(ctx); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := d.Telemetry().Shutdown(ctx); err != nil {
		errs = multierror.Append(errs, err)
	}

	// We don't close the logfile because we probably need
	// to write to if after this if there are issues shutting
	// down everything.
	// if err := d.logfile.Close(); err != nil {
	// 	errs = multierror.Append(errs, err)
	// }

	return errs
}
