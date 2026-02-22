package database

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init opens the SQLite database and runs migrations.
// dbPath is the path to the SQLite file (e.g. from config).
func Init(dbPath string) {
	if dbPath == "" {
		dbPath = "fastchem.db"
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("failed to create database directory", "error", err)
			os.Exit(1)
		}
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	// Enable WAL mode and foreign keys
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA foreign_keys=ON")

	// SQLite only supports a single writer. Limit open connections to avoid
	// SQLITE_BUSY errors under concurrent writes.
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)

	if err := DB.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	migrate()
	slog.Info("database initialized successfully")
}

func migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			total_points INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS scores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			score INTEGER NOT NULL,
			total_answered INTEGER NOT NULL,
			correct_answers INTEGER NOT NULL,
			difficulty TEXT NOT NULL DEFAULT 'easy',
			time_limit INTEGER NOT NULL,
			time_spent REAL NOT NULL DEFAULT 0,
			played_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			difficulty TEXT NOT NULL,
			total_score INTEGER NOT NULL DEFAULT 0,
			best_combo INTEGER NOT NULL DEFAULT 0,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			ended_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS question_attempts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			match_id INTEGER NOT NULL,
			question_snapshot TEXT NOT NULL DEFAULT '{}',
			correct BOOLEAN NOT NULL DEFAULT 0,
			timed_out BOOLEAN NOT NULL DEFAULT 0,
			time_spent_ms INTEGER NOT NULL DEFAULT 0,
			score_awarded INTEGER NOT NULL DEFAULT 0,
			combo_at_answer INTEGER NOT NULL DEFAULT 0,
			combo_multiplier REAL NOT NULL DEFAULT 1.0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (match_id) REFERENCES matches(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_scores_score ON scores(score DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_scores_user_id ON scores(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_scores_played_at ON scores(played_at)`,
		`CREATE INDEX IF NOT EXISTS idx_users_total_points ON users(total_points DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_matches_user_id ON matches(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_matches_ended_at ON matches(ended_at)`,
		`CREATE INDEX IF NOT EXISTS idx_question_attempts_match_id ON question_attempts(match_id)`,

		// ─── Ranked Match Tables ─────────────────────────────────
		`CREATE TABLE IF NOT EXISTS ranked_matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player1_id INTEGER NOT NULL,
			player2_id INTEGER NOT NULL,
			seed INTEGER NOT NULL,
			player1_score INTEGER NOT NULL DEFAULT 0,
			player2_score INTEGER NOT NULL DEFAULT 0,
			winner_id INTEGER,
			status TEXT NOT NULL DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			finished_at DATETIME,
			FOREIGN KEY (player1_id) REFERENCES users(id),
			FOREIGN KEY (player2_id) REFERENCES users(id),
			FOREIGN KEY (winner_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS ranked_question_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ranked_match_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			question_index INTEGER NOT NULL,
			difficulty TEXT NOT NULL,
			topic TEXT NOT NULL DEFAULT '',
			correct BOOLEAN NOT NULL DEFAULT 0,
			time_spent REAL NOT NULL DEFAULT 0,
			score_earned INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (ranked_match_id) REFERENCES ranked_matches(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ranked_matches_player1 ON ranked_matches(player1_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ranked_matches_player2 ON ranked_matches(player2_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ranked_matches_status ON ranked_matches(status)`,
		`CREATE INDEX IF NOT EXISTS idx_ranked_question_results_match ON ranked_question_results(ranked_match_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ranked_question_results_user ON ranked_question_results(user_id)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			slog.Error("migration failed", "error", err, "query", q)
			os.Exit(1)
		}
	}

	// Add difficulty column if it doesn't exist (migration for existing databases)
	DB.Exec("ALTER TABLE scores ADD COLUMN difficulty TEXT NOT NULL DEFAULT 'easy'")

	// Add ranked columns to users table if they don't exist
	DB.Exec("ALTER TABLE users ADD COLUMN rating INTEGER NOT NULL DEFAULT 1200")
	DB.Exec("ALTER TABLE users ADD COLUMN ranked_wins INTEGER NOT NULL DEFAULT 0")
	DB.Exec("ALTER TABLE users ADD COLUMN ranked_losses INTEGER NOT NULL DEFAULT 0")
	DB.Exec("ALTER TABLE users ADD COLUMN highest_rating INTEGER NOT NULL DEFAULT 1200")
}

// Close shuts down the database connection.
func Close() {
	if DB != nil {
		DB.Close()
	}
}
