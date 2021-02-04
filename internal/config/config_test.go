package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openConfig(name string) io.Reader {
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return f
}

func TestLoad(t *testing.T) {
	r := openConfig("test.yaml")

	if c, err := Load(r); assert.NoError(t, err) {
		status := 0
		stdout := "42\n"
		expected := &Config{
			Tests: []Test{
				{
					Command: []string{"echo", "42"},
					Env: []struct {
						Name  string `yaml:"name"`
						Value string `yaml:"value"`
					}{{Name: "ANSWER", Value: "42"}},
					Expect: &struct {
						Status *int    `yaml:"status"`
						Stdout *string `yaml:"stdout"`
						Stderr *string `yaml:"stderr"`
					}{
						Status: &status,
						Stdout: &stdout,
					},
				},
			},
		}

		assert.Equal(t, expected, c)
	}
}
