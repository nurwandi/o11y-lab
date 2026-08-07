# service-go — Data Service

The data-owning backend. A small Go HTTP service (using the standard library only,
plus a Postgres and a Redis driver) that manages a `products` table and caches
single-product reads in Redis.

## Configuration

| Env var | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | Port to listen on |
| `DATABASE_URL` | `postgres://o11y:o11y@localhost:5432/o11y?sslmode=disable` | Postgres DSN |
| `REDIS_ADDR` | `localhost:6379` | Redis address |

## Endpoints

| Method | Path | Behaviour |
|---|---|---|
| GET | `/health` | Liveness |
| GET | `/products` | List all products (Postgres) |
| GET | `/products/{id}` | **Read-through cache:** Redis → miss → Postgres → populate Redis |
| POST | `/products` | Insert a product (Postgres) |

### The cache path (why it's here)

`GET /products/{id}` sets an `X-Cache` response header:

- `X-Cache: HIT` — served from Redis (fast, no DB touch)
- `X-Cache: MISS` — not cached; fetched from Postgres and cached for 60s

This gives us two visibly different request shapes to observe once tracing is added
in Stage 2.

## Design notes

- **Routing:** Go 1.22 `http.ServeMux` with method+path patterns (`GET /products/{id}`)
  — no third-party router needed.
- **Schema:** created idempotently on startup (`CREATE TABLE IF NOT EXISTS`) and
  seeded with a few rows on first run. Startup retries while Postgres finishes booting.
- **Money:** `price` is stored as integer **cents** to avoid float rounding bugs.
- **Dependencies:** `jackc/pgx` (Postgres) and `redis/go-redis` — resolved at Docker
  build time via `go mod tidy`, so no `go.sum` is committed for this lab service.

## Run standalone

```bash
go mod tidy
DATABASE_URL=... REDIS_ADDR=... go run .
```

Usually run via the app's [`docker compose`](../README.md#run-it-locally).

> **☁️ AWS Equivalent** — Postgres → **RDS**; Redis → **ElastiCache**. Same roles,
> managed by AWS. We stay on plain containers here to keep the focus on observability.
