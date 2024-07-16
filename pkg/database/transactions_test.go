package database_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/loggingfakes"
	"github.com/stretchr/testify/assert"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/database/databasefakes"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type deps struct{}

func (d deps) Logger() logging.Logger {
	return &loggingfakes.FakeLogger{}
}

func TestTransactWithRetries(t *testing.T) {
	t.Parallel()

	fakeTel, _, err := telemetry.NewTestTelemeter(context.Background(), deps{}, "test", "test", "test", telemetry.Options{
		PrometheusNamespace: "test",
		TraceProbability:    1.0,
	})
	if err != nil {
		t.Fatalf("could not create Telemeter: %v", err)
	}

	type args struct {
		tel    *telemetry.Telemeter
		logger logging.Logger
		db     *databasefakes.FakeConnection
		opts   *sql.TxOptions
		txFunc func(ctx context.Context, tx database.Tx) error
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrAs     interface{}
		wantAttemptCt int
		wantCommit    bool
	}{
		{
			name: "no retries",
			args: args{
				tel:    fakeTel,
				logger: &loggingfakes.FakeLogger{},
				db:     &databasefakes.FakeConnection{},
				opts:   &sql.TxOptions{},
				txFunc: func(ctx context.Context, tx database.Tx) error {
					return nil
				},
			},
			wantErr:       false,
			wantAttemptCt: 1,
			wantCommit:    true,
		},
		{
			name: "full failure",
			args: args{
				tel:    fakeTel,
				logger: &loggingfakes.FakeLogger{},
				db:     &databasefakes.FakeConnection{},
				opts:   &sql.TxOptions{},
				txFunc: func(ctx context.Context, tx database.Tx) error {
					return errors.New("fail")
				},
			},
			wantErr:       true,
			wantErrAs:     new(multierror.Error),
			wantAttemptCt: 3,
			wantCommit:    false,
		},
		{
			name: "non-retryable",
			args: args{
				tel:    fakeTel,
				logger: &loggingfakes.FakeLogger{},
				db:     &databasefakes.FakeConnection{},
				opts:   &sql.TxOptions{},
				txFunc: func(ctx context.Context, tx database.Tx) error {
					return database.NonRetryableError(errors.New("fail"))
				},
			},
			wantErr:       true,
			wantErrAs:     new(multierror.Error),
			wantAttemptCt: 1,
			wantCommit:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			callCt := 0
			txFunc := func(ctx context.Context, tx database.Tx) error {
				callCt++
				return tt.args.txFunc(ctx, tx)
			}

			fakeTx := &databasefakes.FakeTx{}
			tt.args.db.BeginTxReturns(fakeTx, nil)

			err := database.TransactWithRetries(context.Background(), tt.args.tel, tt.args.logger, tt.args.db, tt.args.opts, txFunc)
			if tt.wantErr {
				if assert.Error(t, err) {
					if tt.wantErrAs != nil {
						assert.ErrorAs(t, err, &tt.wantErrAs)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantAttemptCt, callCt)
			assert.Equal(t, tt.wantAttemptCt, tt.args.db.BeginTxCallCount())

			if tt.wantCommit {
				assert.Equal(t, 1, fakeTx.CommitCallCount())
			} else {
				assert.Equal(t, 0, fakeTx.CommitCallCount())
			}
		})
	}
}
