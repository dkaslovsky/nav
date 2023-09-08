package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	name    = "nav"
	version = "0.0.1"
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

	// Terminal coloring.
	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

	// Run the app.
	_, err = tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.exitCode)
}

func parseArgs(args []string, m *model) error {
	var err error

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			usageAndExit()
		case "--version", "-v":
			versionAndExit()
		case "--no-color":
			m.modeColor = false
		case "--follow-symlinks":
			m.modeFollowSymlink = true
		case "--hidden":
			m.modeHidden = true
		case "--list":
			m.modeList = true
		case "--search":
			m.modeSearch = true
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

func usageAndExit() {
	fmt.Println(usage())
	os.Exit(0)
}

func versionAndExit() {
	_, _ = fmt.Fprintf(os.Stderr, "%s (v%s)", name, version)
	os.Exit(0)
}
