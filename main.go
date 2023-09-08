package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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

func usage() string {
	pad := 12

	usageKeyLine := func(key key.Binding) string {
		keyStr := key.Keys()[0]
		return fmt.Sprintf("\"%s\":%s", keyStr, strings.Repeat(" ", pad-len(keyStr)))
	}

	usage := `
	%s (v%s) is an interactive terminal filesystem explorer.

	For interactive ls/cd workflows, it can be used in a bash/zsh function such as
	function lsi {
		cd "$(%s "$@")"
	}
	
	Useful key commands are listed in the status bar.

	------------------------
	| Full list of commands |
	------------------------

	Arrow keys are used to move the cursor.
	Vim navigation is enabled via "h" (left), "j" (down) "k" (up), and "l" (right). 

	%s navigates into the directory under the cursor
	%s navigates back to the previous directory

	%s enters help mode
	%s enters search mode (insert in location bar)
	%s enters debug mode  (view error details)
	%s switches back to normal mode

	%s toggles showing hidden files
	%s toggles listing full file information (ls -l)
	%s toggles following symlinks

	%s dismisses errors

	%s quits the application and outputs the current directory
	%s quits the application and outputs the path to the entry under the cursor
	%s quits the application with no output
	`

	return fmt.Sprintf(usage,
		name, version, name,
		usageKeyLine(keySelect),
		usageKeyLine(keyBack),
		usageKeyLine(keyHelp),
		usageKeyLine(keySearch),
		usageKeyLine(keyDebug),
		usageKeyLine(keyEsc),
		usageKeyLine(keyHidden),
		usageKeyLine(keyList),
		usageKeyLine(keyFollowSymlink),
		usageKeyLine(keyDismissError),
		usageKeyLine(keyQuit),
		usageKeyLine(keyQuitWithSelected),
		usageKeyLine(keyQuitForce),
	)
}

func usageAndExit() {
	fmt.Println(usage())
	os.Exit(0)
}

func versionAndExit() {
	_, _ = fmt.Fprintf(os.Stderr, "%s (v%s)", name, version)
	os.Exit(0)
}
