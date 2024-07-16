package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/telemetry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
)

var ErrBadExporter = errors.New("unsupported exporter")

type Telemeter = telemetry.Telemeter

var (
	KVString                            = telemetry.KVString
	KVStringSlice                       = telemetry.KVStringSlice
	KVInt                               = telemetry.KVInt
	KVInt64                             = telemetry.KVInt64
	KVInt64Slice                        = telemetry.KVInt64Slice
	WithAttributes                      = telemetry.WithAttributes
	WithMetricAttributes                = telemetry.WithMetricAttributes
	HTTPClientAttributesFromHTTPRequest = telemetry.HTTPClientAttributesFromHTTPRequest
	HTTPAttributesFromHTTPStatusCode    = telemetry.HTTPAttributesFromHTTPStatusCode

	CodeError = telemetry.CodeError
	CodeUnset = telemetry.CodeUnset
	CodeOK    = telemetry.CodeOK

	SpanFromContext = telemetry.SpanFromContext
	Handler         = telemetry.Handler
)

type dependencies interface {
	Logger() logging.Logger
}

type Options struct {
	PrometheusNamespace string
	JaegerURL           string
	TraceProbability    float64
}

type promLogger struct {
	logger logging.Logger
}

func (l *promLogger) Println(v ...interface{}) {
	l.logger.Message(fmt.Sprintln(v...))
}

func newPromLogger(logger logging.Logger) promhttp.Logger {
	return &promLogger{logger}
}

func NewTelemeter(ctx context.Context, deps dependencies, appName, version, instanceID string, opts Options) (*telemetry.Telemeter, http.Handler, error) {
	// traceExp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(opts.JaegerURL)))
	// if err != nil {
	// 	return nil, nil, errors.Wrap(err, "could not set up jaeger")
	// }

	traceExp, err := stdouttrace.New()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not set up tracer")
	}

	promReg := prometheus.NewPedanticRegistry()
	if err := promReg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, nil, errors.Wrap(err, "could not register process collector")
	}

	if err := promReg.Register(collectors.NewGoCollector()); err != nil {
		return nil, nil, errors.Wrap(err, "could not register go collector")
	}

	res := telemetry.NewResource(appName, version, instanceID)

	promExporter, err := otelprom.New(otelprom.WithRegisterer(promReg))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create prometheus collector")
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExporter),
	)

	metricHandler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{
		ErrorLog: newPromLogger(deps.Logger()),
	})

	o := telemetry.NewTelemeterFromResource(res, traceExp, meterProvider, opts.TraceProbability)
	return o, metricHandler, nil
}

func NewTestTelemeter(ctx context.Context, deps dependencies, appName, version, instanceID string, opts Options) (*telemetry.Telemeter, http.Handler, error) {
	traceExp, err := stdouttrace.New()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not set up jaeger")
	}

	promReg := prometheus.NewPedanticRegistry()
	if err := promReg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, nil, errors.Wrap(err, "could not register process collector")
	}

	if err := promReg.Register(collectors.NewGoCollector()); err != nil {
		return nil, nil, errors.Wrap(err, "could not register go collector")
	}

	res := telemetry.NewResource(appName, version, instanceID)

	promExporter, err := otelprom.New(otelprom.WithRegisterer(promReg))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create prometheus collector")
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExporter),
	)

	metricHandler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{
		ErrorLog: newPromLogger(deps.Logger()),
	})

	o := telemetry.NewTelemeterFromResource(res, traceExp, meterProvider, opts.TraceProbability)
	return o, metricHandler, nil
}

func StartSpan(ctx context.Context, telemeter *telemetry.Telemeter, pkg, op string, opts ...telemetry.StartSpanOption) (context.Context, telemetry.Span) {
	// requestID, ok := request.GetRequestID(ctx)
	// if ok {
	// 	opts = append(opts, telemetry.WithAttributes(telemetry.KVString(keys.RequestID, requestID)))
	// }
	return telemeter.StartSpan(ctx, pkg, op, opts...)
}

func EndSpan(span telemetry.Span, err *error) { //nolint:gocritic // this needs to be a *error
	if err != nil && *err != nil {
		span.SetStatus(telemetry.CodeError, (*err).Error())
		span.RecordError(*err)
	}
	span.End()
}

func SetSpanError(span telemetry.Span, errmsg string) {
	if errmsg != "" {
		span.SetStatus(telemetry.CodeError, errmsg)
	}
}
