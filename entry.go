package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type entry struct {
	fs.DirEntry
	mode entryMode
	info fs.FileInfo
}

func newEntry(dirEntry fs.DirEntry) (*entry, error) {
	e := &entry{
		DirEntry: dirEntry,
	}

	var err error
	e.info, err = dirEntry.Info()
	if err != nil {
		return e, err
	}

	e.setMode()
	return e, nil
}

const maskExec = 0o111

func (e *entry) setMode() {
	e.mode = entryModeNone

	// Determine if e represents a hidden file.
	// This check might not be applicable cross-platform.
	if strings.HasPrefix(e.Name(), ".") {
		e.mode = e.mode | entryModeHidden
	}

	var isExec bool
	if fi, err := e.Info(); err == nil {
		// Set e to be a symlink even if it is also a directory or file.
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			e.mode = e.mode | entryModeSymlink
			return
		}
		// Check if e is executable but do not set this mode until after confirming it is a file below.
		if fi.Mode()&maskExec == maskExec {
			isExec = true
		}
	}

	if e.IsDir() {
		e.mode = e.mode | entryModeDir
		return
	}

	// Set e to be a file since it is not a symlink or a directory.
	e.mode = e.mode | entryModeFile
	if isExec {
		e.mode = e.mode | entryModeExec
	}
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
	entryModeExec
)

func (mode entryMode) has(tgt entryMode) bool {
	return mode&tgt == tgt
}

type symlink struct {
	absPath string
	info    fs.FileInfo
}

func followSymlink(path string, e *entry) (*symlink, error) {
	if !e.hasMode(entryModeSymlink) {
		return nil, fmt.Errorf("cannot follow non-symlink entry: %s", e.Name())
	}
	followed, err := filepath.EvalSymlinks(filepath.Join(path, e.Name()))
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(followed)
	if err != nil {
		return nil, err
	}
	return &symlink{
		absPath: followed,
		info:    info,
	}, nil
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
