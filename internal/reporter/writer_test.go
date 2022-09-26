package reporter

import (
	"bytes"
	"fmt"

	g "github.com/onsi/ginkgo/v2" // Reporter are duplicated
	. "github.com/onsi/gomega"
)

var _ = g.Describe("Writer", func() {
	g.Describe("UseColor()", func() {
		g.It("wrap the specified color code and reset code around text when called", func() {
			buf := &bytes.Buffer{}
			w := newWriter(buf, true)
			fmt.Fprintf(buf, "one")
			w.UseColor(Red, func() {
				fmt.Fprintf(buf, "two")
			})
			fmt.Fprintf(buf, "three")

			Expect(buf.String()).To(Equal("one\033[31mtwo\033[0mthree"))
		})
	})
})
