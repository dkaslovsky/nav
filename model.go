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

	displayIndex map[int]int // Map displayed entry index to entry index.

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

		displayIndex: make(map[int]int),

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

func (m *model) selected() (*entry, bool) {
	idx, found := m.displayIndex[index(m.c, m.r, m.rows)]
	if !found || idx > len(m.entries) {
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

func index(c int, r int, rows int) int {
	return r + (c * rows)
}
