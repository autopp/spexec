package spec

import (
	"encoding/json"
	"time"

	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Validator", func() {
	var v *Validator

	BeValidationError := func(message interface{}) types.GomegaMatcher {
		matcher, ok := message.(types.GomegaMatcher)
		if !ok {
			matcher = Equal(message)
		}

		return And(HaveOccurred(), WithTransform(func(err error) string { return err.Error() }, matcher))
	}

	JustBeforeEach(func() {
		v, _ = NewValidator("")
	})

	Describe("AddViolation() and Error()", func() {
		Context("with no AddViolation() call", func() {
			It("makes to Error() to return nil", func() {
				Expect(v.Error()).To(Succeed())
			})
		})

		Context("with a AddViolation() call", func() {
			It("makes to Error() to returns error which contains path and formatted violation message", func() {
				v.AddViolation("answer is %d", 42)
				Expect(v.Error()).To(BeValidationError("$: answer is 42"))
			})
		})

		Context("with multiple AddViolation() calls", func() {
			It("make to Error() to return error which contains given violations joined by newline", func() {
				v.AddViolation("error1")
				v.AddViolation("error2")
				Expect(v.Error()).To(BeValidationError("$: error1\n$: error2"))
			})
		})
	})

	Describe("InPath()", func() {
		It("appends path prefix in callback", func() {
			v.InPath(":prefix1", func() {
				v.AddViolation("error1")
				v.InPath(":prefix2", func() {
					v.AddViolation("error2")
				})
			})
			v.AddViolation("error3")

			Expect(v.Error()).To(BeValidationError("$:prefix1: error1\n$:prefix1:prefix2: error2\n$: error3"))
		})
	})

	Describe("InField()", func() {
		It(`be equivalent to InPath("."+field)`, func() {
			v.InField("field", func() {
				v.AddViolation("error")
			})

			Expect(v.Error()).To(BeValidationError("$.field: error"))
		})
	})

	Describe("InIndex()", func() {
		It(`be equivalent to InPath("["+i+"]")`, func() {
			v.InIndex(1, func() {
				v.AddViolation("error")
			})

			Expect(v.Error()).To(BeValidationError("$[1]: error"))
		})
	})

	Describe("MayBeMap()", func() {
		Context("with a Map", func() {
			It("returns the given Map and true", func() {
				given := make(Map)
				m, b := v.MayBeMap(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not Map", func() {
			It("returns something and false", func() {
				_, b := v.MayBeMap(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeMap()", func() {
		Context("with a Map", func() {
			It("returns the given Map and true", func() {
				given := make(Map)
				m, b := v.MustBeMap(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not Map", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeMap(42)

				Expect(v.Error()).To(BeValidationError("$: should be map, but is int"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeSeq()", func() {
		Context("with a Seq", func() {
			It("returns the given Seq and true", func() {
				given := make(Seq, 0)
				m, b := v.MustBeSeq(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not Seq", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeSeq(42)

				Expect(v.Error()).To(BeValidationError("$: should be seq, but is int"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MayBeString()", func() {
		Context("with a string", func() {
			It("returns the given string and true", func() {
				given := "hello"
				m, b := v.MayBeString(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not string", func() {
			It("returns something and false", func() {
				_, b := v.MayBeString(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeString()", func() {
		Context("with a string", func() {
			It("returns the given string and true", func() {
				given := "hello"
				m, b := v.MustBeString(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not string", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeString(42)

				Expect(v.Error()).To(BeValidationError("$: should be string, but is int"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeStringExpr", func() {
		Context("with a string", func() {
			It("returns literalString and true", func() {
				given := "hello"
				actual, b := v.MustBeStringExpr(given)

				Expect(actual).To(Equal(model.NewLiteralStringExpr(given)))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a map which contains type='env' and name='MESSAGE'", func() {
			It("returns literalString and true", func() {
				given := Map{"type": "env", "name": "MESSAGE"}
				actual, b := v.MustBeStringExpr(given)

				Expect(actual).To(Equal(model.NewEnvStringExpr("MESSAGE")))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a map which contains type='env' and dose not contain name", func() {
			It("adds violation and returns something and false", func() {
				given := Map{"type": "env"}
				_, b := v.MustBeStringExpr(given)

				Expect(b).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(`$: should have .name as string`))
			})
		})

		Context("with a map which contains unknown type", func() {
			It("adds violation and returns something and false", func() {
				given := Map{"type": "unknown"}
				_, b := v.MustBeStringExpr(given)

				Expect(b).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(`$.type: unknown type "unknown"`))
			})
		})

		Context("with a map which dose not contain type", func() {
			It("adds violation and returns something and false", func() {
				given := Map{}
				_, b := v.MustBeStringExpr(given)

				Expect(b).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(`$: should have .type as string`))
			})
		})

		Context("with not string nor map", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeStringExpr(42)

				Expect(b).To(BeFalse())
				Expect(v.Error()).To(BeValidationError("$: should be string or map, but is int"))
			})
		})
	})

	Describe("MustBeInt()", func() {
		Context("with a int", func() {
			It("returns the given int and true", func() {
				given := 42
				m, b := v.MustBeInt(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not int", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeInt("hello")

				Expect(v.Error()).To(BeValidationError("$: should be int, but is string"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeBool()", func() {
		Context("with a bool", func() {
			It("returns the given bool and true", func() {
				given := true
				m, b := v.MustBeBool(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not bool", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeBool("hello")

				Expect(v.Error()).To(BeValidationError("$: should be bool, but is string"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeDuration()", func() {
		DescribeTable("with a duration string or positive integer (success)",
			func(given interface{}, expected time.Duration) {
				d, b := v.MustBeDuration(given)

				Expect(d).To(Equal(expected))
				Expect(b).To(BeTrue())
				Expect(v.Error()).To(BeNil())
			},
			Entry(`given: "3s"`, "3s", 3*time.Second),
			Entry(`given: "1m"`, "1m", 1*time.Minute),
			Entry(`given: "500ms"`, "500ms", 500*time.Millisecond),
			Entry(`given: 10`, 10, 10*time.Second),
		)
	})

	Describe("MayHave()", func() {
		Context("when the given map has specified field", func() {
			It("calls the callback with the value in patch of the field and return it and true", func() {
				contained := Seq{42, "hello"}
				var passed interface{}
				x, exists := v.MayHave(Map{"field": contained}, "field", func(x interface{}) {
					passed = x
					v.AddViolation("error")
				})

				Expect(passed).To(BeEquivalentTo(contained))
				Expect(v.Error()).To(BeValidationError("$.field: error"))
				Expect(x).To(Equal(contained))
				Expect(exists).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("dose not call the callback and return something and false", func() {
				_, exists := v.MayHave(make(Map), "field", func(x interface{}) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
			})
		})
	})

	Describe("MayHaveMap()", func() {
		Context("when the given map has specified field which is a Map", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true, true", func() {
				contained := make(Map)
				var passed Map
				m, exists, ok := v.MayHaveMap(Map{"field": contained}, "field", func(m Map) {
					passed = m
					v.AddViolation("error")
				})

				Expect(passed).To(Equal(contained))
				Expect(v.Error()).To(BeValidationError("$.field: error"))
				Expect(m).To(Equal(contained))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("dose not call the callback and returns something, false, true", func() {
				_, exists, ok := v.MayHaveMap(make(Map), "field", func(Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a Map", func() {
			It("dose not call the callback, adds violation, and returns something, false, false", func() {
				_, exists, ok := v.MayHaveMap(Map{"field": "hello"}, "field", func(Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be map, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveSeq()", func() {
		Context("when the given map has specified field which is a Seq", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true, true", func() {
				contained := make(Seq, 0)
				var passed Seq
				m, exists, ok := v.MayHaveSeq(Map{"field": contained}, "field", func(s Seq) {
					passed = s
					v.AddViolation("error")
				})

				Expect(passed).To(Equal(contained))
				Expect(v.Error()).To(BeValidationError("$.field: error"))
				Expect(m).To(Equal(contained))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("dose not call the callback and returns something, false, true", func() {
				_, exists, ok := v.MayHaveSeq(make(Map), "field", func(Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a Seq", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveSeq(Map{"field": "hello"}, "field", func(Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be seq, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MustHaveSeq()", func() {
		Context("when the given map has specified field which is a Seq", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true", func() {
				contained := make(Seq, 0)
				var passed Seq
				m, ok := v.MustHaveSeq(Map{"field": contained}, "field", func(s Seq) {
					passed = s
					v.AddViolation("error")
				})

				Expect(passed).To(Equal(contained))
				Expect(v.Error()).To(BeValidationError("$.field: error"))
				Expect(m).To(Equal(contained))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("dose not call the callback and returns something, false", func() {
				_, ok := v.MustHaveSeq(make(Map), "field", func(Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$: should have .field as seq"))
				Expect(ok).To(BeFalse())
			})
		})

		Context("when the given map has specified field which is not a Seq", func() {
			It("dose not call the callback, add violation and returns something, false", func() {
				_, ok := v.MustHaveSeq(Map{"field": "hello"}, "field", func(Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be seq, but is string"))
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("ForInSeq()", func() {
		It("calls callback with each index and element in path of it", func() {
			v.ForInSeq(Seq{42, "hello"}, func(i int, x interface{}) bool {
				v.AddViolation("%d:%#v", i, x)
				return true
			})

			Expect(v.Error()).To(BeValidationError("$[0]: 0:42\n$[1]: 1:\"hello\""))
		})

		It("stops calling callback when it returns false", func() {
			calls := make([]int, 0)
			v.ForInSeq(Seq{"a", "b", "c", "d"}, func(i int, x interface{}) bool {
				calls = append(calls, i)
				return i < 2
			})

			Expect(calls).To(Equal([]int{0, 1, 2}))
		})
	})

	Describe("MayHaveString()", func() {
		Context("when the given map has specified field which is a string", func() {
			It("returns the it, true, true", func() {
				s, exists, ok := v.MayHaveString(Map{"field": "hello"}, "field")

				Expect(s).To(Equal("hello"))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveString(Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveString(Map{"field": 42}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be string, but is int"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MustHaveString()", func() {
		Context("when the given map has specified field which is a string", func() {
			It("returns the it, true, true", func() {
				s, ok := v.MustHaveString(Map{"field": "hello"}, "field")

				Expect(s).To(Equal("hello"))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, ok := v.MustHaveString(Map{}, "field")

				Expect(v.Error()).To(BeValidationError("$: should have .field as string"))
				Expect(ok).To(BeFalse())
			})
		})

		Context("when the given map has specified field which is not a string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, ok := v.MustHaveString(Map{"field": 42}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be string, but is int"))
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveInt()", func() {
		Context("when the given map has specified field which is a int", func() {
			It("returns the it, true, true", func() {
				s, exists, ok := v.MayHaveInt(Map{"field": 42}, "field")

				Expect(s).To(Equal(42))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveInt(Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a int", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveInt(Map{"field": "hello"}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be int, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveDuration()", func() {
		Context("when the given map has specified field which is a duration string", func() {
			It("returns the duration, true, true", func() {
				d, exists, ok := v.MayHaveDuration(Map{"field": "1s"}, "field")

				Expect(d).To(Equal(1 * time.Second))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveDuration(Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a duration string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveDuration(Map{"field": "666"}, "field")

				Expect(v.Error()).To(BeValidationError(MatchRegexp(`\$\.field: should be positive integer or duration string, but cannot parse`)))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveEnvSeq", func() {
		Context("when the given map has specified field which is a seq of key-value", func() {
			It("returns parsed array", func() {
				e, exists, ok := v.MayHaveEnvSeq(Map{"env": Seq{Map{"name": "a", "value": "foo"}, Map{"name": "b", "value": "bar"}}}, "env")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(Equal([]util.StringVar{
					{Name: "a", Value: "foo"},
					{Name: "b", Value: "bar"},
				}))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns nil", func() {
				e, exists, ok := v.MayHaveEnvSeq(Map{}, "env")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		DescribeTable("adds vioration",
			func(env interface{}, prefix string) {
				e, _, ok := v.MayHaveEnvSeq(Map{"env": env}, "env")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", Map{}, "$.env:"),
			Entry("when the field contains invalid key-value (name is not string)",
				Seq{Map{"name": "a", "value": "foo"}, Map{"name": 0, "value": "foo"}},
				"$.env[1].name:",
			),
			Entry("when the field contains invalid key-value (name is not var name)",
				Seq{Map{"name": "a", "value": "foo"}, Map{"name": "0a", "value": "foo"}},
				"$.env[1].name:",
			),
			Entry("when the field contains invalid key-value (name is missing)",
				Seq{Map{"name": "a", "value": "foo"}, Map{"value": "foo"}},
				"$.env[1]:",
			),
			Entry("when the field contains invalid key-value (value is not string)",
				Seq{Map{"name": "a", "value": "foo"}, Map{"name": "a", "value": 42}},
				"$.env[1].value:",
			),
			Entry("when the field contains invalid key-value (value is missing)",
				Seq{Map{"name": "a", "value": "foo"}, Map{"name": "a"}},
				"$.env[1]:",
			),
		)
	})

	Describe("MayHaveCommand", func() {
		Context("when the given map has specified field which is array of string", func() {
			It("returns parsed array", func() {
				e, exists, ok := v.MayHaveCommand(Map{"command": Seq{"sh", "-c", "true"}}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(Equal([]model.StringExpr{model.NewLiteralStringExpr("sh"), model.NewLiteralStringExpr("-c"), model.NewLiteralStringExpr("true")}))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns nil", func() {
				e, exists, ok := v.MayHaveCommand(Map{}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		DescribeTable("adds violation",
			func(command interface{}, prefix string) {
				e, _, ok := v.MayHaveCommand(Map{"command": command}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", Map{}, "$.command:"),
			Entry("when the field contains not string", Seq{"sh", 1}, "$.command[1]:"),
			Entry("when the field is empty seq", Seq{}, "$.command"),
		)
	})

	Describe("MustHaveCommand", func() {
		Context("when the given map has specified field which is array of string", func() {
			It("returns parsed array", func() {
				e, ok := v.MustHaveCommand(Map{"command": Seq{"sh", "-c", "true"}}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(Equal([]model.StringExpr{model.NewLiteralStringExpr("sh"), model.NewLiteralStringExpr("-c"), model.NewLiteralStringExpr("true")}))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("adds ", func() {
				e, ok := v.MustHaveCommand(Map{}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix("$: should have .command as command seq")))
			})
		})

		DescribeTable("adds violation",
			func(command interface{}, prefix string) {
				e, ok := v.MustHaveCommand(Map{"command": command}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", Map{}, "$.command:"),
			Entry("when the field contains not string", Seq{"sh", 1}, "$.command[1]:"),
			Entry("when the field is empty seq", Seq{}, "$.command"),
		)
	})
})

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
)
