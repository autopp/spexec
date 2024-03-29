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

package model

import (
	"encoding/json"
	"fmt"
)

type Map = map[string]any
type Seq = []any

type Type int

const (
	TypeNil Type = iota
	TypeInt
	TypeBool
	TypeString
	TypeSeq
	TypeMap
	TypeUnkown
)

var typeNames = map[Type]string{
	TypeNil:    "nil",
	TypeInt:    "int",
	TypeBool:   "bool",
	TypeString: "string",
	TypeSeq:    "seq",
	TypeMap:    "map",
}

func TypeOf(x any) Type {
	if x == nil {
		return TypeNil
	}

	switch casted := x.(type) {
	case int:
		return TypeInt
	case json.Number:
		if _, err := casted.Int64(); err == nil {
			return TypeInt
		}
	case bool:
		return TypeBool
	case string:
		return TypeString
	case Seq:
		return TypeSeq
	case Map:
		return TypeMap
	}

	return TypeUnkown
}

func TypeNameOf(x any) string {
	if name, ok := typeNames[TypeOf(x)]; ok {
		return name
	}

	return fmt.Sprintf("%T", x)
}
