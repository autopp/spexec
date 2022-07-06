package stream

import (
	"time"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
	"github.com/autopp/spexec/internal/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SatisfyMatcher", func() {
	var m *SatisfyMatcher
	JustBeforeEach(func() {
		m = &SatisfyMatcher{
			Command: []string{"bash", "-c", `test "$(cat -)" = hello`},
			Cleanup: func() []error { return nil },
		}
	})

	DescribeTable("Match",
		func(given string, expectedMatched bool, expectedMessage string) {
			matched, message, err := m.Match([]byte(given))
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
	var env *model.Env

	JustBeforeEach(func() {
		v, _ = spec.NewValidator("")
		r = matcher.NewStreamMatcherRegistry()
		env = model.NewEnv(nil)
	})

	DescribeTable("success cases",
		func(given any, expectedCommand []string, expectedEnv []util.StringVar, expectedTimeout time.Duration) {
			m := ParseSatisfyMatcher(env, v, r, given)

			Expect(v.Error()).To(BeNil())
			Expect(m).NotTo(BeNil())

			var satisfyMatcher *SatisfyMatcher
			Expect(m).To(BeAssignableToTypeOf(satisfyMatcher))
			satisfyMatcher = m.(*SatisfyMatcher)
			Expect(satisfyMatcher.Command).To(Equal(expectedCommand))
			Expect(satisfyMatcher.Env).To(Equal(expectedEnv))
			Expect(satisfyMatcher.Timeout).To(Equal(expectedTimeout))
		},
		Entry("with full field",
			model.Map{
				"command": model.Seq{"test.sh"},
				"env":     model.Seq{model.Map{"name": "MSG", "value": "hello"}},
				"timeout": 1,
			},
			[]string{"test.sh"}, []util.StringVar{{Name: "MSG", Value: "hello"}}, 1*time.Second,
		),
		Entry("without env",
			model.Map{
				"command": model.Seq{"test.sh"},
				"timeout": 1,
			},
			[]string{"test.sh"}, nil, 1*time.Second,
		),
		Entry("without timeout",
			model.Map{
				"command": model.Seq{"test.sh"},
				"env":     model.Seq{model.Map{"name": "MSG", "value": "hello"}},
			},
			[]string{"test.sh"}, []util.StringVar{{Name: "MSG", Value: "hello"}}, 5*time.Second,
		),
	)

	DescribeTable("failure cases",
		func(given any, prefix string) {
			m := ParseSatisfyMatcher(env, v, r, given)

			Expect(m).To(BeNil())
			err := v.Error()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(prefix))
		},
		Entry("with not map", 42, "$: "),
		Entry("without command", model.Map{
			"env":     model.Seq{model.Map{"name": "MSG", "value": "hello"}},
			"timeout": 1,
		}, "$: "),
		Entry("with invalid command", model.Map{
			"command": "test.sh",
		}, "$.command: "),
		Entry("with empty command", model.Map{
			"command": model.Seq{},
		}, "$.command: "),
		Entry("with invalid env", model.Map{
			"command": model.Seq{"test.sh"},
			"env":     42,
		}, "$.env: "),
		Entry("with invalid timeout", model.Map{
			"command": model.Seq{"test.sh"},
			"timeout": "foo",
		}, "$.timeout: "),
		Entry("with invalid string expr", model.Map{
			"command": model.Seq{model.Map{"type": "env", "name": "unknown"}},
		}, "$.command: error occured at parsing command"),
	)
})
