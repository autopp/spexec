package runner

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec", func() {
	DescribeTable("Run()",
		// TODO: test about of signal
		func(e *Exec, status int, stdout, stderr string) {
			if e.Timeout == 0 {
				e.Timeout = defaultTimeout
			}
			er := e.Run()
			Expect(er.Stdout).To(Equal([]byte(stdout)))
			Expect(er.Stderr).To(Equal([]byte(stderr)))

			st, sig, err := er.WaitStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(st).To(Equal(status))
			Expect(sig).To(BeNil())
		},
		Entry("with `echo -n 42`",
			&Exec{
				Command: []string{"echo", "-n", "42"},
			},
			0, "42", "",
		),
	)
})
