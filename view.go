package main

import (
	"fmt"
	"strings"
)

func (m *model) normalView() string {
	var (
		// Cache for storing the current state as it is constructed.
		updateCache = newCacheItem()

		displayNames    = []*displayName{}
		displayNameOpts = m.displayNameOpts()
		displayed       = 0
	)

	// Construct display names from filtered entries and populate a new cache mapping between them.
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
		displayNames = append(displayNames, newDisplayName(ent, displayNameOpts...))

		updateCache.displayToEntityIndex[displayed] = entryIdx
		updateCache.entityToDisplayIndex[entryIdx] = displayed
		displayed++
	}

	// Grid layout for display.
	var (
		width     = m.width
		height    = m.height - 2 // Account for location and status bars.
		gridNames [][]string
		layout    gridLayout
	)
	if m.modeList {
		gridNames, layout = gridSingleColumn(displayNames, width, height)
	} else {
		gridNames, layout = gridMultiColumn(displayNames, width, height)
	}

	// Retrieve cached cursor position and index mappings to set cursor position for current state.
	updateCursorPosition := &position{c: 0, r: 0}
	if cache, found := m.viewCache[m.path]; found && cache.hasIndexes() {
		// Lookup the entry index using the cached cursor (display) position.
		if entryIdx, entryFound := cache.displayToEntityIndex[cache.cursorPosition.index(cache.rows)]; entryFound {
			// Use the entry index to get the current display index.
			if dispIdx, dispFound := updateCache.entityToDisplayIndex[entryIdx]; dispFound {
				// Set the cursor position using the current display index and layout.
				updateCursorPosition = newPositionFromIndex(dispIdx, layout.rows)
			}
		}
	}

	// Update the model.
	m.displayed = displayed
	m.columns = layout.columns
	m.rows = layout.rows
	m.setCursor(updateCursorPosition)
	if m.c >= m.columns || m.r > m.rows {
		m.resetCursor()
	}

	// Update the cache.
	updateCache.cursorPosition = updateCursorPosition
	updateCache.columns = layout.columns
	updateCache.rows = layout.rows
	m.viewCache[m.path] = updateCache

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

	// Construct the final view.
	output := []string{m.locationBar()}
	output = append(output, gridOutput...)
	return strings.Join(output, "\n")
}

func (m *model) debugView() string {
	output := barRendererOK.Render("No errors")
	if m.errorStatus != "" && m.error != nil {
		output = fmt.Sprintf(
			"%s\n %s\n\n%s\n %v",
			barRendererError.Render("Error Message"),
			m.errorStatus,
			barRendererError.Render("Error"),
			m.error,
		)
	}
	return fmt.Sprintf("%s\n\n", output)
}

func (m *model) statusBar() string {
	var (
		mode string
		cmds []string
	)
	if m.modeSearch {
		mode = "SEARCH"
		cmds = []string{
			`"/": cancel search`,
			`"esc": force exit`,
		}
	} else if m.modeHelp {
		mode = "HELP"
		cmds = []string{
			`"h": cancel help`,
			`"Q": quit`,
		}
	} else if m.modeDebug {
		mode = "DEBUG"
		cmds = []string{
			`"d": cancel debug`,
			`"Q": force exit`,
		}
	} else {
		mode = "NORMAL"
		cmds = []string{
			`"/": search`,
			`"d": debug`,
			`"h": help`,
			`"q": quit`,
			`"Q": force exit`,
		}
	}

	status := strings.Join([]string{
		"  " + name,
		fmt.Sprintf("%s MODE", mode),
		strings.Join(cmds, " | "),
		"",
	}, "\t")

	err := ""
	if m.errorStatus != "" {
		err = fmt.Sprintf("ERROR: %s\t", m.errorStatus)
	}

	return barRendererStatus.Render(status) + barRendererError.Render(err)
}

func (m *model) locationBar() string {
	locationBar := barRendererLocation.Render(m.location())
	if m.modeSearch {
		locationBar += barRendererSearch.Render(fileSeparator + m.search)
	}
	return locationBar
}
