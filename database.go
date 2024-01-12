package main

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func GetConnection(filename string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())
		panic(err)
	}
	return db
}

func Migrate(db *sqlx.DB) {
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

func InsertRequest(db *sqlx.DB, req RequestData) {
	res, err := json.Marshal(req.Response)
	if err != nil {
		slog.Error("Error marshalling response", "error", err.Error())
	}

	_, err = db.Exec(`
	INSERT INTO requests (id, path, response)
	VALUES ($1, $2, $3);
	`, uuid.New(), req.Path, res)
	if err != nil {
		slog.Error("Error inserting request", "error", err.Error())
	}
}

func GetRequests(db *sqlx.DB) ([]Request, error) {
	query := `SELECT id, created_at, path, response FROM requests;`
	var reqs []Request
	err := db.Select(&reqs, query)
	return reqs, err
}
