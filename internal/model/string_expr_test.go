package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringLiteralExpr", func() {
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
