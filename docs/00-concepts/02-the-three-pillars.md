# 02 — The Three Pillars

Metrics, logs, and traces each answer a different question. Learn *which question*
each one answers and you'll always know which to reach for.

## 📊 Metrics — "Is there a problem?"

Numbers measured over time (time-series): request rate, error rate, latency, CPU,
queue depth. They're **aggregated and cheap**, so they're perfect for **dashboards
and alerts**.

- **Strength:** cheap to store, fast to query, great for "trend over time" and alerting.
- **Weakness:** aggregated — a metric tells you error rate jumped to 5%, but not
  *which* requests failed or *why*.
- **Example:** *"`checkout` p99 latency went from 200ms to 900ms at 14:03."*

> **☁️ AWS Equivalent** — CloudWatch **Metrics** (managed Prometheus: **AMP**).

## 📝 Logs — "What exactly happened?"

Timestamped, discrete event records — ideally **structured** (JSON with fields), not
just free-text. Logs carry the *detail* a metric can't.

- **Strength:** rich, specific context for a single event.
- **Weakness:** high volume and cost if you log everything; hard to see the big
  picture from logs alone.
- **Example:** *"`user=8f21 checkout failed: dial tcp redis:6379: i/o timeout`."*

> **☁️ AWS Equivalent** — CloudWatch **Logs** (log groups + Logs Insights for query).

## 🔍 Traces — "Where, and why, is it slow?"

A **trace** follows one request as it travels across services. It's made of **spans**
— each span is one unit of work (an HTTP handler, a DB query) with a start, a
duration, and a parent. Together the spans form a tree that shows exactly where time
went.

The magic that makes this work across services is **context propagation**: when
service A calls service B, it passes a trace ID (and span ID) along in the request
headers, so B's spans attach to the *same* trace. That's how a single trace can span
Node.js → Go → database.

- **Strength:** pinpoints *where* in a distributed call chain the latency or error is.
- **Weakness:** usually **sampled** at scale (you can't afford to keep every trace).
- **Example:** *"That 900ms checkout? 750ms of it was a single Postgres query in the
  Go service."*

> **☁️ AWS Equivalent** — AWS **X-Ray**.

## The debugging loop (how they work together)

This is the whole point. In a real incident you move through the pillars:

```
  📊 Metric alert fires        "error rate on checkout is up"
        │  (something is wrong)
        ▼
  🔍 Open traces               find a slow / errored request, see WHICH span failed
        │  (where is it)
        ▼
  📝 Jump to that span's logs   read the exact error for that request
           (why is it)
```

A metric tells you *there's a fire*, a trace tells you *which room*, a log tells you
*what's burning*. A good tool (Grafana — [lesson 04](04-lgtm-stack.md)) lets you hop
between them with a click instead of copy-pasting IDs between five consoles.

## Quick reference

| Pillar | Question | Best for | Watch out for | AWS |
|---|---|---|---|---|
| Metrics | Is there a problem? | dashboards, alerts | aggregated, no detail | CloudWatch Metrics / AMP |
| Logs | What happened? | root-cause detail | volume & cost | CloudWatch Logs |
| Traces | Where & why? | distributed latency/errors | sampling | X-Ray |

---

Next: [OpenTelemetry →](03-opentelemetry.md)
