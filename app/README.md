# app — The Observed Application

This is the application that everything else in `o11y-lab` observes. It's small on
purpose: two services, a database, and a cache — just enough moving parts to produce
*interesting* telemetry (a cross-service call, a DB query, a cache hit vs. miss)
without drowning you in code.

> **Stage 1 scope:** get the app running and talking to itself. There is **no
> observability instrumentation yet** — that's [Stage 2](../docs/00-concepts/README.md).
> Right now we just want a working system to instrument.

## Components

| Service | Language | Role |
|---|---|---|
| [`api-node`](api-node/) | Node.js (Express) | Edge API. Owns no data — proxies to `service-go`. |
| [`service-go`](service-go/) | Go | Owns the data. Reads/writes Postgres, caches in Redis. |
| Postgres | — | System of record for `products`. |
| Redis | — | Read-through cache for single-product lookups. |

## Request flow

```
  client
    │  GET /api/products/1
    ▼
 [api-node] ──HTTP──► [service-go]
                          ├─ check Redis  ── HIT ─► return (X-Cache: HIT)
                          └─ MISS ─► query Postgres ─► store in Redis ─► return (X-Cache: MISS)
```

The `X-Cache` response header tells you which path a request took. Why it matters:
in Stage 2 this becomes visible in a trace — a cache HIT is a short, shallow trace;
a MISS is a longer trace with a Postgres span inside it. That contrast is the whole
point of the app.

## Run it locally

Requires Docker.

```bash
docker compose up --build
```

Then, in another terminal:

```bash
# health
curl localhost:3000/health

# list products
curl localhost:3000/api/products

# first read = cache MISS (watch the header)
curl -i localhost:3000/api/products/1
# same read again = cache HIT
curl -i localhost:3000/api/products/1

# create a product
curl -X POST localhost:3000/api/products \
  -H 'content-type: application/json' \
  -d '{"name":"Webcam","price":5999}'
```

Tear down with `docker compose down` (add `-v` to wipe the database volume).

## Endpoints (via api-node, port 3000)

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Liveness of api-node |
| GET | `/api/products` | List all products |
| GET | `/api/products/:id` | Get one product (cached) |
| POST | `/api/products` | Create a product |

> **☁️ AWS Equivalent** — In production on AWS this shape would typically be
> containers on **ECS/EKS** behind an **ALB**, with **RDS** (Postgres) and
> **ElastiCache** (Redis). Here we keep it as plain containers so the focus stays on
> observability, not infrastructure. We deploy to real AWS in Stage 5.

## Tracing (Stage 2)

Both services are instrumented with **OpenTelemetry**:

- **api-node** — zero-code auto-instrumentation (loaded via `node --import`). Traces
  Express and the outgoing `fetch`, and injects the `traceparent` header.
- **service-go** — `otelhttp` (server span + reads that header), `otelpgx` (a span per
  Postgres query), and `redisotel` (a span per Redis command).

By default (`docker compose up`) spans are printed to **stdout** — no backend needed.
To *see* the trace as a waterfall, bring up the Jaeger overlay:

```bash
docker compose -f docker-compose.yml -f docker-compose.otel.yml up --build
# generate some traffic
curl localhost:3000/api/products/1        # cache MISS -> includes a Postgres span
curl localhost:3000/api/products/1        # cache HIT  -> no Postgres span
# then open the UI
open http://localhost:16686               # search service "api-node"
```

You'll see a **single trace spanning both services**. A MISS has a Postgres span; a
HIT doesn't — the same contrast the app was designed around, now visible end-to-end.

> Jaeger here is a temporary viewer. Stage 3 replaces it with the OpenTelemetry
> Collector fanning out to the LGTM stack (Tempo for traces).

## Container images (GHCR)

The images are standardized to **GitHub Container Registry**:

```
ghcr.io/nurwandi/o11y-lab/api-node
ghcr.io/nurwandi/o11y-lab/service-go
```

They're built and pushed automatically by
[`.github/workflows/build-images.yml`](../.github/workflows/build-images.yml) on every
push to `main` (using the built-in `GITHUB_TOKEN`). Locally, `docker compose build`
tags images with these same names.

## What's next

**Stage 3 is built** — see [`platform/`](../platform/) for the full OpenTelemetry
Collector + LGTM stack. Bring it up (all three signals land in Grafana):

```bash
docker compose -f docker-compose.yml -f docker-compose.platform.yml up --build
# Grafana: http://localhost:3001
```

**Stage 4 — Dashboards & correlation:** Grafana dashboards, alerts, SLOs, and the
metric → trace → log drill-down.
