package reporter

import (
	"bytes"
	"fmt"

	g "github.com/onsi/ginkgo/v2" // Reporter are duplicated
	. "github.com/onsi/gomega"
)

var _ = g.Describe("Writer", func() {
	g.Describe("UseColor()", func() {
		g.It("wraps the specified color code and reset code around text when called", func() {
			buf := &bytes.Buffer{}
			w := newWriter(buf, true)
			fmt.Fprintf(buf, "one")
			w.UseColor(Red, func() {
				fmt.Fprintf(buf, "two")
			})
			fmt.Fprintf(buf, "three")

			Expect(buf.String()).To(Equal("one\033[31mtwo\033[0mthree"))
		})

		g.It("wraps the specified color codes around text when nested", func() {
			buf := &bytes.Buffer{}
			w := newWriter(buf, true)
			fmt.Fprintf(buf, "one")
			w.UseColor(Red, func() {
				fmt.Fprintf(buf, "two")
				w.UseColor(Blue, func() {
					fmt.Fprintf(buf, "three")
				})
				fmt.Fprintf(buf, "four")
			})
			fmt.Fprintf(buf, "five")

			Expect(buf.String()).To(Equal("one\033[31mtwo\033[34mthree\033[31mfour\033[0mfive"))
		})

		g.It("dose not wrap the specified color code and reset code around text when color mode is false", func() {
			buf := &bytes.Buffer{}
			w := newWriter(buf, false)
			fmt.Fprintf(buf, "one")
			w.UseColor(Red, func() {
				fmt.Fprintf(buf, "two")
			})
			fmt.Fprintf(buf, "three")

			Expect(buf.String()).To(Equal("onetwothree"))
		})
	})
})
