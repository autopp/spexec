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
	"io/ioutil"

	"github.com/autopp/spexec/internal/test"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Tests []Test `yaml:"tests"`
}

type Test struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
	Env     []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"env"`
	Expect *struct {
		Status *int    `yaml:"status"`
		Stdout *string `yaml:"stdout"`
		Stderr *string `yaml:"stderr"`
	} `yaml:"expect"`
}

func Load(r io.Reader) ([]*test.Test, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	tests := make([]*test.Test, len(c.Tests))
	for i, t := range c.Tests {
		tests[i] = newTest(&t)
	}

	return tests, nil
}

// NewTest creates new Test instance from config
func newTest(c *Test) *test.Test {
	t := &test.Test{
		Name:    c.Name,
		Command: c.Command,
		Env:     make(map[string]string),
	}

	for _, kv := range c.Env {
		t.Env[kv.Name] = kv.Value
	}

	if c.Expect != nil {
		t.Status = c.Expect.Status
		t.Stdout = c.Expect.Stdout
		t.Stderr = c.Expect.Stderr
	}

	return t
}
