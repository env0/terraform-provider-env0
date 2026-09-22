package utils

import (
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

// recordingReporter is a TestReporter whose Errorf is captured instead of failing the test, so a
// spec can assert on what gomock reported. Everything else, Cleanup included, is the real thing.
type recordingReporter struct {
	*TestReporter
	errors []string
}

func (r *recordingReporter) Errorf(format string, args ...any) {
	r.errors = append(r.errors, format)
}

// A minimal gomock mock. The generated client mocks live in package client, which this package
// cannot import, so the protocol is spelled out by hand.
type fakeMock struct {
	ctrl *gomock.Controller
}

type fakeMockRecorder struct {
	mock *fakeMock
}

func (m *fakeMock) EXPECT() *fakeMockRecorder {
	return &fakeMockRecorder{mock: m}
}

func (m *fakeMock) Do() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Do")
}

func (r *fakeMockRecorder) Do() *gomock.Call {
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Do", reflect.TypeOf((*fakeMock)(nil).Do))
}

func TestReporterCleanupDefersTheCallback(t *testing.T) {
	ran := false

	t.Run("callback", func(t *testing.T) {
		reporter := &TestReporter{T: t}
		reporter.Cleanup(func() { ran = true })

		if ran {
			t.Fatal("Cleanup ran the callback immediately instead of deferring it to the end of the test")
		}
	})

	if !ran {
		t.Fatal("Cleanup never ran the callback")
	}
}

// gomock registers ctrl.Finish through Cleanup. Running that callback on the spot finishes the
// controller before a single EXPECT is recorded, which is what made every Times(N) in the
// acceptance suite pass without the call ever happening.
func TestReporterReportsAnUnmetExpectation(t *testing.T) {
	t.Run("unmet expectation", func(t *testing.T) {
		reporter := &recordingReporter{TestReporter: &TestReporter{T: t}}

		// Cleanups run last-in-first-out, so registering this before the controller's own
		// cleanup is what puts the assertion after gomock has finished.
		t.Cleanup(func() {
			if len(reporter.errors) == 0 {
				t.Error("gomock reported no missing call for an expectation that was never met")
			}
		})

		ctrl := gomock.NewController(reporter)
		mock := &fakeMock{ctrl: ctrl}
		mock.EXPECT().Do().Times(1)
	})
}

func TestReporterAcceptsAMetExpectation(t *testing.T) {
	reporter := &recordingReporter{TestReporter: &TestReporter{T: t}}

	t.Cleanup(func() {
		if len(reporter.errors) > 0 {
			t.Errorf("gomock reported %v for an expectation that was met", reporter.errors)
		}
	})

	ctrl := gomock.NewController(reporter)
	mock := &fakeMock{ctrl: ctrl}
	mock.EXPECT().Do().Times(1)
	mock.Do()
}
