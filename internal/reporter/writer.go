// Copyright (C) 2021 Akira Tanimura (@autopp)
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

import "io"

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
)

func (c Color) isValid() bool {
	return c >= Black && c <= White
}

// Writer extends io.Writer
type Writer struct {
	io.Writer
	colorStack []Color
}

func newWriter(w io.Writer) *Writer {
	return &Writer{w, make([]Color, 0)}
}
