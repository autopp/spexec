package status

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuccessMatcher", func() {
	DescribeTable("MatchStatus",
		func(expected bool, given int, expectedMatched bool, expectedMessage string) {
			m := &SuccessMatcher{expected: expected}
			matched, message, err := m.MatchStatus(given)
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when expectation is success and actual is 0, returns true", true, 0, true, "should not succeed, but succeeded (status is 0)"),
		Entry("when expectation is success and actual is 1, returns false", true, 1, false, "should succeed, but not succeeded (status is 1)"),
		Entry("when expectation is failure and actual is 0, returns false", false, 0, false, "should not succeed, but succeeded (status is 0)"),
		Entry("when expectation is failure and actual is 0, returns false", false, 1, true, "should succeed, but not succeeded (status is 1)"),
	)
})

var _ = Describe("ParseSuccessMatcher", func() {
	var v *spec.Validator
	var r *matcher.StatusMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStatusMatcherRegistry()
	})

	Describe("with bool", func() {
		It("returns matcher", func() {
			m := ParseSuccessMatcher(v, r, true)

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var success *SuccessMatcher = m.(*SuccessMatcher)
			Expect(success.expected).To(Equal(true))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}) {
			m := ParseSuccessMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not bool", 1),
	)
})
