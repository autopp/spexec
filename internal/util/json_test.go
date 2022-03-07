package util

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnmarshalJSON", func() {
	DescribeTable("succeess cases",
		func(given string, expected interface{}) {
			var actual interface{}
			Expect(UnmarshalJSON([]byte(given), &actual)).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		},
		Entry(
			`with object, stores map[string]interface{}`,
			`{"message": "hello"}`, map[string]interface{}{"message": "hello"},
		),
		Entry(
			`with array, stores []interface{}`,
			`[true, false]`, []interface{}{true, false},
		),
		Entry(
			`with number, stores json.Number`,
			`42`, json.Number("42"),
		),
	)

	DescribeTable("failure cases",
		func(given string) {
			var actual interface{}
			Expect(UnmarshalJSON([]byte(given), &actual)).To(HaveOccurred())
		},
		Entry("with invalid format, returns err", `{message: "hello"}`),
		Entry("with extra character, returns err", `{"message": "hello"}}`),
	)
})
