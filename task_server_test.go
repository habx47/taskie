package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTaskServer(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "EOF",
			input: "",
			want:  "",
		},
		{
			name:  "invalid",
			input: "invalid command\n",
			want:  "Please provide a valid command\n",
		},
		{
			name:  "exit",
			input: "exit\n",
			want:  "Goodbye!\nShutting down...\n",
		},
		{
			name:  "add",
			input: "add first task\n",
			want:  "Task added: \"first task\"\n",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			stubTaskStore := NewInMemoryTaskStore(map[int]string{})
			r := strings.NewReader(tt.input)
			var buffer bytes.Buffer
			testTaskServer := NewTaskServer(&buffer, stubTaskStore)

			testTaskServer.CaptureCommand(r)

			got := (&buffer).String()
			want := tt.want

			assertString(t, got, want)
		})
	}

	t.Run("failure to read", func(t *testing.T) {
		stubTaskStore := NewInMemoryTaskStore(map[int]string{})
		r := iotest.ErrReader(errors.New("this will always error"))
		testTaskServer := NewTaskServer(io.Discard, stubTaskStore)

		err := testTaskServer.CaptureCommand(r)

		if err == nil {
			t.Fatal("expected an err but got nil")
		}
	})
}

func TestListTasks(t *testing.T) {
	var buffer bytes.Buffer
	stubTaskStore := NewInMemoryTaskStore(map[int]string{
		1: "first task",
		2: "second task",
	})
	taskServer := NewTaskServer(&buffer, stubTaskStore)
	taskServer.ListTasks()

	got := (&buffer).String()
	want := "1. first task\n2. second task\n"

	assertString(t, got, want)
}

func TestAddTask(t *testing.T) {
	var buffer bytes.Buffer
	stubTaskStore := NewInMemoryTaskStore(map[int]string{})
	taskServer := NewTaskServer(&buffer, stubTaskStore)
	taskServer.AddTask("testing add tasks")

	got := (&buffer).String()
	want := "Task added: \"testing add tasks\"\n"

	assertString(t, got, want)
}

func TestDeleteTask(t *testing.T) {

	t.Run("delete with valid id", func(t *testing.T) {
		var buffer bytes.Buffer
		stubTaskStore := NewInMemoryTaskStore(map[int]string{
			1: "first task",
		})
		taskServer := NewTaskServer(&buffer, stubTaskStore)
		taskServer.DeleteTask("1")

		got := buffer.String()
		want := "Task deleted\n"

		assertString(t, got, want)
	})

	t.Run("delete with invalid id", func(t *testing.T) {
		var buffer bytes.Buffer
		stubTaskStore := NewInMemoryTaskStore(map[int]string{
			1: "first task",
		})
		taskServer := NewTaskServer(&buffer, stubTaskStore)
		taskServer.DeleteTask("invalid_id")

		got := buffer.String()
		want := "Operation failed. Invalid id: invalid_id\n"

		assertString(t, got, want)
	})

	t.Run("delete with non-existant id", func(t *testing.T) {
		var buffer bytes.Buffer
		stubTaskStore := NewInMemoryTaskStore(map[int]string{
			1: "first task",
		})
		taskServer := NewTaskServer(&buffer, stubTaskStore)
		taskServer.DeleteTask("2")

		got := buffer.String()
		want := "Error deleting, id: 2 not found\n"

		assertString(t, got, want)
	})
}
