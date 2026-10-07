package main

import (
	"fmt"
	"os"
)

func main() {
	taskServer := NewTaskServer(os.Stdout, NewInMemoryTaskStore(map[int]string{}))
	fmt.Println("Welcome to Taskie")
	taskServer.CaptureCommand(os.Stdin)
}
