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
	keyQuitForce         = key.NewBinding(key.WithKeys("ctrl+c"))
	keyQuitWithDirectory = key.NewBinding(key.WithKeys("ctrl+d"))
	keyQuitWithSelected  = key.NewBinding(key.WithKeys("ctrl+x"))

	keyEsc           = key.NewBinding(key.WithKeys("esc"))
	keySelect        = key.NewBinding(key.WithKeys("enter"))
	keyBack          = key.NewBinding(key.WithKeys("backspace"))
	keyTab           = key.NewBinding(key.WithKeys("tab"))
	keyFileSeparator = key.NewBinding(key.WithKeys("/"))

	keyExtraEsc = key.NewBinding(key.WithKeys(os.Getenv(envEscRemap)))

	keyUp    = key.NewBinding(key.WithKeys("up", "k"))
	keyDown  = key.NewBinding(key.WithKeys("down", "j"))
	keyLeft  = key.NewBinding(key.WithKeys("left", "h"))
	keyRight = key.NewBinding(key.WithKeys("right", "l"))

	keyDebugMode  = key.NewBinding(key.WithKeys("d"))
	keyHelpMode   = key.NewBinding(key.WithKeys("H"))
	keySearchMode = key.NewBinding(key.WithKeys("i"))

	keyToggleFollowSymlink = key.NewBinding(key.WithKeys("f"))
	keyToggleHidden        = key.NewBinding(key.WithKeys("a"))
	keyToggleList          = key.NewBinding(key.WithKeys("L"))

	keyDismissError = key.NewBinding(key.WithKeys("e"))
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	var view string
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
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			m.exitCode = 2
			return m, tea.Quit
		}

		// Help mode

		if m.modeHelp {
			if key.Matches(msg, keyEsc, keyExtraEsc) {
				m.modeHelp = false
			}

			return m, nil
		}

		// Debug mode

		if m.modeDebug {
			if key.Matches(msg, keyEsc, keyExtraEsc) {
				m.modeDebug = false
			}

			return m, nil
		}

		// Search mode

		if m.modeSearch {
			switch {

			case key.Matches(msg, keyEsc, keyExtraEsc):
				m.clearError()
				m.clearSearch()
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
					// TODO: Handle error.
					return m, tea.Quit
				}
				m.path = path

				err = m.list()
				if err != nil {
					// TODO: Improve error handling rather than quitting the application.
					return m, tea.Quit
				}

				return m, nil

			case key.Matches(msg, keySelect):
				if selected, ok := m.selected(); ok && selected.hasMode(entryModeDir) {
					m.path = m.path + "/" + selected.Name()
					// TODO: encapsulate this fix
					if strings.HasPrefix(m.path, "//") {
						m.path = m.path[1:]
					}
					m.search = ""
					err := m.list()
					if err != nil {
						// TODO: Improve error handling rather than quitting the application.
						return m, tea.Quit
					}
				}
				return m, nil

			case key.Matches(msg, keyTab):
				if m.displayed != 1 {
					return m, nil
				}
				// TODO: this is the same as slash, consider encapsulation.
				if selected, ok := m.selected(); ok && selected.hasMode(entryModeDir) {
					m.path = m.path + "/" + selected.Name()
					// TODO: encapsulate this fix
					if strings.HasPrefix(m.path, "//") {
						m.path = m.path[1:]
					}
					m.search = ""
					err := m.list()
					if err != nil {
						// TODO: Improve error handling rather than quitting the application.
						return m, tea.Quit
					}
				}
				return m, nil

			case key.Matches(msg, keyFileSeparator):
				if m.displayed != 1 {
					m.search += keyString(keyFileSeparator)
					return m, nil
				}
				// TODO: this is the same as tab, consider encapsulation.
				if selected, ok := m.selected(); ok && selected.hasMode(entryModeDir) {
					m.path = m.path + "/" + selected.Name()
					// TODO: encapsulate this fix
					if strings.HasPrefix(m.path, "//") {
						m.path = m.path[1:]
					}
					m.search = ""
					err := m.list()
					if err != nil {
						// TODO: Improve error handling rather than quitting the application.
						return m, tea.Quit
					}
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
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			fmt.Println(sanitizeOutputPath(m.path))
			return m, tea.Quit

		case key.Matches(msg, keyQuitWithSelected):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			current, ok := m.selected()
			if !ok {
				return m, nil
			}
			fmt.Println(sanitizeOutputPath(filepath.Join(m.path, current.Name())))
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
					m.errorStatus = "failed to evaluate symlink"
					m.error = err
					return m, nil
				}
				info, err := os.Stat(followed)
				if err != nil {
					m.errorStatus = "failed to evaluate symlink"
					m.error = err
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
					m.errorStatus = "failed to evaluate path"
					m.error = err
					return m, nil
				}
				m.path = path
			}

			err := m.list()
			if err != nil {
				// TODO: Improve error handling rather than quitting the application.
				return m, tea.Quit
			}

			m.clearSearch()
			m.clearError()

			// Return to ensure the cursor is not re-saved using the updated path.
			return m, nil

		case key.Matches(msg, keyBack):
			m.saveCursor()

			path, err := filepath.Abs(filepath.Join(m.path, ".."))
			if err != nil {
				// TODO: Handle error.
				return m, tea.Quit
			}
			m.path = path

			err = m.list()
			if err != nil {
				// TODO: Improve error handling rather than quitting the application.
				return m, tea.Quit
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

		// Dismiss Error
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
