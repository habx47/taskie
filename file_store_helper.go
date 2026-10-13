package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var FatalErr = errors.New("Fatal error")

type JSONPayload struct {
	LatestID int            `json:"latest_id"`
	Tasks    map[int]string `json:"tasks"`
}

func getDataFileName(name string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	projectDir := filepath.Join(homeDir, "projects", "taskie")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(projectDir, name), nil
}

func CreateFileStore(file_store *FileStore) (bool, error) {
	_, statErr := os.Stat(file_store.filepath)

	if os.IsNotExist(statErr) {
		return true, SaveToFile(file_store)
	} else if statErr != nil {
		return false, fmt.Errorf("%w: Error checking file stats: %v", FatalErr, statErr)
	}
	return false, nil
}

func LoadFileStore(file_store *FileStore) error {
	bytes, err := os.ReadFile(file_store.filepath)
	if err != nil {
		return fmt.Errorf("%w: couldn't read the json data file: %v", FatalErr, err)
	}

	if !json.Valid(bytes) {
		return fmt.Errorf("%w: Data file corrupted. Manually try to correct it.", FatalErr)
	}

	var payload JSONPayload
	if jsonErr := json.Unmarshal(bytes, &payload); jsonErr != nil {
		return fmt.Errorf("%w: couldn't unmarshall the json data file: %v", FatalErr, jsonErr)
	}

	file_store.latest_id = payload.LatestID
	file_store.store = payload.Tasks
	return nil
}

func SaveToFile(file_store *FileStore) error {
	file_store.mu.RLock()
	payload := JSONPayload{
		LatestID: file_store.latest_id,
		Tasks:    file_store.store,
	}

	bytes, err := json.MarshalIndent(payload, "", " ")
	file_store.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("Error converting to json for saving: %w", err)
	}

	dir := filepath.Dir(file_store.filepath)
	tempFile, err := os.CreateTemp(dir, "tasks_temp_*.json")
	if err != nil {
		return fmt.Errorf("Error creating temp file: %w", err)
	}

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(bytes); err != nil {
		tempFile.Close()
		return fmt.Errorf("Failed to write json data to temp file :%w", err)
	}

	tempFile.Close()

	renameErr := os.Rename(tempFile.Name(), file_store.filepath)
	if renameErr != nil {
		return fmt.Errorf("Error creating json data file from temp file: %w", renameErr)
	}
	return nil
}
