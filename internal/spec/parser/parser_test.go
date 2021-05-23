package parser

import (
	"path/filepath"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Parser", func() {
	Describe("ParseFile()", func() {
		var eqMatcher *status.EqMatcher
		var p *Parser
		stdout := "42\n"

		JustBeforeEach(func() {
			statusMR := matcher.NewStatusMatcherRegistry()
			statusMR.Add("eq", status.ParseEqMatcher)
			streamMR := matcher.NewStreamMatcherRegistry()
			p = New(statusMR, streamMR)
		})

		DescribeTable("with valid file",
			func(filename string, expected Elements) {
				// Expect(New().ParseFile(filepath.Join("testdata", filename))).To(Equal(expected))
				actual, err := p.ParseFile(filepath.Join("testdata", filename))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(MatchAllElementsWithIndex(IndexIdentity, expected))
			},

			Entry("testdata/test.yaml", "test.yaml", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":          Equal("test_answer"),
					"Command":       Equal([]string{"echo", "42"}),
					"Stdin":         Equal("hello"),
					"Env":           Equal(map[string]string{"ANSWER": "42"}),
					"StatusMatcher": BeAssignableToTypeOf(eqMatcher),
					"Stdout":        Equal(&stdout),
					"Stderr":        BeNil(),
				})),
			}),
			// Entry("testdata/test.json", "test.json", []*model.Test{
			// 	{
			// 		Name:    "test_answer",
			// 		Command: []string{"echo", "42"},
			// 		Stdin:   "hello",
			// 		Env:     map[string]string{"ANSWER": "42"},
			// 		Status:  &status,
			// 		Stdout:  &stdout,
			// 	},
			// }),
		)
	})
})
