# src/components

## Purpose

Reusable UI. Presentational components plus one context provider.

## Ownership

Provider:
- `AuthProvider.tsx` — auth context and the `useAuth()` hook; the only component here that owns global state

Chrome and shared UI:
- `Navbar.tsx`, `ScoreBoard.tsx`, `TimerBar.tsx`, `QuestionCard.tsx`

In-game helpers:
- `PeriodicTable.tsx` — toggleable reference table, reads `@/data/periodicTable`
- `Calculator.tsx` — in-game calculator

Mode-specific UI:
- `RankedQueue.tsx`, `RankedMatchUI.tsx`, `RankedGameOver.tsx`
- `RoomGameOver.tsx`
- `GameOver.tsx` — single player
- `ProfileContent.tsx` — profile view rendered by `app/profile/page.tsx`

## Local Contracts

- **Props in, callbacks out.** Apart from `AuthProvider`, components take data and handlers as props and do not call `@/lib/api` or hold game state. Fetching belongs in `@/hooks` or the page.
- **`AuthProvider` is the single source of auth truth.** It owns the `localStorage` `token` read/write; every consumer uses `useAuth()`. Do not touch `localStorage` for auth elsewhere.
- **Never render or infer the correct answer before the server returns it.** `QuestionCard` shows a choice as correct only from the API's answer result.
- **Default-export one component per file**, named to match the filename. `AuthProvider.tsx` is the exception: it exports both `AuthProvider` and `useAuth` as named exports.
- Prop shapes are declared as a local `interface <Name>Props` above the component.
- Every file here is a client component (`"use client"`).

## Work Guidance

- Style with Tailwind classes; keep dark backgrounds and the existing color language consistent across modes.
- Use `triggerHaptic("correct" | "wrong" | "timeout")` from `@/lib/haptic` for answer feedback, matching how the existing game UI uses it.
- Timers display state — the authoritative countdown is the server's; a component running down to zero must not itself decide the question was missed.
- Before adding a mode-specific variant, check whether an existing `*GameOver` component generalizes with a prop.

## Verification

```bash
npm run lint
npm run build
```

Check rendering at mobile width — the game screens are played on phones, and the periodic table and calculator overlays are the first things to break.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
