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
	store map[int]string
	mu    sync.RWMutex
}

func NewFileStore() (*FileStore, error) {
	filePath := "task_data.json"
	file_store := &FileStore{store: map[int]string{}}

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		file, createErr := os.Create(filePath)
		if createErr != nil {
			return nil, fmt.Errorf("Failed to create a task_data file: %w", createErr)
		}

		_, writeErr := file.Write([]byte("{}"))
		file.Close()
		if writeErr != nil {
			return nil, fmt.Errorf("Failed to write inital json: %w", writeErr)
		}

		return file_store, nil

	} else if err != nil {
		return nil, fmt.Errorf("Error checking file stats: %w", err)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read task_data file: %w", err)
	}

	if jsonErr := json.Unmarshal(bytes, &file_store.store); jsonErr != nil {
		return nil, fmt.Errorf("Failed to parse existing JSON: %w", jsonErr)
	}

	return file_store, nil
}

func (f *FileStore) add(data string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.store[len(f.store)+1] = data
}

func (f *FileStore) delete(id int) error {
	if _, ok := f.store[id]; !ok {
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

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
