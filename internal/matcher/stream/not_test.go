package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/testutil"
	"github.com/autopp/spexec/internal/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotMatcher", func() {
	successMatcher := testutil.NewExampleStreamMatcher(true, "inner", nil)
	failureMatcher := testutil.NewExampleStreamMatcher(false, "inner", nil)

	DescribeTable("Match",
		func(inner model.StreamMatcher, expectedMatched bool, expectedMessage string) {
			m := &NotMatcher{matcher: inner}

			matched, message, err := m.Match([]byte("hello"))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual is not matched to inner matcher, returns true", failureMatcher, true, failureMatcher.FailureMessage()),
		Entry("when actual is matched to inner matcher, returns false", successMatcher, false, successMatcher.SuccessMessage()),
	)
})

var _ = Describe("ParseNotMatcher", func() {
	var v *model.Validator
	var r *matcher.StreamMatcherRegistry
	var parseExampleMatcherParser matcher.StreamMatcherParser
	var calls *testutil.ParserCalls

	JustBeforeEach(func() {
		v, _ = model.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()

		parseExampleMatcherParser, calls = testutil.GenParseExampleStreamMatcher(true, "example", nil)
		r.Add("example", parseExampleMatcherParser)
		failureParser := testutil.GenFailedParseStreamMatcher("failure")
		r.Add("failure", failureParser)
	})

	Describe("with defined matcher", func() {
		It("returns matcher", func() {
			m := ParseNotMatcher(v, r, model.Map{"example": true})

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			Expect(calls.Calls).To(Equal([]any{true}))
		})
	})

	DescribeTable("failure cases",
		func(given any, prefix string) {
			m := ParseNotMatcher(v, r, given)

			Expect(m).To(BeNil())
			err := v.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(prefix))
		},
		Entry("with not map", 42, "$: "),
		Entry("with undefined matcher", model.Map{"foo": 42}, "$: "),
		Entry("with invalid inner matcher parameter", model.Map{"failure": 42}, "$.failure: "),
	)
})
