package util

import (
	"encoding/json"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnmarshalJSON", func() {
	DescribeTable("succeess cases",
		func(given string, expected any) {
			var actual any
			Expect(DecodeJSON(strings.NewReader(given), &actual)).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		},
		Entry(
			`with object, stores map[string]any`,
			`{"message": "hello"}`, map[string]any{"message": "hello"},
		),
		Entry(
			`with array, stores []any`,
			`[true, false]`, []any{true, false},
		),
		Entry(
			`with number, stores json.Number`,
			`42`, json.Number("42"),
		),
	)

	DescribeTable("failure cases",
		func(given string) {
			var actual any
			Expect(DecodeJSON(strings.NewReader(given), &actual)).To(HaveOccurred())
		},
		Entry("with invalid format, returns err", `{message: "hello"}`),
		Entry("with extra character, returns err", `{"message": "hello"}}`),
	)
})
