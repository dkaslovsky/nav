package main

import (
	"fmt"
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

		// Quit

		if key.Matches(msg, keyQuit) {
			m.setExitWithCode("", 2)
			return m, tea.Quit
		}

		// Remapped escape logic

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

		// Error mode

		if m.modeError {

			// Debug mode

			if m.modeDebug {
				if esc || key.Matches(msg, keyEsc) || key.Matches(msg, keyDebugMode) {
					m.modeDebug = false
				}

				if key.Matches(msg, keyDismissError) {
					m.clearError()
					m.modeDebug = false
				}

				return m, nil
			}

			if key.Matches(msg, keyDismissError) {
				m.clearError()
			}

			if key.Matches(msg, keyDebugMode) {
				m.modeDebug = true
			}

			return m, nil
		}

		// Help mode

		if m.modeHelp {
			if esc || key.Matches(msg, keyEsc) || key.Matches(msg, keyHelpMode) {
				m.modeHelp = false
			}

			return m, nil
		}

		// Search mode

		if m.modeSearch {
			if esc || key.Matches(msg, keyEsc) {
				m.clearSearch()
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
					m.setError(err, err.Error())
					return m, nil
				}

				return m, nil

			case key.Matches(msg, keySelect):
				return m.searchSelectAction()

			case key.Matches(msg, keyTab):
				if m.displayed != 1 {
					return m, nil
				}
				return m.searchSelectAction()

			case key.Matches(msg, keyFileSeparator):
				if m.displayed != 1 {
					m.search += keyString(keyFileSeparator)
					return m, nil
				}
				return m.searchSelectAction()

			default:
				if msg.Type == tea.KeyRunes || key.Matches(msg, keySpace) {
					m.search += string(msg.Runes)
					return m, nil
				}

			}
		}

		switch {

		// Return

		case key.Matches(msg, keyReturnDirectory):
			m.setExit(sanitizeOutputPath(m.path))
			if m.modeSubshell {
				fmt.Print(m.exitStr)
			}
			return m, tea.Quit

		case key.Matches(msg, keyReturnSelected):
			selected, err := m.selected()
			if err != nil {
				m.setError(err, "failed to select entry")
				return m, nil
			}

			var path string
			if selected.hasMode(entryModeSymlink) {
				sl, err := followSymlink(m.path, selected)
				if err != nil {
					m.setError(err, "failed to evaluate symlink")
					return m, nil
				}
				path = sl.absPath
			} else {
				path = filepath.Join(m.path, selected.Name())
			}

			m.setExit(sanitizeOutputPath(path))
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
			return m.selectAction()

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
				m.setError(err, err.Error())
				return m, nil
			}

			m.clearSearch()

			// Return to ensure the cursor is not re-saved using the updated path.
			return m, nil

		// Change modes

		case key.Matches(msg, keyHelpMode):
			m.modeHelp = true

		case key.Matches(msg, keySearchMode):
			m.modeSearch = true

		// Toggles

		case key.Matches(msg, keyToggleFollowSymlink):
			m.modeFollowSymlink = !m.modeFollowSymlink

		case key.Matches(msg, keyToggleHidden):
			m.modeHidden = !m.modeHidden

		case key.Matches(msg, keyToggleList):
			m.modeList = !m.modeList

		}
	}

	m.saveCursor()
	return m, nil
}

func sanitizeOutputPath(s string) string {
	return strings.Replace(s, " ", "\\ ", -1)
}
