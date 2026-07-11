package models

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataDir string) {
	dbPath := filepath.Join(dataDir, "drop.db")
	uploadsDir := filepath.Join(dataDir, "uploads")

	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS files (
			uuid TEXT PRIMARY KEY,
			original_name TEXT NOT NULL,
			uploaded_by TEXT NOT NULL,
			size INTEGER NOT NULL,
			uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := DB.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	_, err = DB.Exec(`ALTER TABLE files ADD COLUMN last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		log.Printf("[WARN] Failed to add last_used_at column: %v", err)
	}

	DB.Exec(`UPDATE files SET last_used_at = uploaded_at WHERE last_used_at IS NULL`)

	log.Println("[INFO] SQLite DB initialized")
}
