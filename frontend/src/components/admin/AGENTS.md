# src/components/admin

## Purpose

The operator console's UI. One panel per tab of `app/admin/page.tsx`, plus the
shared primitives they are built from. Nothing here is used outside `/admin`.

## Ownership

- `ui.tsx` — shared primitives and formatters: `Panel`, `Stat`, `Pill`, `Empty`, `Spinner`, `Table`, the two chart forms `BarList` and `Columns`, the `SERIES` / `MAGNITUDE` colour constants, and `formatNumber` / `formatBytes` / `formatDuration` / `formatDateTime`
- `AnalyticsPanel.tsx` — headline stats, the daily series, distribution charts, per-topic table, process health
- `UsersPanel.tsx` — the user table, its sort and search controls, and the edit dialog (rating, points, password, admin flag, delete)
- `AnticheatPanel.tsx` — rule editors and the filtered findings log with breakdowns
- `LivePanel.tsx` — queue, matches in progress, open rooms, recent finished matches
- `LeaderboardPanel.tsx` — both boards, uncached
- `QuestionsPanel.tsx` — topic browser and question preview

## Local Contracts

- **`ui.tsx` is the exception to one-default-export-per-file** that the parent doc sets. None of its pieces is a screen, and splitting seven primitives across seven files would cost more than it explains. Panels keep the default-export rule.
- **Panels are presentational.** They take data, `loading`, and callbacks as props; `app/admin/page.tsx` owns every fetch, every timer, and all tab state. A panel that calls `@/lib/api` is misplaced — `QuestionsPanel` takes even its preview call as an `onPreview` prop.
- **Colour follows the entity, never its rank.** `SERIES` is a fixed order: a chart assigns `SERIES[0]` to its first category and counts up, never cycling for a 6th series and never repainting survivors when a filter changes the count. The five steps are validated against the `#12122a` panel surface for lightness, chroma, contrast and colour-vision separation. Substituting a step by eye breaks that.
- **One measure per chart.** The daily series is four small multiples, not one chart with four lines, because the measures have different units. A second y-axis is never the answer.
- **Every mark is direct-labelled or hoverable.** `BarList` labels every row; `Columns` labels the peak and the latest value, carries a `title` per column, and repeats the series in an `sr-only` list. Identity never rests on colour alone.
- **Destructive actions are typed-confirmation only.** Deleting an account needs the username typed back, and the dialog states that ranked matches disappear from the opponent's history too. Self-demotion and self-deletion are disabled here as well as rejected by the server.
- **The admin flag is a rendering hint.** `user.isAdmin` decides what to draw; it never stands in for authorisation. Every route re-checks server-side, so a forged flag reveals nothing.
- **Tables scroll themselves.** Wide content goes inside `Table`, which owns the `overflow-x-auto` container — the page body must never scroll sideways on a phone.

## Work Guidance

- Thai copy, matching the rest of the app.
- A new tab is a panel here, a `Tab` entry and its fetch in `app/admin/page.tsx`, an API function in `@/lib/api`, and types in `@/types` mirroring `backend/internal/models/admin.go`.
- Reach for `Panel` + `Stat` + `Table` before writing new layout; reach for `BarList` / `Columns` before writing a new chart.
- Show what the server actually returned. An empty list is `Empty`, never a fabricated zero row.

## Verification

```bash
npm run lint
npm run build
```

Then load `/admin` against a running backend as a promoted account, and check
each tab at phone width — the user and findings tables are the first things to
overflow. Confirm a non-admin account sees the locked state rather than an
empty console.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
