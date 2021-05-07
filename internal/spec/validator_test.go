package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Validator", func() {
	var v *Validator

	BeValidationError := func(message string) types.GomegaMatcher {
		return And(HaveOccurred(), WithTransform(func(err error) string { return err.Error() }, Equal(message)))
	}

	JustBeforeEach(func() {
		v = NewValidator()
	})

	Describe("AddViolation() and Error()", func() {
		Context("with no AddViolation() call", func() {
			It("makes to Error() to return nil", func() {
				Expect(v.Error()).To(Succeed())
			})
		})

		Context("with a AddViolation() call", func() {
			It("makes to Error() to returns error which contains path and formatted violation message", func() {
				v.AddViolation("answer is %d", 42)
				Expect(v.Error()).To(BeValidationError("$: answer is 42"))
			})
		})

		Context("with multiple AddViolation() calls", func() {
			It("make to Error() to return error which contains given violations joined by newline", func() {
				v.AddViolation("error1")
				v.AddViolation("error2")
				Expect(v.Error()).To(BeValidationError("$: error1\n$: error2"))
			})
		})
	})

	Describe("InPath()", func() {
		It("appends path prefix in callback", func() {
			v.InPath(":prefix1", func() {
				v.AddViolation("error1")
				v.InPath(":prefix2", func() {
					v.AddViolation("error2")
				})
			})
			v.AddViolation("error3")

			Expect(v.Error()).To(BeValidationError("$:prefix1: error1\n$:prefix1:prefix2: error2\n$: error3"))
		})
	})

	Describe("InField()", func() {
		It(`be equivalent to InPath("."+field)`, func() {
			v.InField("field", func() {
				v.AddViolation("error")
			})

			Expect(v.Error()).To(BeValidationError("$.field: error"))
		})
	})

	Describe("InIndex()", func() {
		It(`be equivalent to InPath("["+i+"]")`, func() {
			v.InIndex(1, func() {
				v.AddViolation("error")
			})

			Expect(v.Error()).To(BeValidationError("$[1]: error"))
		})
	})

	Describe("MustBeMap()", func() {
		Context("with a Map", func() {
			It("returns the given Map and true", func() {
				given := make(Map)
				m, b := v.MustBeMap(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not Map", func() {
			It("returns something and false", func() {
				_, b := v.MustBeMap(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeSeq()", func() {
		Context("with a Seq", func() {
			It("returns the given Seq and true", func() {
				given := make(Seq, 0)
				m, b := v.MustBeSeq(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not Seq", func() {
			It("returns something and false", func() {
				_, b := v.MustBeSeq(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeString()", func() {
		Context("with a string", func() {
			It("returns the given string and true", func() {
				given := "hello"
				m, b := v.MustBeString(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not string", func() {
			It("returns something and false", func() {
				_, b := v.MustBeString(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeInt()", func() {
		Context("with a int", func() {
			It("returns the given int and true", func() {
				given := 42
				m, b := v.MustBeInt(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not int", func() {
			It("returns something and false", func() {
				_, b := v.MustBeInt("hello")

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MayHaveMap()", func() {
		Context("when the given map has specified field which is a Map", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true, true", func() {
				contained := make(Map)
				var passed Map
				m, exists, ok := v.MayHaveMap(Map{"answer": contained}, "answer", func(m Map) {
					passed = m
					v.AddViolation("error")
				})

				Expect(passed).To(Equal(contained))
				Expect(v.Error()).To(BeValidationError("$.answer: error"))
				Expect(m).To(Equal(contained))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("dose not call the callback and returns something, false, true", func() {
				_, exists, ok := v.MayHaveMap(make(Map), "answer", func(Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a Map", func() {
			It("dose not call the callback and returns something, false, false", func() {
				_, exists, ok := v.MayHaveMap(Map{"answer": "hello"}, "answer", func(m Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
