package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const awsCredentialsName = "credential_test"

var _ = Describe("AwsCredentials", func() {
	var apiKey ApiKey
	mockApiKey := ApiKey{
		Id:             "id1",
		Name:           "key1",
		OrganizationId: organizationId,
		Type: "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	mockApiKeySecond := ApiKey{
		Id:             "id2",
		Name:           "key2",
		OrganizationId: organizationId,
		Type: "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	keys := []ApiKey{mockApiKey, mockApiKeySecond}

	Describe("AwsCredentialsCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", AwsCredentialsCreatePayload{
					Name:           awsCredentialsName,
					OrganizationId: organizationId,
					Type: "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value: AwsCredentialsValuePayload{
						RoleArn: "role",
						ExternalId: "external",
					},
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *ApiKey) {
					*response = mockApiKey
				})

			apiKey, _ = apiClient.AwsCredentialsCreate(AwsCredentialsCreatePayload{
				Name: awsCredentialsName,
				Value: AwsCredentialsValuePayload{
					RoleArn: "role",
					ExternalId: "external",
				},
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(apiKey).To(Equal(mockApiKey))
		})
	})

	Describe("AwsCredentialsDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/credentials/" + mockApiKey.Id)
			apiClient.AwsCredentialsDelete(mockApiKey.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})
	})

	Describe("AwsCredentials", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials",  map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApiKey) {
					*response = keys
				})
			apiKey, _ = apiClient.AwsCredentials(mockApiKey.Id)
		})

		It("Should send GET request with project id", func() {
			httpCall.Times(1)
		})

		It("Should return correct key", func() {
			Expect(apiKey).To(Equal(mockApiKey))
		})
	})

	Describe("AwsCredentialsList", func() {
		var credentials []ApiKey

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApiKey) {
					*response = keys
				})
			credentials, _ = apiClient.AwsCredentialsList()
		})

		It("Should send GET request with organization id param", func() {
			httpCall.Times(1)
		})

		It("Should return all credentials", func() {
			Expect(credentials).To(Equal(keys))
		})
	})
})
