package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// initLogger returns an slog.Logger and a shutdown func.
//
// When OTEL_EXPORTER_OTLP_ENDPOINT is set, logs are exported over OTLP (to the
// Collector, then Loki) AND — because we log with *Context methods — each record
// carries the active trace_id. That trace_id is what lets Grafana jump from a log
// line straight to its trace. Without a backend, we just log plain text to stdout.
func initLogger(ctx context.Context) (*slog.Logger, func(context.Context) error, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return slog.New(slog.NewTextHandler(os.Stdout, nil)), noopShutdown, nil
	}

	exp, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, nil, err
	}
	res, err := resource.New(ctx, resource.WithFromEnv(), resource.WithTelemetrySDK())
	if err != nil {
		return nil, nil, err
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(lp)
	return otelslog.NewLogger("service-go"), lp.Shutdown, nil
}

func noopShutdown(context.Context) error { return nil }
