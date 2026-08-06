# internal/database

## Purpose

Owns the SQLite connection and the entire schema. `db.go` is the only file here.

## Ownership

- `database.DB` — the process-wide `*sql.DB`, opened by `Init(dbPath)` and closed by `Close()`
- The `migrate()` function — the single source of truth for every table, index, and column in FastChem

## Local Contracts

- **All schema lives in `migrate()`.** There are no migration files or a version table. Statements are `CREATE TABLE IF NOT EXISTS` / `CREATE INDEX IF NOT EXISTS`, run in order on every boot, so they must stay idempotent.
- **Adding a column to an existing table** uses a bare `DB.Exec("ALTER TABLE ... ADD COLUMN ...")` after the main loop, with its error deliberately ignored — the statement fails harmlessly when the column already exists. Keep new column additions in that trailing block, not in the `queries` slice.
- **Failures in the `queries` loop are fatal** (`slog.Error` + `os.Exit(1)`). Only add a statement there if the server genuinely cannot run without it.
- **Connection limits are intentional:** `SetMaxOpenConns(2)` / `SetMaxIdleConns(2)`, WAL journal, `_busy_timeout=5000`, `foreign_keys=ON`. SQLite has one writer; raising these reintroduces `SQLITE_BUSY`.
- `Init` creates the parent directory of `dbPath`, which is what makes `DB_PATH=/data/fastchem.db` work against a mounted volume.

## Tables

- `users` — credentials, `total_points`, plus ranked columns added by `ALTER`: `rating` (default 1200), `ranked_wins`, `ranked_losses`, `highest_rating`
- `scores` — one row per finished single-player game
- `matches` / `question_attempts` — server-scored single-player match sessions
- `ranked_matches` / `ranked_question_results` — 1v1 results and per-question breakdown

## Work Guidance

- Queries live in the packages that use them (`handlers`, `services`); this package exposes the handle, not a query layer.
- Add an index alongside any new column that will be filtered or ordered on.
- Keep foreign keys declared — `PRAGMA foreign_keys=ON` is enabled.

## Verification

```bash
go build ./...
```

Schema changes are exercised by starting the server against a fresh `DB_PATH` and confirming it reaches `database initialized successfully`, and against an existing DB file to confirm the `ALTER` path is still harmless.

## Child DOX Index

- No child AGENTS.md files. This folder is a single file.
