package spec

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
	var statusEqMatcher *status.EqMatcher
	var streamEqMatcher *stream.EqMatcher
	var p *Parser
	var env *model.Env

	JustBeforeEach(func() {
		statusMR := matcher.NewStatusMatcherRegistry()
		statusMR.Add("eq", status.ParseEqMatcher)
		streamMR := matcher.NewStreamMatcherRegistry()
		streamMR.Add("eq", stream.ParseEqMatcher)
		p = NewParser(statusMR, streamMR, true)
		env = model.NewEnv(nil)
	})

	Describe("ParseFile()", func() {

		DescribeTable("with valid file",
			func(filename string, expected Elements) {
				actual, err := p.ParseFile(env, filepath.Join("testdata", filename))
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
		)

		Describe("with no exist file", func() {
			It("returns err", func() {
				_, err := p.ParseFile(env, filepath.Join("testdata", "unknown.yaml"))
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("loadSpec", func() {
		DescribeTable("success cases",
			func(s any, expected Fields) {
				v, _ := model.NewValidator("testdata/spec.yaml")
				actual, err := p.loadSpec(env, v, s)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(MatchAllElementsWithIndex(IndexIdentity, Elements{
					"0": PointTo(MatchAllFields(expected)),
				}))
			},
			Entry("without .spexec",
				model.Map{
					"tests": model.Seq{
						model.Map{
							"name":    "test_answer",
							"command": model.Seq{"echo", "42"},
							"stdin":   "hello",
							"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
							"timeout": 3,
						},
					},
				},
				Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin":        Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				},
			),
			Entry("with .spexec",
				model.Map{
					"spexec": "v0",
					"tests": model.Seq{
						model.Map{
							"name":    "test_answer",
							"command": model.Seq{"echo", "42"},
							"stdin":   "hello",
							"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
							"timeout": 3,
						},
					},
				},
				Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin":        Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				},
			),
		)

		DescribeTable("failure cases",
			func(s any, expectedErr string) {
				v, _ := model.NewValidator("testdata/spec.yaml")
				_, err := p.loadSpec(env, v, s)
				Expect(err).To(MatchError(expectedErr))
			},
			Entry("with unknown field",
				model.Map{
					"tests": model.Seq{
						model.Map{
							"name":    "test_answer",
							"command": model.Seq{"echo", "42"},
							"stdin":   "hello",
							"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
							"timeout": 3,
						},
					},
					"unknown": 42,
				},
				"$: field .unknown is not expected",
			),
			Entry("with invalid .spexec",
				model.Map{
					"spexec": "invalid",
					"tests": model.Seq{
						model.Map{
							"name":    "test_answer",
							"command": model.Seq{"echo", "42"},
							"stdin":   "hello",
							"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
							"timeout": 3,
						},
					},
				},
				`$.spexec: should be "v0"`,
			),
			Entry("with invalid .tests",
				model.Map{
					"tests": model.Map{
						"name":    "test_answer",
						"command": model.Seq{"echo", "42"},
						"stdin":   "hello",
						"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
						"timeout": 3,
					},
				},
				`$.tests: should be seq, but is map`,
			),
			Entry("with not map",
				model.Seq{
					model.Map{
						"spexec": "invalid",
						"tests": model.Seq{
							model.Map{
								"name":    "test_answer",
								"command": model.Seq{"echo", "42"},
								"stdin":   "hello",
								"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
								"timeout": 3,
							},
						},
					},
				},
				`$: should be map, but is seq`,
			),
		)
	})

	Describe("loadTest", func() {
		DescribeTable("success cases",
			func(test any, expected Fields) {
				v, _ := model.NewValidator("testdata/spec.yaml")
				actual := p.loadTest(env, v, test)
				Expect(v.Error()).NotTo(HaveOccurred())
				Expect(actual).To(PointTo(MatchAllFields(expected)))
			},
			Entry("without any matchers",
				model.Map{
					"name":    "test_answer",
					"command": model.Seq{"echo", "42"},
					"stdin":   "hello",
					"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
					"timeout": 3,
				},
				Fields{
					"Name":         Equal("test_answer"),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command":      Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":          HaveSuffix("/testdata"),
					"Stdin":        Equal([]byte("hello")),
					"Env": Equal([]util.StringVar{
						{Name: "ANSWER", Value: "42"},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				},
			),
			Entry("with matcher",
				model.Map{
					"name":    "test_answer",
					"command": model.Seq{"echo", "42"},
					"expect":  model.Map{"status": model.Map{"eq": 0}},
					"timeout": 3,
				},
				Fields{
					"Name":          Equal("test_answer"),
					"SpecFilename":  HaveSuffix("/testdata/spec.yaml"),
					"Command":       Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":           HaveSuffix("/testdata"),
					"Stdin":         BeNil(),
					"Env":           BeEmpty(),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeAssignableToTypeOf(statusEqMatcher),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				},
			),
			Entry("with TeeStdout",
				model.Map{
					"name":      "test_answer",
					"command":   model.Seq{"echo", "42"},
					"timeout":   3,
					"teeStdout": true,
				},
				Fields{
					"Name":          Equal("test_answer"),
					"SpecFilename":  HaveSuffix("/testdata/spec.yaml"),
					"Command":       Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":           HaveSuffix("/testdata"),
					"Stdin":         BeNil(),
					"Env":           BeEmpty(),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeTrue(),
					"TeeStderr":     BeFalse(),
				},
			),
			Entry("with TeeStderr",
				model.Map{
					"name":      "test_answer",
					"command":   model.Seq{"echo", "42"},
					"timeout":   3,
					"teeStderr": true,
				},
				Fields{
					"Name":          Equal("test_answer"),
					"SpecFilename":  HaveSuffix("/testdata/spec.yaml"),
					"Command":       Equal([]model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("42")}),
					"Dir":           HaveSuffix("/testdata"),
					"Stdin":         BeNil(),
					"Env":           BeEmpty(),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": BeNil(),
					"StdoutMatcher": BeNil(),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeTrue(),
				},
			),
		)

		DescribeTable("failure cases",
			func(test any, expectedErr string) {
				v, _ := model.NewValidator("testdata/spec.yaml")
				p.loadTest(env, v, test)
				Expect(v.Error()).To(MatchError(expectedErr))
			},
			Entry("with not map", 42, "$: should be map, but is int"),
			Entry("with unknown field",
				model.Map{
					"name":    "test_answer",
					"command": model.Seq{"echo", "42"},
					"stdin":   "hello",
					"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
					"timeout": 3,
					"unknown": 42,
				},
				"$: field .unknown is not expected",
			),
			Entry("with invalid timeout",
				model.Map{
					"name":    "test_answer",
					"command": model.Seq{"echo", "42"},
					"stdin":   "hello",
					"env":     model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
					"timeout": false,
				},
				"$.timeout: should be positive integer or duration string, but is bool",
			),
			Entry("with invalid teeStdout",
				model.Map{
					"name":      "test_answer",
					"command":   model.Seq{"echo", "42"},
					"stdin":     "hello",
					"env":       model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
					"teeStdout": 42,
				},
				"$.teeStdout: should be bool, but is int",
			),
			Entry("with invalid teeStderr",
				model.Map{
					"name":      "test_answer",
					"command":   model.Seq{"echo", "42"},
					"stdin":     "hello",
					"env":       model.Seq{model.Map{"name": "ANSWER", "value": "42"}},
					"teeStderr": 42,
				},
				"$.teeStderr: should be bool, but is int",
			),
		)
	})

	Describe("loadCommandStdin", func() {
		DescribeTable("success cases",
			func(stdin any, expected string) {
				v, _ := model.NewValidator("")
				actual := p.loadCommandStdin(v, stdin)
				Expect(v.Error()).NotTo(HaveOccurred())
				Expect(string(actual)).To(Equal(expected))
			},
			Entry("with simple string", "hello", "hello"),
			Entry("with yaml format", model.Map{"format": "yaml", "value": model.Seq{"hello", "world"}}, "- hello\n- world\n"),
		)

		DescribeTable("failure cases",
			func(stdin any, expectedErr string) {
				v, _ := model.NewValidator("")
				Expect(p.loadCommandStdin(v, stdin)).To(BeNil())
				Expect(v.Error()).To(MatchError(expectedErr))
			},
			Entry("with no string nor map", 42, "$: should be a string or map, but is int"),
			Entry("with .format missing map", model.Map{"value": model.Seq{"hello", "world"}}, "$: should have .format as string"),
			Entry("with .value missing map", model.Map{"format": "yaml"}, "$: should have .value"),
			Entry("with invalid .format map", model.Map{"format": 42, "value": 42}, `$.format: should be string, but is int`),
			Entry("with unknown .format map", model.Map{"format": "unknown", "value": 42}, `$.format: should be a "yaml", but is "unknown"`),
			Entry("with unknown field", model.Map{"format": "yaml", "value": 42, "unknown": 42}, `$: field .unknown is not expected`),
		)
	})

	Describe("loadCommandExpect", func() {
		DescribeTable("success cases",
			func(expect model.Map, statusMatcherShouldBeSet bool, stdoutMatcherShouldBeSet, stderrMatcherShouldBeSet bool) {
				v, _ := model.NewValidator("")
				actualStdin, actualStdout, actualStderr := p.loadCommandExpect(env, v, expect)
				Expect(v.Error()).NotTo(HaveOccurred())
				if statusMatcherShouldBeSet {
					Expect(actualStdin).To(BeAssignableToTypeOf(statusEqMatcher))
				} else {
					Expect(actualStdin).To(BeNil())
				}
				if stdoutMatcherShouldBeSet {
					Expect(actualStdout).To(BeAssignableToTypeOf(streamEqMatcher))
				} else {
					Expect(actualStdout).To(BeNil())
				}
				if stderrMatcherShouldBeSet {
					Expect(actualStderr).To(BeAssignableToTypeOf(streamEqMatcher))
				} else {
					Expect(actualStderr).To(BeNil())
				}
			},
			Entry("without any matchers", model.Map{}, false, false, false),
			Entry("with only status", model.Map{"status": model.Map{"eq": 0}}, true, false, false),
			Entry("with only stdout", model.Map{"stdout": model.Map{"eq": ""}}, false, true, false),
			Entry("with only stderr", model.Map{"stderr": model.Map{"eq": ""}}, false, false, true),
			Entry("with all matchers", model.Map{"status": model.Map{"eq": 0}, "stdout": model.Map{"eq": ""}, "stderr": model.Map{"eq": ""}}, true, true, true),
		)

		DescribeTable("failure cases",
			func(expect model.Map, expectedErr string) {
				v, _ := model.NewValidator("")
				p.loadCommandExpect(env, v, expect)
				Expect(v.Error()).To(MatchError(expectedErr))
			},
			Entry("with unknown status", model.Map{"status": model.Map{"unknown": true}}, "$.status: matcher for status unknown is not defined"),
			Entry("with unknown stdout", model.Map{"stdout": model.Map{"unknown": true}}, "$.stdout: matcher for stream unknown is not defined"),
			Entry("with unknown stderr", model.Map{"stderr": model.Map{"unknown": true}}, "$.stderr: matcher for stream unknown is not defined"),
			Entry("with unknown field", model.Map{"unknown": 42}, "$: field .unknown is not expected"),
		)
	})
})
