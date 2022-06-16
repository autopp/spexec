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

package model

import (
	"time"
)

type TemplatableStringVar struct {
	Name  string
	Value Templatable[string]
}

type TestTemplate struct {
	Name          Templatable[string]
	SpecFilename  string
	Dir           Templatable[string]
	Command       []Templatable[StringExpr]
	Stdin         Templatable[string]
	StatusMatcher Templatable[any]
	StdoutMatcher Templatable[any]
	StderrMatcher Templatable[any]
	Env           []TemplatableStringVar
	Timeout       time.Duration
	TeeStdout     Templatable[bool]
	TeeStderr     Templatable[bool]
}
