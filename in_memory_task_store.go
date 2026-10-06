package main

type InMemoryTaskStore struct {
	store map[int]string
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{store: map[int]string{}}
}

func (i *InMemoryTaskStore) Add(id int, data string) {
}

func (i *InMemoryTaskStore) Delete(id int) {
}

func (i *InMemoryTaskStore) List() {
}
