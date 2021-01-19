package test

import (
	"bytes"
	"os/exec"
)

type CommandResult struct {
	Stdout []byte
	Stderr []byte
	Status int
}

type Test struct {
	Command []string
}

func (t *Test) Execute() *CommandResult {
	cmd := exec.Command(t.Command[0], t.Command[1:]...)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Run()

	return &CommandResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		Status: cmd.ProcessState.ExitCode(),
	}
}
