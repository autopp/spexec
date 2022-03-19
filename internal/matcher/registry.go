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

package matcher

import (
	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
)

type matcherParserEntry[T any] struct {
	parser       MatcherParser[T]
	hasDefault   bool
	defaultParam interface{}
}

type matcherParserRegistry[T any] struct {
	target   string
	matchers map[string]*matcherParserEntry[T]
}

func newMatcherParserRegistry[T any](target string) *matcherParserRegistry[T] {
	return &matcherParserRegistry[T]{target: target, matchers: make(map[string]*matcherParserEntry[T])}
}

func (r *matcherParserRegistry[T]) Add(name string, p MatcherParser[T]) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = &matcherParserEntry[T]{
		parser:     p,
		hasDefault: false,
	}
	return nil
}

func (r *matcherParserRegistry[T]) AddWithDefault(name string, p MatcherParser[T], defaultParam interface{}) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = &matcherParserEntry[T]{
		parser:       p,
		hasDefault:   true,
		defaultParam: defaultParam,
	}
	return nil
}

func (r *matcherParserRegistry[T]) get(v *spec.Validator, x interface{}) (string, MatcherParser[T], interface{}) {
	var name string
	var param interface{}
	withParam := false

	switch specifier := x.(type) {
	case string:
		name = specifier
	case spec.Map:
		if len(specifier) != 1 {
			v.AddViolation("matcher specifier should be a matcher name or a map with single key-value (got map with %d key-value)", len(specifier))
			return "", nil, nil
		}

		for name, param = range specifier {
		}
		withParam = true
	default:
		v.AddViolation("matcher specifier should be a matcher name or a map with single key-value (got %s)", spec.TypeNameOf(x))
		return "", nil, nil
	}

	p, ok := r.matchers[name]
	if !ok {
		v.AddViolation("matcher for %s %s is not defined", r.target, name)
		return "", nil, nil
	}

	if !withParam {
		if !p.hasDefault {
			v.InField(name, func() {
				v.AddViolation("parameter is required")
			})
			return "", nil, nil
		}
		param = p.defaultParam
	}

	return name, p.parser, param
}

func (r *matcherParserRegistry[T]) ParseMatcher(v *spec.Validator, x interface{}) model.Matcher[T] {
	name, parser, param := r.get(v, x)
	if parser == nil {
		return nil
	}
	var m model.Matcher[T]
	v.InField(name, func() {
		m = parser(v, r, param)
	})

	return m
}

func (r *matcherParserRegistry[T]) ParseMatchers(v *spec.Validator, x interface{}) []model.Matcher[T] {
	params, ok := v.MustBeSeq(x)
	if !ok {
		return nil
	}

	matchers := make([]model.Matcher[T], len(params))
	ok = v.ForInSeq(params, func(i int, param interface{}) bool {
		m := r.ParseMatcher(v, param)
		if m == nil {
			return false
		}
		matchers[i] = m
		return true
	})

	if !ok {
		return nil
	}
	return matchers
}

type StatusMatcherRegistry = matcherParserRegistry[int]

func NewStatusMatcherRegistry() *StatusMatcherRegistry {
	return newMatcherParserRegistry[int]("status")
}

type StreamMatcherRegistry = matcherParserRegistry[[]byte]

func NewStreamMatcherRegistry() *StreamMatcherRegistry {
	return newMatcherParserRegistry[[]byte]("stream")
}
