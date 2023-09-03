package main

import (
	"os"
	"strings"
)

type model struct {
	path  string
	files []*entry

	width  int // Terminal width.
	height int // Terminal height.

	c int // Cursor column position.
	r int // Cursor row position.

	columns int // Displayed columns.
	rows    int // Displayed columns.

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

	m.files = []*entry{}
	for _, file := range files {
		ent, err := newEntry(file)
		if err != nil {
			return err
		}
		m.files = append(m.files, ent)
	}
	sortEntries(m.files)

	return nil
}

func (m *model) view() string {
	displayNames := []*displayName{}
	for _, file := range m.files {
		// Optionally do not show hidden files.
		if !m.modeHidden && file.hasMode(entryModeHidden) {
			continue
		}
		displayNames = append(displayNames, newDisplayName(file, m.displayNameOpts()...))
	}

	var gridNames [][]string
	var layout gridLayout

	if m.modeList {
		gridNames, layout = gridSingleColumn(displayNames, m.width, m.height)
	} else {
		gridNames, layout = gridMultiColumn(displayNames, m.width, m.height)
	}

	m.columns = layout.columns
	m.rows = layout.rows

	if m.c >= m.columns {
		m.c = 0
	}
	if m.r >= m.rows {
		m.r = 0
	}

	output := make([]string, layout.rows)
	for row := 0; row < layout.rows; row++ {
		for col := 0; col < layout.columns; col++ {
			if col == m.c && row == m.r {
				output[row] += cursor.Render(gridNames[col][row]) + separator[:len(separator)-len(cursorStr)-1]
			} else {
				output[row] += gridNames[col][row] + separator
			}
		}
	}

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
