/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"fmt"
	"os"

	"atheris/requests/tui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// requestsCmd represents the requests command
var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "View and manage requests",
	Long: `View a list of saved requests. You can also view the details of each request,
set a name for them, delete them, and most importantly, compare them to see the differences.`,
	Run: func(_ *cobra.Command, _ []string) {
		m := tui.NewModel([]list.Item{})
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting the program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(requestsCmd)
}
