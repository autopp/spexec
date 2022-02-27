package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeEmptyMatcher", func() {
	DescribeTable("MatchStatus",
		func(expected bool, given string, expectedMatched bool, expectedMessage string) {
			m := &BeEmptyMatcher{expected: expected}
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry(`when expectation is empty and actual is "", returns true`, true, "", true, `should not be empty, but is empty`),
		Entry(`when expectation is empty and actual is "hello", returns false`, true, "hello", false, "should be empty, but is not:\n----------------------------------------\nhello\n----------------------------------------"),
		Entry(`when expectation is not empty and actual is "", returns false`, false, "", false, "should not be empty, but is empty"),
		Entry(`when expectation is not empty and actual is "hello", returns false`, false, "hello", true, "should be empty, but is not:\n----------------------------------------\nhello\n----------------------------------------"),
	)
})

var _ = Describe("ParseBeEmptyMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with bool", func() {
		It("returns matcher", func() {
			m := ParseBeEmptyMatcher(v, r, true)

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var beEmpty *BeEmptyMatcher = m.(*BeEmptyMatcher)
			Expect(beEmpty.expected).To(Equal(true))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}) {
			m := ParseBeEmptyMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not string", 42),
	)
})
