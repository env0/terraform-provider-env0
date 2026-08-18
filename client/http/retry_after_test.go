package http_test

import (
	"net/http"
	"time"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseRetryAfter", func() {
	It("should parse a number of seconds", func() {
		delay, ok := httpModule.ParseRetryAfter("30")

		Expect(ok).To(BeTrue())
		Expect(delay).To(Equal(30 * time.Second))
	})

	It("should parse an HTTP date", func() {
		date := time.Now().UTC().Add(42 * time.Second).Format(http.TimeFormat)

		delay, ok := httpModule.ParseRetryAfter(date)

		Expect(ok).To(BeTrue())
		// The header has a one second resolution, so the parsed delay is 41 or 42 seconds.
		Expect(delay).To(BeNumerically(">", 40*time.Second))
		Expect(delay).To(BeNumerically("<=", 42*time.Second))
	})

	DescribeTable("values that carry no usable delay",
		func(header string) {
			delay, ok := httpModule.ParseRetryAfter(header)

			Expect(ok).To(BeFalse())
			Expect(delay).To(Equal(time.Duration(0)))
		},
		Entry("missing header", ""),
		Entry("unparsable value", "soon"),
		Entry("zero seconds", "0"),
		Entry("negative seconds", "-5"),
		Entry("date in the past", time.Now().UTC().Add(-time.Minute).Format(http.TimeFormat)),
	)
})
