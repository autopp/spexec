package stream

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("SatisfyMatcher", func() {
	var m *SatisfyMatcher
	JustBeforeEach(func() {
		m = &SatisfyMatcher{
			Command: []string{"bash", "-c", `test "$(cat -)" = hello`},
		}
	})

	DescribeTable("MatchStream",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when command with given input via stdin succeeds, returns true", "hello", true, "should make the given command fail"),
		Entry("when command with given input via stdin fails, returns false", "goodbye", false, `should make the given command succeed`),
	)
})
