package util

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

	It("with invalid format, returns err", func() {
		var actual interface{}
		Expect(UnmarshalJSON([]byte(`{"message": "hello"}}`), &actual)).To(HaveOccurred())
	})
})
