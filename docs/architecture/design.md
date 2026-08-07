# o11y-lab — Design Document

**Date:** 2026-08-07
**Status:** Approved (design phase)
**Repo (planned):** `o11y-lab` — public, GitHub `nurwandi`

---

## 1. Purpose & Audience

`o11y-lab` is an educational, portfolio-grade project that teaches **observability
end-to-end**, mimicking a real production environment. A learner clones the repo,
follows a staged path, and ends with a fully instrumented application whose
telemetry (metrics, logs, traces) flows through a modern open-source stack —
deployed to AWS entirely as code.

**Primary audience:** global, general (English narrative). Extra care is taken for
readers with an **AWS background** — every concept and tool is paired with its
**AWS equivalent** so the mental model transfers instantly.

**Success criteria:**
- A reader with no prior observability experience understands the *why* behind each
  of the three pillars, not just the *what*.
- The full environment is reproducible with `terraform apply` and torn down with
  `terraform destroy` (cost-safe).
- Every folder and subfolder has a clear, well-written README.
- The repo reads like a guided course, not a dump of config files.

---

## 2. The Stack (and why)

| Component | Role | ☁️ AWS Equivalent |
|---|---|---|
| **OpenTelemetry** (SDK + Collector) | Generate telemetry in the app; collect & route it. Vendor-neutral. NOT storage. | CloudWatch Agent + X-Ray SDK/daemon; managed: **ADOT** |
| **Prometheus** | Store & query **metrics** (time-series) | CloudWatch Metrics; managed: **AMP** |
| **Loki** | Store & query **logs** | CloudWatch Logs |
| **Tempo** | Store & query **traces** | AWS **X-Ray** |
| **Grafana** | Dashboards, alerting, cross-pillar correlation | CloudWatch Dashboards + Alarms; managed: **AMG** |

> "LGTM" = Loki, Grafana, Tempo, Mimir. On a single EC2 we use plain **Prometheus**
> (lighter) instead of Mimir. Same role: the metrics store.

**Why OpenTelemetry as the hub:** the app only knows "send to the Collector." The
Collector decides where telemetry lands. Swapping a backend later (e.g. Loki →
Elastic) means changing the Collector, not the app. Decoupling = no vendor lock-in.

**The three pillars — when to use which:**
- **Metrics** → *"Is there a problem?"* (aggregate, cheap, alerting)
- **Logs** → *"What exactly happened?"* (detailed events)
- **Traces** → *"Where and why is it slow?"* (one request across services)

The real debugging loop: metric alert → find the slow/errored request in traces →
jump to its exact logs. Grafana ties all three together in one UI.

---

## 3. Architecture

```
Developer
    |
    | terraform apply  (via GitHub Actions + OIDC)
    v
+-------------------------------------+
| AWS: VPC / Public Subnet / IGW / SG |
|                                     |
|   EC2  --(cloud-init)-->  k3s       |
|   Git  --(pull)-->  ArgoCD          |
+-------------------------------------+
    |
    | ArgoCD deploys to cluster
    v
APP (observed) ===============================================

 k6 --load--> Node.js API --trace--> Go service
                  |                     |     |
                  |                     v     v
                  |                 Postgres Redis
                  |
                  +-- metrics / logs / traces
                  |
                  v
      OpenTelemetry Collector  (hub)
                 |
PLATFORM (observability) =====================================

     +-----------+-----------+
     |           |           |
     v           v           v
 Prometheus    Loki       Tempo
     |           |           |
     +-----------+-----------+
                 |
                 v
              Grafana
    (dashboards / correlation / alerts)
```

**Layers (deliberate separation of concerns):**
- **infra** — *where it runs* (Terraform → AWS + k3s + ArgoCD bootstrap)
- **platform** — *the observers* (LGTM + OTel Collector)
- **app** — *what is observed* (Node.js + Go + datastores)
- **gitops** — *how apps get deployed* (ArgoCD Application definitions)

---

## 4. Deployment model — Terraform (infra) + ArgoCD (GitOps)

- **Terraform** provisions infrastructure and bootstraps the cluster:
  VPC → EC2 → install k3s (via cloud-init) → install ArgoCD.
- **ArgoCD** is the single source of truth for workloads: it pulls Kubernetes
  manifests / Helm values from this repo (`platform/` and `app/`, orchestrated by
  `gitops/` app-of-apps) and continuously syncs them into k3s.
- Result: change a manifest → `git push` → the cluster reconciles itself. This is
  GitOps. ☁️ *Comparable to ArgoCD on EKS, or a CodePipeline-driven GitOps flow.*

Because ArgoCD pulls from a **public** GitHub repo, no registry auth is needed.

### CI/CD — Terraform runs in GitHub Actions (OIDC, keyless)

Terraform is **not** run from a laptop. It runs in a **GitHub Actions** pipeline that
authenticates to AWS via **OIDC** — GitHub's identity provider issues a short-lived
token, AWS trusts it and hands back temporary credentials. **No long-lived AWS access
keys are ever stored in the repo.** ☁️ *This is the AWS-recommended pattern for CI/CD
(GitHub OIDC → IAM role AssumeRoleWithWebIdentity).*

- **Region:** `ap-southeast-3` (Jakarta).
- **AWS account:** any AWS account you control, referenced via a local AWS CLI profile.
- **Trust:** an IAM OIDC provider for `token.actions.githubusercontent.com` + an IAM
  role whose trust policy is scoped to your GitHub repo.
- **State:** remote state in **S3** (with native S3 state locking — no DynamoDB
  needed). Required because pipeline runners are ephemeral; local state won't do.

**Bootstrap (chicken-and-egg):** the OIDC provider, the IAM role, and the S3 state
bucket must exist *before* the pipeline can run. These live in `infra/bootstrap/` and
are applied **once, manually, with your local AWS CLI profile**. Everything after that —
the actual environment (VPC, EC2, k3s, ArgoCD) in `infra/terraform/` — runs only
through the pipeline.

---

## 5. Network & cost design (no NAT, minimal billing)

Single-AZ VPC, one **public** subnet. EC2 gets a public IP and reaches the internet
directly through the Internet Gateway — **no NAT Gateway** (which would cost ~$32/mo).

```
+-------------------------------------------------------+
| VPC  (single AZ)                                      |
|                                                       |
|   Public Subnet                                       |
|     Internet Gateway (IGW) .......... free  (egress)  |
|     Route table: 0.0.0.0/0 -> IGW ... free            |
|     Security Group (SSH + UI -> your IP) . free       |
|     EC2 (public IP, runs k3s) ....... BILLED          |
+-------------------------------------------------------+
```

**Billed only while running:**
| Item | Estimate | Note |
|---|---|---|
| EC2 (t3.large) | ~$0.083/hr (~$60/mo if 24/7) | core cost |
| EBS root (gp3 ~30GB) | ~$2.4/mo | attached to the instance |
| Public IPv4 | ~$0.005/hr | AWS charges IPv4 since 2024; gone on destroy |

**Free / $0:** VPC, subnet, route table, IGW, security group.
**Explicitly avoided:** NAT Gateway, ALB/NLB, idle Elastic IP.

**One unavoidable extra (pennies):** the pipeline needs remote Terraform state, so an
**S3 bucket** holds `terraform.tfstate`. At this size it costs well under $0.10/mo
(a few KB of state + a handful of requests), and it persists between runs so it is
*not* destroyed by `terraform destroy`. Using S3-native locking means no DynamoDB
table, keeping the footprint to a single cheap bucket.

`terraform destroy` returns everything to $0. The root README states this prominently.

---

## 6. Repository structure

```
o11y-lab/
├── README.md              # hero: what/why, architecture diagram, quickstart, cost note
├── docs/
│   ├── 00-concepts/       # observability 101, 3 pillars, OTel, LGTM vs CloudWatch
│   └── architecture/      # diagrams + design decisions
├── app/                   # the observed application
│   ├── api-node/          # Node.js API + OTel instrumentation
│   └── service-go/        # Go service + OTel instrumentation
├── platform/              # observability stack (Helm values / manifests)
│   ├── otel-collector/
│   ├── prometheus/
│   ├── loki/
│   ├── tempo/
│   └── grafana/           # provisioned datasources + dashboards
├── gitops/                # ArgoCD app-of-apps + Application CRs (watch app/ + platform/)
├── infra/
│   ├── bootstrap/         # one-time, local (your AWS profile): OIDC provider, IAM role, S3 state bucket
│   └── terraform/         # the environment (VPC, EC2, k3s, ArgoCD) — applied only via pipeline
├── .github/workflows/     # GitHub Actions: terraform plan/apply via OIDC (region ap-southeast-3)
└── load/                  # k6 load generator
```

**README standard (every folder):** clear purpose, how to use it, what it depends on,
and — where relevant — a **"☁️ AWS Equivalent"** callout box.

---

## 7. Learning roadmap (staged)

Each stage is roughly one coaching session and one meaningful commit.

| Stage | Focus | Learner outcome |
|---|---|---|
| **0** | Concepts & foundations | 3 pillars, OTel, why LGTM, LGTM vs CloudWatch |
| **1** | Build the app | Node + Go calling each other, running locally |
| **2** | Instrumentation | OTel SDK → emit traces/metrics/logs; context propagation |
| **3** | Local platform | Collector + LGTM via docker-compose / k3d; see data arrive |
| **4** | Dashboards & correlation | Grafana dashboards, alerts, SLO, metric→trace→log drilldown |
| **5** | IaC to AWS | Terraform: VPC→EC2→k3s→ArgoCD; GitOps sync. **E2E full IaC.** |
| **6** | (optional) Chaos & incident | Inject failures; practice debugging with observability |

**Why local first (1–4) before AWS (5):** fast, free iteration on a laptop to build
understanding, then "lift" to the cloud — the real DevOps workflow. Readers without
an AWS account can still learn through Stage 4.

---

## 8. Out of scope (YAGNI)

- Multi-node / HA cluster (single EC2 by design — cost & clarity).
- Mimir / long-term metrics storage (plain Prometheus is enough at this scale).
- Full OpenTelemetry Demo (~15–20 services — too heavy; we build a small app to
  learn instrumentation ourselves). May be noted as an optional advanced module.
- Managed AWS observability (AMP/AMG/X-Ray) — referenced as equivalents for learning,
  not deployed. A possible future comparison module.

---

## 9. Conventions

- **Repo narrative:** English (global audience).
- **Coaching sessions:** Bahasa Indonesia; explain the *why* at every step.
- Every stack/concept paired with its **AWS equivalent** in the docs.
- Every folder/subfolder has a good README.
- All activity recorded in local memory; work pushed to GitHub (`nurwandi`).
