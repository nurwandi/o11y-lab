# 06 — Golden Signals, SLOs & Correlation

You have three pillars flowing. Now: *what do you actually put on a dashboard, and how
do the pillars work together?*

## What to measure: RED & the Golden Signals

Don't graph everything. Two well-known starter frameworks:

**RED** (great for request-driven services like ours):
- **R**ate — requests per second
- **E**rrors — failed requests per second (or error ratio)
- **D**uration — latency distribution (p50/p95/p99)

**The Four Golden Signals** (Google SRE) add one more:
- Latency, Traffic, Errors, **Saturation** (how "full" a resource is — CPU, memory, queue depth)

Our dashboard (`platform/grafana/.../app-red.json`) is RED: request rate, p95 latency,
rate by route, and an availability stat. Start here; add saturation when you have a
resource that can fill up.

> **Why percentiles, not averages?** An average latency of 100ms can hide that 1% of
> users wait 3 seconds. p95/p99 expose the tail — which is what users actually feel.

## SLI / SLO / SLA

- **SLI** (Indicator) — a number that measures user-visible health.
  *Example: the ratio of non-5xx responses (availability).*
- **SLO** (Objective) — the target for an SLI. *Example: 99.5% availability over 30 days.*
- **SLA** (Agreement) — a contractual promise with consequences if missed. SLOs are
  internal; SLAs are external.

The dashboard's "Availability (non-5xx)" stat is an **SLI**. Pick a target (say 99.5%)
and you have an **SLO**. The gap between your SLO and 100% is your **error budget** —
how much unreliability you're allowed to spend before you stop shipping features and
fix reliability instead.

> ☁️ **AWS Equivalent** — CloudWatch has **ServiceLevelObjective** (Application Signals
> SLOs). Same idea: define an SLI, set a target, track the error budget.

## Correlation: the payoff of one platform

The reason all this lives behind one Grafana is that the pillars **link to each other**:

```
  Dashboard (metric spike)
        │  exemplar / click
        ▼
     Trace (Tempo) ── which span was slow? ──┐
        │  "Logs for this span"              │
        ▼                                     │
     Logs (Loki) ── the exact error ──────────┘  (log's trace_id links back to the trace)
```

This lab wires two of these links (see `datasources.yaml`):
- **Trace → Logs** (`tracesToLogsV2`): from a span, jump to that service's logs.
- **Logs → Trace** (`derivedFields` on `trace_id`): from a log line, jump to its trace.

That round-trip — metric to trace to log and back — is what makes an incident a
two-minute drill-down instead of a two-hour hunt. It only works because service-go's
logs carry the active `trace_id` (that's why we log with `slog` *Context* methods).

## Alerting

An alert is just an SLI query with a threshold and a "for" duration (to avoid paging on
a brief blip). This lab provisions one (`alerting/alerting.yaml`): page if p95 latency
stays above 500ms for 2 minutes. In Stage 6 you can trip it on purpose.

> **Alert on symptoms, not causes.** "p95 latency is high" (what users feel) is a better
> page than "CPU is high" (which may be harmless). Cause-based alerts create noise.

---

Back to [Stage 0 index](README.md). Next stage: [deploy to AWS](../../README.md#learning-roadmap).
