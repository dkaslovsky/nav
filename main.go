package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var err error

	opts := newOptions()
	err = parseArgs(os.Args[1:], &opts)
	if err != nil {
		log.Fatal(err)
	}

	m := &model{
		path:              opts.startPath,
		modeColor:         opts.modeColor,
		modeHidden:        opts.modeHidden,
		modeFollowSymlink: opts.modeFollowSymlink,
		modeTrailing:      opts.modeTrailing,
	}

	err = m.list()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

// options are configuration options set from the command line.
type options struct {
	startPath         string
	modeColor         bool
	modeHidden        bool
	modeFollowSymlink bool
	modeTrailing      bool
}

// newOptions return options with default values.
func newOptions() options {
	return options{
		modeColor:         true,
		modeHidden:        false,
		modeFollowSymlink: false,
		modeTrailing:      true,
	}
}

func parseArgs(args []string, opts *options) error {
	var err error

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			usage()
		case "--version", "-v":
			version()
		case "--no-color":
			opts.modeColor = false
		case "--hidden":
			opts.modeHidden = true
		case "--follow-symlinks":
			opts.modeFollowSymlink = true
		case "--no-trailing":
			opts.modeTrailing = false
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown flag: %s", arg)
			}
			opts.startPath, err = filepath.Abs(arg)
			if err != nil {
				return err
			}
		}
	}

	if opts.startPath == "" {
		opts.startPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	return nil
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage todo...")
	os.Exit(0)
}

func version() {
	_, _ = fmt.Fprintf(os.Stderr, "version todo...")
	os.Exit(0)
}
