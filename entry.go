package main

import (
	"io/fs"
	"sort"
	"strings"
)

type entry struct {
	fs.DirEntry
}

func (e *entry) IsHidden() bool {
	return strings.HasPrefix(e.Name(), ".")
}

// sortEntriesByType performs an in-place sort of a slice of entries by type and alphabetically within
// each type. The ordering of types is:
// - directories
// - files
// - hidden files
func sortEntriesByType(entries []*entry) {
	sort.Slice(entries, func(i, j int) bool {
		iEntry := entries[i]
		jEntry := entries[j]

		if iEntry.IsHidden() {
			if jEntry.IsHidden() {
				if iEntry.IsDir() {
					if jEntry.IsDir() {
						return iEntry.Name() < jEntry.Name()
					}
					return true
				}
				if jEntry.IsDir() {
					return false
				}
				return iEntry.Name() < jEntry.Name()
			}
			return false
		}
		if jEntry.IsHidden() {
			return true
		}

		if iEntry.IsDir() {
			if jEntry.IsDir() {
				return iEntry.Name() < jEntry.Name()
			}
			return true
		}
		if jEntry.IsDir() {
			return false
		}

		return iEntry.Name() < jEntry.Name()
	})
}
