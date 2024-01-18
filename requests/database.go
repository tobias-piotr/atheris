package requests

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: Move error handling up to the caller

func InsertRequest(db *sqlx.DB, req RequestData) {
	res, err := json.Marshal(req.Response)
	if err != nil {
		slog.Error("Error marshalling response", "error", err.Error())
	}

	_, err = db.Exec(`
	INSERT INTO requests (id, method, path, response)
	VALUES ($1, $2, $3, $4);
	`, uuid.New(), req.Method, req.Path, res)
	if err != nil {
		slog.Error("Error inserting request", "error", err.Error())
	}
}

func GetRequests(db *sqlx.DB) ([]Request, error) {
	query := `
	SELECT id, created_at, name, method, path, response
	FROM requests;
	`
	var reqs []Request
	err := db.Select(&reqs, query)
	return reqs, err
}

func SetRequestName(db *sqlx.DB, id uuid.UUID, name string) error {
	_, err := db.Exec(`
	UPDATE requests
	SET name = $1
	WHERE id = $2;
	`, name, id)
	if err != nil {
		slog.Error("Error setting request name", "error", err.Error())
	}
	return err
}

func DeleteRequest(db *sqlx.DB, id uuid.UUID) error {
	_, err := db.Exec(`
	DELETE FROM requests
	WHERE id = $1;
	`, id)
	if err != nil {
		slog.Error("Error deleting request", "error", err.Error())
	}
	return err
}
