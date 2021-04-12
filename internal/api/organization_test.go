package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/env0/terraform-provider-env0/internal/api"
)

var _ = Describe("Organization", func() {
	var organization Organization
	var organizationErr error

	JustBeforeEach(func() {
		organization, organizationErr = apiClient.Organization()
	})

	Describe("Fetch organization data", func() {
		When("Fetching the default organization of given api key", func() {
			It("should not fail", func() {
				Expect(organizationErr).To(BeNil())
			})
			It("Should have id set", func() {
				Expect(organization.Id).ToNot(BeEmpty())
			})
			It("Should have name set", func() {
				Expect(organization.Name).ToNot(BeEmpty())
			})
			It("Should not be self hosted", func() {
				Expect(organization.IsSelfHosted).To(BeFalse())
			})
		})
	})
})
