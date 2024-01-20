package libs

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func GetDBConnection(filename string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %s\n", err.Error())
		os.Exit(1)
	}
	return db
}

func MigrateDB(db *sqlx.DB) {
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
		fmt.Fprintf(os.Stderr, "Error creating table: %s\n", err.Error())
		os.Exit(1)
	}
}
