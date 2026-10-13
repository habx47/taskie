package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

var LoadingFile_ReadErr = errors.New("Couldn't read the json data file")
var LoadingFile_JSONErr = errors.New("Couldn't unmarhsall the json file")

type JSONPayload struct {
	LatestID int            `json:"latest_id"`
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
		return nil, fmt.Errorf("Error related to creating task_data file: %w", err)
	}

	if file_created == true {
		return file_store, nil
	}

	loadErr := LoadFileStore(file_store)
	if loadErr != nil {
		return nil, fmt.Errorf("Error realted to loading task_data file: %w", loadErr)
	}

	return file_store, nil
}

func LoadFileStore(file_store *FileStore) error {
	bytes, err := os.ReadFile(file_store.filename)
	if err != nil {
		return fmt.Errorf("Failure: %w, err: %v", LoadingFile_ReadErr, err)
	}

	var payload JSONPayload
	if jsonErr := json.Unmarshal(bytes, &payload); jsonErr != nil {
		return fmt.Errorf("Failure: %w, err: %v", LoadingFile_JSONErr, jsonErr)
	}

	file_store.latest_id = payload.LatestID
	file_store.store = payload.Tasks
	return nil
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
		return fmt.Errorf("Error converting to json for saving: %w", err)
	}

	if err := os.WriteFile(f.filename, bytes, 0666); err != nil {
		return fmt.Errorf("Failed to write json data to task_data file :%w", err)
	}

	return nil
}

func (f *FileStore) add(data string) error {
	f.mu.Lock()
	f.store[f.latest_id] = data
	f.latest_id++
	f.mu.Unlock()

	err := f.saveToFile()
	if err != nil {
		loadErr := LoadFileStore(f)
		if loadErr != nil {
			return fmt.Errorf("Loading file again on failure to add new task failed: %w", loadErr)
		}
	}

	return nil
}

func (f *FileStore) delete(id int) error {
	f.mu.Lock()
	if _, ok := f.store[id]; !ok {
		f.mu.Unlock()
		return fmt.Errorf("Error deleting, id: %v not found", id)
	}

	delete(f.store, id)
	f.mu.Unlock()

	err := f.saveToFile()
	if err != nil {
		loadErr := LoadFileStore(f)
		if loadErr != nil {
			return fmt.Errorf("Loading file again on failure to add new task failed: %w", loadErr)
		}
	}

	return nil
}

func (f *FileStore) list() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

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
