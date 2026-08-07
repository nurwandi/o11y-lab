// api-node is the edge/front API of the o11y-lab app. It owns no data of its own;
// every request is proxied to service-go over HTTP. This two-service shape is
// deliberate: in Stage 2 it lets us watch a single trace propagate across a
// language boundary (Node.js -> Go).
//
// There is NO OpenTelemetry here yet — that comes in Stage 2.

import express from "express";

const app = express();
app.use(express.json());

const SERVICE_URL = process.env.SERVICE_URL || "http://localhost:8080";
const PORT = process.env.PORT || 3000;

// Proxy a request to service-go and return its raw response. `fetch` is built
// into Node 18+, so no HTTP client dependency is needed.
async function callService(path, init) {
  const res = await fetch(`${SERVICE_URL}${path}`, init);
  return { status: res.status, body: await res.text(), xCache: res.headers.get("x-cache") };
}

// Relay an upstream response to the client, preserving the X-Cache header so the
// cache HIT/MISS signal survives the proxy hop.
function relay(res, r) {
  if (r.xCache) res.set("X-Cache", r.xCache);
  res.status(r.status).type("application/json").send(r.body);
}

app.get("/health", (_req, res) => {
  res.json({ status: "ok", service: "api-node" });
});

app.get("/api/products", async (_req, res, next) => {
  try {
    const r = await callService("/products");
    relay(res, r);
  } catch (err) {
    next(err);
  }
});

app.get("/api/products/:id", async (req, res, next) => {
  try {
    const r = await callService(`/products/${encodeURIComponent(req.params.id)}`);
    relay(res, r);
  } catch (err) {
    next(err);
  }
});

app.post("/api/products", async (req, res, next) => {
  try {
    const r = await callService("/products", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(req.body ?? {}),
    });
    relay(res, r);
  } catch (err) {
    next(err);
  }
});

// If service-go is unreachable, surface a clean 502 instead of crashing.
app.use((err, _req, res, _next) => {
  console.error("upstream error:", err.message);
  res.status(502).json({ error: "upstream service unavailable", detail: err.message });
});

app.listen(PORT, () => {
  console.log(`api-node listening on :${PORT}, proxying to ${SERVICE_URL}`);
});
