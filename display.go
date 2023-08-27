package main

import (
	"fmt"
	"path/filepath"
)

// displayName contains a formatted name and effective length for display in the terminal.
type displayName struct {
	name string
	len  int
}

func (d *displayName) String() string {
	return d.name
}

func (d *displayName) Len() int {
	return d.len
}

// newDisplayName constructs a displayName from an entry and provided functional options.
func newDisplayName(e *entry, opts ...displayNameOption) *displayName {
	c := &displayNameConfig{
		name:      e.Name(),
		nameExtra: "",
		trailing:  "",
		color:     colorGray,
	}

	for _, opt := range opts {
		opt(c, e.mode)
	}

	return &displayName{
		name: fmt.Sprintf("%s%s%s%s%s", c.color, c.name, colorReset, c.trailing, c.nameExtra),
		len:  len(c.name) + len(c.trailing) + len(c.nameExtra),
	}
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
	color     color
	name      string
	nameExtra string
	trailing  string
}

// displayNameOption is a functional option for setting displayNameConfig values.
type displayNameOption func(*displayNameConfig, entryMode)

func displayNameWithColor() displayNameOption {
	return func(c *displayNameConfig, mode entryMode) {
		switch {
		case mode.has(entryModeSymlink):
			c.color = colorMagenta
		case mode.has(entryModeHidden):
			c.color = colorGreen
		case mode.has(entryModeDir):
			c.color = colorCyan
		}
	}
}

func displayNameWithFollowSymlink(path string) displayNameOption {
	return func(c *displayNameConfig, mode entryMode) {
		if !mode.has(entryModeSymlink) {
			return
		}
		if followedName, err := filepath.EvalSymlinks(filepath.Join(path, c.name)); err == nil {
			c.nameExtra = fmt.Sprintf(" -> %s", followedName)
		}
	}
}

func displayNameWithTrailing() displayNameOption {
	return func(c *displayNameConfig, mode entryMode) {
		switch {
		case mode.has(entryModeSymlink):
			c.trailing = "@"
		case mode.has(entryModeDir):
			c.trailing = "/"
		}
	}
}
