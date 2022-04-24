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

package model

type Env struct {
	vars map[string]string
	prev *Env
}

func NewEnv(prev *Env) *Env {
	env := &Env{vars: make(map[string]string), prev: prev}
	return env
}

func (e *Env) Define(name string, value string) bool {
	_, defined := e.vars[name]
	if defined {
		return false
	}

	e.vars[name] = value

	return true
}

func (e *Env) Lookup(name string) (string, bool) {
	if e == nil {
		return "", false
	}

	if v, ok := e.vars[name]; ok {
		return v, true
	}

	return e.prev.Lookup(name)
}