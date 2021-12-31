package model

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("literalStringExpr", func() {
	literal := NewLiteralStringExpr("hello")

	Describe("String()", func() {
		It("returns itself", func() {
			Expect(literal.String()).To(Equal("hello"))
		})
	})

	Describe("Eval()", func() {
		It("returns itself", func() {
			Expect(literal.Eval()).To(Equal("hello"))
		})
	})
})

var _ = Describe("envStringExpr", func() {
	name := "MESSAGE"
	value := "hello"
	os.Setenv(name, value)

	env := NewEnvStringExpr(name)

	Describe("String()", func() {
		It("returns itself with '$' prefix", func() {
			Expect(env.String()).To(Equal("$" + name))
		})
	})

	Describe("Eval()", func() {
		It("returns value of the environment variable", func() {
			Expect(env.Eval()).To(Equal(value))
		})

		It("returns error when given name is not defined", func() {
			_, err := NewEnvStringExpr("SPEXEC_UNDEFINED").Eval()
			Expect(err).To(HaveOccurred())
		})
	})
})
