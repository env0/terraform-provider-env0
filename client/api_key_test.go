package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("ApiKey Client", func() {
	mockApiKey := ApiKey{
		Id:           "id",
		Name:         "name",
		ApiKeyId:     "keyid",
		ApiKeySecret: "keysecret",
	}

	Describe("Get All ApiKeys", func() {
		var returnedApiKeys []ApiKey
		mockApiKeys := []ApiKey{mockApiKey}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/api-keys", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApiKey) {
					*response = mockApiKeys
				})
			returnedApiKeys, _ = apiClient.ApiKeys()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return ApiKeys", func() {
			Expect(returnedApiKeys).To(Equal(mockApiKeys))
		})
	})

	Describe("Create ApiKeys", func() {
		var createdApiKey *ApiKey
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			createApiKeyPayload := ApiKeyCreatePayload{}
			copier.Copy(&createApiKeyPayload, &mockApiKey)

			expectedCreateRequest := ApiKeyCreatePayloadWith{
				ApiKeyCreatePayload: createApiKeyPayload,
				OrganizationId:      organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/api-keys", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *ApiKey) {
					*response = mockApiKey
				})

			createdApiKey, err = apiClient.ApiKeyCreate(createApiKeyPayload)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created ApiKey", func() {
			Expect(*createdApiKey).To(Equal(mockApiKey))
		})
	})

	Describe("Delete ApiKey", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/api-keys/"+mockApiKey.Id, nil)
			apiClient.ApiKeyDelete(mockApiKey.Id)
		})

		It("Should send DELETE request with ApiKey id", func() {
			httpCall.Times(1)
		})
	})
})
