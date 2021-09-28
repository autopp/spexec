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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
)

type EqJSONMatcher struct {
	expected       interface{}
	expectedString string
}

func (m *EqJSONMatcher) MatchStream(actual []byte) (bool, string, error) {
	d := json.NewDecoder(bytes.NewBuffer(actual))
	d.UseNumber()
	var actualBody interface{}
	if err := d.Decode(&actualBody); err != nil {
		return false, fmt.Sprintf("cannot recognize as json: %s", err), nil
	}

	if reflect.DeepEqual(actualBody, m.expected) {
		return true, fmt.Sprintf("should not be %s, but got it", m.expectedString), nil
	}

	return false, fmt.Sprintf("should be %s, but got %s", m.expectedString, string(actual)), nil
}

func ParseEqJSONMatcher(v *spec.Validator, r *matcher.StreamMatcherRegistry, x interface{}) matcher.StreamMatcher {
	expectedBytes, err := json.Marshal(x)
	if err != nil {
		v.AddViolation("parameter is not json value: %s", err)
		return nil
	}

	d := json.NewDecoder(bytes.NewReader(expectedBytes))
	d.UseNumber()
	var expected interface{}
	err = d.Decode(&expected)
	if err != nil {
		v.AddViolation("parameter is not json value: %s", err)
	}

	return &EqJSONMatcher{expected: expected, expectedString: string(expectedBytes)}
}
