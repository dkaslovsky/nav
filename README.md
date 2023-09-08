# nav
Terminal navigator for interactive `ls` and `cd` workflows.

<br/>



`nav` is an interactive terminal filesystem explorer.

For interactive ls/cd workflows, it can be used in a bash/zsh function such as
```bash
function lsi {
	cd "$(nav "$@")"
}
```

<br/>

### Full list of commands

	Arrow keys are used to move the cursor.
	Vim navigation is available using "h" (left), "j" (down) "k" (up), and "l" (right).

	"enter":       navigates into the directory under the cursor
	"backspace":   navigates back to the previous directory

	"ctrl+x":      returns the path to the entry under the cursor and quits
	"ctrl+d":      returns the path to the current directory and quits

	"H":           enters help mode
	"i":           enters search mode (insert in location bar)
	"d":           enters debug mode  (view error details)
	"esc":         switches back to normal mode

	"a":           toggles showing hidden files
	"L":           toggles listing full file information (ls -l)
	"f":           toggles following symlinks

	"e":           dismisses errors
	"ctrl+c":      quits the application with no return
<br/>


`nav` was originally inspired by https://github.com/antonmedv/walk but has deviated significantly and has been written from the ground up to support a different set of features.
