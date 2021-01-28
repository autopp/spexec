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

package test

import "github.com/autopp/spexec/internal/config"

type Test struct {
	Command []string
	Status  *int
	Stdout  *string
	Stderr  *string
}

// NewTest creates new Test instance from config
func NewTest(c *config.Test) *Test {
	t := &Test{
		Command: c.Command,
	}

	if c.Expect != nil {
		t.Status = c.Expect.Status
		t.Stdout = c.Expect.Stdout
		t.Stderr = c.Expect.Stderr
	}

	return t
}

// ToExec converts self to Exec
func (t *Test) ToExec() *Exec {
	return &Exec{
		Command: t.Command,
	}
}
