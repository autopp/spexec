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

package exec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Songmu/timeout"
	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/util"
	"golang.org/x/sys/unix"
)

type ExecResult struct {
	Stdout    []byte
	Stderr    []byte
	Status    int
	Signal    os.Signal
	IsTimeout bool
	Err       error
}

type Exec struct {
	Command []string
	Dir     string
	Stdin   []byte
	Env     []util.StringVar
	Timeout time.Duration
}

const defaultTimeout = 10 * time.Second

type Option interface {
	Apply(e *Exec) error
}

type OptionTimeout time.Duration

func (t OptionTimeout) Apply(e *Exec) error {
	if t > 0 {
		e.Timeout = time.Duration(t)
	}
	return nil
}

func WithTimeout(t time.Duration) Option {
	return OptionTimeout(t)
}

func New(command []string, dir string, stdin []byte, env []util.StringVar, opts ...Option) (*Exec, error) {
	e := &Exec{
		Command: command,
		Dir:     dir,
		Stdin:   stdin,
		Env:     env,
	}

	for _, o := range append([]Option{WithTimeout(defaultTimeout)}, opts...) {
		if err := o.Apply(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Exec) Run() *ExecResult {
	cmd := exec.Command(e.Command[0], e.Command[1:]...)
	cmd.Dir = e.Dir
	cmd.Stdin = bytes.NewReader(e.Stdin)
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	for _, v := range e.Env {
		kv := fmt.Sprintf("%s=%s", v.Name, v.Value)
		cmd.Env = append(cmd.Env, kv)
	}

	tio := &timeout.Timeout{
		Cmd:      cmd,
		Duration: e.Timeout,
		Signal:   unix.SIGKILL,
	}
	ch, err := tio.RunCommand()

	if err != nil {
		return &ExecResult{
			Stdout: stdout.Bytes(),
			Stderr: stderr.Bytes(),
			Err:    err,
		}
	}
	es := <-ch
	es.GetExitCode()
	ps := cmd.ProcessState

	if ps.Exited() {
		return &ExecResult{
			Stdout: stdout.Bytes(),
			Stderr: stderr.Bytes(),
			Status: es.GetExitCode(),
		}
	}

	if es.IsTimedOut() {
		return &ExecResult{
			Stdout:    stdout.Bytes(),
			Stderr:    stderr.Bytes(),
			IsTimeout: true,
		}
	}

	sys, ok := ps.Sys().(syscall.WaitStatus)
	if !ok {
		return &ExecResult{
			Stdout: stdout.Bytes(),
			Stderr: stderr.Bytes(),
			Err:    errors.Errorf(errors.ErrInternalError, "unknown (*ProcessState).Sys() type: %T", ps.Sys()),
		}
	}

	ws := unix.WaitStatus(sys)
	if !ws.Signaled() {
		return &ExecResult{
			Stdout: stdout.Bytes(),
			Stderr: stderr.Bytes(),
			Err:    errors.New(errors.ErrInternalError, "process is neither exited nor signaled"),
		}
	}

	return &ExecResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
		Signal: ws.Signal(),
	}
}
