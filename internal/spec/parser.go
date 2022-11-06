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

package spec

import (
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/model/template"
	"github.com/autopp/spexec/internal/util"
	"gopkg.in/yaml.v3"
)

var evnVarNamePattern = regexp.MustCompile(`^[a-zA-Z_]\w+$`)

type Parser struct {
	statusMR *matcher.StatusMatcherRegistry
	streamMR *matcher.StreamMatcherRegistry
}

func NewParser(statusMR *matcher.StatusMatcherRegistry, streamMR *matcher.StreamMatcherRegistry) *Parser {
	return &Parser{statusMR, streamMR}
}

func (p *Parser) ParseStdin(env *model.Env, v *model.Validator) ([]*template.TestTemplate, error) {
	return p.parseYAML(env, v, "", os.Stdin)
}

func (p *Parser) ParseFile(env *model.Env, v *model.Validator, filename string) ([]*template.TestTemplate, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}
	defer f.Close()

	var tests []*template.TestTemplate
	ext := filepath.Ext(filename)
	if ext == ".yml" || ext == ".yaml" {
		tests, err = p.parseYAML(env, v, filename, f)
	} else {
		tests, err = p.parseJSON(env, v, filename, f)
	}

	return tests, err
}

func (p *Parser) parseYAML(env *model.Env, v *model.Validator, filename string, in io.Reader) ([]*template.TestTemplate, error) {
	return p.load(env, v, filename, in, func(in io.Reader, out any) error {
		return yaml.NewDecoder(in).Decode(out)
	})
}

func (p *Parser) parseJSON(env *model.Env, v *model.Validator, filename string, in io.Reader) ([]*template.TestTemplate, error) {
	return p.load(env, v, filename, in, util.DecodeJSON)
}

func (p *Parser) load(env *model.Env, v *model.Validator, filename string, b io.Reader, unmarshal func(in io.Reader, out any) error) ([]*template.TestTemplate, error) {
	var x any
	err := unmarshal(b, &x)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}

	return p.loadSpec(env, v, x)
}

func (p *Parser) loadSpec(env *model.Env, v *model.Validator, c any) ([]*template.TestTemplate, error) {
	cmap, ok := v.MustBeMap(c)
	if !ok {
		return nil, v.Error()
	}

	ts := make([]*template.TestTemplate, 0)

	v.MustContainOnly(cmap, "spexec", "tests")

	version, exists, ok := v.MayHaveString(cmap, "spexec")
	if ok && exists {
		if version != "v0" {
			v.InField("spexec", func() {
				v.AddViolation(`should be "v0"`)
			})
		}
	}

	v.MustHaveSeq(cmap, "tests", func(tcs model.Seq) {
		v.ForInSeq(tcs, func(i int, tc any) bool {
			t := p.loadTest(env, v, tc)
			ts = append(ts, t)
			return t != nil
		})
	})

	return ts, v.Error()
}

func (p *Parser) loadTest(env *model.Env, v *model.Validator, x any) *template.TestTemplate {
	tc, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	v.MustContainOnly(tc, "name", "command", "stdin", "env", "expect", "timeout", "teeStdout", "teeStderr")

	tt := new(template.TestTemplate)
	tt.SpecFilename = v.Filename
	name, exists, ok := v.MayHaveTemplatableString(tc, "name")
	if exists {
		tt.Name = name
	}

	v.MustHaveSeq(tc, "command", func(seq model.Seq) {
		command := make([]*model.Templatable[any], 0)
		v.ForInSeq(seq, func(i int, x any) bool {
			c, ok := v.MustBeTemplatable(x)
			command = append(command, c)
			return ok
		})

		tt.Command = command
	})

	v.MayHave(tc, "stdin", func(stdin any) {
		tt.Stdin, _ = v.MustBeTemplatable(stdin)
	})

	v.MayHaveSeq(tc, "env", func(seq model.Seq) {
		env := make([]*template.TemplatableStringVar, 0)
		v.ForInSeq(seq, func(i int, x any) bool {
			m, ok := v.MustBeMap(x)
			if !ok {
				return false
			}
			name, ok := v.MustHaveString(m, "name")
			if !ok {
				return false
			}

			value, ok := v.MustHaveTemplatableString(m, "value")
			if !ok {
				return false
			}

			env = append(env, &template.TemplatableStringVar{Name: name, Value: value})
			return true
		})

		tt.Env = env
	})

	if timeout, exists, _ := v.MayHaveDuration(tc, "timeout"); exists {
		tt.Timeout = timeout
	}

	v.MayHaveMap(tc, "expect", func(expect model.Map) {
		tt.StatusMatcher, tt.StdoutMatcher, tt.StderrMatcher = p.loadCommandExpect(env, v, expect)
	})

	if teeStdout, exists, _ := v.MayHaveBool(tc, "teeStdout"); exists {
		tt.TeeStdout = teeStdout
	}

	if teeStderr, exists, _ := v.MayHaveBool(tc, "teeStderr"); exists {
		tt.TeeStderr = teeStderr
	}

	tt.Dir = v.GetDir()

	return tt
}

func (p *Parser) loadCommandExpect(env *model.Env, v *model.Validator, expect model.Map) (*model.Templatable[any], *model.Templatable[any], *model.Templatable[any]) {
	var statusMatcher, stdoutMatcher, stderrMatcher *model.Templatable[any]
	v.MustContainOnly(expect, "status", "stdout", "stderr")

	v.MayHave(expect, "status", func(status any) {
		statusMatcher, _ = v.MustBeTemplatable(status)
	})

	v.MayHave(expect, "stdout", func(stdout any) {
		stdoutMatcher, _ = v.MustBeTemplatable(stdout)
	})

	v.MayHave(expect, "stderr", func(stderr any) {
		stderrMatcher, _ = v.MustBeTemplatable(stderr)
	})

	return statusMatcher, stdoutMatcher, stderrMatcher
}
