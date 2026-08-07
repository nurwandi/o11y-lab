# 04 — The LGTM Stack

Once the OpenTelemetry Collector has your telemetry, it needs somewhere to **store**
each signal and something to **visualize** it. That's the job of the Grafana stack,
nicknamed **LGTM**.

## What LGTM stands for

| Letter | Tool | Stores / does |
|---|---|---|
| **L** | **Loki** | Logs |
| **G** | **Grafana** | Visualization, dashboards, alerting, correlation |
| **T** | **Tempo** | Traces |
| **M** | **Mimir** | Metrics (large-scale Prometheus) |

> **In this lab we use plain Prometheus instead of Mimir.** Mimir is Prometheus
> scaled out for huge, multi-tenant deployments — overkill for a single node. Plain
> Prometheus does the same job (store & query metrics) far more lightly. So our stack
> is really **L·G·T + Prometheus**.

## Each component

**Prometheus — metrics store.** A time-series database with a powerful query
language, **PromQL**. The de-facto standard for metrics in the Kubernetes world.
> ☁️ *AWS Equivalent:* CloudWatch Metrics · managed: **Amazon Managed Prometheus (AMP)**.

**Loki — logs store.** Think "Prometheus, but for logs." It indexes *labels* (not the
full log text), which keeps it cheap. Queried with **LogQL**.
> ☁️ *AWS Equivalent:* CloudWatch **Logs**.

**Tempo — traces store.** A cheap, high-scale trace backend. Queried with **TraceQL**.
> ☁️ *AWS Equivalent:* AWS **X-Ray**.

**Grafana — the single pane of glass.** Connects to all three stores as *data
sources*, renders dashboards, runs alerts, and — crucially — lets you **correlate**:
click a spike on a metric graph → jump to the traces behind it → open the logs for a
specific span. One UI, all three pillars.
> ☁️ *AWS Equivalent:* CloudWatch **Dashboards + Alarms** · managed: **Amazon Managed Grafana (AMG)**.

## How it all fits with OTel

```
  [ Your app ] ──OTLP──► [ OpenTelemetry Collector ] ──┬──► Prometheus ─┐
                                                        ├──► Loki ───────┤
                                                        └──► Tempo ──────┤
                                                                         ▼
                                                                    [ Grafana ]
                                                             dashboards · alerts · correlation
```

The app speaks **one** protocol (OTLP) to the Collector. The Collector fans each
signal out to its store. Grafana reads from all the stores and gives you the unified
view. Clean separation: **generate → route → store → visualize.**

## Why this stack (for this lab)

- **Open-source & self-hostable** — runs comfortably on a single small node.
- **One ecosystem** — the four tools are designed to work together; correlation is
  first-class, not bolted on.
- **Portable** — nothing here is tied to a cloud vendor.
- **Industry-relevant** — Prometheus + Grafana are everywhere in Kubernetes shops,
  *including* teams running on AWS.

The next lesson asks the fair question: *AWS already gives me all this — why bother?*

---

Next: [LGTM vs CloudWatch →](05-lgtm-vs-cloudwatch.md)
