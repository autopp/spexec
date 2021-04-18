package runner

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec", func() {
	DescribeTable("Run()",
		func(command []string, status int, stdout, stderr string) {
			e := Exec{
				Command: command,
			}

			expected := &ExecResult{
				Stdout: []byte(stdout),
				Stderr: []byte(stderr),
				Status: status,
			}
			Expect(e.Run()).To(Equal(expected))
		},
		Entry("with `echo -n 42`", []string{"echo", "-n", "42"}, 0, "42", ""),
	)
})
