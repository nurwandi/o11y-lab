# gitops — ArgoCD (App of Apps)

This is how workloads get onto the cluster. **Terraform never deploys the app or the
platform directly** — it only installs ArgoCD and applies a single *root* Application
(see [`../infra/terraform/user_data.sh.tftpl`](../infra/terraform/user_data.sh.tftpl))
that points here. ArgoCD then syncs everything in this folder.

```
 root Application (created by cloud-init)  --watches-->  gitops/
                                                          ├── platform.yaml  --> platform/k8s  (Collector + LGTM)
                                                          └── app.yaml       --> app/k8s       (Node + Go + PG + Redis)
```

That's the **app-of-apps** pattern: one root app manages child apps; each child app
manages a set of manifests. Change a manifest and `git push` → ArgoCD reconciles the
cluster automatically (`selfHeal` + `prune`).

## The Applications

| File | Deploys | Into namespace |
|---|---|---|
| [`platform.yaml`](platform.yaml) | `platform/k8s` — OpenTelemetry Collector + Tempo + Loki + Prometheus + Grafana | `platform` |
| [`app.yaml`](app.yaml) | `app/k8s` — the observed application | `app` |

## Watching it

```bash
kubectl get applications -n argocd
kubectl -n argocd port-forward svc/argocd-server 8080:443   # ArgoCD UI at https://localhost:8080
```

> ☁️ **AWS Equivalent** — ArgoCD is commonly run on **EKS** for GitOps; the same
> app-of-apps pattern applies. Here it runs on k3s to keep it to a single node.
