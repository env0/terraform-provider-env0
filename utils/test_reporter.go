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

func (r *TestReporter) Error(args ...interface{}) {
	r.T.Error(args...)
}

func (r *TestReporter) Fail() {
	r.T.Fail()
}

func (r *TestReporter) FailNow() {
	r.T.FailNow()
}

func (r *TestReporter) Failed() bool {
	return r.T.Failed()
}

func (r *TestReporter) Fatal(args ...interface{}) {
	r.T.Fatal(args...)
}

func (r *TestReporter) Log(args ...interface{}) {
	r.T.Log(args...)
}

func (r *TestReporter) Logf(format string, args ...interface{}) {
	r.T.Logf(format, args...)
}

func (r *TestReporter) Name() string {
	return r.T.Name()
}

func (r *TestReporter) Parallel() {
	r.T.Parallel()
}

func (r *TestReporter) Skip(args ...interface{}) {
	r.T.Skip(args...)
}

func (r *TestReporter) SkipNow() {
	r.T.SkipNow()
}

func (r *TestReporter) Skipf(format string, args ...interface{}) {
	r.T.Skipf(format, args...)
}

func (r *TestReporter) Skipped() bool {
	return r.T.Skipped()
}

func (r *TestReporter) Helper() {
	r.T.Helper()
}

func (r *TestReporter) Fatalf(format string, args ...interface{}) {
	r.Log(fmt.Sprintf(format, args...))
	r.T.Fail()
	os.Exit(1)
}

func (r *TestReporter) Errorf(format string, args ...interface{}) {
	r.T.Errorf(format, args...)
}
