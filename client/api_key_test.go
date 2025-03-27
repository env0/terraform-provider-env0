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
			mockOrganizationIdCall()
			httpCall = mockHttpClient.EXPECT().
				Get("/api-keys", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]ApiKey) {
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
			mockOrganizationIdCall()

			createApiKeyPayload := ApiKeyCreatePayload{}

			_ = copier.Copy(&createApiKeyPayload, &mockApiKey)

			expectedCreateRequest := ApiKeyCreatePayloadWith{
				ApiKeyCreatePayload: createApiKeyPayload,
				OrganizationId:      organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/api-keys", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request any, response *ApiKey) {
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
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/api-keys/"+mockApiKey.Id, nil)
			err = apiClient.ApiKeyDelete(mockApiKey.Id)
		})

		It("Should send DELETE request with ApiKey id", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Oidc Sub", func() {
		var returnedOidcSub string
		var err error
		mockedOidcSub := "oidc sub 1234"

		BeforeEach(func() {
			mockOrganizationIdCall()
			httpCall = mockHttpClient.EXPECT().
				Get("/api-keys/oidc-sub", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *string) {
					*response = mockedOidcSub
				})
			returnedOidcSub, err = apiClient.OidcSub()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return Oidc sub", func() {
			Expect(returnedOidcSub).To(Equal(mockedOidcSub))
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
