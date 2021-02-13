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
	"io"

	"github.com/autopp/spexec/internal/test"
)

// Reporter is the interface implemented by test reporter
type Reporter interface {
	OnRunStart(w io.Writer)
	OnTestStart(w io.Writer, t *test.Test)
	OnTestComplete(w io.Writer, t *test.Test, tr *TestResult)
	OnRunComplete(w io.Writer, trs []*TestResult)
}

func NewReporter() Reporter {
	return &SimpleReporter{}
}

type SimpleReporter struct{}

func (sr *SimpleReporter) OnRunStart(w io.Writer) {
}

func (sr *SimpleReporter) OnTestStart(w io.Writer, t *test.Test) {
}

func (sr *SimpleReporter) OnTestComplete(w io.Writer, t *test.Test, tr *TestResult) {
	if tr.IsSuccess {
		fmt.Fprint(w, ".")
	} else {
		fmt.Fprint(w, "f")
	}
}

func (sr *SimpleReporter) OnRunComplete(w io.Writer, trs []*TestResult) {
	failures := 0
	for _, tr := range trs {
		if !tr.IsSuccess {
			failures++
		}
	}
	fmt.Fprintf(w, "\n%d examples, %d failures\n", len(trs), failures)
}
