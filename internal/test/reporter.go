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

// Reporter is the interface implemented by test reporter
type Reporter interface {
	OnRunStart()
	OnTestStart(t *Test)
	OnTestComplete(t *Test, tr *TestResult)
	OnRunComplete(trs []*TestResult)
}

func NewReporter() Reporter {
	return &SimpleReporter{}
}

type SimpleReporter struct{}

func (sr *SimpleReporter) OnRunStart() {
}

func (sr *SimpleReporter) OnTestStart(t *Test) {
}

func (sr *SimpleReporter) OnTestComplete(t *Test, tr *TestResult) {
	if tr.IsSuccess {
		fmt.Print(".")
	} else {
		fmt.Print("f")
	}
}

func (sr *SimpleReporter) OnRunComplete(trs []*TestResult) {
	failures := 0
	for _, tr := range trs {
		if !tr.IsSuccess {
			failures++
		}
	}
	fmt.Printf("\n%d examples, %d failures\n", len(trs), failures)
}
