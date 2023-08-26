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

	// Initialize model with defaults.
	m := newModel()

	// Set model options from args.
	err = parseArgs(os.Args[1:], m)
	if err != nil {
		log.Fatal(err)
	}

	// Populate the model.
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

func parseArgs(args []string, m *model) error {
	var err error

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			usage()
		case "--version", "-v":
			version()
		case "--no-color":
			m.modeColor = false
		case "--hidden":
			m.modeHidden = true
		case "--follow-symlinks":
			m.modeFollowSymlink = true
		case "--no-trailing":
			m.modeTrailing = false
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown flag: %s", arg)
			}
			m.path, err = filepath.Abs(arg)
			if err != nil {
				return err
			}
		}
	}

	if m.path == "" {
		m.path, err = os.Getwd()
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
