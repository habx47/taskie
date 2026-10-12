package main

import (
	"fmt"
	"os"
)

func CreateFileStore(name string) (*os.File, error) {
	fileName := name

	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		file, createErr := os.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("Failed to create a task_data file: %w", createErr)
		}
		return WriteFileStore(file)

	} else if err != nil {
		return nil, fmt.Errorf("Error checking file stats: %w", err)
	}
	return nil, nil
}

func WriteFileStore(file *os.File) (*os.File, error) {
	_, writeErr := file.Write([]byte("{}"))
	file.Close()
	if writeErr != nil {
		return nil, fmt.Errorf("Failed to write inital json: %w", writeErr)
	}

	return file, nil
}
