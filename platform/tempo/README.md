# tempo

Stores and queries **traces**, queried with **TraceQL**.

- **Config:** [`tempo.yaml`](tempo.yaml)
- **Ingest:** OTLP/gRPC on `:4317` (the Collector sends here)
- **Storage:** local filesystem (fine for a lab; production uses object storage)
- **API:** http://localhost:3200 (Grafana's Tempo datasource points here)

A trace is one request across services; each span is a unit of work. Tempo is cheap
because it doesn't index span contents — you pivot into it from metrics/logs by
trace ID.

> ☁️ **AWS Equivalent** — AWS **X-Ray**.
