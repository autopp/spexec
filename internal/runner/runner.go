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
	"os"

	"github.com/autopp/spexec/internal/model"
)

type Runner struct{}

type TestResult struct {
	Name      string
	IsSuccess bool
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunTests(tests []*model.Test) []*TestResult {
	results := make([]*TestResult, 0, len(tests))
	reporter := NewReporter()
	w := os.Stdout

	reporter.OnRunStart(w)
	for _, t := range tests {
		reporter.OnTestStart(w, t)
		er := NewExec(t).Run()
		tr := assertResult(t, er)
		reporter.OnTestComplete(w, t, tr)
		results = append(results, tr)
	}
	reporter.OnRunComplete(w, results)

	return results
}

func assertResult(t *model.Test, r *ExecResult) *TestResult {
	return &TestResult{
		Name:      t.Name,
		IsSuccess: (t.Status == nil || *t.Status == r.Status) && (t.Stdout == nil || bytes.Equal([]byte(*t.Stdout), r.Stdout)) && (t.Stderr == nil || bytes.Equal([]byte(*t.Stderr), r.Stderr)),
	}
}
