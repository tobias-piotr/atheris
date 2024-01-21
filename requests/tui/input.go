package tui

import (
	"fmt"
	"os"

	"atheris/requests"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type input struct {
	textInput textinput.Model
}

func NewInput() *input {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &input{
		textInput: ti,
	}
}

func (i input) Init() tea.Cmd {
	return textinput.Blink
}

func (i *input) Update(msg tea.Msg, rd *RequestDetails) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return tea.Quit
		case tea.KeyEsc:
			rd.renaming = false
			return nil
		case tea.KeyEnter:
			err := requests.SetRequestName(rd.db, *rd.reqID, i.textInput.Value())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting request name: %s\n", err.Error())
				os.Exit(1)
			}
			rd.req.Name = requests.NullableString(i.textInput.Value())
			rd.renaming = false
			return nil
		}
	}

	i.textInput, cmd = i.textInput.Update(msg)
	return cmd
}

func (i input) View() string {
	return fmt.Sprintf(
		"Enter a new name for the request:\n\n%s\n\n%s",
		i.textInput.View(),
		"(esc to go back)",
	) + "\n"
}
