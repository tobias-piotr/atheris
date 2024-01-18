package libs

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(filename string) *sqlx.DB {
	slog.Info("Connecting to database", "filename", filename)

	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())
		panic(err)
	}
	return db
}

func MigrateDB(db *sqlx.DB) {
	slog.Info("Migrating database")

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS requests (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name VARCHAR(255),
		method VARCHAT(15) NOT NULL,
		path VARCHAR(255) NOT NULL,
		response JSONB
	);
	`)
	if err != nil {
		slog.Error("Error creating tables", "error", err.Error())
		panic(err)
	}
}
