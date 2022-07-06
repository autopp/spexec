package spec

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/util"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type unmarshalableToYAML struct {
}

func (unmarshalableToYAML) MarshalYAML() (any, error) {
	return nil, errors.New("cannot marshal to YAML")
}

var _ = Describe("Validator", func() {
	var v *Validator

	BeValidationError := func(message any) types.GomegaMatcher {
		matcher, ok := message.(types.GomegaMatcher)
		if !ok {
			matcher = Equal(message)
		}

		return And(HaveOccurred(), WithTransform(func(err error) string { return err.Error() }, matcher))
	}

	JustBeforeEach(func() {
		v, _ = NewValidator("")
	})

	Describe("Filename and GetDir()", func() {
		It("with no filename, returns empty string and current directory", func() {
			Expect(v.Filename).To(BeEmpty())
			wd, err := os.Getwd()
			if err != nil {
				Fail("os.Getwd() fails: " + err.Error())
			}
			Expect(v.GetDir()).To(Equal(wd))
		})

		It("with filename, returns absolute path and directory of it", func() {
			v, _ = NewValidator("validator_test.go")
			filename, _ := filepath.Abs("validator_test.go")
			Expect(v.Filename).To(Equal(filename))
			Expect(v.GetDir()).To(Equal(filepath.Dir(filename)))
		})
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
		Context("with a model.Map", func() {
			It("returns the given model.Map and true", func() {
				given := make(model.Map)
				m, b := v.MayBeMap(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not model.Map", func() {
			It("returns something and false", func() {
				_, b := v.MayBeMap(42)

				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeMap()", func() {
		Context("with a model.Map", func() {
			It("returns the given model.Map and true", func() {
				given := make(model.Map)
				m, b := v.MustBeMap(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not model.Map", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeMap(42)

				Expect(v.Error()).To(BeValidationError("$: should be map, but is int"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeSeq()", func() {
		Context("with a model.Seq", func() {
			It("returns the given model.Seq and true", func() {
				given := make(model.Seq, 0)
				m, b := v.MustBeSeq(given)

				Expect(m).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a not model.Seq", func() {
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

	Describe("MayBeQualified", func() {
		Context("with 1 element map", func() {
			It("returns the qualifier, the value and true", func() {
				given := model.Map{"$": "answer"}
				q, v, b := v.MayBeQualified(given)
				Expect(q).To(Equal("$"))
				Expect(v).To(Equal("answer"))
				Expect(b).To(BeTrue())
			})
		})

		Context("with not map", func() {
			It("returns empty values and false", func() {
				given := "answer"
				_, _, b := v.MayBeQualified(given)
				Expect(b).To(BeFalse())
			})
		})

		Context("with two or more elements map", func() {
			It("returns empty values and false", func() {
				given := model.Map{"$": "answer", "$$": 42}
				_, _, b := v.MayBeQualified(given)
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MayBeVariable", func() {
		Context("with 1 element map which contain '$' and the name", func() {
			It("returns the name and true", func() {
				given := model.Map{"$": "answer"}
				name, b := v.MayBeVariable(given)
				Expect(name).To(Equal("answer"))
				Expect(b).To(BeTrue())
			})
		})

		Context("with 1 elements map contain '$' and not string", func() {
			It("returns something and false", func() {
				given := model.Map{"$": 42}
				_, b := v.MayBeVariable(given)
				Expect(b).To(BeFalse())
			})
		})

		Context("with 1 elements map contain '$' and not variable name", func() {
			It("returns something and false", func() {
				given := model.Map{"$": "foo bar"}
				_, b := v.MayBeVariable(given)
				Expect(b).To(BeFalse())
			})
		})

		Context("with 1 element map without '$'", func() {
			It("returns something and false", func() {
				given := model.Map{"$$": "answer"}
				_, b := v.MayBeVariable(given)
				Expect(b).To(BeFalse())
			})
		})

		Context("with not map", func() {
			It("returns something and false", func() {
				given := "answer"
				_, b := v.MayBeVariable(given)
				Expect(b).To(BeFalse())
			})
		})

		Context("with two or more elements map", func() {
			It("returns something and false", func() {
				given := model.Map{"$": "answer", "$$": 42}
				_, b := v.MayBeVariable(given)
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustBeStringExpr", func() {
		Context("with a string", func() {
			It("returns literalStringExpr and true", func() {
				given := "hello"
				actual, b := v.MustBeStringExpr(given)

				Expect(actual).To(Equal(model.NewLiteralStringExpr(given)))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a map which contains type='env'", func() {
			Context("and name='MESSAGE'", func() {
				It("returns envStringExpr and true", func() {
					given := model.Map{"type": "env", "name": "MESSAGE"}
					actual, b := v.MustBeStringExpr(given)

					Expect(actual).To(Equal(model.NewEnvStringExpr("MESSAGE")))
					Expect(b).To(BeTrue())
				})
			})

			Context("and dose not contain name", func() {
				It("adds violation and returns something and false", func() {
					given := model.Map{"type": "env"}
					_, b := v.MustBeStringExpr(given)

					Expect(b).To(BeFalse())
					Expect(v.Error()).To(BeValidationError(`$: should have .name as string`))
				})
			})
		})

		Context("with a map which contains type='file'", func() {
			Context("and value='hello world'", func() {
				It("returns fileStringExpr and true", func() {
					given := model.Map{"type": "file", "value": "hello"}
					actual, b := v.MustBeStringExpr(given)

					Expect(actual).To(Equal(model.NewFileStringExpr("", "hello")))
					Expect(b).To(BeTrue())
				})
			})

			Context("and value is not string", func() {
				It("adds violation and returns something and false", func() {
					given := model.Map{"type": "file", "value": 42}
					_, b := v.MustBeStringExpr(given)

					Expect(b).To(BeFalse())
					Expect(v.Error()).To(BeValidationError(`$.value: should be string, but is int`))
				})
			})

			Context("and format='yaml'", func() {
				Context("and value is yaml compatible", func() {
					It("returns fileStringExpr and true", func() {
						given := model.Map{"type": "file", "format": "yaml", "value": model.Map{"answer": 42}}
						actual, b := v.MustBeStringExpr(given)

						Expect(actual).To(Equal(model.NewFileStringExpr("*.yaml", "answer: 42\n")))
						Expect(b).To(BeTrue())
						Expect(v.Error()).NotTo(HaveOccurred())
					})
				})

				Context("and value is not yaml compatible", func() {

					It("adds violation and returns something and false", func() {
						given := model.Map{"type": "file", "format": "yaml", "value": model.Map{"answer": unmarshalableToYAML{}}}
						_, b := v.MustBeStringExpr(given)

						Expect(b).To(BeFalse())
						Expect(v.Error()).To(BeValidationError(`$.value: cannot encode to a YAML string: cannot marshal to YAML`))
					})
				})

				Context("and dose not contain value", func() {
					It("adds violation and returns something and false", func() {
						given := model.Map{"type": "file", "format": "yaml"}
						_, b := v.MustBeStringExpr(given)

						Expect(b).To(BeFalse())
						Expect(v.Error()).To(BeValidationError(`$: should have .value`))
					})
				})
			})

			Context("and format='unknown'", func() {
				It("adds violation and returns something and false", func() {
					given := model.Map{"type": "file", "format": "unknown"}
					_, b := v.MustBeStringExpr(given)

					Expect(b).To(BeFalse())
					Expect(v.Error()).To(BeValidationError(`$.format: should be a "raw" or "yaml", but is "unknown"`))
				})
			})

			Context("and format is not string", func() {
				It("adds violation and returns something and false", func() {
					given := model.Map{"type": "file", "format": 42}
					_, b := v.MustBeStringExpr(given)

					Expect(b).To(BeFalse())
					Expect(v.Error()).To(BeValidationError(`$.format: should be string, but is int`))
				})
			})

			Context("and dose not contain value", func() {
				It("adds violation and returns something and false", func() {
					given := model.Map{"type": "file"}
					_, b := v.MustBeStringExpr(given)

					Expect(b).To(BeFalse())
					Expect(v.Error()).To(BeValidationError(`$: should have .value as string`))
				})
			})
		})

		Context("with a map which contains unknown type", func() {
			It("adds violation and returns something and false", func() {
				given := model.Map{"type": "unknown"}
				_, b := v.MustBeStringExpr(given)

				Expect(b).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(`$.type: unknown type "unknown"`))
			})
		})

		Context("with a map which dose not contain type", func() {
			It("adds violation and returns something and false", func() {
				given := model.Map{}
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
				i, b := v.MustBeInt(given)

				Expect(i).To(Equal(given))
				Expect(b).To(BeTrue())
			})
		})

		Context("with a valid json.Number", func() {
			It("returns the int represented and true", func() {
				given := json.Number("42")
				i, b := v.MustBeInt(given)

				Expect(i).To(Equal(42))
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

		Context("with a invalid json.Number", func() {
			It("returns the int represented and true", func() {
				given := json.Number("abc")
				_, b := v.MustBeInt(given)

				Expect(v.Error()).To(BeValidationError(HavePrefix("$: should be int, but is invalid json.Number: ")))
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
			func(given any, expected time.Duration) {
				d, b := v.MustBeDuration(given)

				Expect(d).To(Equal(expected))
				Expect(b).To(BeTrue())
				Expect(v.Error()).To(BeNil())
			},
			Entry(`given: "3s"`, "3s", 3*time.Second),
			Entry(`given: "1m"`, "1m", 1*time.Minute),
			Entry(`given: "500ms"`, "500ms", 500*time.Millisecond),
			Entry(`given: 10`, 10, 10*time.Second),
			Entry(`given: 10 (json.Number)`, json.Number("10"), 10*time.Second),
		)

		Context("when the given value is not integer nor string", func() {
			It("adds violation and returns something and false", func() {
				_, b := v.MustBeDuration(true)

				Expect(v.Error()).To(BeValidationError("$: should be positive integer or duration string, but is bool"))
				Expect(b).To(BeFalse())
			})
		})
	})

	Describe("MustHave()", func() {
		Context("when the given map has specified field", func() {
			It("returns value of the field and true", func() {
				contained := model.Seq{42, "hello"}
				x, ok := v.MustHave(model.Map{"field": contained}, "field")

				Expect(x).To(Equal(contained))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("adds violation and returns something and false", func() {
				_, ok := v.MustHave(make(model.Map), "field")

				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError("$: should have .field"))
			})
		})
	})

	Describe("MayHave()", func() {
		Context("when the given map has specified field", func() {
			It("calls the callback with the value in patch of the field and return it and true", func() {
				contained := model.Seq{42, "hello"}
				var passed any
				x, exists := v.MayHave(model.Map{"field": contained}, "field", func(x any) {
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
				_, exists := v.MayHave(make(model.Map), "field", func(x any) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
			})
		})
	})

	Describe("MayHaveMap()", func() {
		Context("when the given map has specified field which is a model.Map", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true, true", func() {
				contained := make(model.Map)
				var passed model.Map
				m, exists, ok := v.MayHaveMap(model.Map{"field": contained}, "field", func(m model.Map) {
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
				_, exists, ok := v.MayHaveMap(make(model.Map), "field", func(model.Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a model.Map", func() {
			It("dose not call the callback, adds violation, and returns something, false, false", func() {
				_, exists, ok := v.MayHaveMap(model.Map{"field": "hello"}, "field", func(model.Map) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be map, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveSeq()", func() {
		Context("when the given map has specified field which is a model.Seq", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true, true", func() {
				contained := make(model.Seq, 0)
				var passed model.Seq
				m, exists, ok := v.MayHaveSeq(model.Map{"field": contained}, "field", func(s model.Seq) {
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
				_, exists, ok := v.MayHaveSeq(make(model.Map), "field", func(model.Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a model.Seq", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveSeq(model.Map{"field": "hello"}, "field", func(model.Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be seq, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MustHaveSeq()", func() {
		Context("when the given map has specified field which is a model.Seq", func() {
			It("calls the callback with the map in map in path of the field and returns the it, true", func() {
				contained := make(model.Seq, 0)
				var passed model.Seq
				m, ok := v.MustHaveSeq(model.Map{"field": contained}, "field", func(s model.Seq) {
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
				_, ok := v.MustHaveSeq(make(model.Map), "field", func(model.Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$: should have .field as seq"))
				Expect(ok).To(BeFalse())
			})
		})

		Context("when the given map has specified field which is not a model.Seq", func() {
			It("dose not call the callback, add violation and returns something, false", func() {
				_, ok := v.MustHaveSeq(model.Map{"field": "hello"}, "field", func(model.Seq) {
					v.AddViolation("error")
				})

				Expect(v.Error()).To(BeValidationError("$.field: should be seq, but is string"))
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("ForInSeq()", func() {
		It("calls callback with each index and element in path of it", func() {
			v.ForInSeq(model.Seq{42, "hello"}, func(i int, x any) bool {
				v.AddViolation("%d:%#v", i, x)
				return true
			})

			Expect(v.Error()).To(BeValidationError("$[0]: 0:42\n$[1]: 1:\"hello\""))
		})

		It("stops calling callback when it returns false", func() {
			calls := make([]int, 0)
			v.ForInSeq(model.Seq{"a", "b", "c", "d"}, func(i int, x any) bool {
				calls = append(calls, i)
				return i < 2
			})

			Expect(calls).To(Equal([]int{0, 1, 2}))
		})
	})

	Describe("MayHaveString()", func() {
		Context("when the given map has specified field which is a string", func() {
			It("returns the it, true, true", func() {
				s, exists, ok := v.MayHaveString(model.Map{"field": "hello"}, "field")

				Expect(s).To(Equal("hello"))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveString(model.Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveString(model.Map{"field": 42}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be string, but is int"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MustHaveString()", func() {
		Context("when the given map has specified field which is a string", func() {
			It("returns the it, true, true", func() {
				s, ok := v.MustHaveString(model.Map{"field": "hello"}, "field")

				Expect(s).To(Equal("hello"))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, ok := v.MustHaveString(model.Map{}, "field")

				Expect(v.Error()).To(BeValidationError("$: should have .field as string"))
				Expect(ok).To(BeFalse())
			})
		})

		Context("when the given map has specified field which is not a string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, ok := v.MustHaveString(model.Map{"field": 42}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be string, but is int"))
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveInt()", func() {
		Context("when the given map has specified field which is a int", func() {
			It("returns the it, true, true", func() {
				s, exists, ok := v.MayHaveInt(model.Map{"field": 42}, "field")

				Expect(s).To(Equal(42))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveInt(model.Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a int", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveInt(model.Map{"field": "hello"}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be int, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveBool()", func() {
		Context("when the given map has specified field which is a bool", func() {
			It("returns the it, true, true", func() {
				b, exists, ok := v.MayHaveBool(model.Map{"field": true}, "field")

				Expect(b).To(Equal(true))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveBool(model.Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a bool", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveBool(model.Map{"field": "hello"}, "field")

				Expect(v.Error()).To(BeValidationError("$.field: should be bool, but is string"))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveDuration()", func() {
		Context("when the given map has specified field which is a duration string", func() {
			It("returns the duration, true, true", func() {
				d, exists, ok := v.MayHaveDuration(model.Map{"field": "1s"}, "field")

				Expect(d).To(Equal(1 * time.Second))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns something, false, true", func() {
				_, exists, ok := v.MayHaveDuration(model.Map{}, "field")

				Expect(v.Error()).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map has specified field which is not a duration string", func() {
			It("dose not call the callback, add violation and returns something, false, false", func() {
				_, exists, ok := v.MayHaveDuration(model.Map{"field": "666"}, "field")

				Expect(v.Error()).To(BeValidationError(MatchRegexp(`\$\.field: should be positive integer or duration string, but cannot parse`)))
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeFalse())
			})
		})
	})

	Describe("MayHaveEnvSeq", func() {
		Context("when the given map has specified field which is a seq of key-value", func() {
			It("returns parsed array", func() {
				e, exists, ok := v.MayHaveEnvSeq(model.Map{"env": model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"name": "b", "value": "bar"}}}, "env")

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
				e, exists, ok := v.MayHaveEnvSeq(model.Map{}, "env")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		DescribeTable("adds vioration",
			func(env any, prefix string) {
				e, _, ok := v.MayHaveEnvSeq(model.Map{"env": env}, "env")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", model.Map{}, "$.env:"),
			Entry("when the field contains invalid key-value (name is not string)",
				model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"name": 0, "value": "foo"}},
				"$.env[1].name:",
			),
			Entry("when the field contains invalid key-value (name is not var name)",
				model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"name": "0a", "value": "foo"}},
				"$.env[1].name:",
			),
			Entry("when the field contains invalid key-value (name is missing)",
				model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"value": "foo"}},
				"$.env[1]:",
			),
			Entry("when the field contains invalid key-value (value is not string)",
				model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"name": "a", "value": 42}},
				"$.env[1].value:",
			),
			Entry("when the field contains invalid key-value (value is missing)",
				model.Seq{model.Map{"name": "a", "value": "foo"}, model.Map{"name": "a"}},
				"$.env[1]:",
			),
		)
	})

	Describe("MayHaveCommand", func() {
		Context("when the given map has specified field which is array of string", func() {
			It("returns parsed array", func() {
				e, exists, ok := v.MayHaveCommand(model.Map{"command": model.Seq{"sh", "-c", "true"}}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(Equal([]model.StringExpr{model.NewLiteralStringExpr("sh"), model.NewLiteralStringExpr("-c"), model.NewLiteralStringExpr("true")}))
				Expect(exists).To(BeTrue())
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("returns nil", func() {
				e, exists, ok := v.MayHaveCommand(model.Map{}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(ok).To(BeTrue())
			})
		})

		DescribeTable("adds violation",
			func(command any, prefix string) {
				e, _, ok := v.MayHaveCommand(model.Map{"command": command}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", model.Map{}, "$.command:"),
			Entry("when the field contains not string", model.Seq{"sh", 1}, "$.command[1]:"),
			Entry("when the field is empty seq", model.Seq{}, "$.command"),
		)
	})

	Describe("MustHaveCommand", func() {
		Context("when the given map has specified field which is array of string", func() {
			It("returns parsed array", func() {
				e, ok := v.MustHaveCommand(model.Map{"command": model.Seq{"sh", "-c", "true"}}, "command")

				Expect(v.Error()).To(BeNil())
				Expect(e).To(Equal([]model.StringExpr{model.NewLiteralStringExpr("sh"), model.NewLiteralStringExpr("-c"), model.NewLiteralStringExpr("true")}))
				Expect(ok).To(BeTrue())
			})
		})

		Context("when the given map dose not have specified field", func() {
			It("adds ", func() {
				e, ok := v.MustHaveCommand(model.Map{}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix("$: should have .command as command seq")))
			})
		})

		DescribeTable("adds violation",
			func(command any, prefix string) {
				e, ok := v.MustHaveCommand(model.Map{"command": command}, "command")

				Expect(e).To(BeNil())
				Expect(ok).To(BeFalse())
				Expect(v.Error()).To(BeValidationError(HavePrefix(prefix)))
			},
			Entry("when the filed is not seq", model.Map{}, "$.command:"),
			Entry("when the field contains not string", model.Seq{"sh", 1}, "$.command[1]:"),
			Entry("when the field is empty seq", model.Seq{}, "$.command"),
		)
	})

	Describe("MustContainOnly()", func() {
		It("returns true and adds no error when the given map contains only specified", func() {
			m := model.Map{"foo": 1, "baz": "spexec"}
			Expect(v.MustContainOnly(m, "foo", "bar", "baz")).To(BeTrue())
			Expect(v.Error()).To(BeNil())
		})

		It("returns false and adds error when the given map contains not specified field", func() {
			m := model.Map{"foo": 1, "baz": "spexec"}
			Expect(v.MustContainOnly(m, "foo", "bar")).To(BeFalse())
			Expect(v.Error()).To(BeValidationError(Equal(`$: field .baz is not expected`)))
		})
	})
})
