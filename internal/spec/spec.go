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

package spec

import (
	"encoding/json"

	"github.com/autopp/spexec/internal/errors"
	test "github.com/autopp/spexec/internal/model"
	"gopkg.in/yaml.v3"
)

type Map = map[string]interface{}
type Seq = []interface{}

// xxxSchema are structs for documentation, not used
type configSchema struct {
	tests []testSchema
}

type testSchema struct {
	name    string
	command []string
	env     []struct {
		name  string
		value string
	}
	expect *struct {
		status *int
		stdout *string
		stderr *string
	}
}

func LoadYAML(b []byte) ([]*test.Test, error) {
	return load(b, yaml.Unmarshal)
}

func LoadJSON(b []byte) ([]*test.Test, error) {
	return load(b, json.Unmarshal)
}

func load(b []byte, unmarchal func(in []byte, out interface{}) error) ([]*test.Test, error) {
	var x interface{}
	err := unmarchal(b, &x)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, err)
	}

	return parseConfig(x)
}

func parseConfig(c interface{}) ([]*test.Test, error) {
	v := NewValidator()
	cmap, ok := v.MustBeMap(c)
	if !ok {
		return nil, v.Error()
	}

	ts := make([]*test.Test, 0)
	v.MustHaveSeq(cmap, "tests", func(tcs Seq) {
		v.ForInSeq(tcs, func(i int, tc interface{}) {
			t := parseTest(v, tc)
			ts = append(ts, t)
		})
	})

	return ts, v.Error()
}

func parseTest(v *Validator, x interface{}) *test.Test {
	tc, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	t := new(test.Test)
	name, exists, ok := v.MayHaveString(tc, "name")
	if exists {
		t.Name = name
	}

	v.MustHaveSeq(tc, "command", func(command Seq) {
		t.Command = make([]string, len(command))
		v.ForInSeq(command, func(i int, x interface{}) {
			c, _ := v.MustBeString(x)
			t.Command[i] = c
		})

		if len(t.Command) == 0 {
			v.AddViolation("shoud have one ore more elements")
		}
	})

	if stdin, exists, _ := v.MayHaveString(tc, "stdin"); exists {
		t.Stdin = stdin
	}

	v.MayHaveSeq(tc, "env", func(env Seq) {
		t.Env = make(map[string]string)
		v.ForInSeq(env, func(i int, x interface{}) {
			envVar, ok := v.MustBeMap(x)
			if !ok {
				return
			}
			name, nameOk := v.MustHaveString(envVar, "name")
			value, valueOk := v.MustHaveString(envVar, "value")

			if nameOk && valueOk {
				t.Env[name] = value
			}
		})
	})

	v.MayHaveMap(tc, "expect", func(expect Map) {
		status, exists, _ := v.MayHaveInt(expect, "status")
		if exists {
			t.Status = &status
		}

		stdout, exists, _ := v.MayHaveString(expect, "stdout")
		if exists {
			t.Stdout = &stdout
		}

		stderr, exists, _ := v.MayHaveString(expect, "stderr")
		if exists {
			t.Stderr = &stderr
		}
	})

	return t
}