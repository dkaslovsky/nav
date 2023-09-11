package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	var view string
	if m.modeExit {
		if m.modeSubshell || m.exitStr == "" {
			return ""
		}
		return m.exitStr + "\n"
	}
	if m.modeHelp {
		view = commands()
	} else if m.modeDebug {
		view = m.debugView()
	} else {
		view = m.normalView()
	}
	return strings.Join([]string{view, m.statusBar()}, "\n")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:

		// Force Quit
		if key.Matches(msg, keyQuitForce) {
			m.modeExit = true
			m.exitCode = 2
			return m, tea.Quit
		}

		esc := false
		if m.modeSearch || m.modeDebug || m.modeHelp {
			if key.Matches(msg, m.esc.key) {
				if m.esc.triggered() {
					esc = true
				}
			} else {
				m.esc.reset()
			}
		}

		// Help mode

		if m.modeHelp {

			if esc || key.Matches(msg, keyEsc) {
				m.modeHelp = false
			}

			return m, nil
		}

		// Debug mode

		if m.modeDebug {

			if esc || key.Matches(msg, keyEsc) {
				m.modeDebug = false
			}

			return m, nil
		}

		// Search mode

		if m.modeSearch {
			if esc || key.Matches(msg, keyEsc) {
				m.clearSearch()
				if m.error != nil && errors.Is(m.error, ErrNoSearchResults) {
					m.clearError()
				}
				return m, nil
			}

			switch {

			// Do not allow remapped escape key in the search.
			case key.Matches(msg, m.esc.key):
				return m, nil

			case key.Matches(msg, keyBack):
				if len(m.search) > 0 {
					m.search = m.search[:len(m.search)-1]
					return m, nil
				}

				m.saveCursor()

				_, m.search = filepath.Split(m.path)
				path, err := filepath.Abs(filepath.Join(m.path, ".."))
				if err != nil {
					m.setError(err, "failed to evaluate path")
					return m, nil
				}
				m.path = path

				err = m.list()
				if err != nil {
					m.setError(err, "failed to list entries")
					return m, nil
				}

				return m, nil

			// TODO: keySelect, keyTab, and keyFileSeparator share common logic that should be encapsulated.
			case key.Matches(msg, keySelect):
				selected, err := m.selected()
				if err != nil {
					m.setError(err, "failed to select entry")
					m.clearSearch()
					return m, nil
				}
				if !selected.hasMode(entryModeDir) {
					// No-op for non directories
					return m, nil
				}
				m.path = m.path + "/" + selected.Name()
				// TODO: encapsulate this fix
				if strings.HasPrefix(m.path, "//") {
					m.path = m.path[1:]
				}
				m.search = ""
				err = m.list()
				if err != nil {
					m.setError(err, "failed to list entries")
					m.clearSearch()
					return m, nil
				}
				return m, nil

			case key.Matches(msg, keyTab):
				if m.displayed != 1 {
					return m, nil
				}
				selected, err := m.selected()
				if err != nil {
					m.setError(err, "failed to select entry")
					m.clearSearch()
					return m, nil
				}
				if !selected.hasMode(entryModeDir) {
					// No-op for non directories
					return m, nil
				}
				m.path = m.path + "/" + selected.Name()
				// TODO: encapsulate this fix
				if strings.HasPrefix(m.path, "//") {
					m.path = m.path[1:]
				}
				m.search = ""
				err = m.list()
				if err != nil {
					m.setError(err, "failed to list entries")
					m.clearSearch()
					return m, nil
				}
				return m, nil

			case key.Matches(msg, keyFileSeparator):
				if m.displayed != 1 {
					m.search += keyString(keyFileSeparator)
					return m, nil
				}
				selected, err := m.selected()
				if err != nil {
					m.setError(err, "failed to select entry")
					m.clearSearch()
					return m, nil
				}
				if !selected.hasMode(entryModeDir) {
					// No-op for non directories
					return m, nil
				}
				m.path = m.path + "/" + selected.Name()
				// TODO: encapsulate this fix
				if strings.HasPrefix(m.path, "//") {
					m.path = m.path[1:]
				}
				m.search = ""
				err = m.list()
				if err != nil {
					m.setError(err, "failed to list entries")
					m.clearSearch()
					return m, nil
				}
				return m, nil

			default:
				if msg.Type == tea.KeyRunes {
					m.search += string(msg.Runes)
					return m, nil
				}

			}
		}

		switch {

		// Quit

		case key.Matches(msg, keyQuitWithDirectory):
			m.modeExit = true
			m.exitStr = sanitizeOutputPath(m.path)
			if m.modeSubshell {
				fmt.Print(m.exitStr)
			}
			return m, tea.Quit

		case key.Matches(msg, keyQuitWithSelected):
			current, err := m.selected()
			if err != nil {
				m.setError(err, "failed to select entry")
				return m, nil
			}

			m.modeExit = true
			m.exitStr = sanitizeOutputPath(filepath.Join(m.path, current.Name()))
			if m.modeSubshell {
				fmt.Print(m.exitStr)
			}
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
			current, err := m.selected()
			if err != nil {
				m.setError(err, "failed to select entry")
				return m, nil
			}

			m.saveCursor()

			isDir := current.hasMode(entryModeDir)
			isSymlink := current.hasMode(entryModeSymlink)

			if !(isDir || isSymlink) {
				// The selected entry is a file, which is a no-op.
				break
			}

			if isSymlink {
				followed, err := filepath.EvalSymlinks(filepath.Join(m.path, current.Name()))
				if err != nil {
					m.setError(err, "failed to evaluate symlink")
					return m, nil
				}
				info, err := os.Stat(followed)
				if err != nil {
					m.setError(err, "failed to evaluate symlink")
					return m, nil
				}
				if !info.IsDir() {
					// The symlink points to a file, which is a no-op.
					break
				}
				m.path = followed
			} else {
				path, err := filepath.Abs(filepath.Join(m.path, current.Name()))
				if err != nil {
					m.setError(err, "failed to evaluate path")
					return m, nil
				}
				m.path = path
			}

			err = m.list()
			if err != nil {
				m.setError(err, "failed to list entries")
				return m, nil
			}

			m.clearSearch()
			m.clearError()

			// Return to ensure the cursor is not re-saved using the updated path.
			return m, nil

		case key.Matches(msg, keyBack):
			m.saveCursor()

			path, err := filepath.Abs(filepath.Join(m.path, ".."))
			if err != nil {
				m.setError(err, "failed to evaluate path")
				return m, nil
			}
			m.path = path

			err = m.list()
			if err != nil {
				m.setError(err, "failed to list entries")
				return m, nil
			}

			m.clearSearch()
			m.clearError()

			// Return to ensure the cursor is not re-saved using the updated path.
			return m, nil

		// Change modes

		case key.Matches(msg, keyDebugMode):
			m.modeDebug = true

		case key.Matches(msg, keyHelpMode):
			m.modeHelp = true

		case key.Matches(msg, keySearchMode):
			m.modeSearch = true
			m.clearError()

		// Toggles

		case key.Matches(msg, keyToggleFollowSymlink):
			m.modeFollowSymlink = !m.modeFollowSymlink

		case key.Matches(msg, keyToggleHidden):
			m.modeHidden = !m.modeHidden

		case key.Matches(msg, keyToggleList):
			m.modeList = !m.modeList

		case key.Matches(msg, keyDismissError):
			m.clearError()

		}
	}

	m.saveCursor()
	return m, nil
}

func sanitizeOutputPath(s string) string {
	return strings.Replace(s, " ", "\\ ", -1)
}
