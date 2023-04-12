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

package reporter

import (
	"fmt"

	"github.com/autopp/spexec/pkg/model"
)

/*
DocumentationFormatter implements Reporter.

Example of output:

	.F.
	3 examples, 1 failures
*/
type DocumentationFormatter struct{}

// OnRunStart is part of Reporter
func (f *DocumentationFormatter) OnRunStart(w *Writer) error {
	return nil
}

// OnTestStart is part of Reporter
func (f *DocumentationFormatter) OnTestStart(w *Writer, t *model.Test) error {
	return nil
}

// OnTestComplete is part of Reporter
func (f *DocumentationFormatter) OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) error {
	var color Color
	if tr.IsSuccess {
		color = Green
	} else {
		color = Red
	}

	w.UseColor(color, func() {
		fmt.Fprintln(w, t.GetName())
	})

	return nil
}

// OnRunComplete is part of Reporter
func (f *DocumentationFormatter) OnRunComplete(w *Writer, sr *model.SpecResult) error {
	failed := sr.GetFailedTestResults()
	var color Color = Green
	if len(failed) > 0 {
		color = Red
	}
	printFailures(w, failed)
	w.UseColor(color, func() {
		fmt.Fprintf(w, "\n%d examples, %d failures\n", sr.Summary.NumberOfSucceeded, sr.Summary.NumberOfFailed)
	})

	return nil
}
