package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var LoadingFile_ReadErr = errors.New("Couldn't read the json data file")
var LoadingFile_JSONErr = errors.New("Couldn't unmarhsall the json file")

type JSONPayload struct {
	LatestID int            `json:"latest_id"`
	Tasks    map[int]string `json:"tasks"`
}

func CreateFileStore(file_store *FileStore) (bool, error) {
	_, err := os.Stat(file_store.filename)

	if os.IsNotExist(err) {
		file, createErr := os.Create(file_store.filename)
		if createErr != nil {
			return false, fmt.Errorf("Failed to create a task_data file: %w", createErr)
		}
		file.Close()
		return true, SaveToFile(file_store)
	} else if err != nil {
		return false, fmt.Errorf("Error checking file stats: %w", err)
	}
	return false, nil
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

func SaveToFile(file_store *FileStore) error {
	file_store.mu.RLock()
	payload := JSONPayload{
		LatestID: file_store.latest_id,
		Tasks:    file_store.store,
	}
	file_store.mu.RUnlock()

	bytes, err := json.MarshalIndent(payload, "", " ")

	if err != nil {
		return fmt.Errorf("Error converting to json for saving: %w", err)
	}

	if err := os.WriteFile(file_store.filename, bytes, 0666); err != nil {
		return fmt.Errorf("Failed to write json data to task_data file :%w", err)
	}

	return nil
}
