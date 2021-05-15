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

import "github.com/autopp/spexec/internal/errors"

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

func (r *StatusMatcherRegistry) Get(name string) (StatusMatcherParser, error) {
	p, ok := r.matchers[name]

	if !ok {
		return nil, errors.Errorf(errors.ErrInvalidSpec, "matcher %s is not defined", name)
	}

	return p, nil
}

type StreamMatcherRegistry struct {
	matchers map[string]StreamMatcherParser
}

func NewStreamMatcherRegistry() *StreamMatcherRegistry {
	return &StreamMatcherRegistry{matchers: make(map[string]StreamMatcherParser)}
}

func (r *StreamMatcherRegistry) Add(name string, p StreamMatcherParser) error {
	panic("")
}

func (r *StreamMatcherRegistry) Get(name string) (StreamMatcherParser, error) {
	panic("")
}
