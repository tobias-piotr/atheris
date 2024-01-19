package tui

import (
	"fmt"
	"os"

	"atheris/requests"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
)

var (
	appStyle           = lipgloss.NewStyle().Padding(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type model struct {
	list         list.Model
	details      RequestDetails
	delegateKeys *delegateKeyMap
	db           *sqlx.DB
	selectedRq   *string
}

func NewModel(db *sqlx.DB) model {
	rqs, err := requests.GetRequests(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting requests: %v\n", err)
		os.Exit(1)
	}
	items := ItemsFromRequests(rqs)

	delKeys := newDelegateKeyMap()
	delegate := newItemDelegate(delKeys)

	itemList := list.New(items, delegate, 0, 0)
	itemList.Title = "Requests"

	details := NewRequestDetails(db, nil)

	return model{
		list:         itemList,
		details:      details,
		delegateKeys: delKeys,
		db:           db,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Split it into list and details
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.selectedRq == nil {
			// Handle list updates
			if msg.String() == "enter" {
				rq := m.list.SelectedItem().(item)
				m.selectedRq = &rq.id
				return m, nil
			}
		} else {
			// Handle details updates
			return m.details.Update(msg, m)
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.selectedRq == nil {
		return appStyle.Render(m.list.View())
	}
	m.details.reqID = m.selectedRq
	return appStyle.Render(m.details.View())
}
