package exec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec", func() {
	DescribeTable("Run()",
		// TODO: test about of signal
		func(e *Exec, exited bool, status int, signal string, stdout, stderr string) {
			if e.Timeout == 0 {
				e.Timeout = defaultTimeout
			}
			er := e.Run()
			Expect(er.Stdout).To(Equal([]byte(stdout)))
			Expect(er.Stderr).To(Equal([]byte(stderr)))

			Expect(er.Err).NotTo(HaveOccurred())
			if exited {
				Expect(er.Status).To(Equal(status))
			} else {
				Expect(er.Signal.String()).To(Equal(signal))
			}
		},
		Entry("with `echo -n 42`",
			&Exec{
				Command: []string{"echo", "-n", "42"},
			},
			true, 0, "", "42", "",
		),
		Entry("with `echo -n 42 >&2`",
			&Exec{
				Command: []string{"testdata/stderr.sh"},
			},
			true, 0, "", "", "42",
		),
		Entry("with `kill -TERM $$`",
			&Exec{
				Command: []string{"testdata/signal.sh"},
			},
			false, 0, "terminated", "", "",
		),
	)
})
