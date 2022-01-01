package status

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqMatcher", func() {
	var m *EqMatcher
	JustBeforeEach(func() {
		m = &EqMatcher{expected: 0}
	})

	DescribeTable("MatchStatus",
		func(given int, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStatus(given)
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual equals to expected, returns true", 0, true, "should not be 0, but got it"),
		Entry("when actual dose not equal to expected, returns false", 1, false, "should be 0, but got 1"),
	)
})

var _ = Describe("ParseEqMatcher", func() {
	var v *spec.Validator
	var r *matcher.StatusMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStatusMatcherRegistry()
	})

	Describe("with natural number", func() {
		It("returns matcher", func() {
			m := ParseEqMatcher(v, r, 0)

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var eq *EqMatcher = m.(*EqMatcher)
			Expect(eq.expected).To(Equal(0))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}) {
			m := ParseEqMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with negative integer", -1),
		Entry("with float", 0.0),
		Entry("with not number", "0"),
	)
})
