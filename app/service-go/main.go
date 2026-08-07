// Command service-go is the data-owning backend of the o11y-lab app.
//
// It exposes a tiny HTTP API over a "products" table in Postgres, with a
// read-through cache in Redis. There is deliberately NO observability
// instrumentation here yet — that arrives in Stage 2. For now the goal is a
// working service that does real work (a DB query, a cache lookup) so that,
// once instrumented, its traces and metrics are actually interesting.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Product is a single catalog item. Price is stored in integer cents to avoid
// floating-point money bugs.
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var (
	db  *pgxpool.Pool
	rdb *redis.Client
)

func main() {
	ctx := context.Background()

	// Start tracing first so everything below is instrumented.
	shutdown, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("init tracer: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	// Attach the otelpgx tracer to the pool so every query becomes a span.
	cfg, err := pgxpool.ParseConfig(env("DATABASE_URL", "postgres://o11y:o11y@localhost:5432/o11y?sslmode=disable"))
	if err != nil {
		log.Fatalf("parse db config: %v", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()
	db = pool

	rdb = redis.NewClient(&redis.Options{Addr: env("REDIS_ADDR", "localhost:6379")})
	// Emit a span for every Redis command.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatalf("instrument redis: %v", err)
	}

	if err := ensureSchema(ctx); err != nil {
		log.Fatalf("schema: %v", err)
	}

	// Go 1.22's ServeMux understands method + path patterns, so we get simple
	// routing (including path params like {id}) without any third-party router.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /products", handleList)
	mux.HandleFunc("GET /products/{id}", handleGet)
	mux.HandleFunc("POST /products", handleCreate)

	// otelhttp creates a server span per request AND extracts the trace context
	// propagated from api-node, so both services end up in ONE trace.
	handler := otelhttp.NewHandler(mux, "service-go")

	addr := ":" + env("PORT", "8080")
	log.Printf("service-go listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// ensureSchema creates the products table (idempotently) and seeds a few rows
// on first run. It retries because, on a fresh `docker compose up`, Postgres may
// still be starting even after its healthcheck flips.
func ensureSchema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    name  TEXT    NOT NULL,
    price INTEGER NOT NULL
);`

	var err error
	for i := 0; i < 10; i++ {
		if _, err = db.Exec(ctx, ddl); err == nil {
			break
		}
		log.Printf("waiting for postgres... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}

	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec(ctx,
			`INSERT INTO products (name, price) VALUES ($1,$2),($3,$4),($5,$6)`,
			"Keyboard", 4999, "Mouse", 2999, "Monitor", 19999)
	}
	return err
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "service-go"})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(r.Context(), `SELECT id, name, price FROM products ORDER BY id`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
		products = append(products, p)
	}
	writeJSON(w, http.StatusOK, products)
}

// handleGet is a read-through cache: try Redis first, fall back to Postgres on a
// miss, then populate the cache. The X-Cache header exposes which path was taken
// — handy for seeing cache behaviour, and later for correlating with traces.
func handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cacheKey := "product:" + id

	if cached, err := rdb.Get(r.Context(), cacheKey).Result(); err == nil {
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(cached))
		return
	}

	var p Product
	err := db.QueryRow(r.Context(), `SELECT id, name, price FROM products WHERE id=$1`, id).
		Scan(&p.ID, &p.Name, &p.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, errors.New("product not found"))
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	if b, err := json.Marshal(p); err == nil {
		// Cache for 60s. Errors here are non-fatal: a cache miss just costs a DB hit.
		rdb.Set(r.Context(), cacheKey, b, 60*time.Second)
	}
	w.Header().Set("X-Cache", "MISS")
	writeJSON(w, http.StatusOK, p)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var in Product
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if in.Name == "" {
		writeErr(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}
	err := db.QueryRow(r.Context(),
		`INSERT INTO products (name, price) VALUES ($1,$2) RETURNING id`, in.Name, in.Price).
		Scan(&in.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, in)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
