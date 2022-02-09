// Copyright (C) 2021-2022	 Akira Tanimura (@autopp)
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

package model

import (
	"os"

	"github.com/autopp/spexec/internal/errors"
)

type StringExpr interface {
	Eval() (string, func() error, error)
	String() string
	stringExpr()
}

type literalStringExpr string

func NewLiteralStringExpr(v string) StringExpr {
	return literalStringExpr(v)
}

func (e literalStringExpr) Eval() (string, func() error, error) {
	return string(e), nil, nil
}

func (e literalStringExpr) String() string {
	return string(e)
}

func (e literalStringExpr) stringExpr() {}

type envStringExpr string

func NewEnvStringExpr(name string) StringExpr {
	return envStringExpr(name)
}

func (e envStringExpr) Eval() (string, func() error, error) {
	v, ok := os.LookupEnv(string(e))
	if !ok {
		return "", nil, errors.Errorf(errors.ErrInvalidSpec, "envrironment variable $%s is not defined", string(e))
	}
	return v, nil, nil
}

func (e envStringExpr) String() string {
	return "$" + string(e)
}

func (e envStringExpr) stringExpr() {}

func EvalStringExprs(exprs []StringExpr) ([]string, func() []error, error) {
	values := make([]string, len(exprs))
	cleanups := make([]func() error, 0)
	var firstErr error
	for i, expr := range exprs {
		value, cleanup, err := expr.Eval()
		cleanups = append(cleanups, cleanup)
		if err != nil {
			firstErr = err
			break
		}
		values[i] = value
	}

	cleanupAll := func() []error {
		errs := make([]error, 0)
		for _, cleanup := range cleanups {
			if cleanup == nil {
				continue
			}

			if err := cleanup(); err != nil {
				errs = append(errs, err)
			}
		}
		return errs
	}

	if firstErr != nil {
		return nil, cleanupAll, firstErr
	}
	return values, cleanupAll, nil
}
