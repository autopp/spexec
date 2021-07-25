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
	"fmt"

	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/reporter"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunTests(tests []*model.Test, reporter *reporter.Reporter) []*model.TestResult {
	results := make([]*model.TestResult, 0, len(tests))

	reporter.OnRunStart()
	for _, t := range tests {
		reporter.OnTestStart(t)
		er := NewExec(t).Run()
		tr := assertResult(t, er)
		reporter.OnTestComplete(t, tr)
		results = append(results, tr)
	}
	reporter.OnRunComplete(results)

	return results
}

func assertResult(t *model.Test, r *ExecResult) *model.TestResult {
	messages := make([]*model.AssertionMessage, 0)
	var message string
	statusOk := true

	if status, sig, err := r.WaitStatus(); err != nil {
		statusOk = false
		messages = append(messages, &model.AssertionMessage{Name: "status", Message: err.Error()})
	} else if sig != nil {
		statusOk = false
		messages = append(messages, &model.AssertionMessage{Name: "status", Message: fmt.Sprintf("process signaled (%s)", sig.String())})
	} else if t.StatusMatcher != nil {
		statusOk, message, _ = t.StatusMatcher.MatchStatus(status)
		if !statusOk {
			messages = append(messages, &model.AssertionMessage{Name: "status", Message: message})
		}
	}

	stdoutOk := true
	if t.StdoutMatcher != nil {
		stdoutOk, message, _ = t.StdoutMatcher.MatchStream(r.Stdout)
		if !stdoutOk {
			messages = append(messages, &model.AssertionMessage{Name: "stdout", Message: message})
		}
	}

	stderrOk := true
	if t.StderrMatcher != nil {
		stderrOk, message, _ = t.StderrMatcher.MatchStream(r.Stderr)
		if !stderrOk {
			messages = append(messages, &model.AssertionMessage{Name: "stderr", Message: message})
		}
	}

	return &model.TestResult{
		Name:      t.GetName(),
		Messages:  messages,
		IsSuccess: statusOk && stdoutOk && stderrOk,
	}
}
