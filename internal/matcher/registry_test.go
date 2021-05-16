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

type emptyMatcher struct{}

func (*emptyMatcher) MatchStream(actual []byte) (bool, string, error) {
	if len(actual) == 0 {
		return true, "stream should not be empty", nil
	}

	return false, "status should be empty", nil
}

func parseEmptyMatcher(r *StreamMatcherRegistry, fd int, x interface{}) (StreamMatcher, error) {
	return &emptyMatcher{}, nil
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
					m, _ := p(r, 1, []byte{})
					return m
				}
				Expect(r.Get(name)).To(WithTransform(callParser, BeAssignableToTypeOf(&emptyMatcher{})))
			})
		})
	})
})
