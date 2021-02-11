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

	"gopkg.in/yaml.v2"
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

func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	yaml.Unmarshal(b, c)

	return c, nil
}
