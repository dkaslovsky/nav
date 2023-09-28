package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var fileSeparator = string(filepath.Separator)

type model struct {
	path      string
	prevPath  string
	entries   []*entry
	displayed int
	exitCode  int
	exitStr   string
	error     error
	errorStr  string
	esc       *remappedEscKey
	search    string
	viewCache map[string]*cacheItem
	marks     map[int]struct{}

	c       int // Cursor column position.
	r       int // Cursor row position.
	columns int // Displayed columns.
	rows    int // Displayed columns.
	width   int // Terminal width.
	height  int // Terminal height.

	modeColor         bool
	modeDebug         bool
	modeError         bool
	modeExit          bool
	modeFollowSymlink bool
	modeHelp          bool
	modeHidden        bool
	modeList          bool
	modeMarks         bool
	modeSearch        bool
	modeSubshell      bool
	modeTrailing      bool
}

func newModel() *model {
	return &model{
		width:     80,
		height:    60,
		esc:       defaultEscRemapKey(),
		viewCache: make(map[string]*cacheItem),
		marks:     make(map[int]struct{}),

		modeColor:         true,
		modeDebug:         false,
		modeError:         false,
		modeExit:          false,
		modeFollowSymlink: false,
		modeHelp:          false,
		modeHidden:        false,
		modeList:          false,
		modeMarks:         false,
		modeSearch:        false,
		modeSubshell:      false,
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

func (m *model) selected() (*entry, error) {
	cache, ok := m.viewCache[m.path]
	if !ok {
		return nil, fmt.Errorf("cache item not found for %s", m.path)
	}
	idx, found := cache.lookupEntryIndex(m.displayIndex())
	if !found {
		return nil, errors.New("failed to map to valid entry index")
	}
	if idx > len(m.entries) {
		return nil, fmt.Errorf("invalid index %d for entries with length %d", idx, len(m.entries))
	}
	return m.entries[idx], nil
}

func (m *model) location() string {
	location := m.path
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

func (m *model) setPath(path string) {
	m.prevPath = m.path
	m.path = path
}

func (m *model) restorePath() {
	if m.prevPath != "" {
		m.path = m.prevPath
		m.prevPath = ""
	}
}

func (m *model) setError(err error, status string) {
	m.modeError = true
	m.errorStr = status
	m.error = err
}

func (m *model) setExit(exitStr string) {
	m.setExitWithCode(exitStr, 0)
}

func (m *model) setExitWithCode(exitStr string, exitCode int) {
	m.modeExit = true
	m.exitStr = exitStr
	m.exitCode = exitCode
}

func (m *model) clearError() {
	m.modeError = false
	m.errorStr = ""
	m.error = nil
}

func (m *model) clearSearch() {
	m.modeSearch = false
	m.search = ""
}

func index(c int, r int, rows int) int {
	return r + (c * rows)
}
