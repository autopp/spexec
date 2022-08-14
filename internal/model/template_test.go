package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyExpandedValue = "dummyExpandedValue"

type dummyTemplateRef struct{}

func (dummyTemplateRef) Expand(env *Env, v *Validator, value any) (any, bool) {
	return dummyExpandedValue, true
}

var dummyError = "dummyError"

type errorTemplateRef struct{}

func (errorTemplateRef) Expand(env *Env, v *Validator, value any) (any, bool) {
	v.AddViolation(dummyError)
	return nil, false
}

var _ = Describe("TemplateVar", func() {
	Describe("Expand()", func() {
		It("returns bound value when var is defined", func() {
			env := NewEnv(nil)
			env.Define("answer", "42")
			v, _ := NewValidator("")

			tv := NewTemplateVar("answer")
			actual, ok := tv.Expand(env, v, Map{"$": "answer"})
			Expect(ok).To(BeTrue())
			Expect(actual).To(Equal("42"))
		})

		It("returns error when var is not defined", func() {
			env := NewEnv(nil)
			v, _ := NewValidator("")

			tv := NewTemplateVar("answer")
			_, ok := tv.Expand(env, v, Map{"$": "answer"})
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateFieldRef", func() {
	Describe("Expand()", func() {
		var tf *TemplateFieldRef
		field := "answer"
		var env *Env
		var v *Validator

		JustBeforeEach(func() {
			tf = NewTemplateFieldRef(field, dummyTemplateRef{})
			env = NewEnv(nil)
			v, _ = NewValidator("")
		})

		It("returns expanded value when given contains the field", func() {
			given := Map{"answer": Map{"$": "answer"}}

			actual, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeTrue())
			Expect(actual).To(Equal(Map{"answer": dummyExpandedValue}))
		})

		It("returns error when given dose not contain the field", func() {
			given := Map{"notAnswer": Map{"$": "answer"}}

			_, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(HaveOccurred())
		})

		It("returns error when given is not map", func() {
			given := Seq{Map{"$": "answer"}}

			_, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateIndexRef", func() {
	Describe("Expand()", func() {
		var tf *TemplateIndexRef
		var env *Env
		var v *Validator

		JustBeforeEach(func() {
			tf = NewTemplateIndexRef(1, dummyTemplateRef{})
			env = NewEnv(nil)
			v, _ = NewValidator("")
		})

		It("returns expanded value when given contains the element", func() {
			given := Seq{42, Map{"$": "answer"}}

			actual, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeTrue())
			Expect(actual).To(Equal(Seq{42, dummyExpandedValue}))
		})

		It("returns error when given dose not contain the element", func() {
			given := Seq{42}

			_, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(HaveOccurred())
		})

		It("returns error when given is not seq", func() {
			given := Map{"answer": Map{"$": "answer"}}

			_, ok := tf.Expand(env, v, given)
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateValue", func() {
	Describe("Expand()", func() {
		It("returns expanded value", func() {
			tv := NewTemplateValue(
				Map{"foo": Map{"$": "x"}, "bar": Seq{Map{"$": "y"}}},
				[]TemplateRef{
					NewTemplateFieldRef("foo", NewTemplateVar("x")),
					NewTemplateFieldRef("bar", NewTemplateIndexRef(0, NewTemplateVar("y"))),
				},
			)
			env := NewEnv(nil)
			env.Define("x", "hello")
			env.Define("y", "bye")

			actual, err := tv.Expand(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(Map{"foo": "hello", "bar": Seq{"bye"}}))
			Expect(actual).NotTo(Equal(tv.value))
		})

		It("propagate occured error in TemplateRef", func() {
			tv := NewTemplateValue(
				Map{"foo": Map{"$": "x"}, "bar": Seq{Map{"$": "y"}}},
				[]TemplateRef{
					NewTemplateFieldRef("foo", &errorTemplateRef{}),
				},
			)

			_, err := tv.Expand(NewEnv(nil))
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Templatable", func() {
	Describe("Expand()", func() {
		var env *Env
		BeforeEach(func() {
			env = NewEnv(nil)
		})

		It("returns wrapped value, when with simple value", func() {
			t := NewTemplatableFromValue("hello")

			Expect(t.Expand(env)).To(Equal("hello"))
		})

		It("returns expanded value, when with template value", func() {
			t := NewTemplatableFromTemplateValue[string](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{dummyTemplateRef{}}),
			)

			Expect(t.Expand(env)).To(Equal(dummyExpandedValue))
		})

		It("returns error, when with wrong type value", func() {
			t := NewTemplatableFromTemplateValue[int](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{dummyTemplateRef{}}),
			)

			_, err := t.Expand(env)
			Expect(err).To(MatchError("expect int, but got string"))
		})

		It("returns error, when with wrong type value", func() {
			t := NewTemplatableFromTemplateValue[string](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{errorTemplateRef{}}),
			)

			_, err := t.Expand(env)
			Expect(err).To(HaveOccurred())
		})
	})
})
