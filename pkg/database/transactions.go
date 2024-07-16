package database

import (
	"context"
	"database/sql"

	"github.com/hashicorp/go-multierror"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	trace "github.com/kava-forge/eve-alts/lib/telemetry"

	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

func deferRollback(ctx context.Context, logger logging.Logger, tx Tx) {
	if tx == nil {
		return
	}

	err := tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Err("error rolling back transaction", err)
	}
}

func newTransaction(ctx context.Context, db Connection, opts *sql.TxOptions) (Tx, error) {
	tx, err := db.BeginTx(ctx, opts)
	return tx, errors.Wrap(err, "could not start transaction")
}

func TransactWithRetries(ctx context.Context, tel *telemetry.Telemeter, logger logging.Logger, db Connection, opts *sql.TxOptions, txFunc func(context.Context, Tx) error) (err error) {
	ctx, span := tel.StartSpan(ctx, "database", "TransactWithRetries")
	defer telemetry.EndSpan(span, &err)

	var tx Tx
	defer deferRollback(ctx, logger, tx)

	success := false
	var multi *multierror.Error

	var attemptCtx context.Context
	var attemptSpan trace.Span
	for i := 0; i < 3; i++ {
		if tx != nil {
			deferRollback(attemptCtx, logger, tx)
		}

		if attemptSpan != nil {
			telemetry.EndSpan(attemptSpan, &err)
		}

		attemptCtx, attemptSpan = tel.StartSpan(ctx, "database", "TransactWithRetries.Attempt")

		tx, err = newTransaction(ctx, db, opts)
		if err != nil {
			telemetry.EndSpan(attemptSpan, &err)
			return errors.Wrap(err, "could not create new transaction")
		}

		err = txFunc(ctx, tx)
		if err == nil {
			success = true
			break
		}

		multi = multierror.Append(multi, err)

		var noRetry nonRetryableError
		if errors.As(err, &noRetry) {
			break
		}
	}

	if attemptSpan != nil {
		telemetry.EndSpan(attemptSpan, &err)
	}

	if !success {
		return errors.Wrap(multi, "transaction failed")
	}

	return errors.Wrap(tx.Commit(), "could not commit transaction")
}

type nonRetryableError struct {
	error
}

func NonRetryableError(err error) error {
	return nonRetryableError{err}
}
