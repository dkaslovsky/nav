# nav
Terminal navigator for interactive `ls` and `cd` workflows.

<br/>



`nav` is an interactive terminal filesystem explorer.

For interactive ls/cd workflows, it can be used in a bash/zsh function such as
```bash
function nv {
	cd "$(nav --subshell "$@")"
}
```

<br/>

### Full list of commands

	Arrow keys are used to move the cursor.
	Vim navigation is available using "h" (left), "j" (down) "k" (up), and "l" (right).

	"enter":       navigates into the directory or returns the
	               path to the entry under the cursor
	"backspace":   navigates back to the previous directory

	"ctrl+x":      returns the path to the entry under the cursor and quits
	"ctrl+d":      returns the path to the current directory and quits

	"i":           enters search mode (insert into the path)
	"d":           enters debug mode  (view error details)
	"H":           enters help mode
	"esc":         switches back to normal mode

	"a":           toggles showing hidden files
	"f":           toggles following symlinks
	"L":           toggles listing full file information (ls -l)

	"e":           dismisses errors
	"ctrl+c":      quits the application with no return value

<br/>

### Command line flags

	The following flags are available:

	--help, -h, -H:           display help
	--version, -v:            display version

	--search, -s:             start in search mode

	--subshell:               return output suitable for subshell invocation

	--follow-symlinks, -f:    toggle on following symlinks at startup
	--list, -l:               toggle on list mode at startup
	--hidden:                 toggle on showing hidden files at startup

	--no-color:               toggle off color output
	--no-trailing:            toggle off trailing annotators

	--remap-esc:              remap the escape key to the following value, using
	                          repeated values to require multiple presses
<br/>

In the future, `nav` might support a wider range of `ls` options.

`nav` was originally inspired by https://github.com/antonmedv/walk but has deviated significantly and has been written from the ground up to support a different set of features.
