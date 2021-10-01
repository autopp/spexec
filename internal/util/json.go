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

package util

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/autopp/spexec/internal/errors"
)

func UnmarshalJSON(in []byte, out interface{}) error {
	d := json.NewDecoder(bytes.NewBuffer(in))
	d.UseNumber()
	err := d.Decode(out)
	if err != nil {
		return err
	}

	if t, err := d.Token(); t != nil || err != io.EOF {
		return errors.Errorf(errors.ErrInternalError, "invalid json string")
	}
	return nil
}
