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

import "github.com/autopp/spexec/internal/matcher"

func NewStreamMatcherRegistryWithBuiltins() *matcher.StreamMatcherRegistry {
	r := matcher.NewStreamMatcherRegistry()
	r.Add("eq", ParseEqMatcher)
	r.Add("beEmpty", ParseBeEmptyMatcher)
	r.Add("eqJSON", ParseEqJSONMatcher)
	r.Add("contain", ParseContainMatcher)
	r.Add("not", ParseNotMatcher)
	r.Add("any", ParseAnyMatcher)
	return r
}
