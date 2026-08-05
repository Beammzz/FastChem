# src/lib

## Purpose

Cross-cutting helpers: the typed HTTP client and small browser utilities.

## Ownership

- `api.ts` — every REST call to the Go backend, the `API_BASE` constant, the `authHeaders()` helper, and the WebSocket URL builders `getRankedWebSocketURL()` / `getRoomWebSocketURL(action, code?)`
- `haptic.ts` — `triggerHaptic(type)`, a wrapper over the Web Vibration API with `correct` / `wrong` / `timeout` patterns

## Local Contracts

- **`api.ts` is the only place `fetch` is called.** Components, pages, and hooks import functions from here; a raw `fetch` anywhere else is a bug.
- **WebSocket URLs are built here, but connections are not opened here.** The `get*WebSocketURL` helpers derive the origin from `window.location`, pick `wss:` on HTTPS, and append the token as a query parameter. Opening, handling, and closing the socket belongs to `@/hooks`.
- **`API_BASE = process.env.NEXT_PUBLIC_API_URL || ""`** — the empty default is a relative base, correct for the single-origin production deployment. This constant is defined once, here.
- **Every function is typed against `@/types`** for both request and response. No `any`, and no inline response shapes.
- **Auth headers come from `authHeaders()`**, which reads the `token` from `localStorage` behind a `typeof window !== "undefined"` guard and returns `{}` when absent.
- **Non-OK responses throw an `Error`.** Callers handle failure; these functions never return a partial or fabricated result on error.
- **`haptic` must degrade silently** where `navigator.vibrate` is unavailable — most desktop browsers and iOS Safari.

## Work Guidance

- A new endpoint means: a function here, its types in `@/types/index.ts`, and the matching route in `backend/cmd/server/main.go`.
- Keep function names aligned with the endpoint (`fetchQuestion`, `submitAnswer`, `startMatch`) and group them under the existing section comments.
- Send `cache: "no-store"` on anything that must be fresh — question fetches already do.
- Keep helpers here dependency-free and framework-agnostic; anything needing React state belongs in `@/hooks`.

## Verification

```bash
npm run lint
npm run build
```

Exercise a new call against the running backend (`../../run.sh`) and confirm both the success path and the thrown-error path.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
