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

package template

import (
	"time"

	"github.com/autopp/spexec/internal/model"
)

type TemplatableStringVar struct {
	Name  string
	Value model.Templatable[string]
}

type TestTemplate struct {
	Name          model.Templatable[string]
	SpecFilename  string
	Dir           model.Templatable[string]
	Command       []model.Templatable[model.StringExpr]
	Stdin         model.Templatable[string]
	StatusMatcher model.Templatable[any]
	StdoutMatcher model.Templatable[any]
	StderrMatcher model.Templatable[any]
	Env           []TemplatableStringVar
	Timeout       time.Duration
	TeeStdout     model.Templatable[bool]
	TeeStderr     model.Templatable[bool]
}
