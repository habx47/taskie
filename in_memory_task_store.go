package main

import (
	"fmt"
	"io"
)

type InMemoryTaskStore struct {
	store map[int]string
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{store: map[int]string{}}
}

func (i *InMemoryTaskStore) Add(id int, data string) {
	i.store[id] = data
}

func (i *InMemoryTaskStore) Delete(id int) error {
	if _, ok := i.store[id]; !ok {
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	delete(i.store, id)
	return nil
}

func (i *InMemoryTaskStore) List(w io.Writer) {
	for index, value := range i.store {
		fmt.Fprintf(w, "%v. %v\n", index, value)
	}
}
