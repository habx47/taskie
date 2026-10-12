package main

import (
	"fmt"
	"slices"
	"strings"
)

type FileStore struct {
	id    int
	store map[int]string
}

func NewFileStore() *FileStore {

	return &FileStore{id: 1, store: map[int]string{}}
}

func (f *FileStore) add(data string) {
	f.store[f.id] = data
	f.id++
}

func (f *FileStore) delete(id int) error {
	if _, ok := f.store[id]; !ok {
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	delete(f.store, id)
	return nil
}

func (f *FileStore) list() string {
	if len(f.store) == 0 {
		return "No tasks found. Use \"add\" command to add tasks.\n"
	}
	ids_slice := make([]int, len(f.store))

	idx := 0
	for id := range f.store {
		ids_slice[idx] = id
		idx++
	}

	slices.Sort(ids_slice)
	var sb strings.Builder

	for _, id := range ids_slice {
		fmt.Fprintf(&sb, "%v. %v\n", id, f.store[id])
	}

	return sb.String()
}
