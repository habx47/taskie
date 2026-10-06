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
			testTaskStore := NewInMemoryTaskStore()
			r := strings.NewReader(tt.input)
			var buffer bytes.Buffer
			testTaskServer := NewTaskServer(&buffer, testTaskStore)

			testTaskServer.CaptureCommand(r)

			got := (&buffer).String()
			want := tt.want

			assertString(t, got, want)
		})
	}
}

func TestListTasks(t *testing.T) {
	var buffer bytes.Buffer
	stubTaskStore := InMemoryTaskStore{store: map[int]string{
		1: "first task",
		2: "second task",
	}}
	taskServer := NewTaskServer(&buffer, &stubTaskStore)
	taskServer.ListTasks()

	got := (&buffer).String()
	want := "1. first task\n2. second task\n"

	assertString(t, got, want)
}
