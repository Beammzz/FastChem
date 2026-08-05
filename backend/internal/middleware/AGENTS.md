# internal/middleware

## Purpose

Gin middleware: JWT authentication, token issuing, rate limiting, and request IDs.

## Ownership

- `auth.go` — `InitJWT`, `GenerateToken`, `ParseToken`, `AuthRequired()`, `OptionalAuth()`, `GetJWTSecret()`
- `ratelimit.go` — `RateLimiter` (per user ID) and `IPRateLimiter` (per client IP)
- `requestid.go` — `RequestID()` and the `X-Request-Id` header

## Local Contracts

- **`InitJWT(secret, production)` must run before any request is served.** It is called once from `cmd/server/main.go`. In production (`GIN_MODE=release`) an empty `JWT_SECRET` is fatal; outside production it falls back to a development secret.
- **Tokens expire after `TokenExpiry` (7 days).** Change that constant here, not in handlers.
- **`AuthRequired()` is the only place that rejects a request for missing or invalid auth.** It parses `Authorization: Bearer <token>` and sets the user ID and username on the Gin context; downstream handlers read those keys and never re-parse the token.
- **WebSocket endpoints do not use `AuthRequired()`** — browsers cannot set headers on a WS handshake. `/api/ranked/ws` and `/api/room/ws` authenticate from a query parameter by calling `ParseToken` directly in the handler. Keep both paths in sync when token handling changes.
- **`OptionalAuth()` never rejects.** It populates the context when a valid token is present and continues otherwise.
- **Rate limiters are token buckets** constructed in `main.go` with `(rate, capacity, interval)`. Both types run their own background `cleanup()` to drop idle buckets; a new limiter type must do the same or it leaks memory per user/IP.

## Work Guidance

- Reject with a status code and `gin.H{"error": ...}`, then `c.Abort()` — never continue the chain after writing an error.
- Rate limit responses use `429`; auth failures use `401`.
- Keep middleware free of database access; it operates on tokens and in-memory buckets only.

## Verification

```bash
go build ./...
go vet ./...
```

Manual check: a request without a token to an `AuthRequired()` route returns 401; bursting `/api/match/answer` past 3 requests per second returns 429.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
