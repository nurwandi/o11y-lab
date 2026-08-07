# loki

Stores and queries **logs**, queried with **LogQL**. "Prometheus, but for logs" — it
indexes labels, not full log text, which keeps it cheap.

- **Config:** [`loki.yaml`](loki.yaml)
- **Ingest:** native OTLP at `:3100/otlp` (the Collector sends here)
- **Storage:** local filesystem (lab); production uses object storage
- **API:** http://localhost:3100

`allow_structured_metadata: true` lets Loki keep the `trace_id` that rides along with
each log, which is what powers log → trace correlation in Grafana.

Example query: `{service_name="service-go"}`

> ☁️ **AWS Equivalent** — CloudWatch **Logs**.
