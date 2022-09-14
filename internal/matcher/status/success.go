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

package status

import (
	"fmt"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
)

type SuccessMatcher struct {
	expected bool
}

func (m *SuccessMatcher) Match(actual int) (bool, string, error) {
	succeeded := actual == 0
	unexpectedSuccessFormat := "should not succeed, but succeeded (status is %d)"
	unexpectedFailureFormat := "should succeed, but not succeeded (status is %d)"
	if succeeded == m.expected {
		if m.expected {
			return true, fmt.Sprintf(unexpectedSuccessFormat, actual), nil
		}
		return true, fmt.Sprintf(unexpectedFailureFormat, actual), nil
	}

	if m.expected {
		return false, fmt.Sprintf(unexpectedFailureFormat, actual), nil
	}

	return false, fmt.Sprintf(unexpectedSuccessFormat, actual), nil
}

func ParseSuccessMatcher(v *model.Validator, r *matcher.StatusMatcherRegistry, x any) model.StatusMatcher {
	expected, ok := v.MustBeBool(x)
	if !ok {
		return nil
	}

	return &SuccessMatcher{expected: expected}
}
