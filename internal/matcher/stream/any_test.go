package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/testutil"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"

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
		Entry("when actual is matched to any inner matcher, returns true", []model.StreamMatcher{failureExampleMatcher, successExampleMatcher}, true, `successExampleMatcher message`),
		Entry("when actual is not matched to all of inner matchers, returns false", []model.StreamMatcher{failureExampleMatcher, failureExampleMatcher}, false, `should satisfy any of [failureExampleMatcher message], [failureExampleMatcher message]`),
	)
})

var _ = Describe("ParseAnyMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry
	var env *model.Env
	var parseExampleMatcher matcher.StreamMatcherParser
	var parserCalls *testutil.ParserCalls

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
		parseExampleMatcher, parserCalls = testutil.GenParseExampleStreamMatcher(true, "message", nil)
		r.Add("example", parseExampleMatcher)
		r.Add("failure", testutil.GenFailedParseStreamMatcher("parse error"))
		env = model.NewEnv(nil)
	})

	Describe("with defined matchers", func() {
		It("returns matcher", func() {
			m := ParseAnyMatcher(env, v, r, model.Seq{model.Map{"example": "hello"}, model.Map{"example": "hi"}})

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			Expect(parserCalls.Calls).To(Equal([]any{"hello", "hi"}))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}, prefix string) {
			m := ParseAnyMatcher(env, v, r, given)

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
