package config

import (
	"os"
	"path/filepath"

	"github.com/autopp/spexec/internal/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func openConfig(name string) []byte {
	path := filepath.Join("testdata", name)
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

var _ = Describe("LoadYAML()", func() {
	It("returns loaded []*Test", func() {
		b := openConfig("test.yaml")
		status := 0
		stdout := "42\n"
		expected := []*model.Test{
			{
				Name:    "test_answer",
				Command: []string{"echo", "42"},
				Stdin:   "hello",
				Env:     map[string]string{"ANSWER": "42"},
				Status:  &status,
				Stdout:  &stdout,
			},
		}
		Expect(LoadYAML(b)).To(Equal(expected))
	})
})
