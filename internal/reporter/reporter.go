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
	"io"
	"os"

	"github.com/autopp/spexec/internal/model"
)

// Reporter provides formatted test reporter
type Reporter struct {
	w  io.Writer
	rf ReportFormatter
}

// ReportFormatter is the interface implemented by report formatter
type ReportFormatter interface {
	OnRunStart(w io.Writer)
	OnTestStart(w io.Writer, t *model.Test)
	OnTestComplete(w io.Writer, t *model.Test, tr *model.TestResult)
	OnRunComplete(w io.Writer, trs []*model.TestResult)
}

// New returns a new Reporter
func New() *Reporter {
	return &Reporter{w: os.Stdout, rf: &SimpleFormatter{}}
}

// OnRunStart should be called before all test execution
func (r *Reporter) OnRunStart() {
	r.rf.OnRunStart(r.w)
}

// OnTestStart should be called before each test execution
func (r *Reporter) OnTestStart(t *model.Test) {
	r.rf.OnTestStart(r.w, t)
}

// OnTestComplete should be called after each test execution
func (r *Reporter) OnTestComplete(t *model.Test, tr *model.TestResult) {
	r.rf.OnTestComplete(r.w, t, tr)
}

// OnRunComplete should be called afterall test execution
func (r *Reporter) OnRunComplete(trs []*model.TestResult) {
	r.rf.OnRunComplete(r.w, trs)
}
