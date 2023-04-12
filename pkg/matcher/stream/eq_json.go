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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"
	"github.com/autopp/spexec/pkg/util"
)

type EqJSONMatcher struct {
	expected       any
	expectedString string
}

func (m *EqJSONMatcher) Match(actual []byte) (bool, string, error) {
	var actualBody any
	if err := util.DecodeJSON(bytes.NewReader(actual), &actualBody); err != nil {
		return false, fmt.Sprintf("cannot recognize as json: %s", err), nil
	}

	if reflect.DeepEqual(actualBody, m.expected) {
		return true, fmt.Sprintf("should not be %s, but got it", m.expectedString), nil
	}

	return false, fmt.Sprintf("should be %s, but got %s", m.expectedString, string(actual)), nil
}

func ParseEqJSONMatcher(v *model.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	expectedBytes, err := json.Marshal(x)
	if err != nil {
		v.AddViolation("parameter is not json value: %s", err)
		return nil
	}

	var expected any
	err = util.DecodeJSON(bytes.NewReader(expectedBytes), &expected)
	if err != nil {
		v.AddViolation("parameter is not json value: %s", err)
	}

	return &EqJSONMatcher{expected: expected, expectedString: string(expectedBytes)}
}
