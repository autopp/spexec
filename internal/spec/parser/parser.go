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

package parser

import (
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
	"github.com/autopp/spexec/internal/util"
	"gopkg.in/yaml.v3"
)

// xxxSchema are structs for documentation, not used
type specSchema struct {
	spexec string
	tests  []testSchema
}

type testSchema struct {
	name    string
	command []string
	stdin   string
	env     []util.StringVar
	expect  *struct {
		status model.StatusMatcher
		stdout model.StreamMatcher
		stderr model.StreamMatcher
	}
	timeout string
}

var evnVarNamePattern = regexp.MustCompile(`^[a-zA-Z_]\w+$`)

type Parser struct {
	statusMR *matcher.StatusMatcherRegistry
	streamMR *matcher.StreamMatcherRegistry
	isStrict bool
}

func New(statusMR *matcher.StatusMatcherRegistry, streamMR *matcher.StreamMatcherRegistry, isStrict bool) *Parser {
	return &Parser{statusMR, streamMR, isStrict}
}

func (p *Parser) ParseStdin(env *model.Env) ([]*model.Test, error) {
	return p.parseYAML(env, "", os.Stdin)
}

func (p *Parser) ParseFile(env *model.Env, filename string) ([]*model.Test, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}
	defer f.Close()

	var tests []*model.Test
	ext := filepath.Ext(filename)
	if ext == ".yml" || ext == ".yaml" {
		tests, err = p.parseYAML(env, filename, f)
	} else {
		tests, err = p.parseJSON(env, filename, f)
	}

	return tests, err
}

func (p *Parser) parseYAML(env *model.Env, filename string, in io.Reader) ([]*model.Test, error) {
	return p.load(env, filename, in, func(in io.Reader, out interface{}) error {
		return yaml.NewDecoder(in).Decode(out)
	})
}

func (p *Parser) parseJSON(env *model.Env, filename string, in io.Reader) ([]*model.Test, error) {
	return p.load(env, filename, in, util.DecodeJSON)
}

func (p *Parser) load(env *model.Env, filename string, b io.Reader, unmarshal func(in io.Reader, out interface{}) error) ([]*model.Test, error) {
	var x interface{}
	err := unmarshal(b, &x)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}

	v, err := spec.NewValidator(filename)
	if err != nil {
		return nil, err
	}

	return p.loadSpec(env, v, x)
}

func (p *Parser) loadSpec(env *model.Env, v *spec.Validator, c interface{}) ([]*model.Test, error) {
	cmap, ok := v.MustBeMap(c)
	if !ok {
		return nil, v.Error()
	}

	ts := make([]*model.Test, 0)

	if p.isStrict {
		v.MustContainOnly(cmap, "spexec", "tests")
	}

	version, exists, ok := v.MayHaveString(cmap, "spexec")
	if ok && exists {
		if version != "v0" {
			v.InField("spexec", func() {
				v.AddViolation(`should be "v0"`)
			})
		}
	}

	v.MustHaveSeq(cmap, "tests", func(tcs model.Seq) {
		v.ForInSeq(tcs, func(i int, tc interface{}) bool {
			t := p.loadTest(env, v, tc)
			ts = append(ts, t)
			return t != nil
		})
	})

	return ts, v.Error()
}

func (p *Parser) loadTest(env *model.Env, v *spec.Validator, x interface{}) *model.Test {
	tc, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	if p.isStrict {
		v.MustContainOnly(tc, "name", "command", "stdin", "env", "expect", "timeout", "teeStdout", "teeStderr")
	}

	t := new(model.Test)
	t.SpecFilename = v.Filename
	name, exists, ok := v.MayHaveString(tc, "name")
	if exists {
		t.Name = name
	}

	t.Command, _ = v.MustHaveCommand(tc, "command")

	v.MayHave(tc, "stdin", func(stdin interface{}) {
		t.Stdin = p.loadCommandStdin(v, stdin)
	})

	t.Env, _, _ = v.MayHaveEnvSeq(tc, "env")

	if timeout, exists, _ := v.MayHaveDuration(tc, "timeout"); exists {
		t.Timeout = timeout
	}

	v.MayHaveMap(tc, "expect", func(expect model.Map) {
		t.StatusMatcher, t.StdoutMatcher, t.StderrMatcher = p.loadCommandExpect(env, v, expect)
	})

	if teeStdout, exists, _ := v.MayHaveBool(tc, "teeStdout"); exists {
		t.TeeStdout = teeStdout
	}

	if teeStderr, exists, _ := v.MayHaveBool(tc, "teeStderr"); exists {
		t.TeeStderr = teeStderr
	}

	t.Dir = v.GetDir()

	return t
}

func (p *Parser) loadCommandStdin(v *spec.Validator, stdin interface{}) []byte {
	if stdinString, ok := v.MayBeString(stdin); ok {
		return []byte(stdinString)
	} else if stdinMap, ok := v.MayBeMap(stdin); ok {
		if p.isStrict && !v.MustContainOnly(stdinMap, "format", "value") {
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

func (p *Parser) loadCommandExpect(env *model.Env, v *spec.Validator, expect model.Map) (model.StatusMatcher, model.StreamMatcher, model.StreamMatcher) {
	var statusMatcher model.StatusMatcher
	var stdoutMatcher, stderrMatcher model.StreamMatcher
	if p.isStrict {
		v.MustContainOnly(expect, "status", "stdout", "stderr")
	}

	v.MayHave(expect, "status", func(status interface{}) {
		statusMatcher = p.statusMR.ParseMatcher(env, v, status)
	})

	v.MayHave(expect, "stdout", func(stdout interface{}) {
		stdoutMatcher = p.streamMR.ParseMatcher(env, v, stdout)
	})

	v.MayHave(expect, "stderr", func(stderr interface{}) {
		stderrMatcher = p.streamMR.ParseMatcher(env, v, stderr)
	})

	return statusMatcher, stdoutMatcher, stderrMatcher
}
