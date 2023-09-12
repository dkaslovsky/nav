package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) selectAction() (*model, tea.Cmd) {
	selected, err := m.selected()
	if err != nil {
		m.setError(err, "failed to select entry")
		return m, nil
	}

	m.saveCursor()

	if selected.hasMode(entryModeFile) {
		m.setExit(sanitizeOutputPath(filepath.Join(m.path, selected.Name())))
		if m.modeSubshell {
			fmt.Print(m.exitStr)
		}
		return m, tea.Quit
	}
	if selected.hasMode(entryModeSymlink) {
		sl, err := followSymlink(m.path, selected)
		if err != nil {
			m.setError(err, "failed to evaluate symlink")
			return m, nil
		}
		if !sl.info.IsDir() {
			// The symlink points to a file.
			m.setExit(sanitizeOutputPath(sl.absPath))
			if m.modeSubshell {
				fmt.Print(m.exitStr)
			}
			return m, tea.Quit
		}
		m.path = sl.absPath
	} else if selected.hasMode(entryModeDir) {
		path, err := filepath.Abs(filepath.Join(m.path, selected.Name()))
		if err != nil {
			m.setError(err, "failed to evaluate path")
			return m, nil
		}
		m.path = path
	} else {
		m.setError(
			errors.New("selection is not a file, directory, or symlink"),
			"unexpected file type",
		)
		return m, nil
	}

	err = m.list()
	if err != nil {
		m.setError(err, err.Error())
		return m, nil
	}

	m.clearSearch()
	m.clearError()

	// Return to ensure the cursor is not re-saved using the updated path.
	return m, nil
}

func (m *model) searchSelectAction() (*model, tea.Cmd) {
	selected, err := m.selected()
	if err != nil {
		m.setError(err, "failed to select entry")
		m.clearSearch()
		return m, nil
	}

	if selected.hasMode(entryModeFile) {
		m.setExit(sanitizeOutputPath(filepath.Join(m.path, selected.Name())))
		if m.modeSubshell {
			fmt.Print(m.exitStr)
		}
		return m, tea.Quit
	}
	if selected.hasMode(entryModeSymlink) {
		sl, err := followSymlink(m.path, selected)
		if err != nil {
			m.setError(err, "failed to evaluate symlink")
			return m, nil
		}
		if !sl.info.IsDir() {
			// The symlink points to a file.
			m.setExit(sanitizeOutputPath(sl.absPath))
			if m.modeSubshell {
				fmt.Print(m.exitStr)
			}
			return m, tea.Quit
		}
		m.path = sl.absPath
	} else if selected.hasMode(entryModeDir) {
		m.path = m.path + "/" + selected.Name()
	} else {
		m.setError(
			errors.New("selection is not a file, directory, or symlink"),
			"unexpected file type",
		)
		return m, nil
	}

	// TODO: encapsulate this fix
	if strings.HasPrefix(m.path, "//") {
		m.path = m.path[1:]
	}

	m.search = ""
	err = m.list()
	if err != nil {
		m.setError(err, err.Error())
		m.clearSearch()
		return m, nil
	}
	return m, nil
}
