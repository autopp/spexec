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

import (
	"bytes"
	"encoding/gob"

	"github.com/autopp/spexec/internal/errors"
)

type TemplateRef interface {
	Expand(value interface{}, env *Env) (interface{}, error)
}

type TemplateVar struct {
	name string
}

func (tv *TemplateVar) Expand(value interface{}, env *Env) (interface{}, error) {
	value, ok := env.Lookup(tv.name)
	if !ok {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "undefined var: %s", tv.name)
	}

	return value, nil
}

type TemplateFieldRef struct {
	field string
	next  TemplateRef
}

func (tf *TemplateFieldRef) Expand(value interface{}, env *Env) (interface{}, error) {
	m, ok := value.(Map)

	if !ok {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "expect to be map, but got %s", TypeNameOf(value))
	}

	v, ok := m[tf.field]
	if !ok {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "expect to contains .%s", tf.field)
	}

	expanded, err := tf.next.Expand(v, env)
	if err != nil {
		return nil, err
	}

	m[tf.field] = expanded

	return m, nil
}

type TemplateIndexRef struct {
	index int
	next  TemplateRef
}

func (ti *TemplateIndexRef) Expand(value interface{}, env *Env) (interface{}, error) {
	s, ok := value.(Seq)

	if !ok {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "expect to be seq, but got %s", TypeNameOf(value))
	}

	if ti.index >= len(s) {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "expect to have %d items", ti.index)
	}

	expanded, err := ti.next.Expand(s[ti.index], env)
	if err != nil {
		return nil, err
	}

	s[ti.index] = expanded

	return s, nil
}

type TemplateValue struct {
	refs  []TemplateRef
	value interface{}
}

func (tv *TemplateValue) Expand(env *Env) (interface{}, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(tv.value); err != nil {
		return nil, err
	}

	var copied interface{}
	if err := gob.NewDecoder(buf).Decode(&copied); err != nil {
		return nil, err
	}

	for _, ref := range tv.refs {
		var err error
		copied, err = ref.Expand(copied, env)
		if err != nil {
			return nil, err
		}
	}

	return copied, nil
}

func init() {
	gob.Register(Seq{})
	gob.Register(Map{})
}
