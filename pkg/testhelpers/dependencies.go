package testhelpers

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"strings"
	"testing"

	"github.com/kava-forge/eve-alts/lib/http"
	"github.com/kava-forge/eve-alts/lib/http/httpfakes"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/loggingfakes"
	"go.opentelemetry.io/otel/attribute"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/database/databasefakes"
	"github.com/kava-forge/eve-alts/pkg/repository/repositoryfakes"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type TestDependencies struct {
	TestLogger       *loggingfakes.FakeLogger
	configuredLogger logging.Logger
	TestDB           *databasefakes.FakeConnection
	TestStaticDB     *databasefakes.FakeConnection
	TestAppRepo      *repositoryfakes.FakeAppData
	TestStaticRepo   *repositoryfakes.FakeStaticData
	TestHTTPClient   *httpfakes.FakeClient

	TestTelemetry *telemetry.Telemeter
	TestStats     *telemetry.Stats
	statsHandler  stdhttp.Handler
}

const (
	AppName      = "CORVEE_TEST"
	BuildVersion = "TEST"
	Host         = "TEST_HOST"
	JaegerURL    = "http://jaeger.local/api/traces"
)

var skipKeys = map[string]bool{
	"caller":               true,
	"timestamp":            true,
	"request_id":           true,
	"net.sock.peer.addr":   true,
	"net.sock.peer.port":   true,
	"net.protocol.version": true,
	"net.host.name":        true,
	"http.scheme":          true,
}

func NewTestDependencies(t *testing.T) *TestDependencies {
	t.Helper()

	ctx := context.Background()
	var err error

	deps := &TestDependencies{}

	deps.TestLogger = &loggingfakes.FakeLogger{}
	deps.TestLogger.LogCalls(func(keyvals ...interface{}) error {
		b := &strings.Builder{}
		for i := 0; i < len(keyvals); i += 2 {
			k := keyvals[i]

			if ks, ok := k.(string); ok && skipKeys[ks] {
				continue
			}

			if ks, ok := k.(fmt.Stringer); ok && skipKeys[ks.String()] {
				continue
			}

			if ks, ok := k.(attribute.Key); ok && skipKeys[string(ks)] {
				continue
			}

			var v interface{}
			if len(keyvals) > i+1 {
				v = keyvals[i+1]
			}
			_, _ = fmt.Fprintf(b, "\n%s=%v", k, v)
		}
		t.Log(b.String())
		return nil
	})
	logger := logging.Logger(deps.TestLogger)
	logger = logging.WithLevel(logger, "info")
	logger = logging.With(logger, "timestamp", logging.DefaultTimestampUTC, "caller", logging.DefaultCaller)
	deps.configuredLogger = logger

	logging.PatchStdLib(deps.configuredLogger)

	if deps.TestTelemetry, deps.statsHandler, err = telemetry.NewTestTelemeter(ctx, deps, AppName, BuildVersion, Host, telemetry.Options{
		PrometheusNamespace: "test",
		TraceProbability:    1.0,
	}); err != nil {
		t.Fatalf("could not create Telemeter: %v", err)
	}
	deps.TestTelemetry.MakeDefault()

	if deps.TestStats, err = telemetry.NewStats(AppName, deps.TestTelemetry); err != nil {
		t.Fatalf("could not create Stats: %v", err)
	}

	deps.TestHTTPClient = &httpfakes.FakeClient{}

	deps.TestDB = &databasefakes.FakeConnection{}
	deps.TestDB.BeginTxReturns(&databasefakes.FakeTx{}, nil)

	return deps
}

func (d *TestDependencies) Telemetry() *telemetry.Telemeter { return d.TestTelemetry }
func (d *TestDependencies) Stats() *telemetry.Stats         { return d.TestStats }
func (d *TestDependencies) DB() database.Connection         { return d.TestDB }
func (d *TestDependencies) Logger() logging.Logger          { return d.configuredLogger }
func (d *TestDependencies) StatsHandler() stdhttp.Handler   { return d.statsHandler }
func (d *TestDependencies) HTTPClient() http.Client         { return d.TestHTTPClient }
