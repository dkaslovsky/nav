package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

func usage() string {
	usage := `
	%s (v%s) is an interactive terminal filesystem explorer.

	For interactive ls/cd workflows, it can be used in a bash/zsh function such as
	function %s {
		cd "$(%s %s "$@")"
	}
	
	Useful key commands are listed in the status bar.
`

	return fmt.Sprintf(usage,
		name, version, "nv", name, flagSubshell,
	)
}

func commands() string {
	pad := 12

	usageKeyLine := func(key key.Binding, text string) string {
		keyStr := keyString(key)
		return fmt.Sprintf("\t\"%s\":%s%s", keyStr, strings.Repeat(" ", pad-len(keyStr)), text)
	}

	usage := `
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
		usageKeyLine(keyReturnSelected, "returns the path to the entry under the cursor and quits"),
		usageKeyLine(keyReturnDirectory, "returns the path to the current directory and quits"),
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
		usageKeyLine(keyQuit, "quits the application with no return"),
	}

	return fmt.Sprintf(usage, strings.Join(cmds, "\n"))
}

func flags() string {
	pad := 25

	usageFlagLine := func(text string, flagSet ...string) string {
		flagStr := strings.Join(flagSet, ", ")
		return fmt.Sprintf("\t%s:%s%s", flagStr, strings.Repeat(" ", pad-len(flagStr)), text)
	}

	usage := `
	----------------------
	| Command Line Flags |
	----------------------

	The following flags are available:

%s
`
	flags := []string{
		usageFlagLine("display help", flagHelp, flagHelpShort, flagHelpShortCaps),
		usageFlagLine("display version", flagVersion, flagVersionShort),
		"",
		usageFlagLine("start in search mode", flagSearch, flagSearchShort),
		"",
		usageFlagLine("return output suitable for subshell invocation", flagSubshell),
		"",
		usageFlagLine("toggle on showing hidden files at startup", flagHidden),
		usageFlagLine("toggle on list mode at startup", flagList, flagListShort),
		usageFlagLine("toggle on following symlinks at startup", flagFollowSymlinks, flagFollowSymlinksShort),
		"",
		usageFlagLine("toggle off color output", flagNoColor),
		usageFlagLine("toggle off trailing annotators", flagNoTrailing),
		"",
		usageFlagLine("remap the escape key to the following value", flagRemapEsc),
	}
	return fmt.Sprintf(usage, strings.Join(flags, "\n"))
}
