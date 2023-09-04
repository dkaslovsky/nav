package main

import (
	"os"
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
	displayNames := []*displayName{}
	for _, ent := range m.entries {
		// Optionally do not show hidden files.
		if !m.modeHidden && ent.hasMode(entryModeHidden) {
			continue
		}
		displayNames = append(displayNames, newDisplayName(ent, m.displayNameOpts()...))
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
				output[row] += render(cursorStyle, gridNames[col][row])
			} else {
				output[row] += render(normalStyle, gridNames[col][row])
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
