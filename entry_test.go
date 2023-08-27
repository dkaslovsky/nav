package main

import (
	"io/fs"
	"testing"
	"time"
)

// Equality is asserted by a name check so all entries must have unique names
func testEntrySliceEqual(entries1, entries2 []*entry) bool {
	testEntryEqual := func(entry1, entry2 *entry) bool {
		return entry1.Name() == entry2.Name()
	}

	if len(entries1) != len(entries2) {
		return false
	}
	for i, ent := range entries1 {
		if !testEntryEqual(ent, entries2[i]) {
			return false
		}
	}
	return true
}

func TestSortEntriesByType(t *testing.T) {
	tests := map[string]struct {
		entries []*entry
		want    []*entry
	}{
		"no_entries": {
			entries: []*entry{},
			want:    []*entry{},
		},
		"entries": {
			entries: []*entry{
				newEntry(&mockDirEntry{name: ".hidden1", mode: fs.ModeSymlink}),
				newEntry(&mockDirEntry{name: ".hidden2", mode: fs.ModeDir}),
				newEntry(&mockDirEntry{name: "file2", mode: fs.ModeSymlink}),
				newEntry(&mockDirEntry{name: "dir2", mode: fs.ModeDir}),
				newEntry(&mockDirEntry{name: "file1", mode: fs.ModeIrregular}),
				newEntry(&mockDirEntry{name: "dir1", mode: fs.ModeDir}),
			},
			want: []*entry{
				newEntry(&mockDirEntry{name: "dir1", mode: fs.ModeDir}),
				newEntry(&mockDirEntry{name: "dir2", mode: fs.ModeDir}),
				newEntry(&mockDirEntry{name: "file1", mode: fs.ModeIrregular}),
				newEntry(&mockDirEntry{name: "file2", mode: fs.ModeSymlink}),
				newEntry(&mockDirEntry{name: ".hidden2", mode: fs.ModeDir}),
				newEntry(&mockDirEntry{name: ".hidden1", mode: fs.ModeSymlink}),
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(tt *testing.T) {
			// Make a copy because the function operates in-place.
			entries := []*entry{}
			for _, ent := range test.entries {
				ent := ent // Avoid aliasing.
				entries = append(entries, ent)
			}

			sortEntriesByType(entries)
			if !testEntrySliceEqual(entries, test.want) {
				tt.Fatal("incorrect sort order for entries")
			}
		})
	}
}

// mockDirEntry provides a mock implementation of the fs.DirEntry interface for testing.
type mockDirEntry struct {
	name string
	mode fs.FileMode
}

func (de *mockDirEntry) Name() string               { return de.name }
func (de *mockDirEntry) IsDir() bool                { return de.mode&fs.ModeDir == fs.ModeDir }
func (de *mockDirEntry) Info() (fs.FileInfo, error) { return &mockFileInfo{mode: de.mode}, nil }
func (de *mockDirEntry) Type() fs.FileMode          { return fs.FileMode(0) } // Unused.

// mockFileInfo provides a mock implementation of the fs.FileInfo interface for testing.
// The Mode() method is the only relevant implementation for the tests.
type mockFileInfo struct {
	mode fs.FileMode
}

func (fi *mockFileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi *mockFileInfo) Name() string       { return "" }          // Unused.
func (fi *mockFileInfo) Size() int64        { return 0 }           // Unused.
func (fi *mockFileInfo) ModTime() time.Time { return time.Time{} } // Unused.
func (fi *mockFileInfo) IsDir() bool        { return false }       // Unused.
func (fi *mockFileInfo) Sys() any           { return nil }         // Unused.
