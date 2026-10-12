package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file_store, err := NewFileStore()
	if err != nil {
		log.Fatalf("Could not initialize the data: %v", err)
	}

	taskServer := NewTaskServer(os.Stdout, file_store)
	fmt.Println("Welcome to Taskie")
	taskServer.CaptureCommand(os.Stdin)
}
