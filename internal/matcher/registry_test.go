package matcher

import (
	"errors"

	"github.com/autopp/spexec/internal/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type zeroMatcher struct{}

func (*zeroMatcher) MatchStatus(actual int) (bool, string, error) {
	if actual == 0 {
		return true, "status shoud not be zero", nil
	}
	return false, "status should be zero", nil
}

func parseZeroMatcher(_ *spec.Validator, r *StatusMatcherRegistry, x interface{}) (StatusMatcher, error) {
	return &zeroMatcher{}, nil
}

const violationMessage = "syntax error"

func parseViolationStatusMatcher(v *spec.Validator, _ *StatusMatcherRegistry, _ interface{}) (StatusMatcher, error) {
	v.AddViolation(violationMessage)
	return &zeroMatcher{}, nil
}

const errMessage = "some error"

func parseErrStatusMatcher(_ *spec.Validator, _ *StatusMatcherRegistry, _ interface{}) (StatusMatcher, error) {
	return nil, errors.New(errMessage)
}

type emptyMatcher struct{}

func (*emptyMatcher) MatchStream(actual []byte) (bool, string, error) {
	if len(actual) == 0 {
		return true, "stream should not be empty", nil
	}

	return false, "status should be empty", nil
}

func parseEmptyMatcher(_ *spec.Validator, r *StreamMatcherRegistry, fd int, x interface{}) (StreamMatcher, error) {
	return &emptyMatcher{}, nil
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

	Describe("Get()", func() {
		Context("when not registered", func() {
			It("returns error", func() {
				_, err := r.Get(zeroName)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when registered", func() {
			It("returns registerd parser", func() {
				r.Add(zeroName, parseZeroMatcher)

				callParser := func(m StatusMatcherParser) StatusMatcher {
					p, _ := m(nil, r, nil)
					return p
				}

				Expect(r.Get(zeroName)).To(WithTransform(callParser, BeAssignableToTypeOf(&zeroMatcher{})))
			})
		})
	})

	Describe("ParseMatcher", func() {
		var v *spec.Validator
		violationName := "violation"
		errName := "error"

		JustBeforeEach(func() {
			r.Add(zeroName, parseZeroMatcher)
			r.Add(violationName, parseViolationStatusMatcher)
			r.Add(errName, parseErrStatusMatcher)
			v = spec.NewValidator()
		})

		Context("when the given name is registered and it returns matcher", func() {
			It("returns the parsed matcher", func() {
				m, err := r.ParseMatcher(v, spec.Map{zeroName: nil})
				Expect(err).NotTo(HaveOccurred())

				Expect(m).To(BeAssignableToTypeOf(&zeroMatcher{}))
				Expect(v.Error()).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered and it adds violations", func() {
			It("cascades violations", func() {
				_, err := r.ParseMatcher(v, spec.Map{violationName: nil})
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Error()).To(MatchError(ContainSubstring(violationMessage)))
			})
		})

		Context("when the given name is registered and it returns error", func() {
			It("returns error", func() {
				_, err := r.ParseMatcher(v, spec.Map{errName: nil})
				Expect(err).To(MatchError(errMessage))
			})
		})

		Context("when the given name is not registered", func() {
			It("adds violations", func() {
				_, err := r.ParseMatcher(v, spec.Map{"unknown": nil})
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when size of the given map is not one", func() {
			It("adds violations", func() {
				_, err := r.ParseMatcher(v, spec.Map{zeroName: nil, violationName: nil})
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Error()).To(HaveOccurred())
			})
		})

		Context("when the given is not a map", func() {
			It("adds violations", func() {
				_, err := r.ParseMatcher(v, zeroName)
				Expect(err).NotTo(HaveOccurred())
				Expect(v.Error()).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("StreamMatcherRegistry", func() {
	var r *StreamMatcherRegistry
	name := "empty"

	JustBeforeEach(func() {
		r = NewStreamMatcherRegistry()
	})

	Describe("Add()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.Add(name, parseEmptyMatcher)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.Add(name, parseEmptyMatcher)
				err := r.Add(name, parseEmptyMatcher)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Get()", func() {
		Context("when not registered", func() {
			It("returns error", func() {
				_, err := r.Get(name)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when registered", func() {
			It("returns registerd parser", func() {
				r.Add(name, parseEmptyMatcher)

				callParser := func(p StreamMatcherParser) StreamMatcher {
					m, _ := p(nil, r, 1, []byte{})
					return m
				}
				Expect(r.Get(name)).To(WithTransform(callParser, BeAssignableToTypeOf(&emptyMatcher{})))
			})
		})
	})
})
