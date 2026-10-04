package http_test

import (
	httpModule "github.com/env0/terraform-provider-env0/client/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("PerPathRequestLimit", func() {
	limitFor := httpModule.PerPathRequestLimit(950, 190)

	DescribeTable("limit for a key",
		func(key string, expected int) {
			Expect(limitFor(key)).To(Equal(expected))
		},
		Entry("listing environments", "GET /environments", 190),
		Entry("reading an environment", "GET /environments/abc", 950),
		Entry("reading a deployment", "GET /environments/deployments/abc", 950),
		Entry("creating an environment", "POST /environments", 190),
		Entry("deploying an environment", "POST /environments/abc/deployments", 190),
		Entry("updating a deployment", "PUT /environments/deployments/abc", 190),
		Entry("updating an environment", "PUT /environments/abc", 950),
		Entry("deleting on the environments path", "DELETE /environments", 950),
		Entry("creating an environment without a template", "POST /environments/without-template", 190),
		Entry("updating a nested deployment", "PUT /x/environments/deployments/abc", 190),
		Entry("posting to a path that starts with environments", "POST /environmentsX", 190),
		Entry("listing a path that starts with environments", "GET /environmentsX", 950),
		Entry("any other path", "GET /projects/abc", 950),
	)
})
