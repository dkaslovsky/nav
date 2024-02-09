package main

type cacheItem struct {
	cursorPosition *position
	entryToDisplay map[int]int
	displayToEntry map[int]int
	columns        int
	rows           int
}

func newCacheItem() *cacheItem {
	return newCacheItemWithPosition(nil)
}

func newCacheItemWithPosition(pos *position) *cacheItem {
	return &cacheItem{
		cursorPosition: pos,
		entryToDisplay: make(map[int]int),
		displayToEntry: make(map[int]int),
	}
}

type indexPair struct {
	entry   int
	display int
}

func (ci *cacheItem) addIndexPair(pair *indexPair) {
	ci.entryToDisplay[pair.entry] = pair.display
	ci.displayToEntry[pair.display] = pair.entry
}

func (ci *cacheItem) setPosition(pos *position) {
	ci.cursorPosition = pos
}

func (ci *cacheItem) setColumns(c int) {
	ci.columns = c
}

func (ci *cacheItem) setRows(r int) {
	ci.rows = r
}

func (ci *cacheItem) hasIndexes() bool {
	return len(ci.entryToDisplay) > 0 && len(ci.displayToEntry) > 0
}

func (ci *cacheItem) lookupEntryIndex(displayIdx int) (int, bool) {
	entryIdx, ok := ci.displayToEntry[displayIdx]
	return entryIdx, ok
}

func (ci *cacheItem) lookupDisplayIndex(entryIdx int) (int, bool) {
	displayIdx, ok := ci.entryToDisplay[entryIdx]
	return displayIdx, ok
}

// cursorIndex returns the display index of cursor
func (ci *cacheItem) cursorIndex() int {
	return ci.cursorPosition.index(ci.rows)
}
