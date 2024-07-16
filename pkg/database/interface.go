package database

import (
	"context"
	"database/sql"
	"embed"
)

//go:generate counterfeiter -generate

type (
	Row  = sql.Row
	Rows = sql.Rows
)

//counterfeiter:generate . Connection
type Connection interface {
	BeginTx(context.Context, *sql.TxOptions) (Tx, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Migrate(context.Context, embed.FS) error
	Close(ctx context.Context) error
}

//counterfeiter:generate . Tx
type Tx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Commit() error
	Rollback() error
}
