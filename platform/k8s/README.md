# platform/k8s — Kubernetes manifests for the observability stack

Plain Kubernetes manifests for the OpenTelemetry Collector + LGTM, synced by ArgoCD
([`../../gitops/platform.yaml`](../../gitops/platform.yaml)) into the `platform`
namespace.

| File | Resource |
|---|---|
| `namespace.yaml` | the `platform` namespace |
| `otel-collector.yaml` | Collector (ConfigMap + Deployment + Service) |
| `tempo.yaml` | Tempo (traces) |
| `loki.yaml` | Loki (logs) |
| `prometheus.yaml` | Prometheus (metrics) |
| `grafana.yaml` | Grafana + provisioned datasources, alerting |
| `grafana-dashboards.yaml` | the `App — RED` dashboard as a ConfigMap |

**These configs mirror the local compose stack** ([`../`](../)) — same service names,
same content — just packaged as Kubernetes ConfigMaps. The Collector still exports to
`tempo:4317`, `loki:3100`, and exposes metrics on `:8889` for Prometheus to scrape,
exactly as it does locally.

Storage is `emptyDir` (ephemeral) — fine for a lab; use PVCs + object storage for
anything real.

Reach Grafana:

```bash
kubectl -n platform port-forward svc/grafana 3001:3000   # http://localhost:3001
```

> ☁️ **AWS Equivalent** — the managed versions would be ADOT + AMP + Amazon Managed
> Grafana + CloudWatch Logs. See [`docs/00-concepts/05-lgtm-vs-cloudwatch.md`](../../docs/00-concepts/05-lgtm-vs-cloudwatch.md).
