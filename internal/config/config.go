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

package config

import (
	"io"

	"github.com/autopp/spexec/internal/test"
	"gopkg.in/yaml.v3"
)

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

func Load(r io.Reader) ([]*test.Test, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var x interface{}
	err = yaml.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}

	ts, err := parseConfig(x)
	return ts, err
}

func parseConfig(c interface{}) ([]*test.Test, error) {
	v := newValidator()
	cmap, ok := v.MustBeMap(c)
	if !ok {
		return nil, v.Error()
	}

	ts := make([]*test.Test, 0)
	v.MustHaveSeq(cmap, "tests", func(tcs configSeq) {
		v.ForInSeq(tcs, func(i int, tc interface{}) {
			t := parseTest(v, tc)
			ts = append(ts, t)
		})
	})

	return ts, v.Error()
}

func parseTest(v *validator, x interface{}) *test.Test {
	tc, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	t := new(test.Test)
	name, exists, ok := v.MayHaveString(tc, "name")
	if exists {
		t.Name = name
	}

	v.MustHaveSeq(tc, "command", func(command configSeq) {
		t.Command = make([]string, len(command))
		v.ForInSeq(command, func(i int, x interface{}) {
			c, _ := v.MustBeString(x)
			t.Command[i] = c
		})
	})

	v.MayHaveSeq(tc, "env", func(env configSeq) {
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

	v.MayHaveMap(tc, "expect", func(expect configMap) {
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
