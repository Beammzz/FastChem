# src/hooks

## Purpose

Game state. Each hook owns one play mode's lifecycle and hands components a ready-to-render view plus actions.

## Ownership

- `useGame.ts` — single player: question fetch, answer submission, timer, combo, score, game over
- `useRankedGame.ts` — ranked 1v1 over `/api/ranked/ws`: queue, match found, per-question flow, opponent progress, result
- `useRoomGame.ts` — custom rooms over `/api/room/ws`: create/join by code, lobby, match; exports the `RoomPhase` union

## Local Contracts

- **One hook per mode, one WebSocket per hook.** The socket is opened in a `useEffect` and closed in its cleanup. Never leave a connection open across unmount, and never share one socket between hooks.
- **State machines are explicit.** Phase unions such as `RoomPhase` drive rendering; components switch on the phase rather than deriving it from a pile of booleans.
- **The server is authoritative.** Scores, correctness, combo multipliers, timeouts, and match end all come from the API or a `WSMessage`. Local timers exist for display and must not decide the outcome.
- **A failed submit leaves the question unresolved.** The answer is not in the question payload, so there is nothing to mark against offline: `useGame` sets `isCorrect` to null, leaves `revealedIndex` null, and moves on without touching the score. Do not reintroduce a local-scoring fallback — it would disagree with what the server recorded, which `endMatch` then contradicts at the results screen.
- **Handle every `WSEventType` the backend can send** — `MATCH_FOUND`, `MATCH_START`, `QUESTION_START`, `MATCH_END`, `ERROR`, `PING`/`PONG`, `QUEUE_JOINED`, plus the room events. An unhandled event silently stalls the match.
- **Answer `PING` with `PONG`** to keep the connection alive.
- **URLs come from `@/lib/api`, lifecycle stays here.** `getRankedWebSocketURL()` and `getRoomWebSocketURL(action, code?)` build the URL — origin-derived, `wss:` when the page is HTTPS, token in the query string because the handshake cannot send an `Authorization` header. Hooks own opening, message handling, and closing; they do not assemble URLs.
- **Guard browser APIs.** `window`, `localStorage`, and `WebSocket` are touched only inside effects or behind `typeof window !== "undefined"`.
- **Return a stable object** of state plus memoized callbacks (`useCallback`) so consumers do not re-render on every tick.

## Work Guidance

- HTTP calls go through `@/lib/api`; only the socket lifecycle and message handling live here.
- Clear every timer and interval in the effect cleanup — an orphaned countdown outlives the match and corrupts the next one.
- On disconnect, surface the state to the UI rather than reconnecting silently; the backend has its own disconnect timeout and force-win path.
- A new mode gets a new hook rather than another branch inside an existing one.

## Verification

```bash
npm run lint
npm run build
```

Multiplayer changes need two clients against a running backend: confirm both advance on the same question, the score matches the server's `MATCH_END` payload, and closing one tab resolves the other cleanly.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
