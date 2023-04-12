package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpecResult", func() {
	Describe("NewSpecResult()", func() {
		It("returns SpecResult with summary", func() {
			trs := []*TestResult{
				{
					Name:      "test1",
					IsSuccess: true,
				},
				{
					Name:      "test2",
					IsSuccess: false,
				},
				{
					Name:      "test3",
					IsSuccess: false,
				},
			}
			sr := NewSpecResult("test.yaml", trs)
			Expect(sr.Summary).To(Equal(SpecSummary{
				NumberOfTests:     3,
				NumberOfSucceeded: 1,
				NumberOfFailed:    2,
			}))
		})
	})

	Describe("GetFailedTestResults()", func() {
		It("returns failed test results only", func() {
			trs := []*TestResult{
				{
					Name:      "test1",
					IsSuccess: true,
				},
				{
					Name:      "test2",
					IsSuccess: false,
				},
				{
					Name:      "test3",
					IsSuccess: false,
				},
			}
			sr := NewSpecResult("test.yaml", trs)
			Expect(sr.GetFailedTestResults()).To(Equal([]*TestResult{trs[1], trs[2]}))
		})
	})
})
