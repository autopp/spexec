package stream

import (
	"regexp"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchRegexpMatcher", func() {
	var m *MatchRegexpMatcher
	r := regexp.MustCompile("ab+c")
	JustBeforeEach(func() {
		m = &MatchRegexpMatcher{expected: r}
	})

	DescribeTable("Match",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.Match([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual matches to expected, returns true", "xabbcx", true, `should not match to "ab+c", but match`),
		Entry("when actual dose not match to expected, returns false", "xabx", false, `should match to "ab+c", but got "xabx"`),
	)
})

var _ = Describe("ParseMatchRegexpMatcher", func() {
	var v *model.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = model.NewValidator("", true)
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with string", func() {
		It("returns matcher", func() {
			m := ParseMatchRegexpMatcher(v, r, "ab+c")

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var rm *MatchRegexpMatcher = m.(*MatchRegexpMatcher)
			Expect(rm.expected.String()).To(Equal("ab+c"))
		})
	})

	DescribeTable("failure cases",
		func(given any) {
			m := ParseMatchRegexpMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not string", 42),
		Entry("with invalid regexp source", "[a"),
	)
})
