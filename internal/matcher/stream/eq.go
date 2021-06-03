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
	"bytes"
	"fmt"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
)

type EqMatcher struct {
	expected string
}

func (m *EqMatcher) MatchStream(actual []byte) (bool, string, error) {
	if bytes.Equal(actual, []byte(m.expected)) {
		return true, fmt.Sprintf("should not be %q, but got it", m.expected), nil
	}

	return false, fmt.Sprintf("should be %q, but got %q", m.expected, string(actual)), nil
}

func ParseEqMatcher(v *spec.Validator, r *matcher.StatusMatcherRegistry, x interface{}) matcher.StreamMatcher {
	expected, ok := v.MustBeString(x)
	if !ok {
		return nil
	}

	return &EqMatcher{expected: expected}
}
