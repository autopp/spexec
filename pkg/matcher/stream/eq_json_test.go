package stream

import (
	"encoding/json"

	"github.com/autopp/spexec/pkg/matcher"
	"github.com/autopp/spexec/pkg/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("EqJSONMatcher", func() {
	var m *EqJSONMatcher
	JustBeforeEach(func() {
		m = &EqJSONMatcher{
			expected:       model.Map{"code": json.Number("200"), "messages": model.Seq{"hello"}},
			expectedString: `{"code":200,"messages":["hello"]}`,
		}
	})

	DescribeTable("Match",
		func(given string, expectedMatched bool, messageMatcher types.GomegaMatcher) {
			matched, message, err := m.Match([]byte(given))
			Expect(err).NotTo(HaveOccurred())
			Expect(matched).To(Equal(expectedMatched))
			Expect(message).To(messageMatcher)
		},
		Entry("when actual equals to expected as JSON, returns true", `{"messages": ["hello"], "code": 200}`, true, Equal(`should not be {"code":200,"messages":["hello"]}, but got it`)),
		Entry("when actual dose not equal to expected as JSON, returns false", `{"messages": ["hi"], "code": 200}`, false, Equal(`should be {"code":200,"messages":["hello"]}, but got {"messages": ["hi"], "code": 200}`)),
		Entry("when actual is not valid json, returns false", `{"messages": ["hi"], "code": 200`, false, HavePrefix("cannot recognize as json: ")),
	)
})

var _ = Describe("ParseEqJSONMatcher", func() {
	var v *model.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v, _ = model.NewValidator("", true)
		r = matcher.NewStreamMatcherRegistry()
	})

	Describe("with object", func() {
		It("returns matcher", func() {
			m := ParseEqJSONMatcher(v, r, model.Map{"code": 200, "messages": model.Seq{"hello"}})

			Expect(m).NotTo(BeNil())
			Expect(v.Error()).To(BeNil())

			var eq *EqJSONMatcher = m.(*EqJSONMatcher)
			Expect(eq.expected).To(Equal(model.Map{"code": json.Number("200"), "messages": model.Seq{"hello"}}))
		})
	})

	Describe("with json incompatible", func() {
		It("adds violation and returns nil", func() {
			m := ParseEqJSONMatcher(v, r, map[any]any{1: 42})

			Expect(m).To(BeNil())
			Expect(v.Error()).To(MatchError(HavePrefix("$: parameter is not json value: ")))
		})
	})
})
