# infra/bootstrap — One-time setup

This runs **once, locally**, before the pipeline can work. It creates the pieces that
can't provision themselves (chicken-and-egg):

1. **GitHub OIDC provider** — lets GitHub Actions get temporary AWS credentials with
   no stored access keys.
2. **CI IAM role** (`o11y-lab-ci`) — trusted only by the `nurwandi/o11y-lab` repo,
   assumed by the pipeline.
3. **S3 state bucket** — remote backend for the main stack (`../terraform`).

Because it bootstraps the backend itself, this stack keeps **local** state.

## Run it

```bash
cd infra/bootstrap
cp terraform.tfvars.example terraform.tfvars   # then edit state_bucket_name
terraform init
terraform apply
```

> If the account already has a GitHub OIDC provider (only one is allowed per
> account), set `create_github_oidc_provider = false` first.

## After apply

Note the outputs and set them as **GitHub repo variables** (Settings → Secrets and
variables → Actions → Variables) so the pipeline can use them:

| Output | GitHub repo variable |
|---|---|
| `ci_role_arn` | `AWS_CI_ROLE_ARN` |
| `state_bucket` | `TF_STATE_BUCKET` |

You'll also set `OPERATOR_CIDR` (your public IP as `x.x.x.x/32`) and `EC2_KEY_NAME`
(an existing EC2 key pair) — see [`../terraform`](../terraform/).

## Permissions note

The CI role attaches `AmazonEC2FullAccess` for simplicity. In production, scope this
to least privilege. ☁️ *This whole file is the "GitHub OIDC → IAM role" pattern AWS
recommends over long-lived access keys.*
