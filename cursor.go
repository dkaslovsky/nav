package main

import "github.com/charmbracelet/lipgloss"

const cursorStr = ">"

var cursor = lipgloss.NewStyle().Bold(true).SetString(cursorStr)

type position struct {
	c int
	r int
}

func (m *model) moveUp() {
	m.r--
	if m.r < 0 {
		m.r = m.rows - 1
		m.c--
	}
	if m.c < 0 {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.entries))
		m.c = m.columns - 1
	}
}

func (m *model) moveDown() {
	m.r++
	if m.r >= m.rows {
		m.r = 0
		m.c++
	}
	if m.c >= m.columns {
		m.c = 0
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.entries) {
		m.r = 0
		m.c = 0
	}
}

func (m *model) moveLeft() {
	m.c--
	if m.c < 0 {
		m.c = m.columns - 1
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.entries) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.entries))
		m.c = m.columns - 1
	}
}

func (m *model) moveRight() {
	m.c++
	if m.c >= m.columns {
		m.c = 0
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.entries) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.entries))
		m.c = m.columns - 1
	}
}

func (m *model) resetCursor() {
	m.c = 0
	m.r = 0
}

func (m *model) setCursor(pos *position) {
	m.c = pos.c
	m.r = pos.r
}

func (m *model) saveCursor() {
	m.cursorCache[m.path] = &position{c: m.c, r: m.r}
}

func (m *model) current() (*entry, bool) {
	idx := m.r + (m.c * m.rows)
	if idx > len(m.entries) {
		return nil, false
	}
	return m.entries[idx], true
}
