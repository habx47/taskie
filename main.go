package main

import (
	"fmt"
	"os"
)

func main() {
	taskServer := NewTaskServer(os.Stdout, NewFileStore())
	fmt.Println("Welcome to Taskie")
	taskServer.CaptureCommand(os.Stdin)
}
