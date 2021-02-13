package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecRun(t *testing.T) {
	cases := []struct {
		name    string
		command []string
		status  int
		stdout  string
		stderr  string
	}{
		{
			name:    "with `echo -n 42`",
			command: []string{"echo", "-n", "42"},
			status:  0,
			stdout:  "42",
			stderr:  "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			e := Exec{
				Command: tt.command,
			}

			expected := &ExecResult{
				Stdout: []byte(tt.stdout),
				Stderr: []byte(tt.stderr),
				Status: tt.status,
			}
			assert.Equal(t, expected, e.Run())
		})
	}
}
