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

package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"
	"gopkg.in/yaml.v3"
)

var envVarNamePattern = regexp.MustCompile(`^[a-zA-Z_]\w*$`)

type violation struct {
	path    string
	message string
}

type Validator struct {
	Filename   string
	dir        string
	paths      []string
	violations []violation
}

func NewValidator(filename string) (*Validator, error) {
	var dir string
	if len(filename) == 0 {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	} else {
		absFilename, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		filename = absFilename
		dir = filepath.Dir(filename)
	}

	return &Validator{
		Filename:   filename,
		dir:        dir,
		paths:      []string{"$"},
		violations: make([]violation, 0),
	}, nil
}

func (v *Validator) GetDir() string {
	return v.dir
}

func (v *Validator) pushPath(path string) {
	v.paths = append(v.paths, path)
}

func (v *Validator) popPath() {
	if len(v.paths) < 2 {
		panic("pop empty validator.paths ")
	}
	v.paths = v.paths[:len(v.paths)-1]
}

func (v *Validator) InPath(path string, f func()) {
	v.pushPath(path)
	defer v.popPath()
	f()
}

func (v *Validator) InField(field string, f func()) {
	v.InPath("."+field, f)
}

func (v *Validator) InIndex(index int, f func()) {
	v.InPath(fmt.Sprintf("[%d]", index), f)
}

func (v *Validator) AddViolation(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	v.violations = append(v.violations, violation{path: strings.Join(v.paths, ""), message: message})
}

func (v *Validator) MayBeMap(x interface{}) (model.Map, bool) {
	m, ok := x.(model.Map)
	return m, ok
}

func (v *Validator) MustBeMap(x interface{}) (model.Map, bool) {
	if m, ok := v.MayBeMap(x); ok {
		return m, true
	}
	v.AddViolation("should be map, but is %s", model.TypeNameOf(x))
	return nil, false
}

func (v *Validator) MustBeSeq(x interface{}) (model.Seq, bool) {
	if s, ok := x.(model.Seq); ok {
		return s, true
	}
	v.AddViolation("should be seq, but is %s", model.TypeNameOf(x))
	return nil, false
}

func (v *Validator) MayBeString(x interface{}) (string, bool) {
	s, ok := x.(string)
	return s, ok
}

func (v *Validator) MustBeString(x interface{}) (string, bool) {
	s, ok := v.MayBeString(x)
	if !ok {
		v.AddViolation("should be string, but is %s", model.TypeNameOf(x))
	}

	return s, ok
}

func (v *Validator) MayBeQualified(x interface{}) (string, interface{}, bool) {
	qv, ok := v.MayBeMap(x)
	if !ok {
		return "", nil, false
	}

	if len(qv) != 1 {
		return "", nil, false
	}

	for q, v := range qv {
		return q, v, true
	}

	panic("UNREACHABLE CODE")
}

var variablePattern = regexp.MustCompile(`^[_a-zA-Z]\w*$`)

func (v *Validator) MayBeVariable(x interface{}) (string, bool) {
	q, value, ok := v.MayBeQualified(x)
	if !ok || q != "$" {
		return "", false
	}

	name, ok := v.MayBeString(value)
	if !ok || !variablePattern.MatchString(name) {
		return "", false
	}

	return name, true
}

func (v *Validator) MustBeStringExpr(x interface{}) (model.StringExpr, bool) {
	if s, ok := v.MayBeString(x); ok {
		return model.NewLiteralStringExpr(s), true
	}

	m, ok := v.MayBeMap(x)
	if !ok {
		v.AddViolation("should be string or map, but is %s", model.TypeNameOf(x))
		return nil, false
	}

	t, ok := v.MustHaveString(m, "type")
	if !ok {
		return nil, false
	}

	switch t {
	case "env":
		name, ok := v.MustHaveString(m, "name")
		if !ok {
			return nil, false
		}
		return model.NewEnvStringExpr(name), true
	case "file":
		format, exists, ok := v.MayHaveString(m, "format")
		if !ok {
			return nil, false
		}
		if !exists {
			format = "raw"
		}

		switch format {
		case "raw":
			value, ok := v.MustHaveString(m, "value")
			if !ok {
				return nil, false
			}
			return model.NewFileStringExpr("", value), true
		case "yaml":
			value, ok := v.MustHave(m, "value")
			if !ok {
				return nil, false
			}
			marshaled, err := yaml.Marshal(value)
			if err != nil {
				v.InField("value", func() {
					v.AddViolation("cannot encode to a YAML string: %s", err)
				})
				return nil, false
			}

			return model.NewFileStringExpr("*.yaml", string(marshaled)), true
		default:
			v.InField("format", func() {
				v.AddViolation(`should be a "raw" or "yaml", but is %q`, format)
			})

			return nil, false
		}
	default:
		v.InField("type", func() {
			v.AddViolation("unknown type %q", t)
		})
		return nil, false
	}
}

func (v *Validator) MustBeInt(x interface{}) (int, bool) {
	switch n := x.(type) {
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			v.AddViolation("should be int, but is invalid json.Number: %s", err)
		}
		return int(i), err == nil
	default:
		v.AddViolation("should be int, but is %s", model.TypeNameOf(x))
		return 0, false
	}
}

func (v *Validator) MustBeBool(x interface{}) (bool, bool) {
	b, ok := x.(bool)
	if !ok {
		v.AddViolation("should be bool, but is %s", model.TypeNameOf(x))
	}

	return b, ok
}

func (v *Validator) MustBeDuration(x interface{}) (time.Duration, bool) {
	n, ok := toInt(x)
	if ok {
		return time.Duration(n) * time.Second, true
	}

	s, ok := x.(string)
	if !ok {
		v.AddViolation("should be positive integer or duration string, but is %s", model.TypeNameOf(x))
		return 0, false
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		v.AddViolation("should be positive integer or duration string, but cannot parse (%s)", err)
		return 0, false
	}

	return d, true
}

func (v *Validator) MustHave(m model.Map, key string) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		v.AddViolation("should have .%s", key)
	}
	return x, ok
}

func (v *Validator) MayHave(m model.Map, key string, f func(interface{})) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false
	}

	v.InField(key, func() {
		f(x)
	})

	return x, true
}

func (v *Validator) MayHaveMap(m model.Map, key string, f func(model.Map)) (model.Map, bool, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false, true
	}

	var inner model.Map
	v.InField(key, func() {
		inner, ok = v.MustBeMap(x)
		if ok {
			f(inner)
		}
	})

	return inner, ok, ok
}

func (v *Validator) MayHaveSeq(m model.Map, key string, f func(model.Seq)) (model.Seq, bool, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false, true
	}

	var s model.Seq
	v.InField(key, func() {
		s, ok = v.MustBeSeq(x)
		if ok {
			f(s)
		}
	})

	return s, ok, ok
}

func (v *Validator) MustHaveSeq(m model.Map, key string, f func(model.Seq)) (model.Seq, bool) {
	s, exists, ok := v.MayHaveSeq(m, key, f)

	if !exists && ok {
		v.AddViolation("should have .%s as seq", key)
	}

	return s, exists && ok
}

func (v *Validator) ForInSeq(s model.Seq, f func(i int, x interface{}) bool) bool {
	ok := true
	for i, x := range s {
		v.InIndex(i, func() {
			ok = f(i, x)
		})

		if !ok {
			break
		}
	}

	return ok
}

func (v *Validator) MayHaveString(m model.Map, key string) (string, bool, bool) {
	x, ok := m[key]
	if !ok {
		return "", false, true
	}

	var s string
	v.InField(key, func() {
		s, ok = v.MustBeString(x)
	})

	return s, ok, ok
}

func (v *Validator) MustHaveString(m model.Map, key string) (string, bool) {
	s, exists, ok := v.MayHaveString(m, key)

	if !exists && ok {
		v.AddViolation("should have .%s as string", key)
	}

	return s, exists && ok
}

func (v *Validator) MayHaveInt(m model.Map, key string) (int, bool, bool) {
	x, ok := m[key]
	if !ok {
		return 0, false, true
	}

	var i int
	v.InField(key, func() {
		i, ok = v.MustBeInt(x)
	})

	return i, ok, ok
}

func (v *Validator) MayHaveBool(m model.Map, key string) (bool, bool, bool) {
	x, ok := m[key]
	if !ok {
		return false, false, true
	}

	var b bool
	v.InField(key, func() {
		b, ok = v.MustBeBool(x)
	})

	return b, ok, ok
}

func (v *Validator) MayHaveDuration(m model.Map, key string) (time.Duration, bool, bool) {
	x, ok := m[key]
	if !ok {
		return 0, false, true
	}

	var d time.Duration
	v.InField(key, func() {
		d, ok = v.MustBeDuration(x)
	})

	return d, ok, ok
}

func (v *Validator) MayHaveEnvSeq(m model.Map, key string) ([]util.StringVar, bool, bool) {
	var ret []util.StringVar
	ok := true
	_, _, isSeq := v.MayHaveSeq(m, "env", func(env model.Seq) {
		ret = []util.StringVar{}
		v.ForInSeq(env, func(i int, x interface{}) bool {
			var envVar model.Map
			envVar, ok = v.MustBeMap(x)
			if !ok {
				return false
			}
			name, nameOk := v.MustHaveString(envVar, "name")
			value, valueOk := v.MustHaveString(envVar, "value")
			if !nameOk || !valueOk {
				ok = false
				return false
			}

			if !envVarNamePattern.MatchString(name) {
				v.InField("name", func() {
					v.AddViolation("environment variable name shoud be match to /%s/", envVarNamePattern.String())
				})
				ok = false
				return false
			}
			ret = append(ret, util.StringVar{Name: name, Value: value})
			return true
		})
	})

	if !isSeq || !ok {
		return nil, false, false
	}

	return ret, ret != nil, ok
}

func (v *Validator) MayHaveCommand(m model.Map, key string) ([]model.StringExpr, bool, bool) {
	var ret []model.StringExpr
	ok := true
	_, _, isSeq := v.MayHaveSeq(m, key, func(command model.Seq) {
		ret = make([]model.StringExpr, len(command))
		v.ForInSeq(command, func(i int, x interface{}) bool {
			var c model.StringExpr
			c, ok = v.MustBeStringExpr(x)
			ret[i] = c
			return ok
		})

		if len(ret) == 0 {
			v.AddViolation("should have one ore more elements")
			ok = false
		}
	})

	if !isSeq || !ok {
		return nil, false, false
	}

	return ret, ret != nil, ok
}

func (v *Validator) MustHaveCommand(m model.Map, key string) ([]model.StringExpr, bool) {
	c, exists, ok := v.MayHaveCommand(m, key)

	if !exists && ok {
		v.AddViolation("should have .%s as command seq", key)
	}

	return c, exists && ok
}

func (v *Validator) MustContainOnly(m model.Map, keys ...string) bool {
	dict := map[string]struct{}{}
	for _, key := range keys {
		dict[key] = struct{}{}
	}

	ok := true
	for field := range m {
		if _, allowed := dict[field]; !allowed {
			v.AddViolation("field .%s is not expected", field)
			ok = false
		}
	}

	return ok
}

func (v *Validator) Error() error {
	if len(v.violations) == 0 {
		return nil
	}

	messages := make([]string, len(v.violations))
	for i, violation := range v.violations {
		messages[i] = violation.path + ": " + violation.message
	}

	return errors.New(errors.ErrInvalidSpec, strings.Join(messages, "\n"))
}

func toInt(x interface{}) (int, bool) {
	switch n := x.(type) {
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}
