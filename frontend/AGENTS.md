# frontend

## Purpose

Next.js 14 App Router client, TypeScript + Tailwind. Builds to a **static export** (`output: "export"`) in `frontend/out/`, which the Go backend serves.

## Ownership

Owned directly by this doc:

- `next.config.mjs` — `output: "export"`, `trailingSlash: true`
- `package.json`, `package-lock.json` — scripts `dev`, `build`, `start`, `lint`
- `tsconfig.json` (the `@/*` → `./src/*` alias), `tailwind.config.ts`, `postcss.config.mjs`, `.eslintrc.json`
- `src/types/index.ts` — the TypeScript mirror of the backend JSON contract
- `src/data/periodicTable.ts` — periodic table reference data for the in-game helper

Delegated to children: `src/app`, `src/components`, `src/hooks`, `src/lib`.

## Local Contracts

- **Static export constraints are absolute.** No server components fetching at request time, no Route Handlers, no `next/image` optimization, no middleware, no server actions. Anything that needs a server belongs in the Go backend. `npm run build` fails loudly if this is violated.
- **`trailingSlash: true` shapes the output.** Each route emits `<route>/index.html`, which is what the backend's static fallback resolves against. Dynamic user routes are handled client-side: `/profile/<username>/` is served by `profile/index.html` (the backend walks up to the nearest parent index) and the page reads the username from `window.location.pathname`, because the export cannot pre-render unknown paths.
- **`src/types/index.ts` must match `backend/internal/models`** field for field, camelCase on both sides. Changing one without the other is an incomplete change.
- **The API base URL comes from `NEXT_PUBLIC_API_URL` and defaults to `""`** — a relative base, correct for single-origin production. It is read only in `src/lib/api.ts`. `NEXT_PUBLIC_*` values are baked in at build time and are public; never put a secret there.
- **The auth token lives in `localStorage` under `token`.** Any code reading it must guard `typeof window !== "undefined"` — the export pre-renders on the server.
- **Import through the `@/` alias** (`@/lib/api`, `@/types`), not deep relative paths.

## Work Guidance

- Function components only; `"use client"` at the top of any file using state, effects, or browser APIs.
- Style with Tailwind utility classes in the markup. `src/app/globals.css` is for base layers and shared primitives, not per-component rules.
- User-facing copy is Thai (`<html lang="th">`); write new strings to match.
- Never recompute authoritative scores or reveal a correct answer client-side — render what the API returns.
- Keep game state in the hooks under `src/hooks/`; components stay presentational.

## Verification

```bash
npm run lint
npm run build
```

`npm run build` type-checks and produces `out/`; it is the gate that must pass before any frontend change is considered done. Serve the full stack with `../run.sh` to check the export against the real backend.

## Child DOX Index

- `src/app/AGENTS.md` — routes, layout, global styles
- `src/components/AGENTS.md` — presentational and provider components
- `src/hooks/AGENTS.md` — game state and WebSocket session hooks
- `src/lib/AGENTS.md` — API client and browser helpers
