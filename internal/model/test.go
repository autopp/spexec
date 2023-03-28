// Copyright (C) 2021-2023	 Akira Tanimura (@autopp)
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

package model

import (
	"fmt"
	"time"

	"github.com/Wing924/shellwords"
	"github.com/autopp/spexec/internal/exec"
	"github.com/autopp/spexec/internal/util"
)

type Test struct {
	Name          string
	SpecFilename  string
	Dir           string
	Command       []StringExpr
	Stdin         []byte
	StatusMatcher StatusMatcher
	StdoutMatcher StreamMatcher
	StderrMatcher StreamMatcher
	Env           []util.StringVar
	Timeout       time.Duration
	TeeStdout     bool
	TeeStderr     bool
}

func (t *Test) GetName() string {
	if len(t.Name) != 0 {
		return t.Name
	}

	envStr := ""
	for _, v := range t.Env {
		envStr += v.Name + "=" + v.Value + " "
	}

	command := make([]string, len(t.Command))
	for i, x := range t.Command {
		command[i] = x.String()
	}

	return envStr + shellwords.Join(command)
}

func (t *Test) Run() (*TestResult, error) {
	command, cleanup, err, _ := EvalStringExprs(t.Command)
	// FIXME: error handling
	defer cleanup()
	if err != nil {
		return nil, err
	}

	e, err := exec.New(command, t.Dir, t.Stdin, t.Env, exec.WithTimeout(t.Timeout), exec.WithTeeStdout(t.TeeStdout), exec.WithTeeStderr(t.TeeStderr))
	if err != nil {
		return nil, err
	}

	r := e.Run()
	messages := make([]*AssertionMessage, 0)
	var message string
	statusOk := true

	if r.Err != nil {
		statusOk = false
		messages = append(messages, &AssertionMessage{Name: "status", Message: r.Err.Error()})
	} else if r.IsTimeout {
		statusOk = false
		messages = append(messages, &AssertionMessage{Name: "status", Message: fmt.Sprintf("process was timeout")})
	} else if r.Signal != nil {
		statusOk = false
		messages = append(messages, &AssertionMessage{Name: "status", Message: fmt.Sprintf("process was signaled (%s)", r.Signal.String())})
	} else if t.StatusMatcher != nil {
		statusOk, message, _ = t.StatusMatcher.Match(r.Status)
		if !statusOk {
			messages = append(messages, &AssertionMessage{Name: "status", Message: message})
		}
	}

	stdoutOk := true
	if t.StdoutMatcher != nil {
		stdoutOk, message, _ = t.StdoutMatcher.Match(r.Stdout)
		if !stdoutOk {
			messages = append(messages, &AssertionMessage{Name: "stdout", Message: message})
		}
	}

	stderrOk := true
	if t.StderrMatcher != nil {
		stderrOk, message, _ = t.StderrMatcher.Match(r.Stderr)
		if !stderrOk {
			messages = append(messages, &AssertionMessage{Name: "stderr", Message: message})
		}
	}

	return &TestResult{
		Name:      t.GetName(),
		Messages:  messages,
		IsSuccess: statusOk && stdoutOk && stderrOk,
	}, nil
}
