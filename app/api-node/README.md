# api-node — Edge API

The front door of the app. A thin [Express](https://expressjs.com/) service that
receives client requests and proxies them to `service-go`. It intentionally owns no
database or cache of its own.

**Why a separate front service?** So that in Stage 2 we can watch a trace cross a
language boundary (Node.js → Go) via **context propagation** — the single most
important idea in distributed tracing.

## Configuration

| Env var | Default | Purpose |
|---|---|---|
| `PORT` | `3000` | Port to listen on |
| `SERVICE_URL` | `http://localhost:8080` | Base URL of `service-go` |

## Endpoints

| Method | Path | Proxies to |
|---|---|---|
| GET | `/health` | — (local liveness) |
| GET | `/api/products` | `GET service-go/products` |
| GET | `/api/products/:id` | `GET service-go/products/{id}` |
| POST | `/api/products` | `POST service-go/products` |

If `service-go` is unreachable, the API responds with `502` rather than crashing.

## Run standalone

```bash
npm install
SERVICE_URL=http://localhost:8080 npm start
```

Usually you'll run it via the app's [`docker compose`](../README.md#run-it-locally)
instead, which wires it to `service-go` automatically.

## Notes

- Uses Node's built-in `fetch` (Node 18+) — no HTTP client dependency.
- ES modules (`"type": "module"`).
