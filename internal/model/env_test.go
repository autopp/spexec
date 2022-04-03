package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {
	var env *Env

	JustBeforeEach(func() {
		env = NewEnv(nil)
	})

	Describe("Define()", func() {
		It("returns true when not defined", func() {
			Expect(env.Define("command", "spexec")).To(BeTrue())
		})

		It("returns false when already defined", func() {
			env.Define("command", "go")
			Expect(env.Define("command", "spexec")).To(BeFalse())
		})

		It("returns true when already defined and new frame is pushed", func() {
			env.Define("command", "go")
			env = NewEnv(env)
			Expect(env.Define("command", "spexec")).To(BeTrue())
		})
	})

	Describe("Lookup()", func() {
		It("returns value and true when defined", func() {
			env.Define("command", "spexec")
			v, ok := env.Lookup("command")
			Expect(v).To(Equal("spexec"))
			Expect(ok).To(BeTrue())
		})

		It("returns false when not defined", func() {
			_, ok := env.Lookup("command")
			Expect(ok).To(BeFalse())
		})

		It("returns value and true when already defined in previous frame", func() {
			env.Define("command", "spexec")
			env = NewEnv(env)
			v, ok := env.Lookup("command")
			Expect(v).To(Equal("spexec"))
			Expect(ok).To(BeTrue())
		})

		It("returns current value and true when already defined in both previous and current frame", func() {
			env.Define("command", "go")
			env = NewEnv(env)
			env.Define("command", "spexec")
			v, ok := env.Lookup("command")
			Expect(v).To(Equal("spexec"))
			Expect(ok).To(BeTrue())
		})
	})
})
