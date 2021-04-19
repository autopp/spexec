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

package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/autopp/spexec/internal/model"
)

type ExecResult struct {
	Stdout []byte
	Stderr []byte
	Status int
}

type Exec struct {
	Command []string
	Stdin   string
	Env     map[string]string
}

func NewExec(t *model.Test) *Exec {
	return &Exec{
		Command: t.Command,
		Stdin:   t.Stdin,
		Env:     t.Env,
	}
}

func (e *Exec) Run() *ExecResult {
	cmd := exec.Command(e.Command[0], e.Command[1:]...)
	cmd.Stdin = strings.NewReader(e.Stdin)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	for name, v := range e.Env {
		kv := fmt.Sprintf("%s=%s", name, v)
		cmd.Env = append(cmd.Env, kv)
	}
	cmd.Run()

	return &ExecResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		Status: cmd.ProcessState.ExitCode(),
	}
}
