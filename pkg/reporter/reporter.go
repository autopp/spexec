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
	"io"
	"os"

	"github.com/autopp/spexec/pkg/model"
)

// Reporter provides formatted test reporter
type Reporter struct {
	w  *Writer
	rf ReportFormatter
}

type Config struct {
	colorMode bool
	w         io.Writer
	rf        ReportFormatter
}

// ReportFormatter is the interface implemented by report formatter
type ReportFormatter interface {
	OnRunStart(w *Writer) error
	OnTestStart(w *Writer, t *model.Test) error
	OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) error
	OnRunComplete(w *Writer, sr *model.SpecResult) error
}

// Option is functional option of New
type Option func(c *Config) error

// WithWriter is a option of New to specify output writer
func WithWriter(w io.Writer) Option {
	return func(c *Config) error {
		c.w = w
		return nil
	}
}

func WithColor(colorMode bool) Option {
	return func(c *Config) error {
		c.colorMode = colorMode
		return nil
	}
}

func WithFormatter(formatter ReportFormatter) Option {
	return func(c *Config) error {
		c.rf = formatter
		return nil
	}
}

// New returns a new Reporter
func New(opts ...Option) (*Reporter, error) {
	r := &Reporter{}

	c := Config{
		w:         os.Stdout,
		colorMode: false,
		rf:        &SimpleFormatter{},
	}
	for _, o := range append([]Option{WithWriter(os.Stdout)}, opts...) {
		if err := o(&c); err != nil {
			return nil, err
		}
	}
	r.w = newWriter(c.w, c.colorMode)
	r.rf = c.rf

	return r, nil
}

// OnRunStart should be called before all test execution
func (r *Reporter) OnRunStart() error {
	return r.rf.OnRunStart(r.w)
}

// OnTestStart should be called before each test execution
func (r *Reporter) OnTestStart(t *model.Test) error {
	return r.rf.OnTestStart(r.w, t)
}

// OnTestComplete should be called after each test execution
func (r *Reporter) OnTestComplete(t *model.Test, tr *model.TestResult) error {
	return r.rf.OnTestComplete(r.w, t, tr)
}

// OnRunComplete should be called afterall test execution
func (r *Reporter) OnRunComplete(sr *model.SpecResult) error {
	return r.rf.OnRunComplete(r.w, sr)
}
