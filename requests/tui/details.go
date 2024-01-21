package tui

import (
	"encoding/json"
	"fmt"
	"os"

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
	Back   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Rename, k.Delete, k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var keys = keyMap{
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
}

type RequestDetails struct {
	textArea    *textarea.Model
	help        *help.Model
	db          *sqlx.DB
	reqID       *string
	req         *requests.Request
	renaming    bool
	renameInput *input
}

func NewRequestDetails(db *sqlx.DB, reqID *string) RequestDetails {
	ta := textarea.New()
	h := help.New()
	return RequestDetails{
		textArea: &ta,
		help:     &h,
		db:       db,
		reqID:    reqID,
		renaming: false,
	}
}

func (rd *RequestDetails) Update(msg tea.KeyMsg, m *model) (tea.Model, tea.Cmd) {
	if rd.renaming {
		cmd := rd.renameInput.Update(msg, rd)
		return m, cmd
	}

	switch msg.String() {
	case "r":
		rd.renaming = true
		rd.renameInput = NewInput()
		return m, nil
	case "d":
		err := requests.DeleteRequest(rd.db, *rd.reqID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting request: %s\n", err.Error())
			os.Exit(1)
		}
		m.list.items.RemoveItem(m.list.items.Cursor())
		m.selectedRq = nil
		rd.reqID = nil
		rd.req = nil
		return m, nil
	case "b":
		m.selectedRq = nil
		rd.reqID = nil
		rd.req = nil
		return m, nil
	}
	return m, nil
}

func (rd RequestDetails) View() string {
	if rd.renaming {
		return rd.renameInput.View()
	}

	// TODO: This could also be cached
	reqText, err := json.MarshalIndent(rd.req.Response, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	rd.textArea.SetValue(string(reqText))

	return fmt.Sprintf(
		"%s\n\n%s",
		rd.req.ID,
		rd.textArea.View(),
	) + "\n\n" + rd.help.View(keys)
}
