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

package matcher

import (
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
)

type MatcherParser[T any] func(v *spec.Validator, r *matcherParserRegistry[T], x interface{}) model.Matcher[T]

type StatusMatcherParser = MatcherParser[model.StatusMatcher]
type StreamMatcherParser = MatcherParser[model.StreamMatcher]
