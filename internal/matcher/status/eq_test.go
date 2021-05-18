package status

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
