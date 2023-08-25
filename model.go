package main

import (
	"os"
)

type model struct {
	path  string
	files []*entry
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
