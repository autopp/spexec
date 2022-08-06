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
				env := model.NewEnv(nil)
				v, _ := model.NewValidator("")
				statusMR := matcher.NewStatusMatcherRegistry()
				statusMatcherParser, _ := testutil.GenParseExampleStatusMatcher(true, "message", nil)
				statusMR.Add("statusExample", statusMatcherParser)

				streamMR := matcher.NewStreamMatcherRegistry()
				streamMatcherParser, _ := testutil.GenParseExampleStreamMatcher(true, "message", nil)
				streamMR.Add("streamExample", streamMatcherParser)

				Expect(tt.Expand(env, v, statusMR, streamMR)).To(Equal(expected))
			},
			Entry("with no placeholder",
				&TestTemplate{
					Name:         model.NewTemplatableFromValue("sample test"),
					SpecFilename: "sample.yaml",
					Dir:          model.NewTemplatableFromValue(""),
					Command: []*model.Templatable[model.StringExpr]{
						model.NewTemplatableFromValue(model.NewLiteralStringExpr("echo")),
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
		)
	})
})
