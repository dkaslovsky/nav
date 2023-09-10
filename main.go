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

const (
	flagHelp          = "--help"
	flagHelpShort     = "-h"
	flagHelpShortCaps = "-H"
	flagVersion       = "--version"
	flagVersionShort  = "-v"

	flagSearch      = "--search"
	flagSearchShort = "-s"

	flagFollowSymlinks      = "--follow-symlinks"
	flagFollowSymlinksShort = "-f"
	flagHidden              = "--hidden"
	flagList                = "--list"
	flagListShort           = "-l"

	flagNoColor    = "--no-color"
	flagNoTrailing = "--no-trailing"

	flagRemapEsc = "--remap-esc"
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

	// for _, arg := range args {

	i := 0
	for i < len(args) {
		arg := args[i]

		switch arg {
		case flagHelp, flagHelpShort, flagHelpShortCaps:
			usageAndExit()
		case flagVersion, flagVersionShort:
			versionAndExit()
		case flagHidden:
			m.modeHidden = true
		case flagList, flagListShort:
			m.modeList = true
		case flagSearch, flagSearchShort:
			m.modeSearch = true
		case flagFollowSymlinks, flagFollowSymlinksShort:
			m.modeFollowSymlink = true
		case flagNoColor:
			m.modeColor = false
		case flagNoTrailing:
			m.modeTrailing = false
		case flagRemapEsc:
			if i > len(args)-2 {
				return fmt.Errorf("%s must be followed by a string value", flagRemapEsc)
			}
			err := m.setEscRemapKey(args[i+1])
			if err != nil {
				return err
			}
			i += 2
			continue
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown flag: %s", arg)
			}
			m.path, err = filepath.Abs(arg)
			if err != nil {
				return err
			}
		}

		i++
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
	fmt.Printf("%s\n%s\n%s\n", usage(), commands(), flags())
	os.Exit(0)
}

func versionAndExit() {
	_, _ = fmt.Fprintf(os.Stderr, "%s (v%s)", name, version)
	os.Exit(0)
}
