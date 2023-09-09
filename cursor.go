package main

// position intentionally does not have a constructor to avoid potential inversion of column and row.
// It should always be instantiated explicitly: pos := &position{c: 0, r: 0}
type position struct {
	c int
	r int
}

func newPositionFromIndex(idx int, rows int) *position {
	return &position{
		c: int(float64(idx) / float64(rows)),
		r: idx % rows,
	}
}

func (p *position) index(rows int) int {
	return index(p.c, p.r, rows)
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
	pos := &position{c: m.c, r: m.r}
	if cache, ok := m.viewCache[m.path]; ok {
		cache.setPosition(pos)
		return
	}
	m.viewCache[m.path] = newCacheItemWithPosition(pos)
}

func (m *model) moveUp() {
	m.r--
	if m.r < 0 {
		m.r = m.rows - 1
		m.c--
	}
	if m.c < 0 {
		m.r = m.rows - 1 - (m.columns*m.rows - m.displayed)
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
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= m.displayed {
		m.r = 0
		m.c = 0
	}
}

func (m *model) moveLeft() {
	m.c--
	if m.c < 0 {
		m.c = m.columns - 1
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= m.displayed {
		m.r = m.rows - 1 - (m.columns*m.rows - m.displayed)
		m.c = m.columns - 1
	}
}

func (m *model) moveRight() {
	m.c++
	if m.c >= m.columns {
		m.c = 0
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= m.displayed {
		m.r = m.rows - 1 - (m.columns*m.rows - m.displayed)
		m.c = m.columns - 1
	}
}
