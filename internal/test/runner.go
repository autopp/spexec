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

import "fmt"

type Runner struct{}

type TestResult struct {
	IsSuccess bool
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunTests(tests []*Test) []*TestResult {
	results := make([]*TestResult, 0, len(tests))

	for _, t := range tests {
		er := t.ToExec().Run()
		tr := assertResult(t, er)
		if tr.IsSuccess {
			fmt.Print(".")
		} else {
			fmt.Print("F")
		}
		results = append(results, tr)
	}
	fmt.Println("")

	return results
}

func assertResult(t *Test, r *ExecResult) *TestResult {
	return &TestResult{
		IsSuccess: (t.Status == nil || *t.Status == r.Status) && (t.Stdout == nil || *t.Stdout == string(r.Stdout)) && (t.Stderr == nil || *t.Stderr == string(r.Stderr)),
	}
}
