package matcher

import (
	"fmt"

	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo"
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

func parseZeroMatcher(_ *spec.Validator, r *StatusMatcherRegistry, x interface{}) StatusMatcher {
	return &zeroMatcher{expected: x.(bool)}
}

const violationMessage = "syntax error"

func parseViolationStatusMatcher(v *spec.Validator, _ *StatusMatcherRegistry, _ interface{}) StatusMatcher {
	v.AddViolation(violationMessage)
	return &zeroMatcher{}
}

type emptyMatcher struct{}

func (*emptyMatcher) MatchStream(actual []byte) (bool, string, error) {
	if len(actual) == 0 {
		return true, "stream should not be empty", nil
	}

	return false, "status should be empty", nil
}

func parseEmptyMatcher(_ *spec.Validator, r *StreamMatcherRegistry, x interface{}) StreamMatcher {
	return &emptyMatcher{}
}

func parseViolationStreamMatcher(v *spec.Validator, _ *StreamMatcherRegistry, x interface{}) StreamMatcher {
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
			v = spec.NewValidator()
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

		Context("when the given is not a map", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, zeroName)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("StreamMatcherRegistry", func() {
	var r *StreamMatcherRegistry
	emptyName := "empty"

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

	Describe("ParseMatcher", func() {
		var v *spec.Validator
		violationName := "violation"

		JustBeforeEach(func() {
			r.Add(emptyName, parseEmptyMatcher)
			r.Add(violationName, parseViolationStreamMatcher)
			v = spec.NewValidator()
		})

		Context("when the given name is registered and it returns matcher", func() {
			It("returns the parsed matcher", func() {
				m := r.ParseMatcher(v, spec.Map{emptyName: nil})

				Expect(m).To(BeAssignableToTypeOf(&emptyMatcher{}))
				Expect(v.Error()).NotTo(HaveOccurred())
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
				m := r.ParseMatcher(v, spec.Map{emptyName: nil, violationName: nil})
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a map", func() {
			It("adds violations", func() {
				m := r.ParseMatcher(v, emptyName)
				Expect(m).To(BeNil())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})
