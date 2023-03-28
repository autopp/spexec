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

package runner

import (
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/reporter"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) RunTests(name string, tests []*model.Test, reporter *reporter.Reporter) ([]*model.TestResult, error) {
	results := make([]*model.TestResult, 0, len(tests))

	if err := reporter.OnRunStart(); err != nil {
		return nil, err
	}
	for _, t := range tests {
		if err := reporter.OnTestStart(t); err != nil {
			return nil, err
		}

		tr, err := t.Run()
		if err != nil {
			return nil, err
		}

		if err := reporter.OnTestComplete(t, tr); err != nil {
			return nil, err
		}
		results = append(results, tr)
	}
	sr := model.NewSpecResult(name, results)
	if err := reporter.OnRunComplete(sr); err != nil {
		return nil, err
	}

	return results, nil
}
