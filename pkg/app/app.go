package app

import (
	"context"
	stdhttp "net/http"
	"net/http/pprof"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/http"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/migrations"
	"github.com/kava-forge/eve-alts/pkg/background"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/panics"
	"github.com/kava-forge/eve-alts/pkg/repository"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type dependencies interface {
	DB() database.Connection
	StaticDB() database.Connection
	Logger() logging.Logger
	HTTPClient() http.Client
	ESIClient() esi.Client
	ESICallbackServer() *esi.CallbackServer

	Telemetry() *telemetry.Telemeter
	Stats() *telemetry.Stats
	StatsHandler() stdhttp.Handler

	AppRepo() repository.AppData
	StaticRepo() repository.StaticData
}

type App struct {
	deps dependencies
	conf Config

	app    fyne.App
	window fyne.Window
}

type SkillLevel struct {
	Name  string
	Level int
}

func New(appName string, deps dependencies, conf Config) (*App, error) {
	a := app.NewWithID("org.evogames.eve-alts")
	a.Settings().SetTheme(&AppTheme{})
	w := NewMainWindow(deps, a)

	return &App{
		deps:   deps,
		conf:   conf,
		app:    a,
		window: w,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.With(a.deps.Logger(), keys.Component, "App")

	level.Debug(logger).Message("running migrations")
	if err := a.deps.DB().Migrate(ctx, migrations.Migrations); err != nil {
		return errors.Wrap(err, "could not run database migrations")
	}

	done := make(chan error)
	level.Debug(logger).Message("starting background")
	go a.startBackground(ctx, done)
	level.Debug(logger).Message("starting app")
	a.window.ShowAndRun()
	level.Debug(logger).Message("app done")
	cancel()
	level.Debug(logger).Message("waiting for background")
	bgerr := <-done
	level.Debug(logger).Message("background done")

	return errors.Wrap(bgerr, "background tasks ended unexpectedly")
}

func (a *App) startBackground(ctx context.Context, done chan error) {
	logger := logging.With(a.deps.Logger(), keys.Component, "App.startBackground")

	defer panics.Handler(logger)
	defer close(done)

	servers := make([]background.Serverlike, 0, 3)
	servers = append(servers, a.deps.ESICallbackServer())

	ch := a.deps.StatsHandler()
	if ch != nil {
		level.Debug(logger).Message("configuring metrics")

		mux := stdhttp.NewServeMux()
		mux.Handle("/metrics", ch)

		servers = append(servers, &stdhttp.Server{
			Addr:         a.conf.Serving.PrometheusHostPort,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      mux,
		})
	}

	if a.conf.PProf.Enabled {
		level.Debug(logger).Message("configuring pprof")

		mux := stdhttp.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		servers = append(servers, &stdhttp.Server{
			Addr:         a.conf.Serving.PProfHostPort,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      mux,
		})
	}

	g, _ := background.RunBackground(ctx, logger, servers...)

	done <- g.Wait()
}
