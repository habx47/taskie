package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	file_store, err := NewFileStore()
	if errors.Is(err, FatalErr) {
		log.Fatalf("Could not initialize: %v", err)
	}

	taskServer := NewTaskServer(os.Stdout, file_store)
	fmt.Println("Welcome to Taskie")
	taskServer.CaptureCommand(os.Stdin)
}
