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

func (p *Parser) ParseStdin() ([]*model.Test, error) {
	f, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}

	tests, err := p.parseYAML("", f)
	return tests, err
}

func (p *Parser) ParseFile(filename string) ([]*model.Test, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}
	var tests []*model.Test
	ext := filepath.Ext(filename)
	if ext == ".yml" || ext == ".yaml" {
		tests, err = p.parseYAML(filename, f)
	} else {
		tests, err = p.parseJSON(filename, f)
	}

	return tests, err
}

func (p *Parser) parseYAML(filename string, b []byte) ([]*model.Test, error) {
	return p.load(filename, b, yaml.Unmarshal)
}

func (p *Parser) parseJSON(filename string, b []byte) ([]*model.Test, error) {
	return p.load(filename, b, util.UnmarshalJSON)
}

func (p *Parser) load(filename string, b []byte, unmarchal func(in []byte, out interface{}) error) ([]*model.Test, error) {
	var x interface{}
	err := unmarchal(b, &x)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidSpec, err)
	}

	return p.loadSpec(filename, x)
}

func (p *Parser) loadSpec(filename string, c interface{}) ([]*model.Test, error) {
	v, err := spec.NewValidator(filename)
	if err != nil {
		return nil, err
	}
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

	v.MustHaveSeq(cmap, "tests", func(tcs spec.Seq) {
		v.ForInSeq(tcs, func(i int, tc interface{}) bool {
			t := p.loadTest(v, tc)
			ts = append(ts, t)
			return t != nil
		})
	})

	return ts, v.Error()
}

func (p *Parser) loadTest(v *spec.Validator, x interface{}) *model.Test {
	tc, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	if p.isStrict {
		v.MustContainOnly(tc, "name", "command", "stdin", "env", "expect", "timeout")
	}

	t := new(model.Test)
	name, exists, ok := v.MayHaveString(tc, "name")
	if exists {
		t.Name = name
	}

	t.Command, _ = v.MustHaveCommand(tc, "command")

	v.MayHave(tc, "stdin", func(stdin interface{}) {
		if stdinString, ok := v.MayBeString(stdin); ok {
			t.Stdin = []byte(stdinString)
		} else if stdinMap, ok := v.MayBeMap(stdin); ok {
			if p.isStrict {
				v.MustContainOnly(stdinMap, "type", "value")
			}

			stdinType, typeOk := v.MustHaveString(stdinMap, "type")
			stdinValue, valueOk := v.MustHave(stdinMap, "value")
			if !typeOk || !valueOk {
				return
			}

			switch stdinType {
			case "yaml":
				value, err := yaml.Marshal(stdinValue)
				if err != nil {
					v.InField("value", func() {
						v.AddViolation(`cannot encode to a YAML string: %s`, err)
					})
					return
				}
				t.Stdin = value
			default:
				v.InField("type", func() {
					v.AddViolation(`should be a "yaml", but is %q`, stdinType)
				})
			}
		} else {
			v.AddViolation("should be a string or map, but is %s", spec.TypeNameOf(stdin))
		}
	})

	t.Env, _, _ = v.MayHaveEnvSeq(tc, "env")

	if timeout, exists, _ := v.MayHaveDuration(tc, "timeout"); exists {
		t.Timeout = timeout
	}

	v.MayHaveMap(tc, "expect", func(expect spec.Map) {
		if p.isStrict {
			v.MustContainOnly(expect, "status", "stdout", "stderr")
		}

		v.MayHave(expect, "status", func(status interface{}) {
			t.StatusMatcher = p.statusMR.ParseMatcher(v, status)
		})

		v.MayHave(expect, "stdout", func(stdout interface{}) {
			t.StdoutMatcher = p.streamMR.ParseMatcher(v, stdout)
		})

		v.MayHave(expect, "stderr", func(stderr interface{}) {
			t.StderrMatcher = p.streamMR.ParseMatcher(v, stderr)
		})
	})

	t.Dir = v.GetDir()

	return t
}
