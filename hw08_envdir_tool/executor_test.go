package main

import (
	"bytes"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		env          Environment
		expectedCode int
		expectedOut  string
	}{
		{
			name:         "Test command runs and produces expected output",
			cmd:          []string{"bash", "-c", "echo $MYVAR"},
			env:          Environment{"MYVAR": EnvValue{Value: "Hello, world!", NeedRemove: false}},
			expectedCode: 0,
			expectedOut:  "Hello, world!\n",
		},
		{
			name:         "Must return error",
			cmd:          []string{"invalid_command"},
			env:          nil,
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			rc := RunCmd(tt.cmd, tt.env, &buf)

			out := buf.String() // get output from stdout

			if rc != tt.expectedCode {
				t.Errorf("Expected return code %v, got %v", tt.expectedCode, rc)
			}

			if out != tt.expectedOut && tt.expectedOut != "" {
				t.Errorf("Expected output %q, got %q", tt.expectedOut, out)
			}
		})
	}
}
