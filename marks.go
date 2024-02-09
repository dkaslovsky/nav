package main

import "errors"

func (m *model) marked() bool {
	return m.markedIndex(m.displayIndex())
}

func (m *model) markedIndex(dispIdx int) bool {
	_, marked := m.marks[dispIdx]
	return marked
}

func (m *model) reloadMarks() error {
	if len(m.marks) == 0 {
		return nil
	}
	newMarks := make(map[int]int)
	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	for _, entryIdx := range m.marks {
		if newDisplayIdx, ok := cache.lookupDisplayIndex(entryIdx); ok {
			newMarks[newDisplayIdx] = entryIdx
		}
	}
	m.marks = newMarks
	return nil
}

func (m *model) toggleMark() error {
	idx := m.displayIndex()
	if m.markedIndex(idx) {
		delete(m.marks, idx)
		m.modeMarks = len(m.marks) != 0
		return nil
	}

	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	entryIdx, ok := cache.lookupEntryIndex(idx)
	if !ok {
		return errors.New("failed to find entry index")
	}
	m.marks[idx] = entryIdx
	m.modeMarks = true
	return nil
}

func (m *model) toggleMarkAll() error {
	// Check if all displayed entries are marked to determine toggle behavior.
	allMarked := true
	for i := 0; i < m.displayed; i++ {
		if _, marked := m.marks[i]; !marked {
			allMarked = false
			break
		}
	}

	if allMarked {
		m.clearMarks()
		return nil
	}
	return m.markAll()
}

func (m *model) markAll() error {
	m.marks = make(map[int]int)
	cache, ok := m.pathCache[m.path]
	if !ok || !cache.hasIndexes() {
		return errors.New("failed to load page cache indexes")
	}
	for i := 0; i < m.displayed; i++ {
		if entryIdx, ok := cache.lookupEntryIndex(i); ok {
			m.marks[i] = entryIdx
		}
	}
	m.modeMarks = len(m.marks) != 0
	return nil
}

func (m *model) clearMarks() {
	m.marks = make(map[int]int)
	m.modeMarks = false
}
