package main

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type APIData struct {
	Prefix string
	URL    string
	Secret bool
}

type API struct {
	ID        uuid.UUID
	Timestamp time.Time
	Prefix    string
	URL       string
	Secret    bool
}

type RequestData struct {
	Prefix   string
	Path     string
	Response Response
}

type Request struct {
	ID        uuid.UUID
	Timestamp time.Time
	Path      string
	Response  Response
}

type Response struct {
	Status  string
	Content map[string]any
}

func ResolvePrefix(path string) string {
	apiMap := map[string]string{
		"aichat": "http://localhost:8080",
	}
	components := strings.Split(path, "/")
	return apiMap[components[1]]
}

func main() {
	e := echo.New()

	e.Any("/*", func(c echo.Context) error {
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

		c.Response().WriteHeader(res.StatusCode)
		c.Response().Write(body)
		return nil
	})

	e.Logger.Fatal(e.Start(":8888"))
}
