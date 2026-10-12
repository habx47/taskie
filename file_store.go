package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

type FileStore struct {
	filename string
	store    map[int]string
	mu       sync.RWMutex
}

func NewFileStore() (*FileStore, error) {
	file_store := &FileStore{filename: "task_data.json", store: map[int]string{}}

	file, err := CreateFileStore(file_store.filename)

	if err != nil {
		return nil, fmt.Errorf("Error related to task_data file: %w", err)
	}

	if file != nil {
		return file_store, nil
	}

	bytes, err := os.ReadFile(file_store.filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to read task_data file: %w", err)
	}

	if jsonErr := json.Unmarshal(bytes, &file_store.store); jsonErr != nil {
		return nil, fmt.Errorf("Failed to parse existing JSON: %w", jsonErr)
	}

	return file_store, nil
}

func (f *FileStore) saveToFile() error {
	f.mu.RLock()
	bytes, err := json.MarshalIndent(f.store, "", " ")
	f.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("Error converting to json: %w", err)
	}

	if err := os.WriteFile(f.filename, bytes, 0666); err != nil {
		return fmt.Errorf("Failed to write data to task_data file :%w", err)
	}

	return nil
}

func (f *FileStore) add(data string) error {
	f.mu.Lock()
	f.store[len(f.store)+1] = data
	f.mu.Unlock()

	return f.saveToFile()
}

func (f *FileStore) delete(id int) error {
	f.mu.Lock()
	if _, ok := f.store[id]; !ok {
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	delete(f.store, id)
	f.mu.Unlock()

	return f.saveToFile()
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
