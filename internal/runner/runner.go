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
	statusOk := true
	if t.StatusMatcher != nil {
		statusOk, _, _ = t.StatusMatcher.MatchStatus(r.Status)
	}

	stdoutOk := true
	if t.StdoutMatcher != nil {
		stdoutOk, _, _ = t.StdoutMatcher.MatchStream(r.Stdout)
	}

	stderrOk := true
	if t.StderrMatcher != nil {
		stderrOk, _, _ = t.StderrMatcher.MatchStream(r.Stderr)
	}

	return &model.TestResult{
		Name:      t.Name,
		IsSuccess: statusOk && stdoutOk && stderrOk,
	}
}
