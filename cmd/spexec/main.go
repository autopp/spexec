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

package main

import (
	stderrors "errors"
	"fmt"
	"os"

	"github.com/autopp/spexec/pkg/cmd"
	"github.com/autopp/spexec/pkg/errors"
)

var version = "HEAD"

var statuses = map[errors.Code]int{
	errors.ErrTestFailed:    1,
	errors.ErrInvalidSpec:   2,
	errors.ErrInternalError: 3,
}

func main() {
	err := cmd.Main(version, os.Stdin, os.Stdout, os.Stderr, os.Args[1:])

	if err == nil {
		return
	}

	var e *errors.Error
	var status int
	if stderrors.As(err, &e) {
		var ok bool
		if status, ok = statuses[e.Code]; !ok {
			status = statuses[errors.ErrInternalError]
		}
	} else {
		// Assume command line error by via cobra
		status = 4
	}

	if status == statuses[errors.ErrInternalError] {
		fmt.Fprint(os.Stderr, "Internal Error: ")
	}

	fmt.Fprintf(os.Stderr, "%s\n", err.Error())

	os.Exit(status)
}
