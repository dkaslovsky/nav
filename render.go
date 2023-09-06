package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorRendererSelected = newCursorRenderer(lipgloss.NewStyle().Bold(true).SetString(">"))
	cursorRendererNormal   = newCursorRenderer(lipgloss.NewStyle().SetString(" "))

	barRendererLocation = lipgloss.NewStyle().Background(lipgloss.Color("#5C5C5C")).Foreground(lipgloss.Color("#FFFFFF"))
	barRendererSearch   = lipgloss.NewStyle().Background(lipgloss.Color("#499F1C")).Foreground(lipgloss.Color("#FFFFFF"))
)

type cursorRenderer struct {
	style lipgloss.Style
	pad   string
}

func newCursorRenderer(style lipgloss.Style) *cursorRenderer {
	pad := ""
	padLen := columnSeparatorLen - len(style.Value()) - 1
	if padLen > 0 {
		pad = columnSeparator[:padLen]
	}

	return &cursorRenderer{
		style: style,
		pad:   pad,
	}
}

func (r *cursorRenderer) Render(name string) string {
	return r.style.Render(name) + r.pad
}
