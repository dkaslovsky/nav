# nav
Terminal navigator for interactive `ls` workflows.

<br/>



`nav` is a terminal filesystem explorer built for interactive ls workflows.

It can be used in bash/zsh functions such as
```bash
# interactive `ls` + `cd`
function nv {
	cd "$(nav --subshell "$@")"
}
```
or 
```bash
# interactive `ls` to copy selected to the clipboard
function ncopy {
	nav --subshell "$@" | pbcopy
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

	"a":           toggles showing hidden files (ls -a)
	"L":           toggles listing full file information (ls -l)
	"f":           toggles following symlinks

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
