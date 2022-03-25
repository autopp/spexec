package parser

import (
	"path/filepath"
	"time"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/status"
	"github.com/autopp/spexec/internal/matcher/stream"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"

	. "github.com/onsi/ginkgo/v2"
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
			p = New(statusMR, streamMR, false)
		})

		DescribeTable("with valid file",
			func(filename string, expected Elements) {
				actual, err := p.ParseFile(filepath.Join("testdata", filename))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(MatchAllElementsWithIndex(IndexIdentity, expected))
			},

			Entry("testdata/test.yaml", "test.yaml", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("testdata/test.yaml"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin":        Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeAssignableToTypeOf(statusEqMatcher),
					"StdoutMatcher": BeAssignableToTypeOf(streamEqMatcher),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				})),
			}),
			Entry("testdata/test.json", "test.json", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("testdata/test.json"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin":        Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeAssignableToTypeOf(statusEqMatcher),
					"StdoutMatcher": BeAssignableToTypeOf(streamEqMatcher),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				})),
			}),
			Entry("testdata/yaml-stdin.yaml", "yaml-stdin.yaml", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("testdata/yaml-stdin.yaml"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin": Equal([]byte(`array:
    - 1
    - true
    - hello
`)),
					"Env":           BeNil(),
					"Timeout":       Equal(0 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				})),
			}),
		)

		DescribeTable("with invalid file",
			func(filename string, expectedErr string) {
				_, err := p.ParseFile(filepath.Join("testdata", filename))
				Expect(err).To(MatchError(expectedErr))
			},
			Entry("testdata/root-is-not-map.yaml", "root-is-not-map.yaml", "$: should be map, but is seq"),
			Entry("testdata/spexec-version-is-invalid.yaml", "spexec-version-is-invalid.yaml", `$.spexec: should be "v0"`),
			Entry("testdata/spexec-version-is-not-string.yaml", "spexec-version-is-not-string.yaml", `$.spexec: should be string, but is int`),
			Entry("testdata/test-is-not-map.yaml", "test-is-not-map.yaml", "$.tests[0]: should be map, but is seq"),
			Entry("testdata/command-stdin-is-invalid.yaml", "command-stdin-is-invalid.yaml", "$.tests[0].stdin: should be a string or map, but is int"),
			Entry("testdata/command-stdin-format-is-invalid.yaml", "command-stdin-format-is-invalid.yaml", "$.tests[0].stdin.format: should be string, but is int"),
			Entry("testdata/command-stdin-value-is-invalid.yaml", "command-stdin-value-is-invalid.yaml", "$.tests[0].stdin: should have .value"),
		)

		Describe("with no exist file", func() {
			It("returns err", func() {
				_, err := p.ParseFile(filepath.Join("testdata", "unknown.yaml"))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
