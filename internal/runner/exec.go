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
	"syscall"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"
	"golang.org/x/sys/unix"
)

type ExecResult struct {
	Stdout []byte
	Stderr []byte
	ps     *os.ProcessState
}

type Exec struct {
	Command []string
	Stdin   string
	Env     []util.StringVar
}

func NewExec(t *model.Test) *Exec {
	return &Exec{
		Command: t.Command,
		Stdin:   t.Stdin,
		Env:     t.Env,
	}
}

func (er *ExecResult) WaitStatus() (int, os.Signal, error) {
	if er.ps.Exited() {
		return er.ps.ExitCode(), nil, nil
	}

	sys, ok := er.ps.Sys().(syscall.WaitStatus)
	if !ok {
		return -1, nil, errors.Errorf(errors.ErrInternalError, "unknown (*ProcessState).Sys() type: %T", er.ps.Sys())
	}

	ws := unix.WaitStatus(sys)
	if !ws.Signaled() {
		return -1, nil, errors.New(errors.ErrInternalError, "process is neither exited nor signaled")
	}

	return -1, ws.Signal(), nil
}

func (e *Exec) Run() *ExecResult {
	cmd := exec.Command(e.Command[0], e.Command[1:]...)
	cmd.Stdin = strings.NewReader(e.Stdin)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	for _, v := range e.Env {
		kv := fmt.Sprintf("%s=%s", v.Name, v.Value)
		cmd.Env = append(cmd.Env, kv)
	}
	cmd.Run()

	return &ExecResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		ps:     cmd.ProcessState,
	}
}
