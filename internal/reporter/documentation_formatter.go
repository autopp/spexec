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

package reporter

import (
	"fmt"

	"github.com/autopp/spexec/internal/model"
)

/*
DocumentationFormatter implements Reporter.

Example of output:

	.F.
	3 examples, 1 failures
*/
type DocumentationFormatter struct{}

// OnRunStart is part of Reporter
func (sr *DocumentationFormatter) OnRunStart(w *Writer) {
}

// OnTestStart is part of Reporter
func (sr *DocumentationFormatter) OnTestStart(w *Writer, t *model.Test) {
}

// OnTestComplete is part of Reporter
func (sr *DocumentationFormatter) OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) {
	var color Color
	if tr.IsSuccess {
		color = Green
	} else {
		color = Red
	}

	w.UseColor(color, func() {
		fmt.Fprintln(w, t.GetName())
	})
}

// OnRunComplete is part of Reporter
func (sr *DocumentationFormatter) OnRunComplete(w *Writer, trs []*model.TestResult) {
	failures := make([]*model.TestResult, 0)
	for _, tr := range trs {
		if !tr.IsSuccess {
			failures = append(failures, tr)
		}
	}
	var color Color = Green
	if len(failures) > 0 {
		color = Red
	}

	printFailures(w, failures)
	w.UseColor(color, func() {
		fmt.Fprintf(w, "\n%d examples, %d failures\n", len(trs), len(failures))
	})
}
