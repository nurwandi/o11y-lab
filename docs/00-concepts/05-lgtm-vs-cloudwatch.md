# 05 — LGTM vs CloudWatch

If you work on AWS, this is the honest question: *"AWS already gives me observability
end-to-end — why build the Grafana stack at all?"* Let's answer it fairly.

## First, the honest part: CloudWatch **can** do observability

CloudWatch + X-Ray is a legitimate, production-grade solution — not a toy. It covers
all three pillars:

- Metrics → CloudWatch Metrics
- Logs → CloudWatch Logs (+ Logs Insights to query)
- Traces → AWS X-Ray
- Dashboards + alerts → CloudWatch Dashboards + Alarms
- Correlation → **ServiceLens** and **Application Signals** (themselves built on OTel)

Plenty of teams run on CloudWatch alone and are perfectly happy. Anyone who tells you
"CloudWatch can't do observability" is wrong.

## But it is **not** the same experience. Where they differ:

| Aspect | CloudWatch + X-Ray | LGTM (Grafana stack) |
|---|---|---|
| **Query power** | Logs/Metrics Insights (AWS-proprietary) | **PromQL / LogQL / TraceQL** — industry-standard, more expressive, portable |
| **Cross-pillar correlation** | Improving (App Signals), but can feel console-hoppy | One-click metric → trace → log in a single UI |
| **Dashboards** | Functional, somewhat rigid | Highly flexible; huge library of community dashboards |
| **Metrics cost model** | Custom metrics billed per-metric (pricey with many labels) | Dimensional labels are cheap; friendlier to high-cardinality |
| **Portability** | Tied to AWS | Runs on any cloud or on-prem |
| **Data sources** | AWS data only | Grafana plugs into many sources — **including CloudWatch itself** |

## Managed AWS equivalents

AWS also offers *managed* versions of this exact open-source stack — proof of where
the industry landed:

| Open-source (this lab) | AWS managed |
|---|---|
| OpenTelemetry Collector | **ADOT** (AWS Distro for OpenTelemetry) |
| Prometheus | **AMP** (Amazon Managed Service for Prometheus) |
| Grafana | **AMG** (Amazon Managed Grafana) |
| Loki | CloudWatch Logs |
| Tempo | AWS X-Ray |

## So which do you pick?

- **CloudWatch** if you're all-in on AWS and want zero infrastructure to run — it's
  convenient and deeply integrated.
- **LGTM** if you want stronger query/visualization, lower cost at scale, or
  portability across clouds/on-prem.
- **Hybrid** — very common in the real world: run **Grafana** as the visualization
  layer and point it at *both* CloudWatch **and** Prometheus. Best of both.

## The takeaway

The reason this stack is worth learning — even for an AWS engineer — comes down to two
words: **flexibility** and **no vendor lock-in**. And notice: AWS built ADOT, AMP,
AMG, and Application Signals precisely to close the gap with this open ecosystem. If
CloudWatch were already identical, AWS wouldn't have bothered.

Learn the open stack, and you understand observability *itself* — not just one
vendor's console. That transfers everywhere, including back to AWS.

---

Back to [Stage 0 index](README.md) · Next stage: [build the app](../../README.md#learning-roadmap).
