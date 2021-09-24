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

package model

type AssertionMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type TestResult struct {
	Name      string              `json:"name"`
	Messages  []*AssertionMessage `json:"messages"`
	IsSuccess bool                `json:"isSuccess"`
}

type SpecSummary struct {
	NumberOfTests     int `json:"numberOfTests"`
	NumberOfSucceeded int `json:"numberOfSucceeded"`
	NumberOfFailed    int `json:"numberOfFailed"`
}

type SpecResult struct {
	TestResults []*TestResult `json:"testResults"`
	Summary     SpecSummary   `json:"summary"`
}

func NewSpecResult(testResults []*TestResult) *SpecResult {
	sr := &SpecResult{
		TestResults: testResults,
	}

	sr.Summary.NumberOfTests = len(testResults)
	sr.Summary.NumberOfFailed = len(sr.GetFailedTestResults())
	sr.Summary.NumberOfSucceeded = sr.Summary.NumberOfTests - sr.Summary.NumberOfFailed
	return sr
}

func (sr *SpecResult) GetFailedTestResults() []*TestResult {
	failures := make([]*TestResult, 0)
	for _, tr := range sr.TestResults {
		if !tr.IsSuccess {
			failures = append(failures, tr)
		}
	}

	return failures
}
