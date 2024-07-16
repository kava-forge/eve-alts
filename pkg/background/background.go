package background

import (
	"context"
	stdhttp "net/http"
	"time"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"golang.org/x/sync/errgroup"

	"github.com/kava-forge/eve-alts/pkg/panics"
)

type Serverlike interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

var _ Serverlike = (*stdhttp.Server)(nil)

func RunBackground(ctx context.Context, logger logging.Logger, servers ...Serverlike) (*errgroup.Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)

	for _, s := range servers {
		g.Go(serverStartFunc(logger, s))
		g.Go(serverShutdownFunc(ctx, logger, s))
	}

	return g, ctx
}

func serverStartFunc(logger logging.Logger, s Serverlike) func() error {
	return func() error {
		defer panics.Handler(logger)

		level.Debug(logger).Message("starting background server")
		err := s.ListenAndServe()
		if err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			return err
		}

		return nil
	}
}

func serverShutdownFunc(ctx context.Context, logger logging.Logger, s Serverlike) func() error {
	return func() error {
		defer panics.Handler(logger)

		<-ctx.Done() // something said we are done

		level.Debug(logger).Message("stopping background server")

		shutdownCtx, cncl := context.WithTimeout(context.Background(), 2*time.Second)
		defer cncl()

		err := s.Shutdown(shutdownCtx)
		if err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			return err
		}
		return nil
	}
}
