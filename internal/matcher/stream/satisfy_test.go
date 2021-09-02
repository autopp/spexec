package stream

import (
	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/spec"
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

var _ = Describe("ParseSatisfyMatcher", func() {
	var v *spec.Validator
	var r *matcher.StreamMatcherRegistry

	JustBeforeEach(func() {
		v = spec.NewValidator()
		r = matcher.NewStreamMatcherRegistry()
	})

	DescribeTable("success cases",
		func(given interface{}) {
			m := ParseSatisfyMatcher(v, r, given)

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			var satisfyMatcher *SatisfyMatcher
			Expect(m).To(BeAssignableToTypeOf(satisfyMatcher))
		},
		Entry("with full field",
			spec.Map{
				"command": spec.Seq{"test.sh"},
				"env":     spec.Seq{spec.Map{"name": "MSG", "value": "hello"}},
				"timeout": 1,
			},
		),
		Entry("without env",
			spec.Map{
				"command": spec.Seq{"test.sh"},
				"timeout": 1,
			},
		),
		Entry("without timeout",
			spec.Map{
				"command": spec.Seq{"test.sh"},
				"env":     spec.Seq{spec.Map{"name": "MSG", "value": "hello"}},
			},
		),
	)

	DescribeTable("failure cases",
		func(given interface{}, prefix string) {
			m := ParseSatisfyMatcher(v, r, given)

			Expect(m).To(BeNil())
			err := v.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(prefix))
		},
		Entry("with not map", 42, "$: "),
		Entry("without command", spec.Map{
			"env":     spec.Seq{spec.Map{"name": "MSG", "value": "hello"}},
			"timeout": 1,
		}, "$: "),
		Entry("with invalid command", spec.Map{
			"command": "test.sh",
		}, "$.command: "),
		Entry("with empty command", spec.Map{
			"command": spec.Seq{},
		}, "$.command: "),
		Entry("with invalid env", spec.Map{
			"command": spec.Seq{"test.sh"},
			"env":     42,
		}, "$.env: "),
		Entry("with invalid timeout", spec.Map{
			"command": spec.Seq{"test.sh"},
			"timeout": "foo",
		}, "$.timeout: "),
	)
})