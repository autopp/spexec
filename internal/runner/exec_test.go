package runner

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec", func() {
	DescribeTable("Run()",
		// TODO: test about of signal
		func(command []string, status int, stdout, stderr string) {
			e := Exec{
				Command: command,
			}

			er := e.Run()
			Expect(er.Stdout).To(Equal([]byte(stdout)))
			Expect(er.Stderr).To(Equal([]byte(stderr)))

			st, sig, err := er.WaitStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(st).To(Equal(status))
			Expect(sig).To(BeNil())
		},
		Entry("with `echo -n 42`", []string{"echo", "-n", "42"}, 0, "42", ""),
	)
})
