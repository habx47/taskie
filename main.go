package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func HandleCommand(cmd []string, w io.Writer) bool {

	if len(cmd) == 0 {
		fmt.Fprintln(w, "No command provided!")
		return true
	}

	switch cmd[0] {
	case "exit":
		fmt.Fprintln(w, "Goodbye!")
		fmt.Fprintln(w, "Shutting down...")
		return false
	case "help":
		fmt.Fprintln(w, "help command")
	case "add":
		fmt.Fprintln(w, "add command")
	case "delete":
		fmt.Fprintln(w, "delete command")
	default:
		fmt.Fprintln(w, "Please provide a valid command")
	}
	return true
}

func CaptureCommand(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	active := true

	for active {
		fmt.Print("taskie> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		commands_args := strings.Fields(line)

		active = HandleCommand(commands_args, w)
	}
}

func main() {
	fmt.Println("Welcome to Taskie")
	CaptureCommand(os.Stdin, os.Stdout)
}
