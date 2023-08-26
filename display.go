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

// toDisplayName builds a displayName from an entry.
func newDisplayName(e *entry, opts ...displayNameOption) *displayName {
	c := &displayNameConfig{
		name:      e.Name(),
		nameExtra: "",
		trailing:  "",
		color:     colorGray,
	}

	for _, opt := range opts {
		opt(c, e)
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
			c.nameExtra = fmt.Sprintf(" -> %s", followedName)
		}
	}
}

func displayNameWithTrailing() displayNameOption {
	return func(c *displayNameConfig, e *entry) {
		if e.IsSymlink() {
			c.trailing = "@"
			return
		}
		if e.IsDir() {
			c.trailing = "/"
			return
		}
	}
}
