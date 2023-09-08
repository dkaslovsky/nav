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

	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

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
	pad := 12
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
	Vim navigation is enabled via "h" (left), "j" (down) "k" (up), and "l" (right). 

	"%s":%snavigates into the directory under the cursor
	"%s":%snavigates back to the previous directory

	"%s":%senters help mode
	"%s":%senters search mode (insert in location bar)
	"%s":%senters debug mode  (view error details)
	"%s":%sswitches back to normal mode

	"%s":%stoggles showing hidden files
	"%s":%stoggles listing full file information (ls -l)
	"%s":%stoggles following symlinks

	"%s":%sdismisses errors

	"%s":%squits the application and outputs the current directory
	"%s":%squits the application and outputs the path to the entry under the cursor
	"%s":%squits the application with no output
	`
	return fmt.Sprintf(usage,
		name, version, name,

		keySelect.Keys()[0],
		strings.Repeat(" ", pad-len(keySelect.Keys()[0])),

		keyBack.Keys()[0],
		strings.Repeat(" ", pad-len(keyBack.Keys()[0])),

		keyHelp.Keys()[0],
		strings.Repeat(" ", pad-len(keyHelp.Keys()[0])),

		keySearch.Keys()[0],
		strings.Repeat(" ", pad-len(keySearch.Keys()[0])),

		keyDebug.Keys()[0],
		strings.Repeat(" ", pad-len(keyDebug.Keys()[0])),

		keyEsc.Keys()[0],
		strings.Repeat(" ", pad-len(keyEsc.Keys()[0])),

		keyHidden.Keys()[0],
		strings.Repeat(" ", pad-len(keyHidden.Keys()[0])),

		keyList.Keys()[0],
		strings.Repeat(" ", pad-len(keyList.Keys()[0])),

		keyFollowSymlink.Keys()[0],
		strings.Repeat(" ", pad-len(keyFollowSymlink.Keys()[0])),

		keyDismissError.Keys()[0],
		strings.Repeat(" ", pad-len(keyDismissError.Keys()[0])),

		keyQuit.Keys()[0],
		strings.Repeat(" ", pad-len(keyQuit.Keys()[0])),

		keyQuitWithSelected.Keys()[0],
		strings.Repeat(" ", pad-len(keyQuitWithSelected.Keys()[0])),

		keyQuitForce.Keys()[0],
		strings.Repeat(" ", pad-len(keyQuitForce.Keys()[0])),
	)
}

func versionAndExit() {
	_, _ = fmt.Fprintf(os.Stderr, "%s (v%s)", name, version)
	os.Exit(0)
}
