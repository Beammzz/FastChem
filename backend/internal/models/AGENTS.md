# internal/models

## Purpose

Plain data structs shared by handlers, services, and the database layer. Defines the wire format for every API response and the WebSocket protocol.

## Ownership

- `question.go` — `Question`, validate/answer request-response pairs, `Element` (including its real `Isotopes` list), `Compound`, `CompoundElement`
- `user.go` — `User`, `RegisterRequest`, `LoginRequest`, `AuthResponse`, `UserPublic`
- `score.go` — `Score`, `SubmitScoreRequest`, `LeaderboardEntry`, `UserStats`, `ProfileResponse`
- `match.go` — single-player match session types, `DefaultQuestionsPerMatch = 10`
- `ranked.go` — ranked match types, `RankedLeaderboardEntry`, the `WSEventType` constants, and every WebSocket payload struct
- `admin.go` — the read models behind `/api/admin/*`: overview, analytics, user rows, findings, rules, live state, recent matches, topics

## Local Contracts

- **This package is the JSON contract.** Every field carries an explicit camelCase `json:` tag, and `frontend/src/types/index.ts` mirrors these structs. Renaming a field or changing its type is a two-sided change.
- **Never expose a password hash.** `User.PasswordHash` is tagged `json:"-"`; user-facing responses use `UserPublic`.
- **`UserPublic.IsAdmin` is set only on `/api/auth/me`,** where the caller is the account itself. It is `omitempty` and tells the client whether to offer the console; it authorises nothing, because every admin route re-reads the flag from the database.
- **The `Admin*` types aggregate across tables on purpose.** They are read models for one operator screen, not the shape of any single row, and nothing but `handlers/admin.go` writes them. They are still the JSON contract and still mirrored in `frontend/src/types/index.ts`.
- **`AdminQuestionPreview` carries `correctIndex` where `Question` hides it.** That is not an exception to the rule — it is a different object. It describes a question generated for inspection and thrown away, never one a player is answering.
- **Never expose an unanswered question's answer.** `Question.CorrectIndex` is tagged `json:"-"` for the same reason: the struct goes out before the player has chosen, so a serialized answer is readable from the network tab. `AnswerResponse`, `MatchAnswerResponse` and `AnswerResultPayload` carry `correctIndex` because they are sent *after* the answer is in — that is where the reveal belongs. `TestQuestionJSONOmitsTheAnswer` pins both halves.
- **No logic and no imports from sibling packages.** Models depend only on the standard library, which keeps `services` and `handlers` free to import them without cycles.
- **Match shape constants live here**, not in services: `DefaultQuestionsPerMatch = 10`, `RankedQuestionsPerMatch = 10`, and the ranked difficulty split `RankedEasyCount = 4`, `RankedMediumCount = 3`, `RankedHardCount = 3`. The counts must sum to `RankedQuestionsPerMatch`.
- **ELO constants live here:** `DefaultRating = 1200`, `EloKFactor = 24`. `services/rating.go` consumes them; `database` defaults the `rating` column to the same 1200.
- **The WebSocket envelope is `WSMessage{Type, Payload}`.** Every message on `/api/ranked/ws` and `/api/room/ws` uses it. Adding an event means adding a `WSEventType` constant here, a payload struct here, and the matching TypeScript type in the frontend.

## WebSocket events

`MATCH_FOUND`, `MATCH_START`, `QUESTION_START`, `SUBMIT_ANSWER`, `MATCH_END`, `ERROR`, `PING`, `PONG`, `QUEUE_JOINED` — plus the room payloads `RoomCreatedPayload` and `RoomJoinedPayload`.

## Work Guidance

- Use `int64` for database IDs and pointer types (`*int64`, `*time.Time`) for nullable columns.
- Put binding rules on request structs with `binding:"..."` tags so Gin validates before handler logic runs.
- Times are `time.Time`; durations measured in the game are `float64` seconds, except `time_spent_ms` on attempts.

## Verification

```bash
go build ./...
```

After changing a struct, confirm the mirrored interface in `frontend/src/types/index.ts` and run the frontend's `npm run build`.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
