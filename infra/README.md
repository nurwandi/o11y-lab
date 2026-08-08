# infra — Infrastructure as Code

Everything needed to run o11y-lab on AWS, as code. Two Terraform stacks with a
deliberate split:

| Folder | What | How it runs |
|---|---|---|
| [`bootstrap/`](bootstrap/) | GitHub OIDC provider, CI IAM role, S3 state bucket | **once, locally** (chicken-and-egg) |
| [`terraform/`](terraform/) | VPC, EC2, k3s + ArgoCD bootstrap | **via the pipeline** (OIDC, keyless) |

## The flow

```
 1. bootstrap  (local, once)  ->  OIDC provider + CI role + state bucket
 2. push to main             ->  GitHub Actions assumes the CI role via OIDC
 3. terraform apply (in CI)  ->  VPC + EC2 + k3s + ArgoCD
 4. ArgoCD                   ->  syncs platform + app from gitops/
```

Region: **ap-southeast-3 (Jakarta)**. No long-lived AWS keys anywhere — the pipeline
uses short-lived OIDC credentials.

☁️ This is the AWS-recommended CI/CD identity pattern (GitHub OIDC → IAM role via
`AssumeRoleWithWebIdentity`).

## Order of operations

1. Run [`bootstrap/`](bootstrap/) and set the outputs as GitHub repo variables.
2. Set `OPERATOR_CIDR` and `EC2_KEY_NAME` repo variables.
3. Push — the pipeline provisions [`terraform/`](terraform/).
4. ArgoCD takes over ([`../gitops`](../gitops/)).
