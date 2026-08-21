# src/app

## Purpose

App Router route tree. Every folder with a `page.tsx` is a URL; the export writes each one as `<route>/index.html`.

## Ownership

- `layout.tsx` — root layout: `lang="th"`, the `Inter` font from `next/font/google`, Thai metadata, and the `AuthProvider` wrapper
- `globals.css` — Tailwind layers and base styles
- `favicon.ico`
- `fonts/` — `GeistVF.woff`, `GeistMonoVF.woff`, unused leftovers from `create-next-app`; nothing imports them

Routes:
- `page.tsx` — landing page
- `play/page.tsx` — mode picker
- `play/singleplayer/page.tsx` — solo game
- `play/ranked/page.tsx` — matchmaking queue and ranked match
- `play/room/page.tsx` — custom room create/join
- `game/page.tsx` — legacy single-player game screen
- `login/page.tsx` — login and registration
- `leaderboard/page.tsx` — global and ranked leaderboards
- `profile/page.tsx` — thin wrapper that parses the username out of the path and renders `ProfileContent`
- `admin/page.tsx` — the operator console: tab state, every admin fetch, and the panels from `@/components/admin`

## Local Contracts

- **Every interactive page is a client component** (`"use client"`). The static export has no request-time server rendering.
- **No dynamic route segments.** `output: "export"` cannot pre-render unknown paths, so `profile/page.tsx` reads the username from `window.location.pathname` inside a `useEffect` and renders a spinner until it resolves. The URL `/profile/<username>/` works because the backend's static fallback walks up to `profile/index.html`. A new user-specific view follows the same pattern — never add a `[param]` folder.
- **`layout.tsx` renders `AuthProvider` only.** `Navbar` is imported per page, so a new page that needs it must render it itself.
- **Pages compose, they do not implement.** Game logic comes from `@/hooks`, rendering from `@/components`, network calls from `@/lib/api`. A page that owns a `useEffect` polling loop or a raw `fetch` is misplaced.
- **Auth state comes from `useAuth()`**, never by reading `localStorage` in a page.
- **`admin/page.tsx` is the one page that owns fetching,** because its panels are presentational by contract. Each tab fetches only when selected, so opening the console does not pull every table at once; the live tab polls on a 5-second interval that is cleared on unmount. The `user.isAdmin` check here decides what to render — authorisation is the server's, on every request.
- **UI copy is Thai.** `<html lang="th">` and page metadata are Thai; new user-facing strings match.
- Navigate with `next/link` and `useRouter` from `next/navigation`; `trailingSlash: true` means internal hrefs should end in `/`.

## Work Guidance

- Adding a route: create `<name>/page.tsx`, mark it `"use client"` if it holds state, and add navigation from `play/page.tsx` or `Navbar` so it is reachable.
- Guard authenticated pages by checking `useAuth()` and redirecting to `/login/`.
- Keep pages thin — layout, routing, and wiring only.

## Verification

```bash
npm run lint
npm run build
```

Confirm the new route appears as `out/<route>/index.html`, then load it through the backend (`../../run.sh`) so the static fallback path is exercised, not just `next dev`.

## Child DOX Index

- No child AGENTS.md files. Route folders hold a single `page.tsx` each and inherit this doc.
