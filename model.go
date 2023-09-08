package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var fileSeparator = string(filepath.Separator)

type model struct {
	path        string
	entries     []*entry
	displayed   int
	exitCode    int
	errorStatus string
	error       error

	width  int // Terminal width.
	height int // Terminal height.

	c int // Cursor column position.
	r int // Cursor row position.

	columns int // Displayed columns.
	rows    int // Displayed columns.

	viewCache map[string]*cacheItem

	search string

	modeColor         bool
	modeDebug         bool
	modeFollowSymlink bool
	modeHelp          bool
	modeHidden        bool
	modeList          bool
	modeSearch        bool
	modeTrailing      bool
}

func newModel() *model {
	return &model{
		width:  80,
		height: 60,

		viewCache: make(map[string]*cacheItem),

		modeColor:         true,
		modeDebug:         false,
		modeFollowSymlink: false,
		modeHelp:          false,
		modeHidden:        false,
		modeList:          false,
		modeSearch:        false,
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
	cache, ok := m.viewCache[m.path]
	if !ok {
		return nil, false
	}
	idx, found := cache.displayToEntityIndex[m.displayIndex()]
	if !found || idx > len(m.entries) {
		return nil, false
	}
	return m.entries[idx], true
}

func (m *model) location() string {
	location := m.path
	// TODO: encapsulate this fix
	if strings.HasPrefix(location, "//") {
		location = location[1:]
	}
	if userHomeDir, err := os.UserHomeDir(); err == nil {
		location = strings.Replace(m.path, userHomeDir, "~", 1)
	}
	if runtime.GOOS == "windows" {
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

func (m *model) displayIndex() int {
	return index(m.c, m.r, m.rows)
}

func (m *model) clearError() {
	m.errorStatus = ""
	m.error = nil
}

func (m *model) clearSearch() {
	m.modeSearch = false
	m.search = ""
}

func index(c int, r int, rows int) int {
	return r + (c * rows)
}
