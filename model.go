package main

import (
	"io/fs"
	"os"
)

type entry struct {
	fs.DirEntry
}

type model struct {
	path  string
	files []*entry
}

func (m *model) list() error {
	files, err := os.ReadDir(m.path)
	if err != nil {
		return err
	}

	m.files = nil
	for _, file := range files {
		m.files = append(m.files, &entry{file})
	}

	return nil
}
