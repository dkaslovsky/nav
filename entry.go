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
}

func (e *entry) IsHidden() bool {
	return strings.HasPrefix(e.Name(), ".")
}

func (e *entry) IsSymlink() bool {
	fi, err := e.Info()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

type color string

const (
	colorCyan    color = "\033[36m"
	colorGreen   color = "\033[32m"
	colorGray    color = "\033[37m"
	colorMagenta color = "\033[35m"
	colorReset   color = "\033[0m"
)

// displayNameConfig contains configuration values for constructing an entry's display name.
type displayNameConfig struct {
	color color
	name  string
}

// displayNameOption is a functional option for setting displayNameConfig values.
type displayNameOption func(*displayNameConfig, *entry)

func displayNameWithColor() displayNameOption {
	return func(c *displayNameConfig, e *entry) {
		if e.IsSymlink() {
			c.color = colorMagenta
			return
		}
		if e.IsHidden() {
			c.color = colorGreen
			return
		}
		if e.IsDir() {
			c.color = colorCyan
			return
		}
	}
}

func displayNameWithFollowSymlink(path string) displayNameOption {
	return func(c *displayNameConfig, e *entry) {
		if !e.IsSymlink() {
			return
		}
		if followedName, err := filepath.EvalSymlinks(filepath.Join(path, e.Name())); err == nil {
			c.name = fmt.Sprintf("%s%s -> %s", e.Name(), colorReset, followedName)
		}
	}
}

// displayName returns a formatted name for display in the terminal.
func (e *entry) displayName(opts ...displayNameOption) string {
	c := &displayNameConfig{
		name:  e.Name(),
		color: colorGray,
	}

	for _, opt := range opts {
		opt(c, e)
	}

	return fmt.Sprintf("%s%s%s", c.color, c.name, colorReset)
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
