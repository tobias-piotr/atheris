package tui

import (
	"fmt"
	"os"

	"atheris/requests"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
)

var appStyle = lipgloss.NewStyle().Padding(1, 2)

type model struct {
	list       RequestList
	details    RequestDetails
	db         *sqlx.DB
	selectedRq *string
}

func NewModel(db *sqlx.DB) model {
	list := NewRequestList(db)
	details := NewRequestDetails(db, nil)

	return model{
		list:    list,
		details: details,
		db:      db,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		if m.selectedRq == nil {
			// Handle item selection
			if msg.String() == "enter" {
				rq := m.list.items.SelectedItem().(item)
				m.selectedRq = &rq.id
				m.details.reqID = m.selectedRq
				req, err := requests.GetRequest(m.db, *m.selectedRq)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting request: %s\n", err.Error())
					os.Exit(1)
				}
				m.details.req = &req
				return m, nil
			}
		} else {
			// Handle details updates
			return m.details.Update(msg, &m)
		}

	case tea.WindowSizeMsg:
		m.list.Resize(msg)
	}

	// Default to list update
	return m, m.list.Update(msg)
}

func (m model) View() string {
	if m.selectedRq == nil {
		return appStyle.Render(m.list.View())
	}
	return appStyle.Render(m.details.View())
}
