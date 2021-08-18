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

package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
)

type NotMatcher struct {
	matcher matcher.StreamMatcher
}

func (m *NotMatcher) MatchStream(actual []byte) (bool, string, error) {
	matched, message, err := m.matcher.MatchStream(actual)

	return !matched, message, err
}

func ParseNotMatcher(v *spec.Validator, r *matcher.StreamMatcherRegistry, x interface{}) matcher.StreamMatcher {
	m := r.ParseMatcher(v, x)
	if m == nil {
		return nil
	}

	return &NotMatcher{matcher: m}
}