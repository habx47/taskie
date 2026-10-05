package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Taskie")

	for {

		fmt.Print("taskie> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		command_args := strings.Fields(input)

		if len(command_args) == 0 {
			fmt.Println("No command provided!")
			continue
		}

		command := command_args[0]
		switch command {
		case "exit":
			fmt.Println("Shutting down...")
			fmt.Println("Goodbye!")
			return
		case "help":
			fmt.Println("help command")
		case "add":
			fmt.Println("add command")
		case "delete":
			fmt.Println("delete command")
		default:
			fmt.Println("Please provide a valid command")
		}
	}
}
