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
	"io"

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
func (sr *SimpleFormatter) OnRunStart(w io.Writer) {
}

// OnTestStart is part of Reporter
func (sr *SimpleFormatter) OnTestStart(w io.Writer, t *model.Test) {
}

// OnTestComplete is part of Reporter
func (sr *SimpleFormatter) OnTestComplete(w io.Writer, t *model.Test, tr *model.TestResult) {
	if tr.IsSuccess {
		fmt.Fprint(w, ".")
	} else {
		fmt.Fprint(w, "f")
	}
}

// OnRunComplete is part of Reporter
func (sr *SimpleFormatter) OnRunComplete(w io.Writer, trs []*model.TestResult) {
	failures := 0
	for _, tr := range trs {
		if !tr.IsSuccess {
			failures++
		}
	}
	fmt.Fprintf(w, "\n%d examples, %d failures\n", len(trs), failures)
}
