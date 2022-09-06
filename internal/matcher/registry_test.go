package matcher_test

import (
	"fmt"

	"github.com/autopp/spexec/internal/matcher"
	"github.com/autopp/spexec/internal/matcher/testutil"
	"github.com/autopp/spexec/internal/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const violationMessage = "syntax error"

var _ = Describe("MatcherRegistry", func() {
	var r *matcher.MatcherParserRegistry[int]
	var parseExampleMatcher matcher.MatcherParser[int]
	name := "example"
	withDefaultName := "exampleWithDefault"

	JustBeforeEach(func() {
		r = matcher.NewMatcherParserRegistry[int]("int")
		parseExampleMatcher, _ = testutil.GenParseExampleStatusMatcher(true, "matcher message", nil)
	})

	Describe("Add()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.Add(name, parseExampleMatcher)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.Add(name, parseExampleMatcher)
				err := r.Add(name, parseExampleMatcher)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("AddWithDefault()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.AddWithDefault(withDefaultName, parseExampleMatcher, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.AddWithDefault(withDefaultName, parseExampleMatcher, true)
				err := r.AddWithDefault(withDefaultName, parseExampleMatcher, true)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ParseMatcher()", func() {
		var v *model.Validator
		var parserCalls *testutil.ParserCalls
		var parseExampleMatcher matcher.MatcherParser[int]
		var parseExampleMatcherWithDefault matcher.MatcherParser[int]
		var parserWithDefaultCalls *testutil.ParserCalls
		var failedParseMatcher matcher.MatcherParser[int]

		name := "example"
		withDefaultName := "exampleWithDefault"
		failedName := "violation"

		JustBeforeEach(func() {
			parseExampleMatcher, parserCalls = testutil.GenParseExampleStatusMatcher(true, "matcher message", nil)
			r.Add(name, parseExampleMatcher)

			parseExampleMatcherWithDefault, parserWithDefaultCalls = testutil.GenParseExampleStatusMatcher(true, "matcher message", nil)
			r.AddWithDefault(withDefaultName, parseExampleMatcherWithDefault, 42)

			failedParseMatcher = testutil.GenFailedParseStatusMatcher(violationMessage)
			r.Add(failedName, failedParseMatcher)

			v, _ = model.NewValidator("", true)
		})

		Context("for matcher without default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, model.Map{name: 42})

					Expect(m).To(BeAssignableToTypeOf(&testutil.ExampleStatusMatcher{}))
					Expect(parserCalls.Calls).To(Equal([]any{42}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("adds violation", func() {
					r.ParseMatcher(v, name)

					Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: parameter is required", name)))
				})
			})
		})

		Context("for matcher with default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, model.Map{withDefaultName: false})

					Expect(m).To(BeAssignableToTypeOf(&testutil.ExampleStatusMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, withDefaultName)

					Expect(m).To(BeAssignableToTypeOf(&testutil.ExampleStatusMatcher{}))
					Expect(parserWithDefaultCalls.Calls).To(Equal([]any{42}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				r.ParseMatcher(v, model.Map{failedName: nil})
				Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: %s", failedName, violationMessage)))
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, model.Map{"unknown": nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when size of the given map is not one", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, model.Map{name: nil, failedName: nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a map and string", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, 42)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})

	Describe("ParseMatchers()", func() {
		var v *model.Validator
		var parseExampleMatcher matcher.MatcherParser[int]
		var parseExampleMatcherWithDefault matcher.MatcherParser[int]
		var failedParseMatcher matcher.MatcherParser[int]

		name := "example"
		withDefaultName := "exampleWithDefault"
		failedName := "violation"

		JustBeforeEach(func() {
			parseExampleMatcher, _ = testutil.GenParseExampleStatusMatcher(true, "matcher message", nil)
			r.Add(name, parseExampleMatcher)

			parseExampleMatcherWithDefault, _ = testutil.GenParseExampleStatusMatcher(true, "matcher message", nil)
			r.AddWithDefault(withDefaultName, parseExampleMatcherWithDefault, 42)

			failedParseMatcher = testutil.GenFailedParseStatusMatcher(violationMessage)
			r.Add(failedName, failedParseMatcher)

			v, _ = model.NewValidator("", true)
		})

		Context("when params are valid", func() {
			It("returns the parsed matchers", func() {
				m := r.ParseMatchers(v, model.Seq{model.Map{name: true}, withDefaultName})

				Expect(m[0]).To(BeAssignableToTypeOf(&testutil.ExampleStatusMatcher{}))
				Expect(m[1]).To(BeAssignableToTypeOf(&testutil.ExampleStatusMatcher{}))
				Expect(v.Error()).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatchers(v, model.Seq{model.Map{"unknown": false}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				m := r.ParseMatchers(v, model.Seq{model.Map{failedName: nil}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a seq", func() {
			It("adds violations", func() {
				m := r.ParseMatchers(v, 42)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("NewStatusMatcherRegistry()", func() {
	It("returns new registry", func() {
		Expect(matcher.NewStatusMatcherRegistry()).NotTo(BeNil())
	})
})

var _ = Describe("NewStreamMatcherRegistry()", func() {
	It("returns new registry", func() {
		Expect(matcher.NewStreamMatcherRegistry()).NotTo(BeNil())
	})
})
