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

package spec

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/util"
)

var envVarNamePattern = regexp.MustCompile(`^[a-zA-Z_]\w*$`)

type violation struct {
	path    string
	message string
}

type Validator struct {
	filename   string
	paths      []string
	violations []violation
}

func NewValidator(filename string) (*Validator, error) {
	if len(filename) != 0 {
		absFilename, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		filename = absFilename
	}
	return &Validator{
		filename:   filename,
		paths:      []string{"$"},
		violations: make([]violation, 0),
	}, nil
}

func (v *Validator) GetDir() string {
	return filepath.Dir(v.filename)
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

func (v *Validator) MustBeMap(x interface{}) (Map, bool) {
	if m, ok := x.(Map); ok {
		return m, true
	}
	v.AddViolation("should be map, but is %s", Typeof(x))
	return nil, false
}

func (v *Validator) MustBeSeq(x interface{}) (Seq, bool) {
	if s, ok := x.(Seq); ok {
		return s, true
	}
	v.AddViolation("should be seq, but is %s", Typeof(x))
	return nil, false
}

func (v *Validator) MustBeString(x interface{}) (string, bool) {
	s, ok := x.(string)
	if !ok {
		v.AddViolation("should be string, but is %s", Typeof(x))
	}

	return s, ok
}

func (v *Validator) MustBeInt(x interface{}) (int, bool) {
	switch n := x.(type) {
	case int:
		return n, true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			v.AddViolation("should be int, but is %s", Typeof(x))
		}
		return int(i), err == nil
	default:
		v.AddViolation("should be int, but is %s", Typeof(x))
		return 0, false
	}
}

func (v *Validator) MustBeBool(x interface{}) (bool, bool) {
	b, ok := x.(bool)
	if !ok {
		v.AddViolation("should be bool, but is %s", Typeof(x))
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
		v.AddViolation("should be positive integer or duration string, but is %s", Typeof(x))
		return 0, false
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		v.AddViolation("should be positive integer or duration string, but cannot parse (%s)", err)
		return 0, false
	}

	return d, true
}

func (v *Validator) mustHave(m Map, key string) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		v.AddViolation("should have .%s", key)
	}
	return x, ok
}

func (v *Validator) MayHave(m Map, key string, f func(interface{})) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false
	}

	v.InField(key, func() {
		f(x)
	})

	return x, true
}

func (v *Validator) MayHaveMap(m Map, key string, f func(Map)) (Map, bool, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false, true
	}

	var inner Map
	v.InField(key, func() {
		inner, ok = v.MustBeMap(x)
		if ok {
			f(inner)
		}
	})

	return inner, ok, ok
}

func (v *Validator) MayHaveSeq(m Map, key string, f func(Seq)) (Seq, bool, bool) {
	x, ok := m[key]
	if !ok {
		return nil, false, true
	}

	var s Seq
	v.InField(key, func() {
		s, ok = v.MustBeSeq(x)
		if ok {
			f(s)
		}
	})

	return s, ok, ok
}

func (v *Validator) MustHaveSeq(m Map, key string, f func(Seq)) (Seq, bool) {
	s, exists, ok := v.MayHaveSeq(m, key, f)

	if !exists && ok {
		v.AddViolation("should have .%s as seq", key)
	}

	return s, exists && ok
}

func (v *Validator) ForInSeq(s Seq, f func(i int, x interface{}) bool) bool {
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

func (v *Validator) MayHaveString(m Map, key string) (string, bool, bool) {
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

func (v *Validator) MustHaveString(m Map, key string) (string, bool) {
	s, exists, ok := v.MayHaveString(m, key)

	if !exists && ok {
		v.AddViolation("should have .%s as string", key)
	}

	return s, exists && ok
}

func (v *Validator) MayHaveInt(m Map, key string) (int, bool, bool) {
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

func (v *Validator) MayHaveDuration(m Map, key string) (time.Duration, bool, bool) {
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

func (v *Validator) MayHaveEnvSeq(m Map, key string) ([]util.StringVar, bool, bool) {
	var ret []util.StringVar
	ok := true
	_, _, isSeq := v.MayHaveSeq(m, "env", func(env Seq) {
		ret = []util.StringVar{}
		v.ForInSeq(env, func(i int, x interface{}) bool {
			var envVar Map
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

func (v *Validator) MayHaveCommand(m Map, key string) ([]string, bool, bool) {
	var ret []string
	ok := true
	_, _, isSeq := v.MayHaveSeq(m, key, func(command Seq) {
		ret = make([]string, len(command))
		v.ForInSeq(command, func(i int, x interface{}) bool {
			var c string
			c, ok = v.MustBeString(x)
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

func (v *Validator) MustHaveCommand(m Map, key string) ([]string, bool) {
	c, exists, ok := v.MayHaveCommand(m, key)

	if !exists && ok {
		v.AddViolation("should have .%s as command seq", key)
	}

	return c, exists && ok
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

func Typeof(x interface{}) string {
	if x == nil {
		return "nil"
	}

	if _, ok := x.(int); ok {
		return "int"
	}

	if i, ok := x.(json.Number); ok {
		if _, err := i.Int64(); err == nil {
			return "int"
		}
	}

	if _, ok := x.(bool); ok {
		return "bool"
	}

	if _, ok := x.(string); ok {
		return "string"
	}

	if _, ok := x.(Seq); ok {
		return "seq"
	}

	if _, ok := x.(Map); ok {
		return "map"
	}

	return fmt.Sprintf("%T", x)
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
