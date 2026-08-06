# backend

## Purpose

Go 1.24 + Gin server (module `github.com/takumi/fastchem`). Serves the JSON API under `/api`, the ranked and room WebSocket endpoints, and — in production — the static frontend export.

## Ownership

Owned directly by this doc:

- `cmd/server/main.go` — the only wiring point: config load, JWT init, DB init, background cleanup goroutines, middleware order, route table, static-file fallback, graceful shutdown
- `internal/config/config.go` — every environment variable read, with defaults
- `go.mod` / `go.sum`

Delegated to children: `internal/database`, `internal/middleware`, `internal/models`, `internal/services`, `internal/handlers`.

## Local Contracts

- **Startup order in `main.go` is load-bearing:** `middleware.InitJWT` runs before any handler, `database.Init` before any query. Both call `os.Exit(1)` on failure — the server never runs half-initialized.
- **Route registration is centralized.** Handlers never register their own routes. Every path is visible in `main.go`'s `api := router.Group("/api")` block.
- **Middleware order:** `RequestID` → CORS → per-group rate limiter → `AuthRequired`. Adding middleware means placing it deliberately in that chain.
- **Dependencies are constructor-injected.** Services are built in `main.go` and passed to `handlers.New*`. Package-level singletons are limited to `services.GlobalQuestionStore`, `services.GlobalMatchStore`, and `database.DB`.
- **Background cleanup goroutines** live in `main.go` on 5-minute tickers: expired questions (10 min), stale matches (15 min), stale ranked matches (30 min), stale rooms (15 min). A new expiring store needs its sweeper added here.
- **CGO is required** — `mattn/go-sqlite3`. Builds must keep `CGO_ENABLED=1`, and the runtime image needs `libsqlite3-0`.
- **Static serving is a fallback layer**, registered after the API routes. It skips `/api/` prefixes and resolves clean URLs against `FRONTEND_DIR` (`<exe>/../frontend/out` when unset). Dynamic client routes such as `/profile/<username>` resolve by walking up to the nearest parent `index.html`.

## Work Guidance

- Log with `log/slog`, structured key-value pairs; no `fmt.Println`.
- Return errors to clients as `gin.H{"error": "..."}` with an explicit status code; never leak SQL or internal error text.
- Read the authenticated user from the context keys set by `middleware.AuthRequired()` — do not re-parse the token in handlers.
- Guard shared state with the mutex already present on the owning struct; do not add a second lock over the same data.

## Verification

```bash
go build ./...
go vet ./...
gofmt -l .        # must print nothing
go test ./... -race
```

Tests currently cover `internal/services` question generation only; other packages have none yet. Add `_test.go` files alongside the code they cover and extend this section as coverage grows.

## Child DOX Index

- `internal/database/AGENTS.md` — SQLite connection, schema, migrations
- `internal/middleware/AGENTS.md` — JWT auth, rate limiting, request IDs
- `internal/models/AGENTS.md` — structs, JSON contract, WebSocket protocol constants
- `internal/services/AGENTS.md` — question generation, scoring, ELO, matchmaking, in-memory stores
- `internal/handlers/AGENTS.md` — HTTP and WebSocket request handling
