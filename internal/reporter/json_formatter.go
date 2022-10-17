// Copyright (C) 2021-2022	 Akira Tanimura (@autopp)
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
	"encoding/json"
	"fmt"

	"github.com/autopp/spexec/internal/model"
)

/*
JSONFormatter implements Reporter.

Example of output:

	{
		"results": [
			{
				"isSuccess": false,
				"messages": [
					{
						"name": "status",
						"message": "expected success, but not"
					}
				]
			},
			...
		],
		"summary": {
			"numberOfTests": 10,
			"numberOfSuccess": 6,
			"numberOfFailure": 4,
		}
	}
*/
type JSONFormatter struct{}

type jsonSummary struct {
	NumberOfTests   int `json:"numberOfTests"`
	NumberOfSuccess int `json:"numberOfSuccess"`
	NumberOfFailure int `json:"numberOfFailure"`
}

type jsonOutput struct {
	Results []*model.TestResult `json:"results"`
	Summary jsonSummary         `json:"summary"`
}

// OnRunStart is part of Reporter
func (f *JSONFormatter) OnRunStart(w *Writer) error {
	return nil
}

// OnTestStart is part of Reporter
func (f *JSONFormatter) OnTestStart(w *Writer, t *model.Test) error {
	return nil
}

// OnTestComplete is part of Reporter
func (f *JSONFormatter) OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) error {
	return nil
}

// OnRunComplete is part of Reporter
func (f *JSONFormatter) OnRunComplete(w *Writer, sr *model.SpecResult) error {
	output, err := json.Marshal(sr)

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, string(output))

	return err
}
