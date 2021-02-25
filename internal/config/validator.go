package config

import (
	"errors"
	"fmt"
	"strings"
)

type configMap map[string]interface{}
type configSeq []interface{}

type violation struct {
	path    string
	message string
}

type validator struct {
	paths      []string
	violations []violation
}

func newValidator() *validator {
	return &validator{
		paths:      []string{"$"},
		violations: make([]violation, 0),
	}
}

func (v *validator) pushPath(path string) {
	v.paths = append(v.paths, path)
}

func (v *validator) popPath() {
	if len(v.paths) < 2 {
		panic("pop empty validator.paths ")
	}
	v.paths = v.paths[:len(v.paths)-1]
}

func (v *validator) InPath(path string, f func()) {
	v.pushPath(path)
	defer v.popPath()
	f()
}

func (v *validator) InField(field string, f func()) {
	v.InPath("."+field, f)
}

func (v *validator) InIndex(index int, f func()) {
	v.InPath(fmt.Sprintf("[%d]", index), f)
}

func (v *validator) AddViolation(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	v.violations = append(v.violations, violation{path: strings.Join(v.paths, ""), message: message})
}

func (v *validator) MustBeMap(x interface{}) (configMap, bool) {
	if m, ok := x.(configMap); ok {
		return m, true
	}
	v.AddViolation("should be map, but is %T", x)
	return nil, false
}

func (v *validator) MustBeSeq(x interface{}) (configSeq, bool) {
	if s, ok := x.(configSeq); ok {
		return s, true
	}
	v.AddViolation("should be seq, but is %T", x)
	return nil, false
}

func (v *validator) MustBeString(x interface{}) (string, bool) {
	s, ok := x.(string)
	if !ok {
		v.AddViolation("should be string but is %T", x)
	}

	return s, ok
}

func (v *validator) mustHave(m configMap, key string) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		v.AddViolation("should have .%s", key)
	}
	return x, ok
}

func (v *validator) MustHaveSeq(m configMap, key string, f func(configSeq)) (configSeq, bool) {
	x, ok := v.mustHave(m, key)
	if !ok {
		return nil, false
	}

	var s configSeq
	v.InField(key, func() {
		s, ok = v.MustBeSeq(x)
		if ok {
			f(s)
		}
	})

	return s, ok
}

func (v *validator) MayHaveString(m configMap, key string) (string, bool, bool) {
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

func (v *validator) MustHaveString(m configMap, key string) (string, bool) {
	s, exists, ok := v.MayHaveString(m, key)
	return s, exists && ok
}

func (v *validator) Error() error {
	if len(v.violations) == 0 {
		return nil
	}

	b := strings.Builder{}
	for _, violation := range v.violations {
		b.WriteString(violation.path + ": " + violation.message + "\n")
	}

	return errors.New(b.String())
}
