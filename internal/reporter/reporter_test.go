package reporter

import (
	"bytes"

	"github.com/autopp/spexec/internal/model"
	g "github.com/onsi/ginkgo/v2" // Reporter are duplicated
	. "github.com/onsi/gomega"
)

type testReportFormatter struct {
	OnRunStartCalled     int
	OnTestStartCalled    int
	OnTestCompleteCalled int
	OnRunCompleteCalled  int
}

func (rf *testReportFormatter) OnRunStart(w *Writer) {
	rf.OnRunStartCalled++
}

func (rf *testReportFormatter) OnTestStart(w *Writer, t *model.Test) {
	rf.OnTestStartCalled++
}

func (rf *testReportFormatter) OnTestComplete(w *Writer, t *model.Test, tr *model.TestResult) {
	rf.OnTestCompleteCalled++
}

func (rf *testReportFormatter) OnRunComplete(w *Writer, sr *model.SpecResult) {
	rf.OnRunCompleteCalled++
}

var _ = g.Describe("Rerporter", func() {
	var rf *testReportFormatter
	var r *Reporter

	g.JustBeforeEach(func() {
		rf = new(testReportFormatter)
		r, _ = New(WithWriter(&bytes.Buffer{}), WithFormatter(rf))
	})

	g.Describe("OnRunStart()", func() {
		g.It("calls OnRunStart() of formatter", func() {
			r.OnRunStart()
			Expect(rf.OnRunStartCalled).To(Equal(1))
		})
	})

	g.Describe("OnTestStart()", func() {
		g.It("calls OnTestStart() of formatter", func() {
			r.OnTestStart(nil)
			Expect(rf.OnTestStartCalled).To(Equal(1))
		})
	})

	g.Describe("OnTestComplete()", func() {
		g.It("calls OnTestComplete() of formatter", func() {
			r.OnTestComplete(nil, nil)
			Expect(rf.OnTestCompleteCalled).To(Equal(1))
		})
	})

	g.Describe("OnRunComplete()", func() {
		g.It("calls OnRunComplete() of formatter", func() {
			r.OnRunComplete(nil)
			Expect(rf.OnRunCompleteCalled).To(Equal(1))
		})
	})
})
