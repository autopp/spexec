package stream

import (
	"fmt"
	"strings"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type prefixMatcher struct {
	prefix string
}

func (m *prefixMatcher) MatchStream(actual []byte) (bool, string, error) {
	if strings.HasPrefix(string(actual), m.prefix) {
		return true, fmt.Sprintf("should not start with %q", m.prefix), nil
	}

	return false, fmt.Sprintf("should start with %q", m.prefix), nil
}

func parsePrefixMatcher(v *spec.Validator, r *matcher.StreamMatcherRegistry, x interface{}) model.StreamMatcher {
	switch prefix := x.(type) {
	case string:
		return &prefixMatcher{prefix: prefix}
	default:
		v.AddViolation("parameter should be string")
		return nil
	}
}

var _ = Describe("AnyMatcher", func() {
	var m *AnyMatcher
	JustBeforeEach(func() {
		m = &AnyMatcher{matchers: []model.StreamMatcher{&prefixMatcher{prefix: "ab"}, &prefixMatcher{prefix: "xy"}}}
	})

	DescribeTable("MatchStream",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual is matched to any inner matcher, returns true", "xyz", true, `should not start with "xy"`),
		Entry("when actual is not matched to all of inner matchers, returns false", "def", false, `should satisfy any of [should start with "ab"], [should start with "xy"]`),
	)
})

var _ = Describe("ParseAnyMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
		r.Add("prefix", parsePrefixMatcher)
	})

	Describe("with defined matchers", func() {
		It("returns matcher", func() {
			m := ParseAnyMatcher(v, r, spec.Seq{spec.Map{"prefix": "hello"}, spec.Map{"prefix": "hello"}})

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			var any *AnyMatcher = m.(*AnyMatcher)
			var prefix *prefixMatcher
			Expect(any.matchers).To(HaveLen(2))
			Expect(any.matchers[0]).To(BeAssignableToTypeOf(prefix))
			Expect(any.matchers[1]).To(BeAssignableToTypeOf(prefix))
		})
	})

	DescribeTable("failure cases",
		func(given interface{}, prefix string) {
			m := ParseAnyMatcher(v, r, given)

			Expect(m).To(BeNil())
			err := v.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(prefix))
		},
		Entry("with not seq", 42, "$: "),
		Entry("with undefined matcher", spec.Seq{spec.Map{"foo": 42}}, "$[0]: "),
		Entry("with invalid inner matcher parameter", spec.Seq{spec.Map{"prefix": 42}}, "$[0].prefix: "),
	)
})
