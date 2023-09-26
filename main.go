package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Name of the application.
const name = "nav"

// Version is set with ldflags.
var version string

const (
	flagHelp                = "--help"
	flagHelpShort           = "-h"
	flagHelpShortCaps       = "-H"
	flagVersion             = "--version"
	flagVersionShort        = "-v"
	flagSearch              = "--search"
	flagSearchShort         = "-s"
	flagPipe                = "--pipe"
	flagFollowSymlinks      = "--follow"
	flagFollowSymlinksShort = "-f"
	flagHidden              = "--hidden"
	flagHiddenShort         = "-a"
	flagList                = "--list"
	flagListShort           = "-l"
	flagNoColor             = "--no-color"
	flagNoTrailing          = "--no-trailing"
	flagRemapEsc            = "--remap-esc"
)

func main() {
	var err error

	// Initialize model with defaults.
	m := newModel()

	// Set model options from args.
	err = parseArgs(os.Args[1:], m)
	if err != nil {
		exit(err, m.exitCode)
	}

	// Populate the model.
	err = m.list()
	if err != nil {
		exit(err, m.exitCode)
	}

	// Terminal coloring.
	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

	// Run the app.
	_, err = tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		exit(err, m.exitCode)
	}

	exit(nil, m.exitCode)
}

func parseArgs(args []string, m *model) error {
	var err error

	i := 0
	for i < len(args) {
		arg := args[i]

		switch arg {
		case flagHelp, flagHelpShort, flagHelpShortCaps:
			usageAndExit()
		case flagVersion, flagVersionShort:
			versionAndExit()
		case flagHidden, flagHiddenShort:
			m.modeHidden = true
		case flagList, flagListShort:
			m.modeList = true
		case flagSearch, flagSearchShort:
			m.modeSearch = true
		case flagPipe:
			m.modeSubshell = true
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

func exit(err error, code int) {
	if err != nil {
		fmt.Printf("fatal: %v", err)
		if code == 0 {
			os.Exit(1)
		}
	}
	os.Exit(code)
}

func usageAndExit() {
	fmt.Printf("%s\n%s\n%s\n", usage(), commands(), flags())
	os.Exit(0)
}

func versionAndExit() {
	fmt.Printf("%s (%s)", name, getVersion())
	os.Exit(0)
}

func getVersion() string {
	if version == "" {
		return "development"
	}
	return "v" + version
}
