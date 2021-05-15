package matcher

import (
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

func parseZeroMatcher(r *StatusMatcherRegistry, x interface{}) (StatusMatcher, error) {
	return &zeroMatcher{}, nil
}

var _ = Describe("StatusMatcherRegistry", func() {
	var r *StatusMatcherRegistry
	name := "zero"

	JustBeforeEach(func() {
		r = NewStatusMatcherRegistry()
	})

	Describe("Add()", func() {
		Context("when the given name is not registered yet", func() {
			It("returns nil", func() {
				err := r.Add(name, parseZeroMatcher)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the given name is registered already", func() {
			It("returns error", func() {
				r.Add(name, parseZeroMatcher)
				err := r.Add(name, parseZeroMatcher)
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
				r.Add(name, parseZeroMatcher)

				callParser := func(m StatusMatcherParser) StatusMatcher {
					p, _ := m(r, nil)
					return p
				}
				Expect(r.Get(name)).To(WithTransform(callParser, BeAssignableToTypeOf(&zeroMatcher{})))
			})
		})
	})
})
