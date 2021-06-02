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

			It("Should return ssh keys", func() {
				Expect(organization).Should(Equal(mockOrganization))
				Expect(err).Should(BeNil())
			})
		})

		Describe("Failure", func() {
			It("On error from server", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations", nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Organization()
				Expect(expectedErr).Should(Equal(err))
			})

			It("On error from server", func() {
				organizationsResult := []Organization{mockOrganization, mockOrganization}
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Organization) {
						*response = organizationsResult
					})

				_, err = apiClient.Organization()
				Expect(err).ShouldNot(BeNil())
			})
		})
	})
})
