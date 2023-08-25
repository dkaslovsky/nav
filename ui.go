package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var keyQuit = key.NewBinding(key.WithKeys("q"))

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyQuit):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible on exit.
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() string {
	output := []string{}
	for _, file := range m.files {
		output = append(output, file.Name())
	}
	return strings.Join(output, "\n")
}
