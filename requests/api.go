package requests

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ResolvePrefix(path string) string {
	apiMap := map[string]string{
		"aichat": "http://localhost:8080",
	}
	components := strings.Split(path, "/")
	return apiMap[components[1]]
}

type APIHandler struct {
	db *sqlx.DB
}

func NewAPIHandler(db *sqlx.DB) *APIHandler {
	return &APIHandler{db}
}

func (h *APIHandler) HandleRequest(c echo.Context) error {
	// Resolve API
	api := ResolvePrefix(c.Request().URL.Path)
	if api == "" {
		return c.String(http.StatusNotFound, "Not Found")
	}

	// Prepare request
	req, err := http.NewRequest(
		c.Request().Method,
		api+c.Request().URL.Path,
		c.Request().Body,
	)
	if err != nil {
		return err
	}
	req.Header = c.Request().Header
	// TODO: Set headers for forwarding

	// Make request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("Request failed", "error", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set response headers
	for k, v := range res.Header {
		c.Response().Header().Set(k, v[0])
	}
	// TODO: Set headers for forwarding

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return err
	}

	// Insert request into database
	go InsertRequest(h.db, RequestData{
		Method: c.Request().Method,
		Prefix: api,
		Path:   c.Request().URL.Path,
		Response: ResponseData{
			Status: res.Status,
			Body:   body,
		},
	})

	c.Response().WriteHeader(res.StatusCode)
	c.Response().Write(body)
	return nil
}
