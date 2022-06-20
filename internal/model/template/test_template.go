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

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
	"github.com/autopp/spexec/internal/util"
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
	TeeStdout     bool
	TeeStderr     bool
}

func (tt *TestTemplate) Expand(env *model.Env, v *spec.Validator, statusMR matcher.StatusMatcherRegistry, streamMR matcher.StreamMatcherRegistry) (*model.Test, error) {
	name, err := tt.Name.Expand(env)
	if err != nil {
		return nil, err
	}

	dir, err := tt.Dir.Expand(env)
	if err != nil {
		return nil, err
	}

	command := make([]model.StringExpr, 0, len(tt.Command))
	for _, ct := range tt.Command {
		c, err := ct.Expand(env)
		if err != nil {
			return nil, err
		}

		command = append(command, c)
	}

	stdin, err := tt.Stdin.Expand(env)
	if err != nil {
		return nil, err
	}

	status, err := tt.StatusMatcher.Expand(env)
	if err != nil {
		return nil, err
	}
	statusMatcher := statusMR.ParseMatcher(env, v, status)

	stdout, err := tt.StdoutMatcher.Expand(env)
	if err != nil {
		return nil, err
	}
	stdoutMatcher := streamMR.ParseMatcher(env, v, stdout)

	stderr, err := tt.StderrMatcher.Expand(env)
	if err != nil {
		return nil, err
	}
	stderrMatcher := streamMR.ParseMatcher(env, v, stderr)

	tEnv := make([]util.StringVar, 0, len(tt.Env))
	for _, tsv := range tt.Env {
		value, err := tsv.Value.Expand(env)
		if err != nil {
			return nil, err
		}

		tEnv = append(tEnv, util.StringVar{Name: tsv.Name, Value: value})
	}

	return &model.Test{
		Name:          name,
		SpecFilename:  tt.SpecFilename,
		Dir:           dir,
		Command:       command,
		Stdin:         []byte(stdin),
		StatusMatcher: statusMatcher,
		StdoutMatcher: stdoutMatcher,
		StderrMatcher: stderrMatcher,
		Env:           tEnv,
		Timeout:       tt.Timeout,
		TeeStdout:     tt.TeeStdout,
		TeeStderr:     tt.TeeStderr,
	}, nil
}
