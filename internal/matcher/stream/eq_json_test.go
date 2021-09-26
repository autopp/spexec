package stream

import (
	"encoding/json"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("EqJSONMatcher", func() {
	var m *EqJSONMatcher
	JustBeforeEach(func() {
		m = &EqJSONMatcher{expected: spec.Map{"code": json.Number("200"), "messages": spec.Seq{"hello"}}}
	})

	DescribeTable("MatchStatus",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.MatchStream([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(Equal(expectedMessage))
		},
		Entry("when actual equals to expected as JSON, returns true", `{"messages": ["hello"], "code": 200}`, true, `should not be {"code":200,"messages":["hello"]}, but got it`),
		Entry("when actual dose not equal to expected as JSON, returns false", `{"messages": ["hi"], "code": 200}`, false, `should be {"code":200,"messages":["hello"]}, but got {"messages": ["hi"], "code": 200}`),
	)
})

var _ = Describe("ParseEqJSONMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with object", func() {
		It("returns matcher", func() {
			m := ParseEqJSONMatcher(v, r, spec.Map{"code": 200, "messages": spec.Seq{"hello"}})

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var eq *EqJSONMatcher = m.(*EqJSONMatcher)
			Expect(eq.expected).To(Equal(spec.Map{"code": json.Number("200"), "messages": spec.Seq{"hello"}}))
		})
	})
})
