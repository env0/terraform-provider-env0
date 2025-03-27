package utils

import (
	"fmt"
	"os"
	"testing"
)

/*
This is a test reporter overriding testing.T's Fatalf to not runtime.Goexit() but rather os.Exit(1)
See issue in: https://github.com/golang/mock/issues/145
*/

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

func (r *TestReporter) Fatalf(format string, args ...any) {
	r.T.Helper()
	r.Log(fmt.Sprintf(format, args...))
	r.T.Fail()
	os.Exit(1)
}

func (r *TestReporter) Errorf(format string, args ...any) {
	r.T.Helper()
	r.T.Errorf(format, args...)
}

func (r *TestReporter) Cleanup(cb func()) {
	r.T.Helper()
	cb()
}
