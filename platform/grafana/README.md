# grafana

The **single pane of glass** — dashboards, alerting, and cross-pillar correlation
over Prometheus, Tempo, and Loki.

- **Provisioning:** [`provisioning/datasources/datasources.yaml`](provisioning/datasources/datasources.yaml)
  wires all three datasources automatically on startup — no clicking required.
- **UI:** http://localhost:3001 (anonymous admin, login disabled — lab only)

Provisioned automatically:

- **Datasources:** Prometheus (default), Tempo, Loki.
- **Correlation:** Tempo → Loki (`tracesToLogsV2`) and Loki → Tempo (`derivedFields`
  on `trace_id`) — click a trace to see its logs, click a log to open its trace.
- **Dashboard:** `App — RED` (`provisioning/dashboards/app-red.json`) — request rate,
  p95 latency, rate by route, and an availability SLI.
- **Alert:** `alerting/alerting.yaml` — fires when p95 latency stays above 500ms.

See [`docs/00-concepts/06-golden-signals-slo.md`](../../docs/00-concepts/06-golden-signals-slo.md)
for the RED / SLO / correlation concepts behind these.

> ☁️ **AWS Equivalent** — CloudWatch Dashboards + Alarms · managed: **Amazon Managed Grafana (AMG)**.
>
> ⚠️ Anonymous admin with no login is for local learning only — never in production.
