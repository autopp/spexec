package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type emptyMatcher struct{}

func (*emptyMatcher) MatchStream(actual []byte) (bool, string, error) {
	if len(actual) == 0 {
		return true, "should not be empty", nil
	}

	return false, "should be empty", nil
}

func parseEmptyMatcher(_ *spec.Validator, r *matcher.StreamMatcherRegistry, x interface{}) matcher.StreamMatcher {
	return &emptyMatcher{}
}

var _ = Describe("NotMatcher", func() {
	var m *NotMatcher
	JustBeforeEach(func() {
		m = &NotMatcher{matcher: &emptyMatcher{}}
	})

	DescribeTable("MatchStatus",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual is matched to inner matcher, returns true", "hello", true, `should be empty`),
		Entry("when actual is not matched to inner matcher, returns false", "", false, `should not be empty`),
	)
})

var _ = Describe("ParseNotMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v = spec.NewValidator()
		r = matcher.NewStreamMatcherRegistry()
		r.Add("empty", parseEmptyMatcher)
	})

	Describe("with defined matcher", func() {
		It("returns matcher", func() {
			m := ParseNotMatcher(v, r, spec.Map{"empty": nil})

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			var not *NotMatcher = m.(*NotMatcher)
			var empty *emptyMatcher
			Expect(not.matcher).To(BeAssignableToTypeOf(empty))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}) {
			m := ParseNotMatcher(v, r, given)

			Expect(m).To(BeNil())
			Expect(v.Error()).To(HaveOccurred())
		},
		Entry("with not map", 42),
		Entry("with undefined matcher", spec.Map{"foo": 42}),
	)
})
