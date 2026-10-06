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

func TestCaptureCommand(t *testing.T) {

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
			r := strings.NewReader(tt.input)
			buffer := bytes.Buffer{}

			CaptureCommand(r, &buffer)

			got := (&buffer).String()
			want := tt.want

			assertString(t, got, want)
		})
	}
}
