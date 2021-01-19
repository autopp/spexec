// Copyright (C) 2021 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
