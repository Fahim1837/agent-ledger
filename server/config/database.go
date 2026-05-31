package config

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const DefaultDBPath = "data/agent-ledger.db"

func DatabasePath() string {
	if path := os.Getenv("AGENT_LEDGER_DB_PATH"); path != "" {
		return path
	}

	return DefaultDBPath
}

func OpenSQLite(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := configureSQLite(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := migrateSQLite(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func configureSQLite(db *sql.DB) error {
	statements := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func migrateSQLite(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	agent_name TEXT NOT NULL,
	model TEXT,
	started_at TEXT NOT NULL DEFAULT (datetime('now')),
	ended_at TEXT
);

CREATE TABLE IF NOT EXISTS token_usage (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	session_id INTEGER,
	prompt_tokens INTEGER NOT NULL DEFAULT 0,
	completion_tokens INTEGER NOT NULL DEFAULT 0,
	total_tokens INTEGER NOT NULL DEFAULT 0,
	estimated_cost_cents INTEGER NOT NULL DEFAULT 0,
	recorded_at TEXT NOT NULL DEFAULT (datetime('now')),
	FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_token_usage_recorded_at
ON token_usage(recorded_at);
`

	_, err := db.Exec(schema)
	return err
}
