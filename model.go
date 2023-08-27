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

	modeColor         bool
	modeHidden        bool
	modeFollowSymlink bool
	modeTrailing      bool
}

func newModel() *model {
	return &model{
		width:  80,
		height: 60,

		modeColor:         true,
		modeHidden:        false,
		modeFollowSymlink: false,
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
		m.files = append(m.files, newEntry(file))
	}
	sortEntriesByType(m.files)

	return nil
}

func (m *model) view() string {
	displayNameOpts := []displayNameOption{}
	if m.modeColor {
		displayNameOpts = append(displayNameOpts, displayNameWithColor())
	}
	if m.modeFollowSymlink {
		displayNameOpts = append(displayNameOpts, displayNameWithFollowSymlink(m.path))
	}
	if m.modeTrailing {
		displayNameOpts = append(displayNameOpts, displayNameWithTrailing())
	}

	displayNames := []*displayName{}
	for _, file := range m.files {
		// Optionally do not show hidden files.
		if !m.modeHidden && file.hasMode(entryModeHidden) {
			continue
		}
		displayNames = append(displayNames, newDisplayName(file, displayNameOpts...))
	}

	gridNames, layout := grid(displayNames, m.width, m.height)

	output := make([]string, layout.rows)
	for row := 0; row < layout.rows; row++ {
		for col := 0; col < layout.columns; col++ {
			output[row] += gridNames[col][row] + separator
		}
	}

	return strings.Join(output, "\n")
}
