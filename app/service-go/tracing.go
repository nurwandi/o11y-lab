package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer wires up the OpenTelemetry trace provider and returns a shutdown
// function that flushes any buffered spans.
//
//   - If OTEL_EXPORTER_OTLP_ENDPOINT is set, spans are exported over OTLP/HTTP
//     (in Stage 3+ that endpoint is the OpenTelemetry Collector).
//   - Otherwise spans are printed to stdout, so the instrumentation is verifiable
//     with no backend at all.
//
// The most important line here is the propagator: installing W3C tracecontext is
// what lets a trace that started in api-node CONTINUE into this service — the
// incoming `traceparent` header is read back into a Go context.
func initTracer(ctx context.Context) (func(context.Context) error, error) {
	var (
		exp sdktrace.SpanExporter
		err error
	)
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		exp, err = otlptracehttp.New(ctx) // reads endpoint from the environment
	} else {
		exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
	if err != nil {
		return nil, err
	}

	// Service name (and any extra attributes) come from OTEL_SERVICE_NAME /
	// OTEL_RESOURCE_ATTRIBUTES via the env detector.
	res, err := resource.New(ctx, resource.WithFromEnv(), resource.WithTelemetrySDK())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	return tp.Shutdown, nil
}
