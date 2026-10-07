package main

import (
	"bytes"
	"strings"
	"testing"
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

	taskServer.ListTasks()

	got := (&buffer).String()
	want := "1. testing add tasks\n"

	assertString(t, got, want)
}

/*func TestDeleteTask(t *testing.T) {
	var buffer bytes.Buffer
	stubTaskStore := NewInMemoryTaskStore(map[int]string{
		1: "first task",
	})
} */
