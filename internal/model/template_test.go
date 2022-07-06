package model

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyExpandedValue = "dummyExpandedValue"

type dummyTemplateRef struct{}

func (dummyTemplateRef) Expand(value any, env *Env) (any, error) {
	return dummyExpandedValue, nil
}

var dummyError = errors.New("dummyError")

type errorTemplateRef struct{}

func (errorTemplateRef) Expand(value any, env *Env) (any, error) {
	return nil, dummyError
}

var _ = Describe("TemplateVar", func() {
	Describe("Expand()", func() {
		It("returns bound value when var is defined", func() {
			env := NewEnv(nil)
			env.Define("answer", "42")

			tv := NewTemplateVar("answer")
			Expect(tv.Expand(Map{"$": "answer"}, env)).To(Equal("42"))
		})

		It("returns error when var is not defined", func() {
			env := NewEnv(nil)

			tv := NewTemplateVar("answer")
			_, err := tv.Expand(Map{"$": "answer"}, env)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateFieldRef", func() {
	Describe("Expand()", func() {
		var tf *TemplateFieldRef
		field := "answer"
		var env *Env

		JustBeforeEach(func() {
			tf = NewTemplateFieldRef(field, dummyTemplateRef{})
			env = NewEnv(nil)
		})

		It("returns expanded value when given contains the field", func() {
			given := Map{"answer": Map{"$": "answer"}}

			Expect(tf.Expand(given, env)).To(Equal(Map{"answer": dummyExpandedValue}))
		})

		It("returns error when given dose not contain the field", func() {
			given := Map{"notAnswer": Map{"$": "answer"}}

			_, err := tf.Expand(given, env)
			Expect(err).To(HaveOccurred())
		})

		It("returns error when given is not map", func() {
			given := Seq{Map{"$": "answer"}}

			_, err := tf.Expand(given, env)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("TemplateIndexRef", func() {
	Describe("Expand()", func() {
		var tf *TemplateIndexRef
		var env *Env

		JustBeforeEach(func() {
			tf = NewTemplateIndexRef(1, dummyTemplateRef{})
			env = NewEnv(nil)
		})

		It("returns expanded value when given contains the element", func() {
			given := Seq{42, Map{"$": "answer"}}

			Expect(tf.Expand(given, env)).To(Equal(Seq{42, dummyExpandedValue}))
		})

		It("returns error when given dose not contain the element", func() {
			given := Seq{42}

			_, err := tf.Expand(given, env)
			Expect(err).To(HaveOccurred())
		})

		It("returns error when given is not seq", func() {
			given := Map{"answer": Map{"$": "answer"}}

			_, err := tf.Expand(given, env)
			Expect(err).To(HaveOccurred())
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
			Expect(err).To(MatchError(dummyError))
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
			Expect(err).To(MatchError(dummyError))
		})
	})
})
