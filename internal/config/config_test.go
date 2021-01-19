package config

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openConfig(name string) io.Reader {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("cannot load config " + name)
	}

	path := filepath.Join(filepath.Dir(filename), "..", "..", "example", name)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return f
}

func TestLoad(t *testing.T) {
	r := openConfig("test.yaml")

	if c, err := Load(r); assert.NoError(t, err) {
		expected := &Config{
			Tests: []Test{
				{Command: []string{"echo"}},
				{Command: []string{"exit", "1"}},
			},
		}

		assert.Equal(t, expected, c)
	}
}
