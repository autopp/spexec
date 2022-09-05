package model

import (
	"github.com/autopp/spexec/internal/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyExpandedValue = "dummyExpandedValue"

type dummyTemplateRef struct{}

func (dummyTemplateRef) Expand(env *Env, v *Validator, value any) (any, bool, error) {
	return dummyExpandedValue, true, nil
}

var dummyError = "dummyError"

type errorTemplateRef struct{}

func (errorTemplateRef) Expand(env *Env, v *Validator, value any) (any, bool, error) {
	return nil, false, errors.New(errors.ErrInternalError, dummyError)
}

var _ = Describe("TemplateVar", func() {
	Describe("Expand()", func() {
		It("returns bound value when var is defined", func() {
			env := NewEnv(nil)
			env.Define("answer", "42")
			v, _ := NewValidator("")

			tv := NewTemplateVar("answer")
			actual, ok, err := tv.Expand(env, v, Map{"$": "answer"})
			Expect(ok).To(BeTrue())
			Expect(v.Error()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal("42"))
		})

		It("returns error when var is not defined", func() {
			env := NewEnv(nil)
			v, _ := NewValidator("")

			tv := NewTemplateVar("answer")
			_, ok, err := tv.Expand(env, v, Map{"$": "answer"})
			Expect(ok).To(BeFalse())
			Expect(v.Error()).To(BeValidationError("$.$answer: is not defined"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateFieldRef", func() {
	Describe("Expand()", func() {
		field := "answer"
		var env *Env
		var v *Validator

		JustBeforeEach(func() {

			env = NewEnv(nil)
			v, _ = NewValidator("")
		})

		It("returns expanded value when given contains the field", func() {
			tf := NewTemplateFieldRef(field, dummyTemplateRef{})
			given := Map{"answer": Map{"$": "answer"}}

			actual, ok, err := tf.Expand(env, v, given)
			Expect(ok).To(BeTrue())
			Expect(v.Error()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(Map{"answer": dummyExpandedValue}))
		})

		It("returns error when given dose not contain the field", func() {
			tf := NewTemplateFieldRef(field, dummyTemplateRef{})
			given := Map{"notAnswer": Map{"$": "answer"}}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(BeValidationError("unexpected template structure: $: should have .answer"))
		})

		It("returns error when given is not map", func() {
			tf := NewTemplateFieldRef(field, dummyTemplateRef{})
			given := Seq{Map{"$": "answer"}}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(BeValidationError("unexpected template structure: $: should be map, but is seq"))
		})

		It("returns error when it occured in nested field", func() {
			tf := NewTemplateFieldRef(field, errorTemplateRef{})
			given := Map{"answer": Map{"$": "answer"}}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(MatchError(dummyError))
		})
	})
})

var _ = Describe("TemplateIndexRef", func() {
	Describe("Expand()", func() {
		var env *Env
		var v *Validator

		JustBeforeEach(func() {
			env = NewEnv(nil)
			v, _ = NewValidator("")
		})

		It("returns expanded value when given contains the element", func() {
			tf := NewTemplateIndexRef(1, dummyTemplateRef{})
			given := Seq{42, Map{"$": "answer"}}

			actual, ok, err := tf.Expand(env, v, given)
			Expect(ok).To(BeTrue())
			Expect(v.Error()).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(Seq{42, dummyExpandedValue}))
		})

		It("returns error when given dose not contain the element", func() {
			tf := NewTemplateIndexRef(1, dummyTemplateRef{})
			given := Seq{42}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(BeValidationError("unexpected template structure: $: should have 2 items"))
		})

		It("returns error when given is not seq", func() {
			tf := NewTemplateIndexRef(1, dummyTemplateRef{})
			given := Map{"answer": Map{"$": "answer"}}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(BeValidationError("unexpected template structure: $: should be seq, but is map"))
		})

		It("returns error when it occured in the element", func() {
			tf := NewTemplateIndexRef(1, errorTemplateRef{})
			given := Seq{42, Map{"$": "answer"}}

			_, _, err := tf.Expand(env, v, given)
			Expect(err).To(MatchError(dummyError))
		})
	})
})

var _ = Describe("TemplateValue", func() {
	Describe("Expand()", func() {
		var v *Validator
		JustBeforeEach(func() {
			v, _ = NewValidator("")
		})

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

			actual, err := tv.Expand(env, v)
			Expect(err).NotTo(HaveOccurred())
			Expect(v.Error()).NotTo(HaveOccurred())
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

			_, err := tv.Expand(NewEnv(nil), v)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Templatable", func() {
	Describe("Expand()", func() {
		var env *Env
		var v *Validator
		BeforeEach(func() {
			env = NewEnv(nil)
			v, _ = NewValidator("")
		})

		It("returns wrapped value, when with simple value", func() {
			t := NewTemplatableFromValue("hello")

			Expect(t.Expand(env, v)).To(Equal("hello"))
			Expect(v.Error()).NotTo(HaveOccurred())
		})

		It("returns expanded value, when with template value", func() {
			t := NewTemplatableFromTemplateValue[string](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{dummyTemplateRef{}}),
			)

			Expect(t.Expand(env, v)).To(Equal(dummyExpandedValue))
			Expect(v.Error()).NotTo(HaveOccurred())
		})

		It("returns error, when with wrong type value", func() {
			t := NewTemplatableFromTemplateValue[int](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{dummyTemplateRef{}}),
			)

			_, err := t.Expand(env, v)
			Expect(err).To(HaveOccurred())
		})

		It("returns error, when with wrong type value", func() {
			t := NewTemplatableFromTemplateValue[string](
				NewTemplateValue(Map{"$": "x"}, []TemplateRef{errorTemplateRef{}}),
			)

			_, err := t.Expand(env, v)
			Expect(err).To(HaveOccurred())
		})
	})
})
