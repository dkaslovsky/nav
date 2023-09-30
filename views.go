package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

func (m *model) normalView() string {
	var (
		updateCache     = newCacheItem() // Cache for storing the current state as it is constructed.
		displayNames    = []*displayName{}
		displayNameOpts = m.displayNameOpts()
		displayed       = 0
		validEntries    = 0
	)

	// Construct display names from filtered entries and populate a new cache mapping between them.
	for entryIdx, ent := range m.entries {
		// Filter hidden files.
		if !m.modeHidden && ent.hasMode(entryModeHidden) {
			continue
		}

		validEntries++

		// Filter for search.
		if m.search != "" {
			if !strings.HasPrefix(ent.Name(), m.search) {
				continue
			}
		}

		displayNames = append(displayNames, newDisplayName(ent, displayNameOpts...))
		updateCache.addIndexPair(&indexPair{entry: entryIdx, display: displayed})
		displayed++
	}

	if validEntries == 0 {
		return m.locationBar() + "\n\n\t(no entries)\n"
	}

	if m.modeSearch || m.search != "" {
		if displayed == 0 && validEntries > 0 {
			return m.locationBar() + "\n\n\t(no matching entries)\n"
		}
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
		if entryIdx, entryFound := cache.lookupEntryIndex(cache.cursorIndex()); entryFound {
			// Use the entry index to get the current display index.
			if dispIdx, dispFound := updateCache.lookupDisplayIndex(entryIdx); dispFound {
				// Set the cursor position using the current display index and layout.
				updateCursorPosition = newPositionFromIndex(dispIdx, layout.rows)
			}
		}
	}

	// Update the cache.
	updateCache.setPosition(updateCursorPosition)
	updateCache.setColumns(layout.columns)
	updateCache.setRows(layout.rows)

	// Update the model.
	m.viewCache[m.path] = updateCache
	m.displayed = displayed
	m.columns = layout.columns
	m.rows = layout.rows
	m.setCursor(updateCursorPosition)
	if m.c >= m.columns || m.r > m.rows {
		m.resetCursor()
	}

	// Render entry names in grid.
	gridOutput := make([]string, layout.rows)
	for row := 0; row < layout.rows; row++ {
		for col := 0; col < layout.columns; col++ {
			if col == m.c && row == m.r {
				if m.marked() {
					gridOutput[row] += cursorRendererSelectedMarked.Render(gridNames[col][row])
				} else {
					gridOutput[row] += cursorRendererSelected.Render(gridNames[col][row])
				}
			} else {
				if m.markedIndex(index(col, row, layout.rows)) {
					gridOutput[row] += cursorRendererMarked.Render(gridNames[col][row])
				} else {
					gridOutput[row] += cursorRendererNormal.Render(gridNames[col][row])
				}
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
	if m.modeError {
		output = fmt.Sprintf(
			"%s\n %s\n\n%s\n %v",
			barRendererError.Render("Error Message"),
			m.errorStr,
			barRendererError.Render("Error"),
			m.error,
		)
	}
	return fmt.Sprintf("%s\n\n", output)
}

type statusBarItem string

func (s statusBarItem) String() string { return string(s) }
func (s statusBarItem) Len() int       { return len(s) }

func (m *model) statusBar() string {
	const rows = 2

	var (
		mode string
		cmds []statusBarItem
	)

	if m.modeDebug {
		mode = "DEBUG"
		cmds = []statusBarItem{
			statusBarItem(fmt.Sprintf(`"%s": dismiss error`, keyString(keyDismissError))),
			statusBarItem(fmt.Sprintf(`"%s": normal mode`, keyString(keyEsc))),
		}
	} else if m.modeSearch {
		mode = "SEARCH"
		cmds = []statusBarItem{
			statusBarItem(fmt.Sprintf(`"%s": complete`, keyString(keyTab))),
			statusBarItem(fmt.Sprintf(`"%s": normal mode`, keyString(keyEsc))),
		}
	} else if m.modeHelp {
		mode = "HELP"
		cmds = []statusBarItem{
			statusBarItem(fmt.Sprintf(`"%s": normal mode`, keyString(keyEsc))),
		}
	} else {
		mode = "NORMAL"
		cmds = []statusBarItem{
			statusBarItem(fmt.Sprintf(`"%s": search`, keyString(keyModeSearch))),
			statusBarItem(fmt.Sprintf(`"%s": help`, keyString(keyModeHelp))),
			statusBarItem(fmt.Sprintf(`"%s": multiselect`, keyString(keyMark))),
		}
	}

	globalCmds := []statusBarItem{
		statusBarItem(fmt.Sprintf(`"%s": quit`, keyString(keyQuit))),
		statusBarItem(fmt.Sprintf(`"%s": return dir`, keyString(keyReturnDirectory))),
	}
	if m.modeMarks {
		globalCmds = append(globalCmds, statusBarItem(fmt.Sprintf(`"%s": return marked`, keyString(keyReturnSelected))))
	} else {
		globalCmds = append(globalCmds, statusBarItem(fmt.Sprintf(`"%s": return cursor`, keyString(keyReturnSelected))))
	}

	columns := max(len(cmds), len(globalCmds))
	for len(cmds) < columns {
		cmds = append(cmds, statusBarItem(""))
	}
	for len(globalCmds) < columns {
		globalCmds = append(globalCmds, statusBarItem(""))
	}
	cmds = append(cmds, globalCmds...)
	gridItems := gridRowMajorFixedLayout(cmds, columns, rows)

	nameAndMode := fmt.Sprintf(" %s   %s MODE  |", name, mode)
	output := strings.Join([]string{
		barRendererStatus.Render(
			fmt.Sprintf("%s\t%s\t",
				nameAndMode,
				strings.Join(gridItems[0], "\t\t"),
			),
		),
		barRendererStatus.Render(
			fmt.Sprintf("%s|\t%s\t",
				strings.Repeat(" ", len(nameAndMode)-1),
				strings.Join(gridItems[1], "\t\t"),
			),
		),
	}, "\n")

	return output
}

func (m *model) locationBar() string {
	err := ""
	if m.modeError {
		err = fmt.Sprintf(
			"\tERROR (\"%s\": dismiss, \"%s\": debug): %s",
			keyString(keyDismissError),
			keyString(keyModeDebug),
			m.errorStr,
		)
		return barRendererError.Render(err + "\t\t")
	}

	locationBar := barRendererLocation.Render(m.location())
	if m.modeSearch || m.search != "" {
		if m.path != fileSeparator {
			locationBar += barRendererSearch.Render(fileSeparator + m.search)
		}
	}
	return locationBar
}

func keyString(key key.Binding) string {
	return key.Keys()[0]
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
