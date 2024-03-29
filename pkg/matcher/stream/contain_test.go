package stream

import (
	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainMatcher", func() {
	var m *ContainMatcher
	JustBeforeEach(func() {
		m = &ContainMatcher{expected: "hello"}
	})

	DescribeTable("Match",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.Match([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual contains expected, returns true", "Message: hello world", true, `should not contain "hello", but contain`),
		Entry("when actual dose not contain expected, returns false", "goodbye\x01", false, `should contain "hello", but got "goodbye\x01"`),
	)
})

var _ = Describe("ParseContainMatcher", func() {
	var v *model.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = model.NewValidator("", true)
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with string", func() {
		It("returns matcher", func() {
			m := ParseContainMatcher(v, r, "hello")

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var eq *ContainMatcher = m.(*ContainMatcher)
			Expect(eq.expected).To(Equal("hello"))
		})
	})

	DescribeTable("failure cases",
		func(given any) {
			m := ParseContainMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not string", 42),
	)
})
