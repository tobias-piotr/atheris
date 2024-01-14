package libs

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(filename string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())
		panic(err)
	}
	return db
}

func MigrateDB(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS requests (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		path VARCHAR(255) NOT NULL,
		response JSONB
	);
	`)
	if err != nil {
		slog.Error("Error creating tables", "error", err.Error())
		panic(err)
	}
}