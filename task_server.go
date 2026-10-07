package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TaskStore interface {
	add(id int, data string)
	delete(id int) error
	list() string
}

type TaskServer struct {
	id    int
	w     io.Writer
	store TaskStore
}

func NewTaskServer(w io.Writer, store TaskStore) *TaskServer {
	t := new(TaskServer)
	t.id = 1
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

		active = t.CommandRouter(commands_args)
	}
}

func (t *TaskServer) CommandRouter(cmd []string) bool {
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
- add -> add [task]
- delete -> delete [task_id]
- list 
- exit 
- help`)

	case "add":
		task := strings.Join(cmd[1:], " ")
		t.AddTask(task)
		fmt.Fprintf(t.w, "Task added: \"%v\"\n", task)

	case "list":
		t.ListTasks()

	case "delete":
		t.DeleteTask(cmd[1])

	default:
		fmt.Fprintln(t.w, "Please provide a valid command")
	}
	return true
}

func (t *TaskServer) AddTask(data string) {
	t.store.add(t.id, data)
	t.id++
}

func (t *TaskServer) DeleteTask(id string) {
	int_id, conv_err := strconv.Atoi(id)
	if conv_err != nil {
		fmt.Fprintf(t.w, "Operation failed. Invalid id: %v\n", id)
		return
	}

	deletion_err := t.store.delete(int_id)
	if deletion_err != nil {
		fmt.Fprintf(t.w, "%v\n", deletion_err)
		return
	}

	fmt.Fprint(t.w, "Task deleted\n")
}

func (t *TaskServer) ListTasks() {
	taskListStr := t.store.list()
	fmt.Fprint(t.w, taskListStr)
}
