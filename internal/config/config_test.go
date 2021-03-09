package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/autopp/spexec/internal/model"
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

	if c, err := Load(r, YAMLFormat); assert.NoError(t, err) {
		status := 0
		stdout := "42\n"
		expected := []*model.Test{
			{
				Name:    "test_answer",
				Command: []string{"echo", "42"},
				Env:     map[string]string{"ANSWER": "42"},
				Status:  &status,
				Stdout:  &stdout,
			},
		}

		assert.Equal(t, expected, c)
	}
}
