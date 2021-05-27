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

package status

import (
	"fmt"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
)

type SuccessMatcher struct {
	expected bool
}

func (m *SuccessMatcher) MatchStatus(actual int) (bool, string, error) {
	succeeded := actual == 0
	if succeeded == m.expected {
		return true, fmt.Sprintf("should not succeed, but success"), nil
	}

	return false, fmt.Sprintf("should succeed, but got %d", actual), nil
}

func ParseSuccessMatcher(v *spec.Validator, r *matcher.StatusMatcherRegistry, x interface{}) matcher.StatusMatcher {
	expected, ok := v.MustBeBool(x)
	if !ok {
		return nil
	}

	return &SuccessMatcher{expected: expected}
}
