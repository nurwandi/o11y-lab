# o11y-lab — Observability, End to End

> Learn production-grade **observability** from the ground up. A small polyglot app
> (Node.js + Go) is instrumented with **OpenTelemetry**, ships its telemetry to the
> **Grafana LGTM** stack, runs on **Kubernetes (k3s)**, and is deployed to **AWS**
> entirely as code with **Terraform** + **ArgoCD (GitOps)**.
>
> Every concept is paired with its **☁️ AWS equivalent**, so if you know AWS, this
> clicks immediately.

**Status:** 🚧 Work in progress — built in the open, one stage at a time.

---

## Why this repo

Most observability tutorials either hand you a `docker-compose up` black box or drop
you into a wall of YAML. This repo is different: it's a **guided course**. You build
the whole thing yourself and, more importantly, you understand **why** each piece
exists — the difference between *monitoring* ("is it up?") and *observability*
("why is it slow, and where?").

By the end you can answer, for any request in a distributed system:
- 📊 **Is there a problem?** → metrics
- 📝 **What exactly happened?** → logs
- 🔍 **Where and why is it slow?** → traces

...and correlate all three in a single Grafana click.

---

## The stack (and its AWS equivalent)

| Component | Role | ☁️ AWS Equivalent |
|---|---|---|
| **OpenTelemetry** (SDK + Collector) | Generate & route telemetry. Vendor-neutral. *Not* storage. | CloudWatch Agent + X-Ray SDK · managed: **ADOT** |
| **Prometheus** | Store & query **metrics** | CloudWatch Metrics · managed: **AMP** |
| **Loki** | Store & query **logs** | CloudWatch Logs |
| **Tempo** | Store & query **traces** | AWS **X-Ray** |
| **Grafana** | Dashboards, alerts, correlation | CloudWatch Dashboards + Alarms · managed: **AMG** |

> **Key idea:** OpenTelemetry is the *hub*. The app only knows "send to the
> Collector" — the Collector decides where data lands. Swap a backend later without
> touching a line of application code. That's how you avoid vendor lock-in.

---

## Architecture

```
   Developer ──terraform apply──►  AWS
                                    │
                   ┌────────────────┴─────────────────┐
                   │  VPC · Public Subnet · IGW · SG   │
                   │   single EC2  (cloud-init → k3s)   │
                   └────────────────┬───────────────────┘
                                    │  k3s (lightweight Kubernetes)
                            [ ArgoCD ] ──pull from Git──► sync
                                    │
        ┌───────────────────────────┼────────────────────────────┐
        │  APP (observed)           │      PLATFORM (observability) │
        │                           │                               │
        │  [Node.js API] ──trace──► [Go service] ──► Postgres        │
        │        │                      │      └────► Redis          │
        │        └── metrics/logs/traces ──┐                         │
        │                                   ▼                        │
        │                        [OpenTelemetry Collector]  ◄─ hub   │
        │                                   │                        │
        │            ┌───────────┬──────────┼──────────┐             │
        │            ▼           ▼          ▼           │            │
        │        Prometheus     Loki      Tempo         │            │
        │            └───────────┴──────────┘           │            │
        │                        ▼                       │            │
        │                    [Grafana] ← dashboards, correlation, alert│
        └────────────────────────────────────────────────────────────┘
                 ▲
          [k6 load generator] ── keeps the dashboards alive
```

Full design & decisions: [`docs/architecture/design.md`](docs/architecture/design.md).

---

## Learning roadmap

Follow the stages in order. Each is self-contained; **stages 0–4 run on your laptop
for free** — you only need AWS for stage 5.

| Stage | Focus | What you'll learn |
|---|---|---|
| **0** | Concepts & foundations | The 3 pillars, OpenTelemetry, why LGTM, LGTM vs CloudWatch |
| **1** | Build the app | Node.js + Go services calling each other |
| **2** | Instrumentation | OTel SDK → emit traces/metrics/logs; context propagation |
| **3** | Local platform | Collector + LGTM locally; watch telemetry arrive |
| **4** | Dashboards & correlation | Grafana dashboards, alerts, SLOs, metric → trace → log drilldown |
| **5** | IaC to AWS | Terraform → EC2 → k3s → ArgoCD; full GitOps deploy |
| **6** | *(optional)* Chaos & incident | Break things on purpose; debug with observability |

---

## Repository layout

| Path | What lives here |
|---|---|
| [`docs/`](docs/) | The course: concepts and architecture |
| `app/` | The observed application (Node.js API + Go service) |
| `platform/` | The observability stack (OTel Collector + LGTM) |
| `gitops/` | ArgoCD Application definitions (app-of-apps) |
| `infra/` | Terraform: AWS network, EC2, k3s + ArgoCD bootstrap |
| `load/` | k6 load generator |

Folders appear as each stage is built. Every folder has its own README.

---

## 💰 Cost (read before stage 5)

Stage 5 runs on a **single EC2** with a deliberately cheap network — **no NAT
Gateway, no load balancer, no idle Elastic IP**.

**Billed only while running:** the EC2 instance (~$0.083/hr for `t3.large`), its EBS
root volume (~$2.4/mo), and one public IPv4 (~$0.005/hr).
**Free:** VPC, subnet, route table, Internet Gateway, security group.

Everything is Infrastructure as Code, so `terraform destroy` returns you to **$0**.
Spin it up to learn, tear it down when done.

---

## License

[MIT](LICENSE) — free to learn from, fork, and build on.
