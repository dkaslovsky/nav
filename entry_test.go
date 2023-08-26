package main

import (
	"io/fs"
	"testing"
)

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
				{&testDirEntry{name: ".hidden1", isDir: false}},
				{&testDirEntry{name: ".hidden2", isDir: true}},
				{&testDirEntry{name: "file2", isDir: false}},
				{&testDirEntry{name: "dir2", isDir: true}},
				{&testDirEntry{name: "file1", isDir: false}},
				{&testDirEntry{name: "dir1", isDir: true}},
			},
			want: []*entry{
				{&testDirEntry{name: "dir1", isDir: true}},
				{&testDirEntry{name: "dir2", isDir: true}},
				{&testDirEntry{name: "file1", isDir: false}},
				{&testDirEntry{name: "file2", isDir: false}},
				{&testDirEntry{name: ".hidden2", isDir: true}},
				{&testDirEntry{name: ".hidden1", isDir: false}},
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

// testDirEntry mocks the fs.DirEntry interface for testing.
type testDirEntry struct {
	name  string
	isDir bool
}

func (tde *testDirEntry) Name() string {
	return tde.name
}

func (tde *testDirEntry) IsDir() bool {
	return tde.isDir
}

func (tde *testDirEntry) Type() fs.FileMode {
	return fs.FileMode(0) // Unused.
}

func (tde *testDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil // Unused.
}

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
