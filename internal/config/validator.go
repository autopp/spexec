package config

import "fmt"

type validator struct {
	paths  []string
	errors []string
}

func newValidator() *validator {
	return &validator{
		paths:  make([]string, 0),
		errors: make([]string, 0),
	}
}

func (v *validator) pushPath(path string) {
	v.paths = append(v.paths, path)
}

func (v *validator) popPath() {
	if len(v.paths) == 0 {
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
