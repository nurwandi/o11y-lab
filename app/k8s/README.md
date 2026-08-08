# app/k8s — Kubernetes manifests for the app

Plain Kubernetes manifests for the observed application, synced by ArgoCD
([`../../gitops/app.yaml`](../../gitops/app.yaml)) into the `app` namespace.

| File | Resource |
|---|---|
| `namespace.yaml` | the `app` namespace |
| `postgres.yaml` | Postgres Deployment + Service |
| `redis.yaml` | Redis Deployment + Service |
| `service-go.yaml` | Go service (GHCR image) + Service |
| `api-node.yaml` | Node API (GHCR image) + Service |

Notes:
- Images come from **GHCR** (built by [`build-images.yml`](../../.github/workflows/build-images.yml)).
  Make those packages **public** so the cluster can pull them without credentials.
- The app sends OTLP to the Collector **across namespaces**:
  `http://otel-collector.platform:4318`.
- Data stores use no PVC — ephemeral, which is fine for a lab. Add PVCs if you want
  data to survive pod restarts.

Reach the API to generate traffic:

```bash
kubectl -n app port-forward svc/api-node 3000:3000
curl localhost:3000/api/products/1
```

> ☁️ **AWS Equivalent** — Postgres → **RDS**, Redis → **ElastiCache** in production;
> here they're in-cluster to stay single-node.
