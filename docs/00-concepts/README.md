# Stage 0 — Concepts & Foundations

Before we build or deploy anything, we build the *mental model*. This is the
difference between copying YAML off the internet and actually knowing what your
telemetry is doing.

Read these in order. They're short on purpose.

| # | Lesson | You'll understand |
|---|---|---|
| 01 | [What is observability?](01-what-is-observability.md) | Monitoring vs observability, and why it matters |
| 02 | [The three pillars](02-the-three-pillars.md) | Metrics, logs, traces — and when to reach for each |
| 03 | [OpenTelemetry](03-opentelemetry.md) | The vendor-neutral "hub" that generates & routes telemetry |
| 04 | [The LGTM stack](04-lgtm-stack.md) | Loki, Grafana, Tempo, Prometheus — who stores what |
| 05 | [LGTM vs CloudWatch](05-lgtm-vs-cloudwatch.md) | How this compares to AWS-native observability |
| 06 | [Golden signals, SLOs & correlation](06-golden-signals-slo.md) | What to dashboard, SLIs/SLOs, and how the pillars link (Stage 4) |

> **☁️ Coming from AWS?** Every lesson pairs concepts with their AWS equivalent, so
> you can map new ideas onto tools you already know (CloudWatch, X-Ray, ADOT, AMP, AMG).

**By the end of Stage 0**, you should be able to explain — to a colleague, in plain
language — what the three pillars are, why OpenTelemetry exists, and why a team might
run the Grafana stack instead of (or alongside) CloudWatch.

Next up: [Stage 1 — build the app](../../README.md#learning-roadmap).
