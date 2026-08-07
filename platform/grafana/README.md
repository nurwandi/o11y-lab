# grafana

The **single pane of glass** — dashboards, alerting, and cross-pillar correlation
over Prometheus, Tempo, and Loki.

- **Provisioning:** [`provisioning/datasources/datasources.yaml`](provisioning/datasources/datasources.yaml)
  wires all three datasources automatically on startup — no clicking required.
- **UI:** http://localhost:3001 (anonymous admin, login disabled — lab only)

Datasources provisioned: **Prometheus** (default), **Tempo**, **Loki**. Dashboards
and the metric→trace→log drill-down are built in Stage 4.

> ☁️ **AWS Equivalent** — CloudWatch Dashboards + Alarms · managed: **Amazon Managed Grafana (AMG)**.
>
> ⚠️ Anonymous admin with no login is for local learning only — never in production.
