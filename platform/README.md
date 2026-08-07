# platform — The Observability Stack

This is the **observers**: the OpenTelemetry Collector (the hub) and the LGTM stack
that stores and visualizes the telemetry the app emits.

```
 app --OTLP--> [ OpenTelemetry Collector ] --+--> Tempo       (traces)
                                             +--> Prometheus  (metrics, scraped)
                                             +--> Loki        (logs)
                                                     |
                                                  [ Grafana ] <- one UI for all three
```

The app only ever talks to the Collector. Each backend can be swapped by editing the
Collector config — never the application. That's the hub pattern in practice.

## Components

| Folder | Component | Role | ☁️ AWS Equivalent |
|---|---|---|---|
| [`otel-collector/`](otel-collector/) | OpenTelemetry Collector | Receive OTLP, fan out per signal | ADOT |
| [`prometheus/`](prometheus/) | Prometheus | Store & query metrics | CloudWatch Metrics / AMP |
| [`loki/`](loki/) | Loki | Store & query logs | CloudWatch Logs |
| [`tempo/`](tempo/) | Tempo | Store & query traces | AWS X-Ray |
| [`grafana/`](grafana/) | Grafana | Dashboards, correlation, alerts | CloudWatch Dashboards / AMG |

## Run it (locally, Stage 3)

The stack runs as a compose overlay on top of the app, so everything shares one
network:

```bash
cd ../app
docker compose -f docker-compose.yml -f docker-compose.platform.yml up --build
```

UIs once it's up:

| Service | URL |
|---|---|
| Grafana | http://localhost:3001 (anonymous admin) |
| Prometheus | http://localhost:9090 |
| Tempo (API) | http://localhost:3200 |
| Loki (API) | http://localhost:3100 |

Generate some traffic, then explore in Grafana (datasources are pre-provisioned):

```bash
for i in $(seq 1 20); do curl -s localhost:3000/api/products/1 >/dev/null; done
```

## How each signal gets there

- **Traces** — app → Collector (OTLP) → **Tempo** (OTLP). A trace spans api-node → Go.
- **Metrics** — app → Collector (OTLP) → Collector re-exposes them in Prometheus
  format on `:8889` → **Prometheus** scrapes it. (service-go's `otelhttp` emits HTTP
  server metrics automatically once a MeterProvider exists.)
- **Logs** — service-go logs via `slog` → Collector (OTLP) → **Loki**. Because the logs
  carry the active `trace_id`, Grafana can jump log → trace.

## Why config lives here (not inline in compose)

Each component's config is a real file under its folder. In Stage 5 these same files
become Kubernetes **ConfigMaps** — so nothing here is throwaway; the compose overlay
and the future k3s deployment share one source of truth.

## Note on the Stage 2 overlay

`app/docker-compose.otel.yml` (Jaeger) from Stage 2 still exists as a minimal
traces-only viewer. This platform overlay supersedes it with the full stack.
