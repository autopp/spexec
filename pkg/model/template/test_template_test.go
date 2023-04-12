package template

import (
	"time"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/matcher/testutil"
	"github.com/autopp/spexec/pkg/model"
	"github.com/autopp/spexec/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestTemplate", func() {
	Describe("Expand()", func() {
		DescribeTable("success cases",
			func(tt *TestTemplate, expected *model.Test) {
				// Arrenge
				env := model.NewEnv(nil)
				env.Define("name", "sample test")
				env.Define("dir", "")
				env.Define("command", "echo")
				env.Define("envMessage", "hello")
				env.Define("statusMatcher", model.Map{"statusExample": nil})
				env.Define("streamMatcher", model.Map{"streamExample": nil})

				v, _ := model.NewValidator("", true)
				statusMR := matcher.NewStatusMatcherRegistry()
				statusMatcherParser, _ := testutil.GenParseExampleStatusMatcher(true, "message", nil)
				statusMR.Add("statusExample", statusMatcherParser)

				streamMR := matcher.NewStreamMatcherRegistry()
				streamMatcherParser, _ := testutil.GenParseExampleStreamMatcher(true, "message", nil)
				streamMR.Add("streamExample", streamMatcherParser)

				// Act & Assert
				Expect(v.Error()).NotTo(HaveOccurred())
				Expect(tt.Expand(env, v, statusMR, streamMR)).To(Equal(expected))
			},
			Entry("with no placeholder",
				&TestTemplate{
					Name:         model.NewTemplatableFromValue("sample test"),
					SpecFilename: "sample.yaml",
					Dir:          "",
					Command: []*model.Templatable[any]{
						model.NewTemplatableFromValue[any]("echo"),
					},
					Stdin: model.NewTemplatableFromValue[any]("stdin"),
					Env: []*TemplatableStringVar{
						{
							Name:  "MESSAGE",
							Value: model.NewTemplatableFromValue("hello"),
						},
					},
					StatusMatcher: model.NewTemplatableFromValue[any](model.Map{"statusExample": nil}),
					StdoutMatcher: model.NewTemplatableFromValue[any](model.Map{"streamExample": nil}),
					StderrMatcher: model.NewTemplatableFromValue[any](model.Map{"streamExample": nil}),
					TeeStdout:     false,
					TeeStderr:     false,
					Timeout:       1 * time.Second,
				},
				&model.Test{
					Name:          "sample test",
					SpecFilename:  "sample.yaml",
					Dir:           "",
					Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
					Stdin:         []byte("stdin"),
					Env:           []util.StringVar{{Name: "MESSAGE", Value: "hello"}},
					StatusMatcher: testutil.NewExampleStatusMatcher(true, "message", nil),
					StdoutMatcher: testutil.NewExampleStreamMatcher(true, "message", nil),
					StderrMatcher: testutil.NewExampleStreamMatcher(true, "message", nil),
					TeeStdout:     false,
					TeeStderr:     false,
					Timeout:       1 * time.Second,
				},
			),
			Entry("with placeholders",
				&TestTemplate{
					Name:         model.NewTemplatableFromVariable[string]("name"),
					SpecFilename: "sample.yaml",
					Dir:          "",
					Command: []*model.Templatable[any]{
						model.NewTemplatableFromVariable[any]("command"),
					},
					Stdin: model.NewTemplatableFromValue[any]("stdin"),
					Env: []*TemplatableStringVar{
						{
							Name:  "MESSAGE",
							Value: model.NewTemplatableFromValue("hello"),
						},
					},
					StatusMatcher: model.NewTemplatableFromVariable[any]("statusMatcher"),
					StdoutMatcher: model.NewTemplatableFromVariable[any]("streamMatcher"),
					StderrMatcher: model.NewTemplatableFromVariable[any]("streamMatcher"),
					TeeStdout:     false,
					TeeStderr:     false,
					Timeout:       1 * time.Second,
				},
				&model.Test{
					Name:          "sample test",
					SpecFilename:  "sample.yaml",
					Dir:           "",
					Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
					Stdin:         []byte("stdin"),
					Env:           []util.StringVar{{Name: "MESSAGE", Value: "hello"}},
					StatusMatcher: testutil.NewExampleStatusMatcher(true, "message", nil),
					StdoutMatcher: testutil.NewExampleStreamMatcher(true, "message", nil),
					StderrMatcher: testutil.NewExampleStreamMatcher(true, "message", nil),
					TeeStdout:     false,
					TeeStderr:     false,
					Timeout:       1 * time.Second,
				},
			),
		)
	})
})

var _ = Describe("evalCommandStdin", func() {
	DescribeTable("success cases",
		func(stdin any, expected string) {
			v, _ := model.NewValidator("", true)
			actual := evalCommandStdin(v, stdin)
			Expect(v.Error()).NotTo(HaveOccurred())
			Expect(string(actual)).To(Equal(expected))
		},
		Entry("with simple string", "hello", "hello"),
		Entry("with yaml format", model.Map{"format": "yaml", "value": model.Seq{"hello", "world"}}, "- hello\n- world\n"),
	)

	DescribeTable("failure cases",
		func(stdin any, expectedErr string) {
			v, _ := model.NewValidator("", true)
			Expect(evalCommandStdin(v, stdin)).To(BeNil())
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
