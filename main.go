package main

import (
	"atheris/libs"
	"atheris/requests"

	"github.com/labstack/echo/v4"
)

func main() {
	db := libs.GetDBConnection("database.db")
	libs.MigrateDB(db)

	handler := requests.NewAPIHandler(db)

	e := echo.New()

	e.Any("/*", handler.HandleRequest)

	// TODO: Replace with slog
	e.Logger.Fatal(e.Start(":8888"))
}
