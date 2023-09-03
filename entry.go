package main

import (
	"io/fs"
	"os"
	"sort"
	"strings"
)

type entry struct {
	fs.DirEntry
	mode entryMode
}

func newEntry(dirEntry fs.DirEntry) *entry {
	e := &entry{
		DirEntry: dirEntry,
		mode:     entryModeNone,
	}

	// Determine if e represents a hidden file.
	// This check might not be applicable cross-platform.
	if strings.HasPrefix(e.Name(), ".") {
		e.mode = e.mode | entryModeHidden
	}

	// Set e to be a symlink even if it is also a directory or file.
	if fi, err := e.Info(); err == nil {
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			e.mode = e.mode | entryModeSymlink
			return e
		}
	}

	if e.IsDir() {
		e.mode = e.mode | entryModeDir
		return e
	}

	// Set e to be a file since it is not a symlink or a directory.
	e.mode = e.mode | entryModeFile
	return e
}

func (e *entry) hasMode(mode entryMode) bool {
	return e.mode.has(mode)
}

type entryMode uint32

const (
	entryModeNone entryMode = 1 << iota
	entryModeDir
	entryModeFile
	entryModeSymlink
	entryModeHidden
)

func (mode entryMode) has(tgt entryMode) bool {
	return mode&tgt == tgt
}

// sortEntries performs an in-place sort of a slice of entries by type (mode) and alphabetically
// within each type (mode). The ordering of types (modes) is:
// - directories
// - files
// - hidden files
func sortEntries(entries []*entry) {
	sort.Slice(entries, func(i, j int) bool {
		iEntry := entries[i]
		jEntry := entries[j]

		if iEntry.hasMode(entryModeHidden) {
			if jEntry.hasMode(entryModeHidden) {
				if iEntry.hasMode(entryModeDir) {
					if jEntry.hasMode(entryModeDir) {
						return iEntry.Name() < jEntry.Name()
					}
					return true
				}
				if jEntry.hasMode(entryModeDir) {
					return false
				}
				return iEntry.Name() < jEntry.Name()
			}
			return false
		}
		if jEntry.hasMode(entryModeHidden) {
			return true
		}

		if iEntry.hasMode(entryModeDir) {
			if jEntry.hasMode(entryModeDir) {
				return iEntry.Name() < jEntry.Name()
			}
			return true
		}
		if jEntry.hasMode(entryModeDir) {
			return false
		}

		return iEntry.Name() < jEntry.Name()
	})
}
