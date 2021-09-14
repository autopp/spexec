package parser

import (
	"path/filepath"
	"time"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/status"
	"github.com/autopp/spexec/internal/matcher/stream"
	"github.com/autopp/spexec/internal/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Parser", func() {
	Describe("ParseFile()", func() {
		var statusEqMatcher *status.EqMatcher
		var streamEqMatcher *stream.EqMatcher
		var p *Parser

		JustBeforeEach(func() {
			statusMR := matcher.NewStatusMatcherRegistry()
			statusMR.Add("eq", status.ParseEqMatcher)
			streamMR := matcher.NewStreamMatcherRegistry()
			streamMR.Add("eq", stream.ParseEqMatcher)
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
					"Name":    Equal("test_answer"),
					"Command": Equal([]string{"echo", "42"}),
					"Dir":     HaveSuffix("/testdata"),
					"Stdin":   Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeAssignableToTypeOf(statusEqMatcher),
					"StdoutMatcher": BeAssignableToTypeOf(streamEqMatcher),
					"StderrMatcher": BeNil(),
				})),
			}),
			Entry("testdata/test.json", "test.json", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":    Equal("test_answer"),
					"Command": Equal([]string{"echo", "42"}),
					"Dir":     HaveSuffix("/testdata"),
					"Stdin":   Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeAssignableToTypeOf(statusEqMatcher),
					"StdoutMatcher": BeAssignableToTypeOf(streamEqMatcher),
					"StderrMatcher": BeNil(),
				})),
			}),
		)
	})
})
