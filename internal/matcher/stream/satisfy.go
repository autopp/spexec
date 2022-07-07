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
	"time"

	"github.com/autopp/spexec/internal/exec"
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
	"github.com/autopp/spexec/internal/util"
)

type SatisfyMatcher struct {
	Command []string
	Dir     string
	Cleanup func() []error
	Env     []util.StringVar
	Timeout time.Duration
}

func (m *SatisfyMatcher) Match(actual []byte) (bool, string, error) {
	// FIXME: error handling
	defer m.Cleanup()

	e, err := exec.New(m.Command, m.Dir, actual, m.Env, exec.WithTimeout(m.Timeout))
	if err != nil {
		return false, "", err
	}

	er := e.Run()
	if er.Err != nil || er.Signal != nil || er.Status != 0 {
		return false, "should make the given command succeed", nil
	}
	return true, "should make the given command fail", nil
}

func ParseSatisfyMatcher(v *spec.Validator, r *matcher.StreamMatcherRegistry, x any) model.StreamMatcher {
	p, ok := v.MustBeMap(x)
	if !ok {
		return nil
	}

	m := &SatisfyMatcher{}
	command, ok := v.MustHaveCommand(p, "command")
	if !ok {
		return nil
	}

	evaled, cleanup, err, _ := model.EvalStringExprs(command)
	if err != nil {
		v.InField("command", func() {
			v.AddViolation("error occured at parsing command: %s", err)
		})
		return nil
	}

	m.Command = evaled
	m.Cleanup = cleanup

	m.Dir = v.GetDir()

	m.Env, _, ok = v.MayHaveEnvSeq(p, "env")
	if !ok {
		return nil
	}

	var exist bool
	m.Timeout, exist, ok = v.MayHaveDuration(p, "timeout")
	if !ok {
		return nil
	}
	if !exist {
		m.Timeout = 5 * time.Second
	}

	return m
}
