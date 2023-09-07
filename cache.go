package main

type cacheItem struct {
	cursorPosition       *position
	entityToDisplayIndex map[int]int
	displayToEntityIndex map[int]int
	columns              int
	rows                 int
}

func newCacheItem() *cacheItem {
	return newCacheItemWithPosition(nil)
}

func newCacheItemWithPosition(pos *position) *cacheItem {
	return &cacheItem{
		cursorPosition:       pos,
		entityToDisplayIndex: make(map[int]int),
		displayToEntityIndex: make(map[int]int),
	}
}

func (ci *cacheItem) hasIndexes() bool {
	return len(ci.entityToDisplayIndex) > 0 && len(ci.displayToEntityIndex) > 0
}
