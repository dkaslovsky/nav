package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

func usage() string {
	pad := 12

	usageKeyLine := func(key key.Binding, text string) string {
		keyStr := keyString(key)
		return fmt.Sprintf("\t\"%s\":%s%s", keyStr, strings.Repeat(" ", pad-len(keyStr)), text)
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
	Vim navigation is available using "h" (left), "j" (down) "k" (up), and "l" (right).

%s
`
	cmds := []string{
		usageKeyLine(keySelect, "navigates into the directory under the cursor"),
		usageKeyLine(keyBack, "navigates back to the previous directory"),
		"",
		usageKeyLine(keyQuitWithSelected, "returns the path to the entry under the cursor and quits"),
		usageKeyLine(keyQuitWithDirectory, "returns the path to the current directory and quits"),
		"",
		usageKeyLine(keyHelpMode, "enters help mode"),
		usageKeyLine(keySearchMode, "enters search mode (insert in location bar)"),
		usageKeyLine(keyDebugMode, "enters debug mode  (view error details)"),
		usageKeyLine(keyEsc, "switches back to normal mode"),
		"",
		usageKeyLine(keyToggleHidden, "toggles showing hidden files"),
		usageKeyLine(keyToggleList, "toggles listing full file information (ls -l)"),
		usageKeyLine(keyToggleFollowSymlink, "toggles following symlinks"),
		"",
		usageKeyLine(keyDismissError, "dismisses errors"),
		usageKeyLine(keyQuitForce, "quits the application with no return"),
	}

	return fmt.Sprintf(usage,
		name, version, name, strings.Join(cmds, "\n"),
	)
}
