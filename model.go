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

func (m *model) list() error {
	files, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}

	m.files = []*entry{}
	for _, file := range files {
		m.files = append(m.files, &entry{file})
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

	output := []string{}
	for _, file := range m.files {
		// Optionally do not show hidden files.
		if !m.modeHidden && file.IsHidden() {
			continue
		}
		output = append(output, file.displayName(displayNameOpts...))
	}
	return strings.Join(output, "\n")
}
