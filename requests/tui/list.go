package tui

import (
	"fmt"
	"os"
	"time"

	"atheris/requests"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
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

// ItemsFromRequests converts a slice of requests to a slice of list items
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

type delegateKeyMap struct {
	choose key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	help := []key.Binding{keys.choose}
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type RequestList struct {
	items list.Model
}

func NewRequestList(db *sqlx.DB) RequestList {
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

	return RequestList{items: itemList}
}

func (l *RequestList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	l.items, cmd = l.items.Update(msg)
	return cmd
}

func (l RequestList) View() string {
	return l.items.View()
}

func (l *RequestList) Resize(msg tea.WindowSizeMsg) {
	h, v := appStyle.GetFrameSize()
	l.items.SetSize(msg.Width-h, msg.Height-v)
}
