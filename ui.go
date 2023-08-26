package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyQuit          = key.NewBinding(key.WithKeys("q"))
	keyHidden        = key.NewBinding(key.WithKeys("a")) // Toggles showing hidden files, (similar to ls -a).
	keyFollowSymlink = key.NewBinding(key.WithKeys("s")) // Toggles showing symlink paths.
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	return m.view()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keyQuit):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			return m, tea.Quit

		case key.Matches(msg, keyHidden):
			m.modeHidden = !m.modeHidden
			return m, nil

		case key.Matches(msg, keyFollowSymlink):
			m.modeFollowSymlink = !m.modeFollowSymlink
			return m, nil
		}
	}

	return m, nil
}
