package matcher

import (
	"fmt"

	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type zeroMatcher struct {
	expected bool
}

func (m *zeroMatcher) Match(actual int) (bool, string, error) {
	if actual == 0 {
		if m.expected {
			return true, "status shoud not be zero", nil
		} else {
			return false, "status should be zero", nil
		}
	}

	if m.expected {
		return false, "status should be zero", nil
	}
	return true, "status shoud not be zero", nil
}

func parseZeroMatcher(_ *model.Env, _ *spec.Validator, r *StatusMatcherRegistry, x interface{}) model.StatusMatcher {
	return &zeroMatcher{expected: x.(bool)}
}

const violationMessage = "syntax error"

func parseViolationMatcher(_ *model.Env, v *spec.Validator, _ *StatusMatcherRegistry, _ interface{}) model.StatusMatcher {
	v.AddViolation(violationMessage)
	return nil
}

var _ = Describe("MatcherRegistry", func() {
	var r *MatcherParserRegistry[int]
	zeroName := "zero"
	zeroWithDefaultName := "zeroWithDefault"

	JustBeforeEach(func() {
		r = newMatcherParserRegistry[int]("int")
	})

	Describe("Add()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.Add(zeroName, parseZeroMatcher)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.Add(zeroName, parseZeroMatcher)
				err := r.Add(zeroName, parseZeroMatcher)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("AddWithDefault()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.AddWithDefault(zeroName, parseZeroMatcher, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.AddWithDefault(zeroName, parseZeroMatcher, true)
				err := r.AddWithDefault(zeroName, parseZeroMatcher, true)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ParseMatcher(env, )", func() {
		var v *spec.Validator
		zeroWithDefaultName := "zeroWithDefault"
		violationName := "violation"
		var env *model.Env

		JustBeforeEach(func() {
			r.Add(zeroName, parseZeroMatcher)
			r.AddWithDefault(zeroWithDefaultName, parseZeroMatcher, true)
			r.Add(violationName, parseViolationMatcher)
			v, _ = spec.NewValidator("")
			env = model.NewEnv(nil)
		})

		Context("for matcher without default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(env, v, model.Map{zeroName: true})

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("adds violation", func() {
					r.ParseMatcher(env, v, zeroName)

					Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: parameter is required", zeroName)))
				})
			})
		})

		Context("for matcher with default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(env, v, model.Map{zeroWithDefaultName: false})

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(env, v, zeroWithDefaultName)

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(m.(*zeroMatcher).expected).To(Equal(true))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				r.ParseMatcher(env, v, model.Map{violationName: nil})
				Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: %s", violationName, violationMessage)))
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(env, v, model.Map{"unknown": nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when size of the given map is not one", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(env, v, model.Map{zeroName: nil, violationName: nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a map and string", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(env, v, 42)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})

	Describe("ParseMatchers", func() {
		var v *spec.Validator
		violationName := "violation"
		var env *model.Env

		JustBeforeEach(func() {
			r.Add(zeroName, parseZeroMatcher)
			r.AddWithDefault(zeroWithDefaultName, parseZeroMatcher, true)
			r.Add(violationName, parseViolationMatcher)
			v, _ = spec.NewValidator("")
			env = model.NewEnv(nil)
		})

		Context("when params are valid", func() {
			It("returns the parsed matchers", func() {
				m := r.ParseMatchers(env, v, model.Seq{model.Map{zeroName: true}, zeroWithDefaultName})

				Expect(m[0]).To(BeAssignableToTypeOf(&zeroMatcher{}))
				Expect(m[1]).To(BeAssignableToTypeOf(&zeroMatcher{}))
				Expect(v.Error()).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatchers(env, v, model.Seq{model.Map{"unknown": false}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				m := r.ParseMatchers(env, v, model.Seq{model.Map{violationName: nil}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a seq", func() {
			It("adds violations", func() {
				m := r.ParseMatchers(env, v, 42)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("NewStatusMatcherRegistry()", func() {
	It("returns new registry", func() {
		Expect(NewStatusMatcherRegistry()).NotTo(BeNil())
	})
})

var _ = Describe("NewStreamMatcherRegistry()", func() {
	It("returns new registry", func() {
		Expect(NewStreamMatcherRegistry()).NotTo(BeNil())
	})
})
