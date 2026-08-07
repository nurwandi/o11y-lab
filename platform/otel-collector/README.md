# otel-collector

The **hub**. Receives OTLP from the apps and routes each signal to its backend.

- **Config:** [`config.yaml`](config.yaml)
- **Receives:** OTLP on `:4317` (gRPC) and `:4318` (HTTP)
- **Exports:** traces → Tempo (OTLP), metrics → Prometheus format on `:8889`
  (scraped), logs → Loki (OTLP)
- **Internal telemetry:** its own metrics on `:8888`

The point of the Collector: the app targets *one* endpoint and never needs to know
which backends exist. Swapping Loki for something else is a change in `config.yaml`,
not in application code.

> ☁️ **AWS Equivalent** — **ADOT** (AWS Distro for OpenTelemetry) is AWS's supported
> build of this exact Collector.
