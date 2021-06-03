package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqMatcher", func() {
	var m *EqMatcher
	JustBeforeEach(func() {
		m = &EqMatcher{expected: "hello"}
	})

	DescribeTable("MatchStatus",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual equals to expected, returns true", "hello", true, `should not be "hello", but got it`),
		Entry("when actual dose not equal to expected, returns false", "goodbye\x01", false, `should be "hello", but got "goodbye\x01"`),
	)
})

var _ = Describe("ParseEqMatcher", func() {
	var v *spec.Validator
	var r *matcher.StatusMatcherRegistry

	JustBeforeEach(func() {
		v = spec.NewValidator()
		r = matcher.NewStatusMatcherRegistry()
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
		func(given interface{}) {
			m := ParseEqMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not string", 42),
	)
})
