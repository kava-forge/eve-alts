package database

import (
	"context"
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kava-forge/eve-alts/lib/errors"
)

type WrappedConnection struct {
	*sql.DB
}

var _ Connection = (*WrappedConnection)(nil)

func (c *WrappedConnection) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return c.DB.BeginTx(ctx, opts)
}

func (c *WrappedConnection) Migrate(ctx context.Context, migrations embed.FS) error {
	sd, err := iofs.New(migrations, ".")
	if err != nil {
		return errors.Wrap(err, "could not create source driver")
	}

	dd, err := sqlite3.WithInstance(c.DB, &sqlite3.Config{})
	if err != nil {
		return errors.Wrap(err, "could not create db driver")
	}

	migrator, err := migrate.NewWithInstance("iofs", sd, "sqlite3", dd)
	if err != nil {
		return errors.Wrap(err, "could not create migrator")
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "could not migrate database")
	}

	return nil
}

func (c *WrappedConnection) Close(ctx context.Context) error {
	return c.DB.Close()
}
