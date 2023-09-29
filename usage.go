package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

func usage() string {
	usage := `
	%s (%s) is a terminal filesystem explorer built for interactive ls workflows.
	
	Useful key commands are listed in the status bar.
`

	return fmt.Sprintf(usage,
		name, getVersion(),
	)
}

func commands() string {
	pad := 12

	usageKeyLine := func(text string, key key.Binding) string {
		keyStr := keyString(key)
		textParts := strings.Split(text, "\n")
		s := []string{
			fmt.Sprintf("\t\"%s\":%s%s", keyStr, strings.Repeat(" ", pad-len(keyStr)), textParts[0]),
		}
		for _, textPart := range textParts[1:] {
			s = append(s, fmt.Sprintf("\t   %s%s", strings.Repeat(" ", pad), textPart))
		}
		return strings.Join(s, "\n")
	}

	usage := `
	------------
	| Commands |
	------------

	Arrow keys are used to move the cursor.
	Vim navigation is available using "h" (left), "j" (down) "k" (up), and "l" (right).

%s
`
	cmds := []string{
		usageKeyLine("navigates into the directory or returns the\npath to the entry under the cursor", keySelect),
		usageKeyLine("navigates back to the previous directory", keyBack),
		"",
		usageKeyLine("returns the path to the entry under the cursor and quits", keyReturnSelected),
		usageKeyLine("returns the path to the current directory and quits", keyReturnDirectory),
		"",
		usageKeyLine("enters search mode (insert into the path)", keySearchMode),
		usageKeyLine("enters debug mode  (view error details)", keyDebugMode),
		usageKeyLine("enters help mode", keyHelpMode),
		usageKeyLine("switches back to normal mode", keyEsc),
		"",
		usageKeyLine("toggles showing hidden files (ls -a)", keyToggleHidden),
		usageKeyLine("toggles listing full file information (ls -l)", keyToggleList),
		usageKeyLine("toggles following symlinks", keyToggleFollowSymlink),
		"",
		usageKeyLine("dismisses errors", keyDismissError),
		usageKeyLine("quits the application with no return value", keyQuit),
	}

	return fmt.Sprintf(usage, strings.Join(cmds, "\n"))
}

func flags() string {
	pad := 25

	usageFlagLine := func(text string, flags ...string) string {
		flagStr := strings.Join(flags, ", ")
		textParts := strings.Split(text, "\n")
		s := []string{
			fmt.Sprintf("\t%s:%s%s", flagStr, strings.Repeat(" ", pad-len(flagStr)), textParts[0]),
		}
		for _, textPart := range textParts[1:] {
			s = append(s, fmt.Sprintf("\t %s%s", strings.Repeat(" ", pad), textPart))
		}
		return strings.Join(s, "\n")
	}

	usage := `
	---------
	| Flags |
	---------

%s
`
	flags := []string{
		usageFlagLine("display help", flagHelp, flagHelpShort, flagHelpShortCaps),
		usageFlagLine("display version", flagVersion, flagVersionShort),
		"",
		usageFlagLine("start in search mode", flagSearch, flagSearchShort),
		"",
		usageFlagLine("return output suitable for pipe and subshell usage", flagPipe),
		"",
		usageFlagLine("toggle on following symlinks at startup", flagFollowSymlinks, flagFollowSymlinksShort),
		usageFlagLine("toggle on showing hidden files at startup", flagHidden, flagHiddenShort),
		usageFlagLine("toggle on list mode at startup", flagList, flagListShort),
		"",
		usageFlagLine("toggle off color output", flagNoColor),
		usageFlagLine("toggle off bottom status bar menu", flagNoStatusBar),
		usageFlagLine("toggle off trailing annotators", flagNoTrailing),
		"",
		usageFlagLine("remap the escape key to the following value, using\nrepeated values to require multiple presses", flagRemapEsc),
	}
	return fmt.Sprintf(usage, strings.Join(flags, "\n"))
}
