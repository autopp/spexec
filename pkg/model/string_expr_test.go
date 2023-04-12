package model

import (
	"errors"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("literalStringExpr", func() {
	literal := NewLiteralStringExpr("hello")

	Describe("String()", func() {
		It("returns itself", func() {
			Expect(literal.String()).To(Equal("hello"))
		})
	})

	Describe("Eval()", func() {
		It("returns itself", func() {
			v, cleanup, err := literal.Eval()
			Expect(v).To(Equal("hello"))
			Expect(cleanup).To(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("envStringExpr", func() {
	name := "MESSAGE"
	value := "hello"

	e := NewEnvStringExpr(name)

	BeforeEach(func() {
		oldValue, setAlready := os.LookupEnv(name)
		os.Setenv(name, value)

		DeferCleanup(func() {
			if setAlready {
				os.Setenv(name, oldValue)
			} else {
				os.Unsetenv(name)
			}
		})
	})

	Describe("String()", func() {
		It("returns itself with '$' prefix", func() {
			Expect(e.String()).To(Equal("$" + name))
		})
	})

	Describe("Eval()", func() {
		It("returns value of the environment variable", func() {
			v, cleanup, err := e.Eval()
			Expect(v).To(Equal(value))
			Expect(cleanup).To(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when given name is not defined", func() {
			_, _, err := NewEnvStringExpr("SPEXEC_UNDEFINED").Eval()
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("fileStringExpr", func() {
	pattern := "*.yaml"
	contents := "hello"

	var file StringExpr
	JustBeforeEach(func() {
		file = NewFileStringExpr(pattern, contents)
	})

	Describe("String()", func() {
		It("returns dummy file path", func() {
			Expect(file.String()).To(SatisfyAll(HavePrefix(os.TempDir()), HaveSuffix(".yaml")))
		})
	})

	Describe("Eval()", func() {
		It("returns the file name contains the given contents and function for remove the file", func() {
			v, cleanup, err := file.Eval()
			Expect(v).NotTo(BeEmpty())
			Expect(cleanup).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			Expect(os.ReadFile(v)).To(Equal([]byte(contents)))

			Expect(cleanup()).NotTo(HaveOccurred())
			_, err = os.Stat(v)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})
	})
})

type testStringExpr struct {
	v              string
	isEvaled       bool
	successEval    bool
	isCleanuped    bool
	successCleanup bool
}

func (e *testStringExpr) Eval() (string, func() error, error) {
	e.isEvaled = true
	if !e.successEval {
		return "", e.getCleanup(), errors.New(e.v)
	}

	return e.v, e.getCleanup(), nil
}

func (e *testStringExpr) getCleanup() func() error {
	return func() error {
		e.isCleanuped = true
		if !e.successCleanup {
			return errors.New(e.v)
		}

		return nil
	}
}

func (e *testStringExpr) String() string {
	return e.v
}

func extractBools(exprs []StringExpr, f func(e *testStringExpr) bool) []bool {
	results := make([]bool, len(exprs))

	for i, e := range exprs {
		results[i] = f(e.(*testStringExpr))
	}
	return results
}

func (e *testStringExpr) stringExpr() {}

var _ = Describe("EvalStringExprs()", func() {
	DescribeTable("returns results, aggregated cleanup function, and errors",
		func(successEvalAndCleanups [][2]bool, expectedValues []string, expectedEvaled []bool, expectedCleanuped []bool, expectedErr string, expectedErrIndex int, expectedCleanupErrs []string) {
			exprs := make([]StringExpr, len(successEvalAndCleanups))
			for i, fields := range successEvalAndCleanups {
				expr := &testStringExpr{v: strconv.Itoa(i), successEval: fields[0], successCleanup: fields[1]}
				exprs[i] = expr
			}

			values, cleanup, err, i := EvalStringExprs(exprs)

			if expectedErr == "" {
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(Equal(expectedValues))
			} else {
				Expect(err).To(MatchError(expectedErr))
				Expect(values).To(BeNil())
			}
			Expect(i).To(Equal(expectedErrIndex))

			Expect(extractBools(exprs, func(e *testStringExpr) bool { return e.isEvaled })).To(Equal(expectedEvaled))
			cleanupErrs := make([]string, 0)
			for _, err := range cleanup() {
				cleanupErrs = append(cleanupErrs, err.Error())
			}
			Expect(cleanupErrs).To(Equal(expectedCleanupErrs))
			Expect(extractBools(exprs, func(e *testStringExpr) bool { return e.isCleanuped })).To(Equal(expectedCleanuped))
		},
		Entry("with 3 exprs all success", [][2]bool{{true, true}, {true, true}, {true, true}}, []string{"0", "1", "2"}, []bool{true, true, true}, []bool{true, true, true}, "", -1, []string{}),
		Entry("with 3 exprs 2nd fails", [][2]bool{{true, true}, {false, true}, {true, true}}, nil, []bool{true, true, false}, []bool{true, true, false}, "1", 1, []string{}),
		Entry("with 3 exprs 1st and 3rd cleanup fails", [][2]bool{{true, false}, {true, true}, {true, false}}, []string{"0", "1", "2"}, []bool{true, true, true}, []bool{true, true, true}, "", -1, []string{"0", "2"}),
	)
})
