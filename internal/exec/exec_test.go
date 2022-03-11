package exec

import (
	"errors"
	"time"

	"github.com/autopp/spexec/internal/util"
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

type testOption struct {
	ret    error
	called bool
}

func (o *testOption) Apply(*Exec) error {
	o.called = true
	return o.ret
}

var _ = Describe("New()", func() {
	Describe("without option", func() {
		It("returns new Exec", func() {
			env := []util.StringVar{{Name: "ANSWER", Value: "42"}}
			e, err := New([]string{"echo", "hello"}, "/tmp", nil, env)

			Expect(err).NotTo(HaveOccurred())
			Expect(e).To(Equal(&Exec{
				Command: []string{"echo", "hello"},
				Dir:     "/tmp",
				Stdin:   nil,
				Env:     env,
				Timeout: defaultTimeout,
			}))
		})
	})

	Describe("with options (all succeeded)", func() {
		It("invokes given options and returns new Exec", func() {
			o1 := &testOption{ret: nil}
			o2 := &testOption{ret: nil}
			env := []util.StringVar{{Name: "ANSWER", Value: "42"}}
			e, err := New([]string{"echo", "hello"}, "/tmp", nil, env, o1, o2)

			Expect(err).NotTo(HaveOccurred())
			Expect(e).To(Equal(&Exec{
				Command: []string{"echo", "hello"},
				Dir:     "/tmp",
				Stdin:   nil,
				Env:     env,
				Timeout: defaultTimeout,
			}))
			Expect(o1.called).To(BeTrue())
			Expect(o2.called).To(BeTrue())
		})
	})

	Describe("with options (failed at first)", func() {
		It("invokes given options and returns new Exec", func() {
			o1 := &testOption{ret: errors.New("option error")}
			o2 := &testOption{ret: nil}
			env := []util.StringVar{{Name: "ANSWER", Value: "42"}}
			_, err := New([]string{"echo", "hello"}, "/tmp", nil, env, o1, o2)

			Expect(err).To(MatchError("option error"))
			Expect(o1.called).To(BeTrue())
			Expect(o2.called).To(BeFalse())
		})
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
