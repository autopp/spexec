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

import (
	"github.com/Wing924/shellwords"
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/util"
)

type Test struct {
	Name          string
	Command       []string
	Stdin         string
	StatusMatcher matcher.StatusMatcher
	StdoutMatcher matcher.StreamMatcher
	StderrMatcher matcher.StreamMatcher
	Env           []util.StringVar
}

func (t *Test) GetName() string {
	if len(t.Name) != 0 {
		return t.Name
	}

	envStr := ""
	for _, v := range t.Env {
		envStr += v.Name + "=" + v.Value + " "
	}
	return envStr + shellwords.Join(t.Command)
}
