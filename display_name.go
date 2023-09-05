package main

import (
	"fmt"
	"io/fs"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
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
		listInfo:  "",
	}

	for _, opt := range opts {
		opt(c, e.mode, e.info)
	}

	return &displayName{
		name: fmt.Sprintf("%s%s%s%s%s%s", c.listInfo, c.color, c.name, colorReset, c.trailing, c.nameExtra),
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
	listInfo  string
}

// displayNameOption is a functional option for setting displayNameConfig values.
type displayNameOption func(*displayNameConfig, entryMode, fs.FileInfo)

func displayNameWithColor() displayNameOption {
	return func(c *displayNameConfig, mode entryMode, info fs.FileInfo) {
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
	return func(c *displayNameConfig, mode entryMode, info fs.FileInfo) {
		if !mode.has(entryModeSymlink) {
			return
		}
		if followedName, err := filepath.EvalSymlinks(filepath.Join(path, c.name)); err == nil {
			c.nameExtra = fmt.Sprintf(" -> %s", followedName)
		}
	}
}

func displayNameWithList() displayNameOption {
	return func(c *displayNameConfig, mode entryMode, info fs.FileInfo) {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return
		}

		usr, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
		if err != nil {
			return
		}

		grp, err := user.LookupGroupId(strconv.FormatUint(uint64(stat.Gid), 10))
		if err != nil {
			return
		}

		c.listInfo = fmt.Sprintf(
			"%10s %8s %8s %8s %14s  ",
			info.Mode(),
			usr.Username,
			grp.Name,
			byteCountSI(info.Size()),
			formatModTime(info.ModTime(), time.Now().Year()),
		)
	}
}

func displayNameWithTrailing() displayNameOption {
	return func(c *displayNameConfig, mode entryMode, info fs.FileInfo) {
		switch {
		case mode.has(entryModeSymlink):
			c.trailing = "@"
		case mode.has(entryModeDir):
			c.trailing = "/"
		}
	}
}

func formatModTime(t time.Time, currentYear int) string {
	layout := "Jan 02 2006"
	if t.Year() == currentYear {
		layout = "Jan 02 15:04"
	}
	return t.Format(layout)
}

// Adapted from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/.
func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}
