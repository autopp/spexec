package template

import (
	"time"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/testutil"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"
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

				v, _ := model.NewValidator("")
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
					Dir:          model.NewTemplatableFromValue(""),
					Command: []*model.Templatable[any]{
						model.NewTemplatableFromValue[any]("echo"),
					},
					Stdin: model.NewTemplatableFromValue(""),
					Env: []TemplatableStringVar{
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
					Stdin:         []byte(""),
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
					Name:         model.NewTemplatableFromTemplateValue[string](model.NewTemplateValue(model.Map{"$": "name"}, []model.TemplateRef{model.NewTemplateVar("name")})),
					SpecFilename: "sample.yaml",
					Dir:          model.NewTemplatableFromValue(""),
					Command: []*model.Templatable[any]{
						model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"$": "command"}, []model.TemplateRef{model.NewTemplateVar("command")})),
					},
					Stdin: model.NewTemplatableFromValue(""),
					Env: []TemplatableStringVar{
						{
							Name:  "MESSAGE",
							Value: model.NewTemplatableFromValue("hello"),
						},
					},
					StatusMatcher: model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"$": "statusMatcher"}, []model.TemplateRef{model.NewTemplateVar("statusMatcher")})),
					StdoutMatcher: model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"$": "streamMatcher"}, []model.TemplateRef{model.NewTemplateVar("streamMatcher")})),
					StderrMatcher: model.NewTemplatableFromTemplateValue[any](model.NewTemplateValue(model.Map{"$": "streamMatcher"}, []model.TemplateRef{model.NewTemplateVar("streamMatcher")})),
					TeeStdout:     false,
					TeeStderr:     false,
					Timeout:       1 * time.Second,
				},
				&model.Test{
					Name:          "sample test",
					SpecFilename:  "sample.yaml",
					Dir:           "",
					Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
					Stdin:         []byte(""),
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
