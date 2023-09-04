package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyQuit = key.NewBinding(key.WithKeys("q"))

	keyUp    = key.NewBinding(key.WithKeys("up"))
	keyDown  = key.NewBinding(key.WithKeys("down"))
	keyLeft  = key.NewBinding(key.WithKeys("left"))
	keyRight = key.NewBinding(key.WithKeys("right"))

	keySelect = key.NewBinding(key.WithKeys("enter"))
	keyBack   = key.NewBinding(key.WithKeys("backspace"))

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

		// Quit

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

		// Selectors

		case key.Matches(msg, keySelect):
			current, ok := m.current()
			if !ok {
				return m, nil
			}
			if current.hasMode(entryModeDir) {
				m.path = filepath.Join(m.path, current.Name())

				m.resetCursor()
				if pos, ok := m.cursorCache[m.path]; ok {
					m.setCursor(pos)
				}

				err := m.list()
				if err != nil {
					// TODO: Improve error handling rather than quitting the application.
					return m, tea.Quit
				}
			}
			// TODO: handle files.

		case key.Matches(msg, keyBack):
			m.path = filepath.Join(m.path, "..")

			m.resetCursor()
			if pos, ok := m.cursorCache[m.path]; ok {
				m.setCursor(pos)
			}

			err := m.list()
			if err != nil {
				// TODO: Improve error handling rather than quitting the application.
				return m, tea.Quit
			}

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

	m.saveCursor()
	return m, nil
}
