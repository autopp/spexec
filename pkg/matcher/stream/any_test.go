package stream

import (
	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/matcher/testutil"
	"github.com/autopp/spexec/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AnyMatcher", func() {
	successExampleMatcher := testutil.NewExampleStreamMatcher(true, "successExampleMatcher message", nil)
	failureExampleMatcher := testutil.NewExampleStreamMatcher(false, "failureExampleMatcher message", nil)

	DescribeTable("Match",
		func(innerMatchers []model.StreamMatcher, expectedMatched bool, expectedMessage string) {
			m := &AnyMatcher{matchers: innerMatchers}
			matched, message, err := m.Match([]byte("given"))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual is matched to any inner matcher, returns true", []model.StreamMatcher{failureExampleMatcher, successExampleMatcher}, true, successExampleMatcher.SuccessMessage()),
		Entry("when actual is not matched to all of inner matchers, returns false", []model.StreamMatcher{failureExampleMatcher, failureExampleMatcher}, false, `should satisfy any of [`+failureExampleMatcher.FailureMessage()+`], [`+failureExampleMatcher.FailureMessage()+`]`),
	)
})

var _ = Describe("ParseAnyMatcher", func() {
	var v *model.Validator
	var r *matcher.StreamMatcherRegistry
	var parseExampleMatcher matcher.StreamMatcherParser
	var parserCalls *testutil.ParserCalls

	JustBeforeEach(func() {
		v, _ = model.NewValidator("", true)
		r = matcher.NewStreamMatcherRegistry()
		parseExampleMatcher, parserCalls = testutil.GenParseExampleStreamMatcher(true, "message", nil)
		r.Add("example", parseExampleMatcher)
		r.Add("failure", testutil.GenFailedParseStreamMatcher("parse error"))
	})

	Describe("with defined matchers", func() {
		It("returns matcher", func() {
			m := ParseAnyMatcher(v, r, model.Seq{model.Map{"example": "hello"}, model.Map{"example": "hi"}})

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			Expect(parserCalls.Calls).To(Equal([]any{"hello", "hi"}))
		})
	})

	DescribeTable("failure cases",
		func(given any, prefix string) {
			m := ParseAnyMatcher(v, r, given)

			Expect(m).To(BeNil())
			err := v.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(prefix))
		},
		Entry("with not seq", 42, "$: "),
		Entry("with undefined matcher", model.Seq{model.Map{"foo": 42}}, "$[0]: "),
		Entry("with invalid inner matcher parameter", model.Seq{model.Map{"failure": 42}}, "$[0].failure: "),
	)
})
