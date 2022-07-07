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

package stream

import (
	"fmt"
	"strings"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
)

type AnyMatcher struct {
	matchers []model.StreamMatcher
}

func (m *AnyMatcher) Match(actual []byte) (bool, string, error) {
	messages := make([]string, 0)
	for _, inner := range m.matchers {
		matched, message, err := inner.Match(actual)
		if err != nil {
			return false, "", err
		}
		if matched {
			return true, message, nil
		}
		messages = append(messages, "["+message+"]")
	}

	return false, fmt.Sprintf("should satisfy any of %s", strings.Join(messages, ", ")), nil
}

func ParseAnyMatcher(v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	matchers := r.ParseMatchers(v, x)
	if matchers == nil {
		return nil
	}

	return &AnyMatcher{matchers: matchers}
}
