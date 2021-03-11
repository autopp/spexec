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

	"github.com/autopp/spexec/internal/model"
)

// Reporter is the interface implemented by test reporter
type Reporter interface {
	OnRunStart(w io.Writer)
	OnTestStart(w io.Writer, t *model.Test)
	OnTestComplete(w io.Writer, t *model.Test, tr *model.TestResult)
	OnRunComplete(w io.Writer, trs []*model.TestResult)
}

func NewReporter() Reporter {
	return &SimpleReporter{}
}
