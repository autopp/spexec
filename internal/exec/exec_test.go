package exec

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithTimeout", func() {
	Context("when the given is positive", func() {
		It("sets .Timeout", func() {
			e := &Exec{}
			err := WithTimeout(42 * time.Second).Apply(e)

			Expect(err).NotTo(HaveOccurred())
			Expect(e.Timeout).To(Equal(42 * time.Second))
		})
	})

	Context("when the given is not positive", func() {
		It("dose not set .Timeout", func() {
			e := &Exec{}
			err := WithTimeout(-42 * time.Second).Apply(e)

			Expect(err).NotTo(HaveOccurred())
			Expect(e.Timeout).To(Equal(time.Duration(0)))
		})
	})
})

var _ = Describe("WithTeeStdout", func() {
	It("sets .TeeStdout", func() {
		e := &Exec{}
		err := WithTeeStdout(true).Apply(e)

		Expect(err).NotTo(HaveOccurred())
		Expect(e.TeeStdout).To(BeTrue())
	})
})

var _ = Describe("WithTeeStderr", func() {
	It("sets .TeeStderr", func() {
		e := &Exec{}
		err := WithTeeStderr(true).Apply(e)

		Expect(err).NotTo(HaveOccurred())
		Expect(e.TeeStderr).To(BeTrue())
	})
})

var _ = Describe("Exec", func() {
	DescribeTable("Run()",
		func(e *Exec, exited bool, status int, signal string, stdout, stderr string, isTimeout bool) {
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
				if !isTimeout {
					Expect(er.Signal.String()).To(Equal(signal))
				}
			}
			Expect(er.IsTimeout).To(Equal(isTimeout))
		},
		Entry("with `echo -n 42`",
			&Exec{
				Command: []string{"echo", "-n", "42"},
			},
			true, 0, "", "42", "", false,
		),
		Entry("with `echo -n 42 >&2`",
			&Exec{
				Command: []string{"testdata/stderr.sh"},
			},
			true, 0, "", "", "42", false,
		),
		Entry("with `echo -n 42 >&2 (in ./testdata)`",
			&Exec{
				Command: []string{"./stderr.sh"},
				Dir:     "testdata",
			},
			true, 0, "", "", "42", false,
		),
		Entry("with `kill -TERM $$`",
			&Exec{
				Command: []string{"testdata/signal.sh"},
			},
			false, 0, "terminated", "", "", false,
		),
		Entry("with `sleep 1` and 1ms timeout",
			&Exec{
				Command: []string{"sleep", "1"},
				Timeout: 1 * time.Millisecond,
			},
			false, 0, "", "", "", true,
		),
	)
})
