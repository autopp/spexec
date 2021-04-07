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

package errors

import (
	"errors"
	"fmt"
)

type Code string

const (
	ErrTestFailed    Code = "test failed"
	ErrInvalidConfig Code = "invalid config"
)

type Error struct {
	Code Code
	err  error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func Errorf(code Code, format string, a ...interface{}) error {
	return Wrap(code, fmt.Errorf(format, a...))
}

func New(code Code, text string) error {
	return Wrap(code, errors.New(text))
}

func Wrap(code Code, err error) error {
	return &Error{Code: code, err: err}
}
