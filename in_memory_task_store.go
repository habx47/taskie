package main

import (
	"fmt"
	"slices"
	"strings"
)

type InMemoryTaskStore struct {
	store map[int]string
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{store: map[int]string{}}
}

func (i *InMemoryTaskStore) add(id int, data string) {
	i.store[id] = data
}

func (i *InMemoryTaskStore) delete(id int) error {
	if _, ok := i.store[id]; !ok {
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	delete(i.store, id)
	return nil
}

func (i *InMemoryTaskStore) list() string {
	if len(i.store) == 0 {
		return "No tasks found. Use \"add\" command to add tasks.\n"
	}
	ids_slice := make([]int, len(i.store))

	idx := 0
	for id := range i.store {
		ids_slice[idx] = id
		idx++
	}

	slices.Sort(ids_slice)
	var sb strings.Builder

	for _, id := range ids_slice {
		fmt.Fprintf(&sb, "%v. %v\n", id, i.store[id])
	}

	return sb.String()
}
