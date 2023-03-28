// Copyright (C) 2021-2023	 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reporter

import (
	"fmt"
	"io"
	"log"
)

// Color is enum type for console color
type Color int

const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Reset
)

func (c Color) isValid() bool {
	return c >= Black && c <= Reset
}

// Writer extends io.Writer
type Writer struct {
	io.Writer
	colorMode  bool
	colorStack []Color
}

func newWriter(w io.Writer, colorMode bool) *Writer {
	return &Writer{w, colorMode, []Color{Reset}}
}

func (w *Writer) UseColor(c Color, f func()) {
	if !w.colorMode {
		f()
		return
	}

	w.writeEscapeSequense(c)
	w.colorStack = append(w.colorStack, c)

	f()

	w.colorStack = w.colorStack[:len(w.colorStack)-1]
	w.writeEscapeSequense(w.colorStack[len(w.colorStack)-1])
}

func (w *Writer) writeEscapeSequense(c Color) {
	switch c {
	case Black:
		fmt.Fprint(w, "\033[30m")
	case Red:
		fmt.Fprint(w, "\033[31m")
	case Green:
		fmt.Fprint(w, "\033[32m")
	case Yellow:
		fmt.Fprint(w, "\033[33m")
	case Blue:
		fmt.Fprint(w, "\033[34m")
	case Magenta:
		fmt.Fprint(w, "\033[35m")
	case Cyan:
		fmt.Fprint(w, "\033[36m")
	case White:
		fmt.Fprint(w, "\033[37m")
	case Reset:
		fmt.Fprint(w, "\033[0m")
	default:
		log.Fatalf("color is not valid: %d", c)
	}
}
