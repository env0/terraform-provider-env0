package utils

import (
	"fmt"
	"os"
	"testing"
)

// TestReporter adapts *testing.T to gomock and to the terraform-plugin-sdk test framework. Only
// Fatalf and Cleanup deviate from testing.T; see those methods for why.
type TestReporter struct {
	T *testing.T
}

func (r *TestReporter) Error(args ...any) {
	r.T.Helper()
	r.T.Error(args...)
}

func (r *TestReporter) Fail() {
	r.T.Helper()
	r.T.Fail()
}

func (r *TestReporter) FailNow() {
	r.T.Helper()
	r.T.FailNow()
}

func (r *TestReporter) Failed() bool {
	return r.T.Failed()
}

func (r *TestReporter) Fatal(args ...any) {
	r.T.Helper()
	r.T.Fatal(args...)
}

func (r *TestReporter) Log(args ...any) {
	r.T.Log(args...)
}

func (r *TestReporter) Logf(format string, args ...any) {
	r.T.Logf(format, args...)
}

func (r *TestReporter) Name() string {
	return r.T.Name()
}

func (r *TestReporter) Parallel() {
	r.T.Parallel()
}

func (r *TestReporter) Skip(args ...any) {
	r.T.Skip(args...)
}

func (r *TestReporter) SkipNow() {
	r.T.SkipNow()
}

func (r *TestReporter) Skipf(format string, args ...any) {
	r.T.Skipf(format, args...)
}

func (r *TestReporter) Skipped() bool {
	return r.T.Skipped()
}

func (r *TestReporter) Helper() {
	r.T.Helper()
}

// Fatalf must not return: gomock's Controller.Call dereferences the match it never got as soon as
// Fatalf hands control back, so the unexpected-call message is replaced by a nil-pointer panic.
// Goexit is not an option either - gomock calls this from the provider's gRPC goroutine, not the
// test's, so it would abandon the test rather than fail it (golang/mock#145, still open on
// go.uber.org/mock). The message goes to stderr because os.Exit skips the test's buffered output,
// which is what used to leave an aborted run printing nothing but "exit status 1".
func (r *TestReporter) Fatalf(format string, args ...any) {
	r.T.Helper()
	fmt.Fprintf(os.Stderr, "%s: %s\n", r.T.Name(), fmt.Sprintf(format, args...))
	r.T.Fail()
	os.Exit(1)
}

func (r *TestReporter) Errorf(format string, args ...any) {
	r.T.Helper()
	r.T.Errorf(format, args...)
}

// Cleanup must register cb rather than run it. gomock.NewController uses this method to defer
// ctrl.Finish() to the end of the test; running cb on the spot finishes the controller before
// a single expectation is recorded, so every Times/MinTimes assertion passes trivially.
func (r *TestReporter) Cleanup(cb func()) {
	r.T.Helper()
	r.T.Cleanup(cb)
}
