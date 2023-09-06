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
	keyHelp          = key.NewBinding(key.WithKeys("h")) // Toggles showing help screen.
	keyHidden        = key.NewBinding(key.WithKeys("a")) // Toggles showing hidden files, (similar to ls -a).
	keyList          = key.NewBinding(key.WithKeys("l")) // Toggles showing file info in list mode.
	keySearch        = key.NewBinding(key.WithKeys("/")) // Toggles search mode.
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	if m.modeHelp {
		return strings.Join([]string{usage(), m.status()}, "\n")
	}

	output := []string{}

	// First row of output is the location bar.
	locationBar := barRendererLocation.Render(m.location())
	if m.modeSearch {
		locationBar += barRendererSearch.Render(fileSeparator + m.search)
	}
	output = append(output, locationBar)

	// Construct display names from filtered entries and populate a local cache mapping between them.
	var (
		displayIdx   = 0
		displayNames = []*displayName{}
		localCache   = newCacheItem(&position{c: 0, r: 0}) // Store local copy of current state.
	)
	for entryIdx, ent := range m.entries {

		// Filter hidden files.
		if !m.modeHidden && ent.hasMode(entryModeHidden) {
			continue
		}
		// Filter for search.
		if m.modeSearch && m.search != "" {
			if !strings.HasPrefix(ent.Name(), m.search) {
				continue
			}
		}

		displayNames = append(displayNames, newDisplayName(ent, m.displayNameOpts()...))

		// Populate local cache.
		localCache.displayToEntityIndex[displayIdx] = entryIdx
		localCache.entityToDisplayIndex[entryIdx] = displayIdx
		displayIdx++
	}
	m.displayed = displayIdx

	// Grid layout for display.
	var (
		width     = m.width
		height    = m.height - 2 // Account for location and status bars
		gridNames [][]string
		layout    gridLayout
	)
	if m.modeList {
		gridNames, layout = gridSingleColumn(displayNames, width, height)
	} else {
		gridNames, layout = gridMultiColumn(displayNames, width, height)
	}
	localCache.columns, localCache.rows = layout.columns, layout.rows
	m.columns, m.rows = layout.columns, layout.rows
	if m.c >= m.columns || m.r > m.rows {
		m.resetCursor()
	}

	// Retrieve cached state for the current page from its previous rendering.
	var cache *cacheItem
	if c, found := m.viewCache[m.path]; found {
		if c.hasIndexes() {
			// Full cache is found.
			cache = c
		} else {
			// Only previous cursor position is cached, use local cache with this cursor position.
			cache = localCache
			cache.cursorPosition = c.cursorPosition
		}
	} else {
		// No cache exists yet, use local cache.
		cache = localCache
	}
	// Lookup the entry index using the cached cursor (display) position.
	if entryIdx, entryFound := cache.displayToEntityIndex[cache.cursorPosition.index(cache.rows)]; entryFound {
		// Use the entry index to get the current (local cache) display index.
		if dispIdx, dispFound := localCache.entityToDisplayIndex[entryIdx]; dispFound {
			// Set the cursor position using the current display index.
			m.setCursor(newPosition(dispIdx, localCache.rows))
		}
	}

	// Update the cache with current state.
	m.viewCache[m.path] = localCache

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

	// Add status bar to output.
	output = append(output, m.status())
	return strings.Join(output, "\n")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:

		// Help mode

		if m.modeHelp {
			switch {

			case key.Matches(msg, keyQuit, keyHelp):
				m.modeHelp = !m.modeHelp

			case key.Matches(msg, keyQuitForce):
				return m, tea.Quit

			}

			return m, nil
		}

		// Search mode

		if m.modeSearch {
			switch {

			case key.Matches(msg, keySearch):
				m.search = ""
				m.modeSearch = false

			case key.Matches(msg, keyBack):
				if len(m.search) > 0 {
					m.search = m.search[:len(m.search)-1]
				}

			default:
				if msg.Type == tea.KeyRunes {
					m.search += string(msg.Runes)
				}
			}

			return m, nil
		}

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
					// TODO: Handle error.
					return m, tea.Quit
				}
				info, err := os.Stat(followed)
				if err != nil {
					// TODO: Handle error.
					return m, tea.Quit
				}
				if !info.IsDir() {
					// The symlink points to a file, which is a no-op.
					break
				}
				m.path = followed
			} else {
				path, err := filepath.Abs(filepath.Join(m.path, current.Name()))
				if err != nil {
					// TODO: Handle error.
					return m, tea.Quit
				}
				m.path = path
			}

			err := m.list()
			if err != nil {
				// TODO: Improve error handling rather than quitting the application.
				return m, tea.Quit
			}

			// Clear search mode
			m.modeSearch = false
			m.search = ""

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

			// Clear search mode
			m.modeSearch = false
			m.search = ""

			// Return to ensure the cursor is not re-saved using the updated path.
			return m, nil

		// Toggles

		case key.Matches(msg, keyFollowSymlink):
			m.modeFollowSymlink = !m.modeFollowSymlink

		case key.Matches(msg, keyHelp):
			m.modeHelp = !m.modeHelp

		case key.Matches(msg, keyHidden):
			m.modeHidden = !m.modeHidden

		case key.Matches(msg, keyList):
			m.modeList = !m.modeList

		case key.Matches(msg, keySearch):
			m.modeSearch = !m.modeSearch

		}
	}

	m.saveCursor()
	return m, nil
}

func (m *model) status() string {
	mode := "NORMAL"
	cmds := []string{
		`"/": search`,
		`"h": help`,
		`"q": quit`,
	}
	if m.modeSearch {
		mode = "SEARCH"
		cmds = []string{
			`"/": cancel search`,
		}
	} else if m.modeHelp {
		mode = "HELP"
		cmds = []string{
			`"q": cancel help`,
			`"esc": exit application`,
		}
	}
	status := strings.Join([]string{
		"  " + name,
		fmt.Sprintf("%s MODE", mode),
		strings.Join(cmds, ", "),
		"",
	}, "\t")

	err := ""
	if m.error != "" {
		err = fmt.Sprintf("ERROR: %s\t", m.error)
	}

	return barRendererStatus.Render(status) + barRendererError.Render(err)
}
