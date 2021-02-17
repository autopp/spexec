package config

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
