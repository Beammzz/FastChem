# internal/handlers

## Purpose

Gin HTTP handlers and the WebSocket session loops. Translates requests into service calls and database queries, then into JSON or `WSMessage` responses.

## Ownership

- `question.go` — `GET /api/question`, `POST /api/validate`, `POST /api/answer`
- `auth.go` — `POST /api/auth/register`, `POST /api/auth/login`, `GET /api/auth/me`
- `leaderboard.go` — `/api/leaderboard`, `/api/profile/:username`, `/api/scores`, `/api/scores/me`, plus the TTL `leaderboardCache`
- `match.go` — `/api/match/start`, `/api/match/answer`, `/api/match/end`
- `ranked.go` — `/api/ranked/ws`, `/api/ranked/stats`, `/api/ranked/history`, `/api/ranked/leaderboard`
- `room.go` — `/api/room/ws`, room create and join flows
- `match_session.go` — `MatchSessionHandler`, the shared in-match loop used by both ranked and room matches

## Local Contracts

- **`GET /api/question` stores the answer, it does not send it.** `GetQuestion` writes `CorrectIndex` into `GlobalQuestionStore` keyed by question ID and serves a payload with the answer stripped; `POST /api/answer` looks it up and returns it alongside the marking. A handler that hands the client the answer up front makes every anti-cheat check downstream decorative.
- **`POST /api/answer` runs the anti-cheat check itself.** The match and ranked paths get theirs inside `services`, but casual play has no match record, so `SubmitAnswer` builds the `Signal` and keys history on `c.ClientIP()`. A new answer path needs its own call or it is unwatched.
- **Handlers do not register routes.** Constructors (`NewQuestionHandler`, `NewRankedHandler`, …) take their dependencies; `cmd/server/main.go` owns the route table.
- **`match_session.go` is shared by ranked and room play.** Both `ranked.go` and `room.go` drive `HandleMatchSession`; per-mode behavior is injected through the `MatchEndFunc` callback. Fix in-match bugs there once rather than in each caller.
- **WebSocket auth is by query parameter**, validated with `middleware.ParseToken` inside the handler, because the handshake carries no `Authorization` header.
- **The upgrader checks origins against the allow-list** passed in from `config.AllowedOriginsSet()`. Never widen it to accept all origins.
- **One reader per connection.** `StartMsgPump` is the only goroutine calling `ReadMessage`; everything else consumes its `RawMsg` channel. Concurrent readers on a `*websocket.Conn` are a data race.
- **Writes go through `sendWSMessage`.** Do not call `conn.WriteJSON` directly from a second goroutine.
- **Question timeouts are server-side.** `startQuestionTimeout` / `startSyncedQuestionTimeout` end a question regardless of what the client sends; a late or absent answer is scored as a timeout.
- **Both players advance together.** `tryAdvanceBothPlayers` gates the next question, and `notifyOpponentProgress` keeps the other side updated.
- **Passwords are bcrypt-hashed in `auth.go`** and never returned. Login failures give one generic message — do not distinguish "no such user" from "wrong password".
- **The leaderboard is cached with a TTL.** Score submission must not assume its write is immediately visible in `GET /api/leaderboard`.

## Work Guidance

- Bind with `c.ShouldBindJSON` and return `400` with `gin.H{"error": ...}` on failure.
- Read the caller's identity from the context keys set by `AuthRequired()`.
- Use parameterized SQL (`?` placeholders) for every query — no string concatenation.
- Status codes: `400` bad input, `401` unauthenticated, `404` missing resource, `409` duplicate username, `429` rate-limited, `500` unexpected.
- Keep chemistry, scoring, and rating math in `internal/services`; a handler that computes a score is misplaced.
- Always `defer conn.Close()` and stop any timers a session started before returning.

## Verification

```bash
go build ./...
go vet ./...
```

Manual smoke test with the server running:

```bash
curl -s localhost:8080/api/health
curl -s "localhost:8080/api/question?difficulty=easy"
```

WebSocket changes need a real two-client run — open a ranked or room match in two browser sessions and confirm both advance in step and see the same final result.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
