# infra/terraform — The Environment

Provisions the running environment and hands off to GitOps:

- **VPC** — one public subnet, Internet Gateway, **no NAT** (cost-conscious).
- **EC2** — a single instance whose cloud-init installs **k3s**, installs **ArgoCD**,
  and applies the root Application pointing at [`../../gitops`](../../gitops/).
- From there, **ArgoCD** syncs the platform + app from Git. Terraform never deploys
  the workloads directly.

This stack runs **only through the pipeline** ([`.github/workflows/terraform.yml`](../../.github/workflows/terraform.yml)),
which authenticates via OIDC. Run [`../bootstrap`](../bootstrap/) once first.

## Inputs (set as GitHub repo variables)

| Variable | Repo variable | Example |
|---|---|---|
| `operator_cidr` | `OPERATOR_CIDR` | `203.0.113.4/32` (your public IP) |
| `key_name` | `EC2_KEY_NAME` | an existing EC2 key pair |
| — | `AWS_CI_ROLE_ARN` | from bootstrap output |
| — | `TF_STATE_BUCKET` | from bootstrap output |

## Run locally (optional, for testing)

```bash
cd infra/terraform
terraform init -backend-config="bucket=<your state bucket>"
terraform apply -var="operator_cidr=$(curl -s ifconfig.me)/32" -var="key_name=<your key>"
```

## After apply

```bash
terraform output          # public_ip, ssh + kubeconfig hints
# grab kubeconfig, then:
kubectl get applications -n argocd   # watch ArgoCD sync the stack
kubectl -n platform port-forward svc/grafana 3001:3000
```

## 💰 Tear down when done

```bash
terraform destroy -var="operator_cidr=..." -var="key_name=..."
```

Everything except the state bucket is destroyed → back to ~$0. See the
[root README cost note](../../README.md#-cost-read-before-stage-5).
