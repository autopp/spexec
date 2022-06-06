package model

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("TypeOf()",
	func(x interface{}, expected Type) {
		Expect(TypeOf(x)).To(Equal(expected))
	},
	Entry(`when 42 given, returns TypeInt`, 42, TypeInt),
	Entry(`when json.Number("42") given, returns TypeInt`, json.Number("42"), TypeInt),
	Entry(`when true given, returns TypeBool`, true, TypeBool),
	Entry(`when "hello" given, returns TypeString`, "hello", TypeString),
	Entry(`when slice given, returns TypeSeq`, Seq{42, true, "hello"}, TypeSeq),
	Entry(`when string key map given, returns TypeMap`, Map{"message": "hello"}, TypeMap),
	Entry(`when nil given, returns TypeNil`, nil, TypeNil),
	Entry(`when other given, returns TypeUnknown`, make(chan int), TypeUnkown),
)

var _ = DescribeTable("TypeNameOf()",
	func(x interface{}, expected string) {
		Expect(TypeNameOf(x)).To(Equal(expected))
	},
	Entry(`when 42 given, returns "int"`, 42, "int"),
	Entry(`when json.Number("42") given, returns "int"`, json.Number("42"), "int"),
	Entry(`when true given, returns "bool"`, true, "bool"),
	Entry(`when "hello" given, returns "string"`, "hello", "string"),
	Entry(`when slice given, returns "seq"`, Seq{42, true, "hello"}, "seq"),
	Entry(`when string key map given, returns "map"`, Map{"message": "hello"}, "map"),
	Entry(`when nil given, returns "nil"`, nil, "nil"),
	Entry(`when other given, returns type name`, make(chan int), "chan int"),
)
