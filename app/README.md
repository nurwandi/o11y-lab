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

## What's next

[Stage 2 — Instrumentation](../docs/00-concepts/README.md): add the OpenTelemetry SDK
to both services and watch a single request produce a trace that spans Node → Go →
Postgres/Redis.
