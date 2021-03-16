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
	w  *ReportWriter
	rf ReportFormatter
}

// ReportWriter extends io.Writer
type ReportWriter struct {
	io.Writer
}

// ReportFormatter is the interface implemented by report formatter
type ReportFormatter interface {
	OnRunStart(w *ReportWriter)
	OnTestStart(w *ReportWriter, t *model.Test)
	OnTestComplete(w *ReportWriter, t *model.Test, tr *model.TestResult)
	OnRunComplete(w *ReportWriter, trs []*model.TestResult)
}

// Option is functional option of New
type Option func(r *Reporter) error

// WithWriter is a option of New to specify output writer
func WithWriter(w io.Writer) Option {
	return func(r *Reporter) error {
		r.w = &ReportWriter{w}
		return nil
	}
}

// New returns a new Reporter
func New(opts ...Option) (*Reporter, error) {
	r := &Reporter{rf: &SimpleFormatter{}}

	for _, o := range append([]Option{WithWriter(os.Stdout)}, opts...) {
		if err := o(r); err != nil {
			return nil, err
		}
	}
	return r, nil
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
