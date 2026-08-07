# 03 — OpenTelemetry (OTel)

## The problem OTel solves

Before OpenTelemetry, every vendor had its own agent and SDK. Want to emit traces?
Use the vendor's tracing library. Switch vendors later? Re-instrument your entire
codebase. With *M* telemetry tools and *N* services, you faced an *N×M* integration
mess — and painful **vendor lock-in**.

**OpenTelemetry** is an open, vendor-neutral standard for **generating, collecting,
and exporting** telemetry (traces, metrics, logs). Instrument your app *once* with
OTel, and you can send that data to *any* backend — Grafana's stack, an AWS service,
a commercial vendor — by changing configuration, not code.

## The two halves of OTel

**1. The SDK / instrumentation** — libraries you add to your application code. They
create spans, record metrics, and emit logs, then hand them off. Instrumentation can
be:
- **Automatic** — drop-in libraries that trace common frameworks (HTTP, gRPC, DB
  drivers) with almost no code changes.
- **Manual** — you add spans/attributes for your own business logic ("this span is
  the payment step").

**2. The Collector** — a standalone service that **receives** telemetry from your
apps, **processes** it (batch, filter, add metadata, sample), and **exports** it to
one or more backends. Its pipeline is simple to picture:

```
   apps ──► [ receivers ] ──► [ processors ] ──► [ exporters ] ──► backends
              (OTLP in)        (batch, filter)     (to Loki,        (Loki,
                                                    Tempo, Prom)     Tempo,
                                                                     Prometheus)
```

## The key idea: OTel is the *hub*, not the *storage*

This is the most common misconception, so let's be explicit:

> **OpenTelemetry does NOT store your telemetry.** It's the *plumbing* — the courier
> that collects and routes. The storage lives in the backends (Loki, Tempo,
> Prometheus).

Your application only needs to know one thing: *"send everything to the Collector."*
The Collector decides where it lands. Want to swap Loki for Elasticsearch tomorrow?
Change one exporter in the Collector config — your application code never moves. That
decoupling is the entire point, and it's why OTel has become the industry standard.

## Common mix-up

> ❌ *"OpenTelemetry is like a CloudWatch log group."*
> ✅ Not quite. A log group is **storage** (that's **Loki's** job). OTel is closer to
> the **CloudWatch agent + X-Ray daemon** — it *collects and forwards*, it doesn't
> keep the data.

> **☁️ AWS Equivalent** — AWS ships **ADOT** (AWS Distro for OpenTelemetry), its own
> supported build of the OTel Collector. That's the tell: AWS didn't replace
> OpenTelemetry — it **adopted** it. Learn OTel once and it works on AWS, on other
> clouds, and on-prem.

---

Next: [The LGTM stack →](04-lgtm-stack.md)
