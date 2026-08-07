# prometheus

Stores and queries **metrics** (time-series numbers), queried with **PromQL**.

- **Config:** [`prometheus.yml`](prometheus.yml)
- **Model:** Prometheus **pulls** — it scrapes targets. Here it scrapes the
  OpenTelemetry Collector's `:8889` (the apps' metrics, re-exposed in Prometheus
  format) and `:8888` (the Collector's own internal metrics).
- **UI/API:** http://localhost:9090

Example query — request rate by route:
`sum by (http_route) (rate(http_server_request_duration_seconds_count[1m]))`

> ☁️ **AWS Equivalent** — CloudWatch Metrics · managed: **Amazon Managed Prometheus (AMP)**.
