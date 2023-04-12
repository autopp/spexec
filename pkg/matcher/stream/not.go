// Copyright (C) 2021-2023	 Akira Tanimura (@autopp)
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
	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"
)

type NotMatcher struct {
	matcher model.StreamMatcher
}

func (m *NotMatcher) Match(actual []byte) (bool, string, error) {
	matched, message, err := m.matcher.Match(actual)

	return !matched, message, err
}

func ParseNotMatcher(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	m := r.ParseMatcher(v, x)
	if m == nil {
		return nil
	}

	return &NotMatcher{matcher: m}
}
