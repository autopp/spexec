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

package matcher

import (
	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/spec"
)

type statusMatcherParserEntry struct {
	parser       StatusMatcherParser
	hasDefault   bool
	defaultParam interface{}
}

type matcherParserEntry struct {
	parser       interface{}
	hasDefault   bool
	defaultParam interface{}
}

type matcherParserRegistry struct {
	matchers map[string]*matcherParserEntry
}

type StatusMatcherRegistry struct {
	registry *matcherParserRegistry
}

func newMatcherParserRegistry() *matcherParserRegistry {
	return &matcherParserRegistry{matchers: make(map[string]*matcherParserEntry)}
}

func (r *matcherParserRegistry) add(name string, p interface{}) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = &matcherParserEntry{
		parser:     p,
		hasDefault: false,
	}
	return nil
}

func (r *matcherParserRegistry) addWithDefault(name string, p interface{}, defaultParam interface{}) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = &matcherParserEntry{
		parser:       p,
		hasDefault:   true,
		defaultParam: defaultParam,
	}
	return nil
}

func (r *matcherParserRegistry) get(v *spec.Validator, x interface{}) (string, interface{}, interface{}) {
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
		v.AddViolation("matcher specifier should be a matcher name or a map with single key-value (got %s)", spec.Typeof(x))
		return "", nil, nil
	}

	p, ok := r.matchers[name]
	if !ok {
		v.AddViolation("matcher for status %s is not defined", name)
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

func NewStatusMatcherRegistry() *StatusMatcherRegistry {
	return &StatusMatcherRegistry{registry: newMatcherParserRegistry()}
}

func (r *StatusMatcherRegistry) Add(name string, p StatusMatcherParser) error {
	return r.registry.add(name, p)
}

func (r *StatusMatcherRegistry) AddWithDefault(name string, p StatusMatcherParser, defaultParam interface{}) error {
	return r.registry.addWithDefault(name, p, defaultParam)
}

func (r *StatusMatcherRegistry) ParseMatcher(v *spec.Validator, x interface{}) StatusMatcher {
	name, parser, param := r.registry.get(v, x)
	if parser == nil {
		return nil
	}
	var m StatusMatcher
	v.InField(name, func() {
		m = parser.(StatusMatcherParser)(v, r, param)
	})

	return m
}

type StreamMatcherRegistry struct {
	registry *matcherParserRegistry
}

func NewStreamMatcherRegistry() *StreamMatcherRegistry {
	return &StreamMatcherRegistry{registry: newMatcherParserRegistry()}
}

func (r *StreamMatcherRegistry) Add(name string, p StreamMatcherParser) error {
	return r.registry.add(name, p)
}

func (r *StreamMatcherRegistry) AddWithDefault(name string, p StreamMatcherParser, defaultParam interface{}) error {
	return r.registry.addWithDefault(name, p, defaultParam)
}

func (r *StreamMatcherRegistry) ParseMatcher(v *spec.Validator, x interface{}) StreamMatcher {
	name, parser, param := r.registry.get(v, x)
	if parser == nil {
		return nil
	}
	var m StreamMatcher
	v.InField(name, func() {
		m = parser.(StreamMatcherParser)(v, r, param)
	})

	return m
}
