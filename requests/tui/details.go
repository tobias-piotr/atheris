package tui

import (
	"encoding/json"
	"fmt"

	"atheris/requests"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type keyMap struct {
	Rename key.Binding
	Delete key.Binding
	Help   key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Rename, k.Delete},
		{k.Back, k.Quit, k.Help},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
	Rename: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Delete: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type RequestDetails struct {
	textArea textarea.Model
	help     help.Model
	db       *sqlx.DB
	reqID    *string
}

func NewRequestDetails(db *sqlx.DB, reqID *string) RequestDetails {
	return RequestDetails{
		textArea: textarea.New(),
		help:     help.New(),
		db:       db,
		reqID:    reqID,
	}
}

func (rd RequestDetails) Update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?":
		rd.help.ShowAll = !rd.help.ShowAll
	case "b":
		m.selectedRq = nil
		return m, nil
	case "r":
		m.selectedRq = nil
		return m, nil
	case "d":
		m.selectedRq = nil
		return m, nil
	}
	return m, nil
}

func (rd RequestDetails) View() string {
	// TODO: Handle error
	req, _ := requests.GetRequest(rd.db, *rd.reqID)
	reqText, _ := json.MarshalIndent(req.Response, "", "  ")

	rd.textArea.SetValue(string(reqText))

	return fmt.Sprintf(
		"%s\n\n%s",
		req.ID,
		rd.textArea.View(),
	) + "\n\n" + rd.help.View(keys)
}
