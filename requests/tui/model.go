package tui

import (
	"fmt"
	"os"
	"time"

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

// Item represents a request on the list
type item struct {
	id        string
	createdAt time.Time
	name      string
	method    string
	path      string
}

func (i item) Title() string {
	if i.name == "" {
		return i.id
	}
	return i.name
}

func (i item) Description() string {
	return fmt.Sprintf("%s %s at %s", i.method, i.path, i.createdAt)
}

func (i item) FilterValue() string {
	if i.name == "" {
		return i.id
	}
	return i.name
}

func ItemsFromRequests(rqs []requests.Request) []list.Item {
	items := make([]list.Item, len(rqs))
	for i, rq := range rqs {
		items[i] = item{
			id:        rq.ID.String(),
			createdAt: rq.CreatedAt,
			name:      string(rq.Name),
			method:    rq.Method,
			path:      rq.Path,
		}
	}
	return items
}

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
	db           *sqlx.DB
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

	return model{
		list:         itemList,
		delegateKeys: delKeys,
		db:           db,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
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
	return appStyle.Render(m.list.View())
}
