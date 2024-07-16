package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type Dependencies struct {
	logger    logging.Logger
	staticDB  database.Connection
	telemetry *telemetry.Telemeter
}

func CreateDependencies(ctx context.Context, appName, buildVersion, databaseFile string) (deps *Dependencies, err error) {
	deps = &Dependencies{}

	var logger logging.Logger
	logger = logging.NewJSONLogger()

	logger = logging.WithLevel(logger, "debug")
	logger = logging.With(logger, "timestamp", logging.DefaultTimestampUTC, "caller", logging.DefaultCaller)
	deps.logger = logger

	logging.PatchStdLib(deps.logger)

	host, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get hostname")
	}

	if deps.telemetry, _, err = telemetry.NewTelemeter(ctx, deps, appName, buildVersion, host, telemetry.Options{
		PrometheusNamespace: "clean-static",
		JaegerURL:           "",
		TraceProbability:    0.0,
	}); err != nil {
		return nil, errors.Wrap(err, "could not create Telemeter")
	}
	deps.telemetry.MakeDefault()

	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database", keys.Path, databaseFile)
	}

	deps.staticDB = &database.WrappedConnection{DB: db}

	return deps, nil
}

func (d *Dependencies) StaticDB() database.Connection   { return d.staticDB }
func (d *Dependencies) Logger() logging.Logger          { return d.logger }
func (d *Dependencies) Telemetry() *telemetry.Telemeter { return d.telemetry }
