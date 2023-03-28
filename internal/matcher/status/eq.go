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

package status

import (
	"fmt"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
)

type EqMatcher struct {
	expected int
}

func (m *EqMatcher) Match(actual int) (bool, string, error) {
	if actual == m.expected {
		return true, fmt.Sprintf("should not be %d, but got it", m.expected), nil
	}

	return false, fmt.Sprintf("should be %d, but got %d", m.expected, actual), nil
}

func ParseEqMatcher(v *model.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
	expected, ok := v.MustBeInt(x)
	if !ok {
		return nil
	}

	if expected < 0 {
		v.AddViolation("should be positive integer")
		return nil
	}

	return &EqMatcher{expected: expected}
}
