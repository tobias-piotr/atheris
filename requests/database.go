package requests

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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
