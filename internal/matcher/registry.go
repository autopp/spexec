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

type StatusMatcherRegistry struct {
	matchers map[string]StatusMatcherParser
}

func NewStatusMatcherRegistry() *StatusMatcherRegistry {
	return &StatusMatcherRegistry{matchers: make(map[string]StatusMatcherParser)}
}

func (r *StatusMatcherRegistry) Add(name string, p StatusMatcherParser) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = p
	return nil
}

func (r *StatusMatcherRegistry) ParseMatcher(v *spec.Validator, x interface{}) StatusMatcher {
	specifier, ok := x.(spec.Map)
	if !ok {
		v.AddViolation("matcher specifier should be a map with single key-value (got %T)", x)
		return nil
	}
	if len(specifier) != 1 {
		v.AddViolation("matcher specifier should be a map with single key-value (got map with %d key-value)", len(specifier))
		return nil
	}

	var name string
	var param interface{}
	for name, param = range specifier {
	}

	p, ok := r.matchers[name]
	if !ok {
		v.AddViolation("matcher for status %s is not defined", name)
		return nil
	}
	return p(v, r, param)
}

type StreamMatcherRegistry struct {
	matchers map[string]StreamMatcherParser
}

func NewStreamMatcherRegistry() *StreamMatcherRegistry {
	return &StreamMatcherRegistry{matchers: make(map[string]StreamMatcherParser)}
}

func (r *StreamMatcherRegistry) Add(name string, p StreamMatcherParser) error {
	_, ok := r.matchers[name]
	if ok {
		return errors.Errorf(errors.ErrInternalError, "matcher %s is already registered", name)
	}
	r.matchers[name] = p
	return nil
}

func (r *StreamMatcherRegistry) ParseMatcher(v *spec.Validator, fd int, x interface{}) StreamMatcher {
	specifier, ok := x.(spec.Map)
	if !ok {
		v.AddViolation("matcher specifier should be a map with single key-value (got %T)", x)
		return nil
	}
	if len(specifier) != 1 {
		v.AddViolation("matcher specifier should be a map with single key-value (got map with %d key-value)", len(specifier))
		return nil
	}

	var name string
	var param interface{}
	for name, param = range specifier {
	}

	p, ok := r.matchers[name]
	if !ok {
		v.AddViolation("matcher for status %s is not defined", name)
		return nil
	}
	return p(v, r, param)
}
