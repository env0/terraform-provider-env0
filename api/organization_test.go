package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/env0/terraform-provider-env0/api"
)

var _ = Describe("Organization", func() {
	var client *ApiClient
	var organization Organization

	BeforeEach(func() {
		var err error
		client, err = NewClientFromEnv()
		Expect(err).To(BeNil())
		Expect(client).ToNot(BeNil())
	})

	JustBeforeEach(func() {
		var err error
		organization, err = client.Organization()
		Expect(err).To(BeNil())
	})

	Describe("Fetch organization data", func() {
		When("Fetching the default organization of given api key", func() {
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
