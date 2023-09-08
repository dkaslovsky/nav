package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

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
