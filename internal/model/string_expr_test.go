package model

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
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
			v, cleanup, err := literal.Eval()
			Expect(v).To(Equal("hello"))
			Expect(cleanup).To(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("envStringExpr", func() {
	name := "MESSAGE"
	value := "hello"

	env := NewEnvStringExpr(name)

	BeforeEach(func() {
		oldValue, setAlready := os.LookupEnv(name)
		os.Setenv(name, value)

		DeferCleanup(func() {
			if setAlready {
				os.Setenv(name, oldValue)
			} else {
				os.Unsetenv(name)
			}
		})
	})

	Describe("String()", func() {
		It("returns itself with '$' prefix", func() {
			Expect(env.String()).To(Equal("$" + name))
		})
	})

	Describe("Eval()", func() {
		It("returns value of the environment variable", func() {
			v, cleanup, err := env.Eval()
			Expect(v).To(Equal(value))
			Expect(cleanup).To(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when given name is not defined", func() {
			_, _, err := NewEnvStringExpr("SPEXEC_UNDEFINED").Eval()
			Expect(err).To(HaveOccurred())
		})
	})
})
