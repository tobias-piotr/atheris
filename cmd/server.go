/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"atheris/libs"
	"atheris/requests"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var (
	dbFilename string
	port       string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts server",
	Long:  "Starts atheris server, that will catch all requests and store them in database",
	Run: func(_ *cobra.Command, _ []string) {
		db := libs.GetDBConnection(dbFilename)
		libs.MigrateDB(db)

		handler := requests.NewAPIHandler(db)

		e := echo.New()

		e.Any("/*", handler.HandleRequest)

		e.Logger.Fatal(e.Start(":" + port))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// SQLite database file
	serverCmd.Flags().StringVarP(&dbFilename, "db", "d", "database.sqlite", "SQLite database file")

	// Server port
	serverCmd.Flags().StringVarP(&port, "port", "p", "8888", "Port to listen on")
}
