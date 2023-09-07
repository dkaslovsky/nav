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
Useful key commands are listed below in the status bar.

<br/>

### Full list of commands

- Arrow keys are used to move the cursor.

<br/>


- "enter":     navigates into the directory under the cursor (no action for files, yet)
- "backspace": navigates back to the previous directory

<br/>

- "h": enters/exits help mode
- "d": enters/exits debug mode
- "/": enters/exits search mode

<br/>


- "a": toggles showing hidden files
- "l": toggles listing full file information (ls -l)
- "s": toggles following symlinks

<br/>


- "q":   quits the application and outputs the current directory
- "c":   quits the application and outputs the path to the entry under the cursor
- "Q":   force quits the application (exit 2) when not in search mode
- "esc": force quits the application (exit 2) in any mode

<br/>


`nav` was originally inspired by https://github.com/antonmedv/walk but has deviated significantly and has been written from the ground up to support a different set of of features.
