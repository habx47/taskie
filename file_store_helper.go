package main

import (
	"fmt"
	"os"
)

func CreateFileStore(file_store *FileStore) (bool, error) {

	_, err := os.Stat(file_store.filename)

	if os.IsNotExist(err) {
		file, createErr := os.Create(file_store.filename)
		if createErr != nil {
			return false, fmt.Errorf("Failed to create a task_data file: %w", createErr)
		}
		file.Close()
		return true, file_store.saveToFile()
	} else if err != nil {
		return false, fmt.Errorf("Error checking file stats: %w", err)
	}
	return false, nil
}
