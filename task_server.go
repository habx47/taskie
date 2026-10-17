package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type TaskStore interface {
	add(data string) error
	delete(id int) error
	list() string
}

type TaskServer struct {
	w     io.Writer
	store TaskStore
}

func NewTaskServer(w io.Writer, store TaskStore) *TaskServer {
	t := new(TaskServer)
	t.w = w
	t.store = store
	return t
}

func (t *TaskServer) CaptureCommand(r io.Reader) error {
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

	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("error reading the input: %v", err)
	}

	return nil
}

func (t *TaskServer) CommandRouter(cmd []string) bool {
	if len(cmd) == 0 {
		fmt.Fprintln(t.w, "No command provided!")
		return true
	}

	switch strings.ToLower(cmd[0]) {

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
		if len(cmd) == 1 {
			fmt.Fprint(t.w, "Need a valid task to add. Currently empty\n")
			return true
		}
		task := strings.Join(cmd[1:], " ")
		t.AddTask(task)

	case "list":
		t.ListTasks()

	case "delete":
		if len(cmd) == 1 {
			fmt.Fprint(t.w, "Need a valid task-id to delete\n")
			return true
		}

		t.DeleteTask(cmd[1])

	default:
		fmt.Fprintln(t.w, "Please provide a valid command")
	}
	return true
}

func (t *TaskServer) AddTask(data string) {
	err := t.store.add(data)
	if err == nil {
		fmt.Fprintf(t.w, "Task added: \"%v\"\n", data)
		return
	}

	if errors.Is(err, FatalErr) {
		log.Fatalf("%v\n", err)
	}

	fmt.Fprintf(t.w, "%v\n", err)
}

func (t *TaskServer) DeleteTask(id string) {
	int_id, conv_err := strconv.Atoi(id)
	if conv_err != nil {
		fmt.Fprintf(t.w, "Operation failed. Invalid id: %v\n", id)
		return
	}

	deletion_err := t.store.delete(int_id)

	if deletion_err == nil {
		fmt.Fprint(t.w, "Task deleted\n")
		return
	}

	if errors.Is(deletion_err, FatalErr) {
		log.Fatalf("%v\n", deletion_err)
	}

	fmt.Fprintf(t.w, "%v\n", deletion_err)
}

func (t *TaskServer) ListTasks() {
	taskListStr := t.store.list()
	fmt.Fprint(t.w, taskListStr)
}
