package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyExpandedValue = "dummyExpandedValue"

type dummyTemplateRef struct{}

func (dummyTemplateRef) Expand(value interface{}, env *Env) (interface{}, error) {
	return dummyExpandedValue, nil
}

var _ = Describe("TemplateVar", func() {
	Describe("Expand()", func() {
		It("returns bound value when var is defined", func() {
			env := NewEnv(nil)
			env.Define("answer", "42")

			tv := &TemplateVar{name: "answer"}
			Expect(tv.Expand(Map{"$": "answer"}, env)).To(Equal("42"))
		})

		It("returns error when var is not defined", func() {
			env := NewEnv(nil)

			tv := &TemplateVar{name: "answer"}
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
			tf = &TemplateFieldRef{field: field, next: dummyTemplateRef{}}
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
			tf = &TemplateIndexRef{index: 1, next: dummyTemplateRef{}}
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
