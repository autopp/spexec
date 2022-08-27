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
	Expand(env *Env, v *Validator, value any) (any, bool, error)
}

type TemplateVar struct {
	name string
}

func NewTemplateVar(name string) *TemplateVar {
	return &TemplateVar{name}
}

func (tv *TemplateVar) Expand(env *Env, v *Validator, value any) (any, bool, error) {
	value, ok := env.Lookup(tv.name)
	if !ok {
		v.InField("$"+tv.name, func() {
			v.AddViolation("is not defined")
		})
		return nil, false, nil
	}

	return value, true, nil
}

type TemplateFieldRef struct {
	field string
	next  TemplateRef
}

func NewTemplateFieldRef(field string, next TemplateRef) *TemplateFieldRef {
	return &TemplateFieldRef{
		field: field,
		next:  next,
	}
}

func (tf *TemplateFieldRef) Expand(env *Env, v *Validator, value any) (any, bool, error) {
	m, ok := v.MustBeMap(value)
	if !ok {
		return nil, false, unexpectedTemplateStructureError(v)
	}

	field, ok := v.MustHave(m, tf.field)
	if !ok {
		return nil, false, unexpectedTemplateStructureError(v)
	}

	var expanded any
	var err error
	v.InField(tf.field, func() {
		expanded, ok, err = tf.next.Expand(env, v, field)
	})

	if err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, false, nil
	}

	m[tf.field] = expanded

	return m, true, nil
}

type TemplateIndexRef struct {
	index int
	next  TemplateRef
}

func NewTemplateIndexRef(index int, next TemplateRef) *TemplateIndexRef {
	return &TemplateIndexRef{
		index: index,
		next:  next,
	}
}

func (ti *TemplateIndexRef) Expand(env *Env, v *Validator, value any) (any, bool, error) {
	s, ok := v.MustBeSeq(value)
	if !ok {
		return nil, false, unexpectedTemplateStructureError(v)
	}

	if ti.index >= len(s) {
		v.AddViolation("should have %d items", ti.index+1)
		return nil, false, unexpectedTemplateStructureError(v)
	}

	var expanded any
	var err error
	v.InIndex(ti.index, func() {
		expanded, ok, err = ti.next.Expand(env, v, s[ti.index])
	})

	if err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, false, nil
	}

	s[ti.index] = expanded

	return s, true, nil
}

type TemplateValue struct {
	value any
	refs  []TemplateRef
}

func NewTemplateValue(value any, refs []TemplateRef) *TemplateValue {
	return &TemplateValue{
		value: value,
		refs:  refs,
	}
}

func (tv *TemplateValue) Expand(env *Env, v *Validator) (any, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(&tv.value); err != nil {
		return nil, err
	}

	var copied any
	if err := gob.NewDecoder(buf).Decode(&copied); err != nil {
		return nil, err
	}

	for _, ref := range tv.refs {
		var ok bool
		var err error
		copied, ok, err = ref.Expand(env, v, copied)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, v.Error()
		}
	}

	return copied, nil
}

type Templatable[T any] struct {
	tv    *TemplateValue
	value T
}

func NewTemplatableFromValue[T any](value T) *Templatable[T] {
	return &Templatable[T]{value: value}
}

func NewTemplatableFromTemplateValue[T any](tv *TemplateValue) *Templatable[T] {
	return &Templatable[T]{tv: tv}
}

func NewTemplatableFromVariable[T any](name string) *Templatable[T] {
	return NewTemplatableFromTemplateValue[T](NewTemplateValue(Map{"$": name}, []TemplateRef{NewTemplateVar(name)}))
}

func (t *Templatable[T]) Expand(env *Env, v *Validator) (T, error) {
	if t.tv == nil {
		return t.value, nil
	}

	expanded, err := t.tv.Expand(env, v)

	if err != nil {
		var defaultV T
		return defaultV, err
	}

	x, ok := expanded.(T)
	if !ok {
		var defaultV T
		return defaultV, errors.Errorf(errors.ErrInvalidSpec, "expect %T, but got %T", defaultV, expanded)
	}

	return x, nil
}

func init() {
	gob.Register([]any{})
	gob.Register(map[string]any{})
}

func unexpectedTemplateStructureError(v *Validator) error {
	return errors.Errorf(errors.ErrInternalError, "unexpected template structure: %s", v.LastViolation())
}
