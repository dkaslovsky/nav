package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyQuit   = key.NewBinding(key.WithKeys("q"))
	keyHidden = key.NewBinding(key.WithKeys("a")) // Toggles showing hidden files, (similar to ls -a).
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	return m.view()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keyQuit):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			return m, tea.Quit

		case key.Matches(msg, keyHidden):
			m.hidden = !m.hidden
			return m, nil
		}
	}

	return m, nil
}
