/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atheris",
	Short: "Atheris is a tool for catching and analyzing requests",
	Long: `With atheris you can catch all requests and store them in database for later analysis.
First, you need to start atheris server, that will catch all requests and store them in database.
Then, you can use atheris cli to analyze requests stored in database.
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
