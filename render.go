package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	rendererCursor = newRenderer(lipgloss.NewStyle().Bold(true).SetString(">"))
	rendererNormal = newRenderer(lipgloss.NewStyle().SetString(" "))
)

type renderer struct {
	style lipgloss.Style
	pad   string
}

func newRenderer(style lipgloss.Style) *renderer {
	pad := ""
	padLen := columnSeparatorLen - len(style.Value()) - 1
	if padLen > 0 {
		pad = columnSeparator[:padLen]
	}

	return &renderer{
		style: style,
		pad:   pad,
	}
}

func (r *renderer) render(name string) string {
	return r.style.Render(name) + r.pad
}
