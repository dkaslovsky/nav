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
	pathCache map[string]*cacheItem // Map path to cached state.
	marks     map[int]int           // Map display index to entry index for marked entries.

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

	hideStatusBar bool
}

func newModel() *model {
	return &model{
		width:     80,
		height:    60,
		esc:       defaultEscRemapKey(),
		pathCache: make(map[string]*cacheItem),
		marks:     make(map[int]int),

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

		hideStatusBar: false,
	}
}

func (m *model) normalMode() bool {
	return !(m.modeSearch || m.modeDebug || m.modeHelp)
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
	cache, ok := m.pathCache[m.path]
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

func (m *model) clearError() {
	m.modeError = false
	m.errorStr = ""
	m.error = nil
}

func (m *model) setExit(exitStr string) {
	m.setExitWithCode(exitStr, 0)
}

func (m *model) setExitWithCode(exitStr string, exitCode int) {
	m.modeExit = true
	m.exitStr = exitStr
	m.exitCode = exitCode
}

func (m *model) marked() bool {
	return m.markedIndex(m.displayIndex())
}

func (m *model) markedIndex(dispIdx int) bool {
	_, marked := m.marks[dispIdx]
	return marked
}

func (m *model) reloadMarks() error {
	newMarks := make(map[int]int)
	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	for _, entryIdx := range m.marks {
		if newDisplayIdx, ok := cache.lookupDisplayIndex(entryIdx); ok {
			newMarks[newDisplayIdx] = entryIdx
		}
	}
	m.marks = newMarks
	return nil
}

func (m *model) toggleMark() error {
	idx := m.displayIndex()
	if m.markedIndex(idx) {
		delete(m.marks, idx)
		m.modeMarks = len(m.marks) != 0
		return nil
	}

	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	entryIdx, ok := cache.lookupEntryIndex(idx)
	if !ok {
		return errors.New("failed to find entry index")
	}
	m.marks[idx] = entryIdx
	m.modeMarks = true
	return nil
}

func (m *model) toggleMarkAll() error {
	allMarked := true
	for i := 0; i < m.displayed; i++ {
		if _, marked := m.marks[i]; !marked {
			allMarked = false
			break
		}
	}

	if allMarked {
		m.clearMarks()
		return nil
	} else {
		return m.markAll()
	}
}

func (m *model) markAll() error {
	m.marks = make(map[int]int)
	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	for i := 0; i < m.displayed; i++ {
		if entryIdx, ok := cache.lookupEntryIndex(i); ok {
			m.marks[i] = entryIdx
		}
	}
	m.modeMarks = len(m.marks) != 0
	return nil
}

func (m *model) clearMarks() {
	m.marks = make(map[int]int)
	m.modeMarks = false
}

func (m *model) clearSearch() {
	m.modeSearch = false
	m.search = ""
}

func index(c int, r int, rows int) int {
	return r + (c * rows)
}
