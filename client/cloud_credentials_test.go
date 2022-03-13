package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudCredentials", func() {
	const credentialsName = "credential_test"
	var apiKey ApiKey
	mockApiKey := ApiKey{
		Id:             "id1",
		Name:           "key1",
		OrganizationId: organizationId,
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	mockApiKeySecond := ApiKey{
		Id:             "id2",
		Name:           "key2",
		OrganizationId: organizationId,
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	keys := []ApiKey{mockApiKey, mockApiKeySecond}

	Describe("AwsCredentialsCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			payloadValue := AwsCredentialsValuePayload{
				RoleArn:    "role",
				ExternalId: "external",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", AwsCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *ApiKey) {
					*response = mockApiKey
				})

			apiKey, _ = apiClient.AwsCredentialsCreate(AwsCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
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

	Describe("GcpCredentialsCreate", func() {
		const gcpRequestType = "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT"
		mockGcpApiKey := mockApiKey
		mockGcpApiKey.Type = gcpRequestType
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			payloadValue := GcpCredentialsValuePayload{
				ProjectId:         "projectId",
				ServiceAccountKey: "serviceAccountKey",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", GcpCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           gcpRequestType,
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *ApiKey) {
					*response = mockGcpApiKey
				})

			apiKey, _ = apiClient.GcpCredentialsCreate(GcpCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  gcpRequestType,
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(apiKey).To(Equal(mockGcpApiKey))
		})
	})

	Describe("CloudCredentialsDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/credentials/" + mockApiKey.Id)
			apiClient.CloudCredentialsDelete(mockApiKey.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})
	})

	Describe("CloudCredentials", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApiKey) {
					*response = keys
				})
			apiKey, _ = apiClient.CloudCredentials(mockApiKey.Id)
		})

		It("Should send GET request with project id", func() {
			httpCall.Times(1)
		})

		It("Should return correct key", func() {
			Expect(apiKey).To(Equal(mockApiKey))
		})
	})

	Describe("CloudCredentialsList", func() {
		var credentials []ApiKey

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApiKey) {
					*response = keys
				})
			credentials, _ = apiClient.CloudCredentialsList()
		})

		It("Should send GET request with organization id param", func() {
			httpCall.Times(1)
		})

		It("Should return all credentials", func() {
			Expect(credentials).To(Equal(keys))
		})
	})
})
