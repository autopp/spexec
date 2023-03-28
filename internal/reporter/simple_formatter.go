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

	"github.com/autopp/spexec/internal/model"
)

/*
SimpleFormatter implements Reporter.

Example of output:

	.F.
	3 examples, 1 failures
*/
type SimpleFormatter struct{}

// OnRunStart is part of Reporter
func (f *SimpleFormatter) OnRunStart(w *Writer) error {
	return nil
}

// OnTestStart is part of Reporter
func (f *SimpleFormatter) OnTestStart(w *Writer, t *model.Test) error {
	return nil
}

// OnTestComplete is part of Reporter
func (f *SimpleFormatter) OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) error {
	if tr.IsSuccess {
		w.UseColor(Green, func() {
			fmt.Fprint(w, ".")
		})
	} else {
		w.UseColor(Red, func() {
			fmt.Fprint(w, "F")
		})
	}
	return nil
}

// OnRunComplete is part of Reporter
func (f *SimpleFormatter) OnRunComplete(w *Writer, sr *model.SpecResult) error {
	printFailures(w, sr.GetFailedTestResults())
	fmt.Fprintf(w, "\n%d examples, %d failures\n", sr.Summary.NumberOfTests, sr.Summary.NumberOfFailed)
	return nil
}
