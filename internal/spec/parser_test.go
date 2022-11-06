package spec

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/status"
	"github.com/autopp/spexec/internal/matcher/stream"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/model/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Parser", func() {
	var p *Parser
	var env *model.Env

	JustBeforeEach(func() {
		statusMR := matcher.NewStatusMatcherRegistry()
		statusMR.Add("eq", status.ParseEqMatcher)
		streamMR := matcher.NewStreamMatcherRegistry()
		streamMR.Add("eq", stream.ParseEqMatcher)
		p = NewParser(statusMR, streamMR)
		env = model.NewEnv(nil)
	})

	Describe("ParseFile()", func() {

		DescribeTable("with valid file",
			func(filename string, expected Elements) {
				v, _ := model.NewValidator(filepath.Join("testdata", filename), true)
				actual, err := p.ParseFile(env, v, filepath.Join("testdata", filename))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(MatchAllElementsWithIndex(IndexIdentity, expected))
			},

			Entry("testdata/test.yaml", "test.yaml", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("testdata/test.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":   HaveSuffix("/testdata"),
					"Stdin": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("hello", []model.TemplateRef{}))),
					"Env": Equal([]*template.TemplatableStringVar{
						{Name: "ANSWER", Value: model.NewTemplatableFromValue("42")},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"eq": 0}, []model.TemplateRef{}))),
					"StdoutMatcher": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"eq": "42\n"}, []model.TemplateRef{}))),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				})),
			}),
			Entry("testdata/test.json", "test.json", Elements{
				"0": PointTo(MatchAllFields(Fields{
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("testdata/test.json"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":   HaveSuffix("/testdata"),
					"Stdin": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("hello", []model.TemplateRef{}))),
					"Env": Equal([]*template.TemplatableStringVar{
						{Name: "ANSWER", Value: model.NewTemplatableFromValue("42")},
					}),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"eq": json.Number("0")}, []model.TemplateRef{}))),
					"StdoutMatcher": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"eq": "42\n"}, []model.TemplateRef{}))),
					"StderrMatcher": BeNil(),
					"TeeStdout":     BeFalse(),
					"TeeStderr":     BeFalse(),
				})),
			}),
		)

		Describe("with no exist file", func() {
			It("returns err", func() {
				v, _ := model.NewValidator(filepath.Join("testdata", "unknown.yaml"), true)
				_, err := p.ParseFile(env, v, filepath.Join("testdata", "unknown.yaml"))
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("loadSpec", func() {
		DescribeTable("success cases",
			func(s any, expected Fields) {
				v, _ := model.NewValidator("testdata/spec.yaml", true)
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":   HaveSuffix("/testdata"),
					"Stdin": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("hello", []model.TemplateRef{}))),
					"Env": Equal([]*template.TemplatableStringVar{
						{Name: "ANSWER", Value: model.NewTemplatableFromValue("42")},
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":   HaveSuffix("/testdata"),
					"Stdin": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("hello", []model.TemplateRef{}))),
					"Env": Equal([]*template.TemplatableStringVar{
						{Name: "ANSWER", Value: model.NewTemplatableFromValue("42")},
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
				v, _ := model.NewValidator("testdata/spec.yaml", true)
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
				v, _ := model.NewValidator("testdata/spec.yaml", true)
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":   HaveSuffix("/testdata"),
					"Stdin": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("hello", []model.TemplateRef{}))),
					"Env": Equal([]*template.TemplatableStringVar{
						{Name: "ANSWER", Value: model.NewTemplatableFromValue("42")},
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
					"Dir":           HaveSuffix("/testdata"),
					"Stdin":         BeNil(),
					"Env":           BeEmpty(),
					"Timeout":       Equal(3 * time.Second),
					"StatusMatcher": Equal(model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"eq": 0}, []model.TemplateRef{}))),
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
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
					"Name":         Equal(model.NewTemplatableFromValue("test_answer")),
					"SpecFilename": HaveSuffix("/testdata/spec.yaml"),
					"Command": Equal([]*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("echo", []model.TemplateRef{})),
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue("42", []model.TemplateRef{}))},
					),
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
				v, _ := model.NewValidator("testdata/spec.yaml", true)
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

	Describe("loadCommandExpect", func() {
		DescribeTable("success cases",
			func(expect model.Map, statusMatcherShouldBeSet bool, stdoutMatcherShouldBeSet, stderrMatcherShouldBeSet bool) {
				v, _ := model.NewValidator("", true)
				actualStdin, actualStdout, actualStderr := p.loadCommandExpect(env, v, expect)
				Expect(v.Error()).NotTo(HaveOccurred())
				if statusMatcherShouldBeSet {
					Expect(actualStdin).NotTo(BeNil())
				} else {
					Expect(actualStdin).To(BeNil())
				}
				if stdoutMatcherShouldBeSet {
					Expect(actualStdout).NotTo(BeNil())
				} else {
					Expect(actualStdout).To(BeNil())
				}
				if stderrMatcherShouldBeSet {
					Expect(actualStderr).NotTo(BeNil())
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
				v, _ := model.NewValidator("", true)
				p.loadCommandExpect(env, v, expect)
				Expect(v.Error()).To(MatchError(expectedErr))
			},
			Entry("with unknown field", model.Map{"unknown": 42}, "$: field .unknown is not expected"),
		)
	})
})
