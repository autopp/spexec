package model_test

import (
	"time"

	"github.com/autopp/spexec/pkg/matcher/testutil"
	"github.com/autopp/spexec/pkg/model"
	"github.com/autopp/spexec/pkg/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test", func() {
	DescribeTable("GetName()",
		func(test *model.Test, expected string) {
			Expect(test.GetName()).To(Equal(expected))
		},
		Entry("Name is not empty", &model.Test{Name: "test of echo", Command: []model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("hello")}}, "test of echo"),
		Entry("Name is empty", &model.Test{Name: "", Command: []model.StringExpr{model.NewLiteralStringExpr("echo"), model.NewLiteralStringExpr("hello")}},
			"echo hello"),
		Entry("Name is empty and Env is given", &model.Test{
			Name:    "",
			Command: []model.StringExpr{model.NewLiteralStringExpr("make"), model.NewLiteralStringExpr("build")},
			Env: []util.StringVar{
				{Name: "GOOS", Value: "linux"},
				{Name: "GOARCH", Value: "amd64"},
			},
		}, "GOOS=linux GOARCH=amd64 make build"),
	)

	Describe("Run()", func() {
		successStatusMatcher := testutil.NewExampleStatusMatcher(true, "status", nil)
		failureStatusMatcher := testutil.NewExampleStatusMatcher(false, "status", nil)
		successStdoutMatcher := testutil.NewExampleStreamMatcher(true, "stdout", nil)
		failureStdoutMatcher := testutil.NewExampleStreamMatcher(false, "stdout", nil)
		successStderrMatcher := testutil.NewExampleStreamMatcher(true, "stderr", nil)
		failureStderrMatcher := testutil.NewExampleStreamMatcher(false, "stderr", nil)

		DescribeTable("succeeded cases",
			func(test *model.Test, expectedMessages []*model.AssertionMessage, expectedIsSuccess bool) {
				tr, err := test.Run()
				Expect(err).NotTo(HaveOccurred())
				Expect(tr).To(Equal(&model.TestResult{Name: test.GetName(), Messages: expectedMessages, IsSuccess: expectedIsSuccess}))
			},
			Entry("no matchers", &model.Test{
				Name:    "no matchers case",
				Command: []model.StringExpr{model.NewLiteralStringExpr("echo")},
			}, []*model.AssertionMessage{}, true),
			Entry("matchers are passed", &model.Test{
				Name:          "matchers are passed",
				Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
				StatusMatcher: successStatusMatcher,
				StdoutMatcher: successStdoutMatcher,
				StderrMatcher: successStderrMatcher,
			}, []*model.AssertionMessage{}, true),
			Entry("status matcher is failed", &model.Test{
				Name:          "status matcher is failed",
				Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
				StatusMatcher: failureStatusMatcher,
				StdoutMatcher: successStdoutMatcher,
				StderrMatcher: successStderrMatcher,
			}, []*model.AssertionMessage{{Name: "status", Message: failureStatusMatcher.FailureMessage()}}, false),
			Entry("stdout matcher is failed", &model.Test{
				Name:          "stdout matcher is failed",
				Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
				StatusMatcher: successStatusMatcher,
				StdoutMatcher: failureStdoutMatcher,
				StderrMatcher: successStderrMatcher,
			}, []*model.AssertionMessage{{Name: "stdout", Message: failureStdoutMatcher.FailureMessage()}}, false),
			Entry("stderr matcher is failed", &model.Test{
				Name:          "stderr matcher is failed",
				Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
				StatusMatcher: successStatusMatcher,
				StdoutMatcher: successStdoutMatcher,
				StderrMatcher: failureStderrMatcher,
			}, []*model.AssertionMessage{{Name: "stderr", Message: failureStderrMatcher.FailureMessage()}}, false),
			Entry("all matchers are failed", &model.Test{
				Name:          "all matchers are failed",
				Command:       []model.StringExpr{model.NewLiteralStringExpr("echo")},
				StatusMatcher: failureStatusMatcher,
				StdoutMatcher: failureStdoutMatcher,
				StderrMatcher: failureStderrMatcher,
			}, []*model.AssertionMessage{{Name: "status", Message: failureStatusMatcher.FailureMessage()}, {Name: "stdout", Message: failureStdoutMatcher.FailureMessage()}, {Name: "stderr", Message: failureStderrMatcher.FailureMessage()}}, false),
			Entry("process is timeout", &model.Test{
				Name:    "process is timeout",
				Command: []model.StringExpr{model.NewLiteralStringExpr("sleep"), model.NewLiteralStringExpr("1")},
				Timeout: 1 * time.Millisecond,
			}, []*model.AssertionMessage{{Name: "status", Message: "process was timeout"}}, false),
			Entry("process is signaled", &model.Test{
				Name:    "process is signaled",
				Command: []model.StringExpr{model.NewLiteralStringExpr("bash"), model.NewLiteralStringExpr("-c"), model.NewLiteralStringExpr("kill -TERM $$")},
			}, []*model.AssertionMessage{{Name: "status", Message: "process was signaled (terminated)"}}, false),
		)

		DescribeTable("failed cases",
			func(test *model.Test, expectedErr string) {
				tr, err := test.Run()
				Expect(tr).To(BeNil())
				Expect(err).To(MatchError(expectedErr))
			},
			Entry("command evaluating is failed", &model.Test{
				Name:    "command evaluating is failed",
				Command: []model.StringExpr{model.NewEnvStringExpr("undefined")},
			}, "environment variable $undefined is not defined"),
		)
	})
})
