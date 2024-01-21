package requests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InsertRequest(db *sqlx.DB, req RequestData) {
	res, err := json.Marshal(req.Response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %s\n", err.Error())
	}

	_, err = db.Exec(`
	INSERT INTO requests (id, method, path, response)
	VALUES ($1, $2, $3, $4);
	`, uuid.New(), req.Method, req.Path, res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting request: %s\n", err.Error())
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

func GetRequest(db *sqlx.DB, id string) (Request, error) {
	query := `
	SELECT id, created_at, name, method, path, response
	FROM requests
	WHERE id = $1;
	`
	var req Request
	err := db.Get(&req, query, id)
	return req, err
}

func SetRequestName(db *sqlx.DB, id string, name string) error {
	_, err := db.Exec(`
	UPDATE requests
	SET name = $1
	WHERE id = $2;
	`, name, id)
	return err
}

func DeleteRequest(db *sqlx.DB, id string) error {
	_, err := db.Exec(`
	DELETE FROM requests
	WHERE id = $1;
	`, id)
	return err
}
