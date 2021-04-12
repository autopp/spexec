package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/autopp/spexec/internal/model"
	"github.com/stretchr/testify/assert"
)

func openConfig(name string) []byte {
	path := filepath.Join("testdata", name)
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

func TestLoadYAML(t *testing.T) {
	b := openConfig("test.yaml")

	if c, err := LoadYAML(b); assert.NoError(t, err) {
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
