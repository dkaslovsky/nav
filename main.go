package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
			m.modeHelp = true
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

func usage() string {
	usage := `
	%s (v%s) is an interactive terminal filesystem explorer.

	For interactive ls/cd workflows, it can be used in a bash/zsh function such as
	function lsi {
		cd "$(%s "$@")"
	}
	
	Useful key commands are listed below in the status bar.

	------------------------
	| Full list of commands |
	------------------------

	Arrow keys are used to move the cursor.

	"enter":     navigates into the directory under the cursor (no action for files, yet)
	"backspace": navigates back to the previous directory

	"h": enters/exits help mode
	"d": enters/exits debug mode
	"/": enters/exits search mod

	"a": toggles showing hidden files
	"l": toggles listing full file information (ls -l)
	"s": toggles following symlinks

	"q":   quits the application and outputs the current directory
	"c":   quits the application and outputs the path to the entry under the cursor
	"Q":   force quits the application (exit 2) when not in search mode
	"esc": force quits the application (exit 2) in any mode
	`
	return fmt.Sprintf(usage, name, version, name)
}

func versionAndExit() {
	_, _ = fmt.Fprintf(os.Stderr, "%s (v%s)", name, version)
	os.Exit(0)
}
