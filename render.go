package main

import "github.com/charmbracelet/lipgloss"

const (
	cursorStr = ">"
	emptyStr  = " "
)

var (
	cursorStyle = lipgloss.NewStyle().Bold(true).SetString(cursorStr)
	normalStyle = lipgloss.NewStyle().SetString(emptyStr)
)

func render(style lipgloss.Style, name string) string {
	// return nonCursor.Render(gridNames[col][row]) + separator[:len(separator)-len(emptyStr)-1]
	return style.Render(name) + separator[:separatorLen-2]
}
