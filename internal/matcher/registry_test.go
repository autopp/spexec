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

func (m *zeroMatcher) MatchStatus(actual int) (bool, string, error) {
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

func parseZeroMatcher(_ *spec.Validator, r *StatusMatcherRegistry, x interface{}) model.StatusMatcher {
	return &zeroMatcher{expected: x.(bool)}
}

const violationMessage = "syntax error"

func parseViolationStatusMatcher(v *spec.Validator, _ *StatusMatcherRegistry, _ interface{}) model.StatusMatcher {
	v.AddViolation(violationMessage)
	return &zeroMatcher{}
}

type emptyMatcher struct {
	expected bool
}

func (m *emptyMatcher) MatchStream(actual []byte) (bool, string, error) {
	if len(actual) == 0 {
		if m.expected {
			return true, "stream should not be empty", nil
		}
		return false, "status should be empty", nil
	}

	if m.expected {
		return false, "status should be empty", nil
	}
	return true, "stream should not be empty", nil
}

func parseEmptyMatcher(_ *spec.Validator, r *StreamMatcherRegistry, x interface{}) model.StreamMatcher {
	return &emptyMatcher{expected: x.(bool)}
}

func parseViolationStreamMatcher(v *spec.Validator, _ *StreamMatcherRegistry, x interface{}) model.StreamMatcher {
	v.AddViolation(violationMessage)
	return &emptyMatcher{}
}

var _ = Describe("StatusMatcherRegistry", func() {
	var r *StatusMatcherRegistry
	zeroName := "zero"

	JustBeforeEach(func() {
		r = NewStatusMatcherRegistry()
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

	Describe("ParseMatcher()", func() {
		var v *spec.Validator
		zeroWithDefaultName := "zeroWithDefault"
		violationName := "violation"

		JustBeforeEach(func() {
			r.Add(zeroName, parseZeroMatcher)
			r.AddWithDefault(zeroWithDefaultName, parseZeroMatcher, true)
			r.Add(violationName, parseViolationStatusMatcher)
			v, _ = spec.NewValidator("")
		})

		Context("for matcher without default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, spec.Map{zeroName: true})

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("adds violation", func() {
					r.ParseMatcher(v, zeroName)

					Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: parameter is required", zeroName)))
				})
			})
		})

		Context("for matcher with default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, spec.Map{zeroWithDefaultName: false})

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, zeroWithDefaultName)

					Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
					Expect(m.(*zeroMatcher).expected).To(Equal(true))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				r.ParseMatcher(v, spec.Map{violationName: nil})
				Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: %s", violationName, violationMessage)))
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, spec.Map{"unknown": nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when size of the given map is not one", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, spec.Map{zeroName: nil, violationName: nil})
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
})

var _ = Describe("StreamMatcherRegistry", func() {
	var r *StreamMatcherRegistry
	emptyName := "empty"
	emptyWithDefaultName := "emptyWithDefault"

	JustBeforeEach(func() {
		r = NewStreamMatcherRegistry()
	})

	Describe("Add()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.Add(emptyName, parseEmptyMatcher)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.Add(emptyName, parseEmptyMatcher)
				err := r.Add(emptyName, parseEmptyMatcher)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("AddWithDefault()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.AddWithDefault(emptyName, parseEmptyMatcher, true)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.AddWithDefault(emptyName, parseEmptyMatcher, true)
				err := r.AddWithDefault(emptyName, parseEmptyMatcher, true)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ParseMatcher", func() {
		var v *spec.Validator
		violationName := "violation"

		JustBeforeEach(func() {
			r.Add(emptyName, parseEmptyMatcher)
			r.AddWithDefault(emptyWithDefaultName, parseEmptyMatcher, true)
			r.Add(violationName, parseViolationStreamMatcher)
			v, _ = spec.NewValidator("")
		})

		Context("for matcher without default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, spec.Map{emptyName: true})

					Expect(m).To(BeAssignableToTypeOf(&emptyMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("adds violation", func() {
					r.ParseMatcher(v, emptyName)

					Expect(v.Error()).To(MatchError(fmt.Sprintf("$.%s: parameter is required", emptyName)))
				})
			})
		})

		Context("for matcher with default parameter", func() {
			Context("when param is passed and it returns matcher", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, spec.Map{emptyWithDefaultName: false})

					Expect(m).To(BeAssignableToTypeOf(&emptyMatcher{}))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})

			Context("when param is not passed", func() {
				It("returns the parsed matcher", func() {
					m := r.ParseMatcher(v, emptyWithDefaultName)

					Expect(m).To(BeAssignableToTypeOf(&emptyMatcher{}))
					Expect(m.(*emptyMatcher).expected).To(Equal(true))
					Expect(v.Error()).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, spec.Map{"unknown": nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when size of the given map is not one", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, spec.Map{emptyName: nil, violationName: nil})
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

	Describe("ParseMatchers", func() {
		var v *spec.Validator
		violationName := "violation"

		JustBeforeEach(func() {
			r.Add(emptyName, parseEmptyMatcher)
			r.AddWithDefault(emptyWithDefaultName, parseEmptyMatcher, true)
			r.Add(violationName, parseViolationStreamMatcher)
			v, _ = spec.NewValidator("")
		})

		Context("when params are valid", func() {
			It("returns the parsed matchers", func() {
				m := r.ParseMatchers(v, spec.Seq{spec.Map{emptyName: true}, emptyWithDefaultName})

				Expect(m[0]).To(BeAssignableToTypeOf(&emptyMatcher{}))
				Expect(m[1]).To(BeAssignableToTypeOf(&emptyMatcher{}))
				Expect(v.Error()).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, spec.Seq{spec.Map{"unknown": false}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				m := r.ParseMatcher(v, spec.Seq{spec.Map{violationName: nil}})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a seq", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, 42)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})
