package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqMatcher", func() {
	var m *EqMatcher
	JustBeforeEach(func() {
		m = &EqMatcher{expected: "hello"}
	})

	DescribeTable("Match",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.Match([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual equals to expected, returns true", "hello", true, `should not be "hello", but got it`),
		Entry("when actual dose not equal to expected, returns false", "goodbye\x01", false, "should be \"hello\", but got:\n----------------------------------------\n\x1b[31mh\x1b[0m\x1b[32mgoodby\x1b[0me\x1b[31mllo\x1b[0m\x1b[32m\x01\x1b[0m\n----------------------------------------"),
	)
})

var _ = Describe("ParseEqMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with string", func() {
		It("returns matcher", func() {
			m := ParseEqMatcher(v, r, "hello")

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var eq *EqMatcher = m.(*EqMatcher)
			Expect(eq.expected).To(Equal("hello"))
		})
	})

	DescribeTable("failure cases",
		func(given any) {
			m := ParseEqMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not string", 42),
	)
})
