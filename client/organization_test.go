package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organization", func() {
	mockOrganization := Organization{
		Id:   organizationId,
		Name: "env0 🦄",
	}

	Describe("Organization", func() {
		var organization Organization
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				organizationsResult := []Organization{mockOrganization}
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Organization) {
						*response = organizationsResult
					})

				organization, err = apiClient.Organization()
			})

			It("Should send GET request once", func() {
				httpCall.Times(1)
			})

			It("Should return organization", func() {
				Expect(organization).Should(Equal(mockOrganization))
				Expect(err).Should(BeNil())
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations", nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Organization()
				Expect(expectedErr).Should(Equal(err))
			})

			It("On too many organizations return error", func() {
				organizationsResult := []Organization{mockOrganization, mockOrganization}
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Organization) {
						*response = organizationsResult
					})

				_, err = apiClient.Organization()
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("server responded with too many organizations"))
			})
		})
	})

	Describe("OrganizationPolicyUpdate", func() {
		hour12 := "12-h"
		t := true
		updatedMockOrganization := mockOrganization
		updatedMockOrganization.DoNotConsiderMergeCommitsForPrPlans = true
		updatedMockOrganization.EnableOidc = true
		updatedMockOrganization.DefaultTtl = &hour12

		var updatedOrganization *Organization
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)
				updateOrganizationPolicyPayload := OrganizationPolicyUpdatePayload{
					DefaultTtl:                          &hour12,
					DoNotConsiderMergeCommitsForPrPlans: &t,
					EnableOidc:                          &t,
				}

				httpCall = mockHttpClient.EXPECT().
					Post("/organizations/"+organizationId+"/policies", updateOrganizationPolicyPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *Organization) {
						*response = updatedMockOrganization
					})

				updatedOrganization, err = apiClient.OrganizationPolicyUpdate(updateOrganizationPolicyPayload)
			})

			It("Should send Post request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return organization received from API", func() {
				Expect(*updatedOrganization).To(Equal(updatedMockOrganization))
			})
		})
	})
})
