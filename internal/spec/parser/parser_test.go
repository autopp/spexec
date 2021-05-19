package parser

import (
	"path/filepath"

	"github.com/autopp/spexec/internal/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("ParseFile()", func() {
		status := 0
		stdout := "42\n"

		DescribeTable("with valid file",
			func(filename string, expected []*model.Test) {
				Expect(New().ParseFile(filepath.Join("testdata", filename))).To(Equal(expected))
			},
			Entry("testdata/test.yaml", "test.yaml", []*model.Test{
				{
					Name:    "test_answer",
					Command: []string{"echo", "42"},
					Stdin:   "hello",
					Env:     map[string]string{"ANSWER": "42"},
					Status:  &status,
					Stdout:  &stdout,
				},
			}),
			Entry("testdata/test.json", "test.json", []*model.Test{
				{
					Name:    "test_answer",
					Command: []string{"echo", "42"},
					Stdin:   "hello",
					Env:     map[string]string{"ANSWER": "42"},
					Status:  &status,
					Stdout:  &stdout,
				},
			}),
		)
	})
})
