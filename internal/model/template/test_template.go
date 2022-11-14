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

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"
	"gopkg.in/yaml.v3"
)

type TemplatableStringVar struct {
	Name  string
	Value *model.Templatable[string]
}

type TestTemplate struct {
	Name          *model.Templatable[string]
	SpecFilename  string
	Dir           string
	Command       []*model.Templatable[any]
	Stdin         *model.Templatable[any]
	StatusMatcher *model.Templatable[any]
	StdoutMatcher *model.Templatable[any]
	StderrMatcher *model.Templatable[any]
	Env           []*TemplatableStringVar
	Timeout       time.Duration
	TeeStdout     bool
	TeeStderr     bool
}

// TODO: set validator path
func (tt *TestTemplate) Expand(env *model.Env, v *model.Validator, statusMR *matcher.StatusMatcherRegistry, streamMR *matcher.StreamMatcherRegistry) (*model.Test, error) {
	name := ""
	if tt.Name != nil {
		var err error
		name, err = tt.Name.Expand(env, v)
		if err != nil {
			return nil, err
		}
	}

	command := make([]model.StringExpr, 0, len(tt.Command))
	for _, ct := range tt.Command {
		x, err := ct.Expand(env, v)
		if err != nil {
			return nil, err
		}

		// TODO: error handling
		c, _ := v.MustBeStringExpr(x)
		command = append(command, c)
	}

	evaledStdin := []byte("")
	if tt.Stdin != nil {
		stdin, err := tt.Stdin.Expand(env, v)
		if err != nil {
			return nil, err
		}
		evaledStdin = evalCommandStdin(v, stdin)
		if evaledStdin == nil {
			// TODO: error handling
			return nil, errors.New(errors.ErrInvalidSpec, "cannot load stdin")
		}
	}

	var statusMatcher model.StatusMatcher
	if tt.StatusMatcher != nil {
		status, err := tt.StatusMatcher.Expand(env, v)
		if err != nil {
			return nil, err
		}
		statusMatcher = statusMR.ParseMatcher(v, status)
	}

	var stdoutMatcher model.StreamMatcher
	if tt.StdoutMatcher != nil {
		stdout, err := tt.StdoutMatcher.Expand(env, v)
		if err != nil {
			return nil, err
		}
		stdoutMatcher = streamMR.ParseMatcher(v, stdout)
	}

	var stderrMatcher model.StreamMatcher
	if tt.StderrMatcher != nil {
		stderr, err := tt.StderrMatcher.Expand(env, v)
		if err != nil {
			return nil, err
		}
		stderrMatcher = streamMR.ParseMatcher(v, stderr)
	}

	tEnv := make([]util.StringVar, 0, len(tt.Env))
	for _, tsv := range tt.Env {
		value, err := tsv.Value.Expand(env, v)
		if err != nil {
			return nil, err
		}

		tEnv = append(tEnv, util.StringVar{Name: tsv.Name, Value: value})
	}

	return &model.Test{
		Name:          name,
		SpecFilename:  tt.SpecFilename,
		Dir:           tt.Dir,
		Command:       command,
		Stdin:         evaledStdin,
		StatusMatcher: statusMatcher,
		StdoutMatcher: stdoutMatcher,
		StderrMatcher: stderrMatcher,
		Env:           tEnv,
		Timeout:       tt.Timeout,
		TeeStdout:     tt.TeeStdout,
		TeeStderr:     tt.TeeStderr,
	}, nil
}

func evalCommandStdin(v *model.Validator, stdin any) []byte {
	if stdinString, ok := v.MayBeString(stdin); ok {
		return []byte(stdinString)
	} else if stdinMap, ok := v.MayBeMap(stdin); ok {
		if !v.MustContainOnly(stdinMap, "format", "value") {
			return nil
		}

		stdinFormat, formatOk := v.MustHaveString(stdinMap, "format")
		stdinValue, valueOk := v.MustHave(stdinMap, "value")
		if !formatOk || !valueOk {
			return nil
		}

		switch stdinFormat {
		case "yaml":
			value, err := yaml.Marshal(stdinValue)
			if err != nil {
				v.InField("value", func() {
					v.AddViolation(`cannot encode to a YAML string: %s`, err)
				})
				return nil
			}
			return value
		default:
			v.InField("format", func() {
				v.AddViolation(`should be a "yaml", but is %q`, stdinFormat)
			})
			return nil
		}
	} else {
		v.AddViolation("should be a string or map, but is %s", model.TypeNameOf(stdin))
		return nil
	}
}
