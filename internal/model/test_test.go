package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test", func() {
	DescribeTable("GetName()",
		func(test *Test, expected string) {
			Expect(test.GetName()).To(Equal(expected))
		},
		Entry("Name is not empty", &Test{Name: "test of echo", Command: []string{"echo", "hello"}}, "test of echo"),
		Entry("Name is empty", &Test{Name: "", Command: []string{"echo", "hello"}},
			"echo hello"),
		Entry("Name is empty and Env is given", &Test{
			Name:    "",
			Command: []string{"make", "build"},
			Env: []struct {
				Name  string
				Value string
			}{
				{Name: "GOOS", Value: "linux"},
				{Name: "GOARCH", Value: "amd64"},
			},
		}, "GOOS=linux GOARCH=amd64 make build"),
	)
})
