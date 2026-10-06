package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TaskStore interface {
	Add(id int, data string)
	Delete(id int)
	List()
}

type TaskServer struct {
	id    int
	w     io.Writer
	store TaskStore
}

func NewTaskServer(w io.Writer, store TaskStore) *TaskServer {
	t := new(TaskServer)
	t.w = w
	t.store = store
	return t
}

func (t *TaskServer) CaptureCommand(r io.Reader) {
	scanner := bufio.NewScanner(r)
	active := true

	for active {
		fmt.Print("taskie> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		commands_args := strings.Fields(line)

		active = t.HandleCommand(commands_args)
	}
}

func (t *TaskServer) HandleCommand(cmd []string) bool {
	if len(cmd) == 0 {
		fmt.Fprintln(t.w, "No command provided!")
		return true
	}

	switch cmd[0] {
	case "exit":
		fmt.Fprintln(t.w, "Goodbye!")
		fmt.Fprintln(t.w, "Shutting down...")
		return false
	case "help":
		fmt.Fprintln(t.w, `Available commands:
- add
- delete
- exit
- help`)
	case "add":
		fmt.Fprintln(t.w, "add command")
	case "delete":
		fmt.Fprintln(t.w, "delete command")
	default:
		fmt.Fprintln(t.w, "Please provide a valid command")
	}
	return true
}

func (t *TaskServer) AddTask(data string) {
}

func (t *TaskServer) DeleteTask(id int) {
}

func (t *TaskServer) ListTasks() {
}
