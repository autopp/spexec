package config

import (
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

func (v *validator) inPath(path string, f func()) {
	v.pushPath(path)
	defer v.popPath()
	f()
}

func (v *validator) inField(field string, f func()) {
	v.inPath("."+field, f)
}

func (v *validator) inIndex(index int, f func()) {
	v.inPath(fmt.Sprintf("[%d]", index), f)
}

func (v *validator) addViolation(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	v.violations = append(v.violations, violation{path: strings.Join(v.paths, ""), message: message})
}

func (v *validator) toBeMap(x interface{}) (configMap, bool) {
	if m, ok := x.(configMap); ok {
		return m, true
	}
	v.addViolation("should be map, but is %T", x)
	return nil, false
}

func (v *validator) toBeSeq(x interface{}) (configSeq, bool) {
	if s, ok := x.(configSeq); ok {
		return s, true
	}
	v.addViolation("should be seq, but is %T", x)
	return nil, false
}
