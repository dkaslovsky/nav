package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyQuit = key.NewBinding(key.WithKeys("q"))

	keyUp    = key.NewBinding(key.WithKeys("up"))
	keyDown  = key.NewBinding(key.WithKeys("down"))
	keyLeft  = key.NewBinding(key.WithKeys("left"))
	keyRight = key.NewBinding(key.WithKeys("right"))

	keyFollowSymlink = key.NewBinding(key.WithKeys("s")) // Toggles showing symlink paths.
	keyHidden        = key.NewBinding(key.WithKeys("a")) // Toggles showing hidden files, (similar to ls -a).
	keyList          = key.NewBinding(key.WithKeys("l")) // Toggles showing file info in list mode.
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

		// Cursor

		case key.Matches(msg, keyUp):
			m.moveUp()

		case key.Matches(msg, keyDown):
			m.moveDown()

		case key.Matches(msg, keyLeft):
			m.moveLeft()

		case key.Matches(msg, keyRight):
			m.moveRight()

		// Toggles

		case key.Matches(msg, keyFollowSymlink):
			m.modeFollowSymlink = !m.modeFollowSymlink
			return m, nil

		case key.Matches(msg, keyHidden):
			m.modeHidden = !m.modeHidden
			return m, nil

		case key.Matches(msg, keyList):
			m.modeList = !m.modeList
			return m, nil
		}
	}

	return m, nil
}
