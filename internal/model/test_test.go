package model

import (
	"time"

	"github.com/autopp/spexec/internal/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var matcherMessage = "failed"

type testStatusMatcher bool

func (m testStatusMatcher) Match(int) (bool, string, error) {
	return bool(m), matcherMessage, nil
}

type testStreamMatcher bool

func (m testStreamMatcher) Match([]byte) (bool, string, error) {
	return bool(m), matcherMessage, nil
}

var _ = Describe("Test", func() {
	DescribeTable("GetName()",
		func(test *Test, expected string) {
			Expect(test.GetName()).To(Equal(expected))
		},
		Entry("Name is not empty", &Test{Name: "test of echo", Command: []StringExpr{NewLiteralStringExpr("echo"), NewLiteralStringExpr("hello")}}, "test of echo"),
		Entry("Name is empty", &Test{Name: "", Command: []StringExpr{NewLiteralStringExpr("echo"), NewLiteralStringExpr("hello")}},
			"echo hello"),
		Entry("Name is empty and Env is given", &Test{
			Name:    "",
			Command: []StringExpr{NewLiteralStringExpr("make"), NewLiteralStringExpr("build")},
			Env: []util.StringVar{
				{Name: "GOOS", Value: "linux"},
				{Name: "GOARCH", Value: "amd64"},
			},
		}, "GOOS=linux GOARCH=amd64 make build"),
	)

	Describe("Run()", func() {
		DescribeTable("succeeded cases",
			func(test *Test, expectedMessages []*AssertionMessage, expectedIsSuccess bool) {
				tr, err := test.Run()
				Expect(err).NotTo(HaveOccurred())
				Expect(tr).To(Equal(&TestResult{Name: test.GetName(), Messages: expectedMessages, IsSuccess: expectedIsSuccess}))
			},
			Entry("no matchers", &Test{
				Name:    "no matchers case",
				Command: []StringExpr{NewLiteralStringExpr("echo")},
			}, []*AssertionMessage{}, true),
			Entry("matchers are passed", &Test{
				Name:          "matchers are passed",
				Command:       []StringExpr{NewLiteralStringExpr("echo")},
				StatusMatcher: testStatusMatcher(true),
				StdoutMatcher: testStreamMatcher(true),
				StderrMatcher: testStreamMatcher(true),
			}, []*AssertionMessage{}, true),
			Entry("status matcher is failed", &Test{
				Name:          "status matcher is failed",
				Command:       []StringExpr{NewLiteralStringExpr("echo")},
				StatusMatcher: testStatusMatcher(false),
				StdoutMatcher: testStreamMatcher(true),
				StderrMatcher: testStreamMatcher(true),
			}, []*AssertionMessage{{Name: "status", Message: matcherMessage}}, false),
			Entry("stdout matcher is failed", &Test{
				Name:          "stdout matcher is failed",
				Command:       []StringExpr{NewLiteralStringExpr("echo")},
				StatusMatcher: testStatusMatcher(true),
				StdoutMatcher: testStreamMatcher(false),
				StderrMatcher: testStreamMatcher(true),
			}, []*AssertionMessage{{Name: "stdout", Message: matcherMessage}}, false),
			Entry("stderr matcher is failed", &Test{
				Name:          "stderr matcher is failed",
				Command:       []StringExpr{NewLiteralStringExpr("echo")},
				StatusMatcher: testStatusMatcher(true),
				StdoutMatcher: testStreamMatcher(true),
				StderrMatcher: testStreamMatcher(false),
			}, []*AssertionMessage{{Name: "stderr", Message: matcherMessage}}, false),
			Entry("all matchers are failed", &Test{
				Name:          "all matchers are failed",
				Command:       []StringExpr{NewLiteralStringExpr("echo")},
				StatusMatcher: testStatusMatcher(false),
				StdoutMatcher: testStreamMatcher(false),
				StderrMatcher: testStreamMatcher(false),
			}, []*AssertionMessage{{Name: "status", Message: matcherMessage}, {Name: "stdout", Message: matcherMessage}, {Name: "stderr", Message: matcherMessage}}, false),
			Entry("process is timeout", &Test{
				Name:    "process is timeout",
				Command: []StringExpr{NewLiteralStringExpr("sleep"), NewLiteralStringExpr("1")},
				Timeout: 1 * time.Millisecond,
			}, []*AssertionMessage{{Name: "status", Message: "process was timeout"}}, false),
			Entry("process is signaled", &Test{
				Name:    "process is signaled",
				Command: []StringExpr{NewLiteralStringExpr("bash"), NewLiteralStringExpr("-c"), NewLiteralStringExpr("kill -TERM $$")},
			}, []*AssertionMessage{{Name: "status", Message: "process was signaled (terminated)"}}, false),
		)
	})

	DescribeTable("failed cases",
		func(test *Test, expectedErr string) {
			tr, err := test.Run()
			Expect(tr).To(BeNil())
			Expect(err).To(MatchError(expectedErr))
		},
		Entry("command evaluating is failed", &Test{
			Name:    "command evaluating is failed",
			Command: []StringExpr{NewEnvStringExpr("undefined")},
		}, "envrironment variable $undefined is not defined"),
	)
})
