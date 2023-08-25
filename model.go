package main

import (
	"os"
	"strings"
)

type model struct {
	path   string
	files  []*entry
	hidden bool
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
	output := []string{}
	for _, file := range m.files {
		// Optionally do not show hidden files.
		if !m.hidden && file.IsHidden() {
			continue
		}
		output = append(output, file.displayName())
	}
	return strings.Join(output, "\n")
}
