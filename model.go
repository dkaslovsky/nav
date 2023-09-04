package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type model struct {
	path    string
	entries []*entry

	width  int // Terminal width.
	height int // Terminal height.

	c int // Cursor column position.
	r int // Cursor row position.

	columns int // Displayed columns.
	rows    int // Displayed columns.

	cursorCache map[string]*position

	exitCode int

	modeColor         bool
	modeFollowSymlink bool
	modeHidden        bool
	modeList          bool
	modeTrailing      bool
}

func newModel() *model {
	return &model{
		width:  80,
		height: 60,

		cursorCache: make(map[string]*position),

		modeColor:         true,
		modeFollowSymlink: false,
		modeHidden:        false,
		modeList:          false,
		modeTrailing:      true,
	}
}

func (m *model) list() error {
	files, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}

	m.entries = []*entry{}
	for _, file := range files {
		ent, err := newEntry(file)
		if err != nil {
			return err
		}
		m.entries = append(m.entries, ent)
	}
	sortEntries(m.entries)

	return nil
}

func (m *model) view() string {
	output := []string{
		// First row of output is the location bar.
		barRendererLocation.Render(m.location()),
	}

	displayNames := []*displayName{}
	for _, ent := range m.entries {
		// Optionally do not show hidden files.
		if !m.modeHidden && ent.hasMode(entryModeHidden) {
			continue
		}
		displayNames = append(displayNames, newDisplayName(ent, m.displayNameOpts()...))
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

func (m *model) displayNameOpts() []displayNameOption {
	opts := []displayNameOption{}
	if m.modeColor {
		opts = append(opts, displayNameWithColor())
	}
	if m.modeFollowSymlink {
		opts = append(opts, displayNameWithFollowSymlink(m.path))
	}
	if m.modeList {
		opts = append(opts, displayNameWithList())
	}
	if m.modeTrailing {
		opts = append(opts, displayNameWithTrailing())
	}
	return opts
}

func (m *model) selected() (*entry, bool) {
	idx := index(m.c, m.r, m.rows)
	if idx > len(m.entries) {
		return nil, false
	}
	return m.entries[idx], true
}

func (m *model) location() string {
	location := m.path
	if userHomeDir, err := os.UserHomeDir(); err == nil {
		location = strings.Replace(m.path, userHomeDir, "~", 1)
	}
	if runtime.GOOS == "windows" {
		fileSeparator := string(filepath.Separator)
		location = strings.ReplaceAll(strings.Replace(location, "\\/", fileSeparator, 1), "/", fileSeparator)
	}
	return location
}

func index(c int, r int, rows int) int {
	return r + (c * rows)
}
