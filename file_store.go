package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

type JSONPayload struct {
	LatestID int            `json:"latest_id:"`
	Tasks    map[int]string `json:"tasks"`
}

type FileStore struct {
	filename  string
	latest_id int
	store     map[int]string
	mu        sync.RWMutex
}

func NewFileStore() (*FileStore, error) {
	file_store := &FileStore{filename: "task_data.json", latest_id: 1, store: map[int]string{}}

	file_created, err := CreateFileStore(file_store)

	if err != nil {
		return nil, fmt.Errorf("Error related to task_data file: %w", err)
	}

	if file_created == true {
		return file_store, nil
	}

	bytes, err := os.ReadFile(file_store.filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to read task_data file: %w", err)
	}

	var payload JSONPayload
	if jsonErr := json.Unmarshal(bytes, &payload); jsonErr != nil {
		return nil, fmt.Errorf("Failed to parse existing JSON: %w", jsonErr)
	}

	file_store.latest_id = payload.LatestID
	file_store.store = payload.Tasks

	return file_store, nil
}

func (f *FileStore) saveToFile() error {
	f.mu.RLock()
	payload := JSONPayload{
		LatestID: f.latest_id,
		Tasks:    f.store,
	}
	f.mu.RUnlock()

	bytes, err := json.MarshalIndent(payload, "", " ")

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
	f.store[f.latest_id] = data
	f.latest_id++
	f.mu.Unlock()

	return f.saveToFile()
}

func (f *FileStore) delete(id int) error {
	f.mu.Lock()
	if _, ok := f.store[id]; !ok {
		f.mu.Unlock()
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
