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
	"fmt"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"
)

type BeEmptyMatcher struct {
	expected bool
}

func (m *BeEmptyMatcher) Match(actual []byte) (bool, string, error) {
	empty := len(actual) == 0
	unexpectedEmptyFormat := "should not be empty, but is empty"
	unexpectedNotEmptyFormat := "should be empty, but is not:\n----------------------------------------\n%s\n----------------------------------------"
	if empty == m.expected {
		if m.expected {
			return true, fmt.Sprintf(unexpectedEmptyFormat), nil
		}
		return true, fmt.Sprintf(unexpectedNotEmptyFormat, string(actual)), nil
	}

	if m.expected {
		return false, fmt.Sprintf(unexpectedNotEmptyFormat, string(actual)), nil
	}

	return false, fmt.Sprintf(unexpectedEmptyFormat), nil
}

func ParseBeEmptyMatcher(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	expected, ok := v.MustBeBool(x)
	if !ok {
		return nil
	}

	return &BeEmptyMatcher{expected: expected}
}
