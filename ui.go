package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keyQuit        = key.NewBinding(key.WithKeys("q"))
	keyQuitCurrent = key.NewBinding(key.WithKeys("c"))
	keyQuitForce   = key.NewBinding(key.WithKeys("esc"))

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
	output := []string{
		// First row of output is the location bar.
		barRendererLocation.Render(m.location()),
	}

	displayNames := []*displayName{}
	for idx, ent := range m.entries {
		// Optionally do not show hidden files.
		if !m.modeHidden && ent.hasMode(entryModeHidden) {
			continue
		}
		displayNames = append(displayNames, newDisplayName(ent, m.displayNameOpts()...))
		m.displayIndex[len(displayNames)-1] = idx
	}

	// Grid layout for display.
	var (
		width     = m.width
		height    = m.height - 1 // Account for location bar
		gridNames [][]string
		layout    gridLayout
	)
	if m.modeList {
		gridNames, layout = gridSingleColumn(displayNames, width, height)
	} else {
		gridNames, layout = gridMultiColumn(displayNames, width, height)
	}
	m.columns = layout.columns
	m.rows = layout.rows
	if m.c >= m.columns {
		m.c = 0
	}
	if m.r >= m.rows {
		m.r = 0
	}

	// Render entry names in grid.
	gridOutput := make([]string, layout.rows)
	for row := 0; row < layout.rows; row++ {
		for col := 0; col < layout.columns; col++ {
			if col == m.c && row == m.r {
				gridOutput[row] += cursorRendererSelected.Render(gridNames[col][row])
			} else {
				gridOutput[row] += cursorRendererNormal.Render(gridNames[col][row])
			}
		}
	}
	output = append(output, gridOutput...)

	return strings.Join(output, "\n")
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
			fmt.Println(m.path)
			return m, tea.Quit

		case key.Matches(msg, keyQuitCurrent):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			current, ok := m.selected()
			if !ok {
				return m, nil
			}
			fmt.Println(filepath.Join(m.path, current.Name()))
			return m, tea.Quit

		case key.Matches(msg, keyQuitForce):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			m.exitCode = 2
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
			current, ok := m.selected()
			if !ok {
				return m, nil
			}

			isDir := current.hasMode(entryModeDir)
			isSymlink := current.hasMode(entryModeSymlink)

			if !(isDir || isSymlink) {
				// The selected entry is a file, which is a no-op.
				break
			}

			if isSymlink {
				followed, err := filepath.EvalSymlinks(filepath.Join(m.path, current.Name()))
				if err != nil {
					return m, nil
				}
				info, err := os.Stat(followed)
				if err != nil {
					return m, nil
				}
				if !info.IsDir() {
					// The symlink points to a file, which is a no-op.
					break
				}
				m.path = followed
			} else {
				m.path = filepath.Join(m.path, current.Name())
			}

			m.resetCursor()
			if pos, ok := m.cursorCache[m.path]; ok {
				m.setCursor(pos)
			}

			err := m.list()
			if err != nil {
				// TODO: Improve error handling rather than quitting the application.
				return m, tea.Quit
			}

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

		case key.Matches(msg, keyHidden):
			m.modeHidden = !m.modeHidden

		case key.Matches(msg, keyList):
			m.modeList = !m.modeList

		}
	}

	m.saveCursor()
	return m, nil
}
