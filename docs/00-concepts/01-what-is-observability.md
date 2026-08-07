# 01 — What is Observability?

## Monitoring vs observability

**Monitoring** answers questions you *knew to ask in advance*. You decide up front —
"alert me if CPU > 80%" or "page me if the site is down" — and you watch those known
signals. Monitoring is great for **known-unknowns**: failure modes you can predict.

**Observability** is the ability to ask questions you *didn't* anticipate, after the
fact, from the data your system already emits. It's built for **unknown-unknowns** —
the weird, novel failure at 2am that no dashboard was pre-built for.

> Monitoring: *"Is the thing I'm watching okay?"*
> Observability: *"Something's weird — let me interrogate the system and find out why."*

In a single monolith you can often get by with monitoring. In a distributed system —
many services, calling each other, deployed independently — a request can fail in a
hundred ways you never predicted. That's where observability earns its keep.

## The definition

A system is **observable** when you can understand its internal state purely from the
signals it emits, without shipping new code to answer a new question. Those signals
come in three complementary forms — the **three pillars**:

- 📊 **Metrics** — *"Is there a problem?"*
- 📝 **Logs** — *"What exactly happened?"*
- 🔍 **Traces** — *"Where, across all my services, and why?"*

The next lesson digs into each. The important idea for now: they're not competitors,
they're a **team**. Metrics tell you *something* is wrong, traces tell you *where*,
and logs tell you the *exact detail*.

## Why it matters (the payoff)

The goal isn't pretty dashboards. It's **shrinking the time between "something's
wrong" and "here's the root cause"** — what teams call MTTR (mean time to
resolution). Good observability turns a two-hour blind hunt into a two-minute
drill-down.

> **☁️ AWS Equivalent** — On AWS, observability is delivered by **CloudWatch**
> (metrics, logs, dashboards, alarms) plus **X-Ray** (traces), with **CloudWatch
> Application Signals** stitching them together. Same three pillars — different
> packaging. We'll build the open-source version so you understand the machinery
> underneath, then map it back to AWS in [lesson 05](05-lgtm-vs-cloudwatch.md).

---

Next: [The three pillars →](02-the-three-pillars.md)
