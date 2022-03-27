package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("APIKey Client", func() {
	mockAPIKey := APIKey{
		Id:           "id",
		Name:         "name",
		APIKeyId:     "keyid",
		APIKeySecret: "keysecret",
	}

	Describe("Get All APIKeys", func() {
		var returnedAPIKeys []APIKey
		mockAPIKeys := []APIKey{mockAPIKey}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/api-keys", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]APIKey) {
					*response = mockAPIKeys
				})
			returnedAPIKeys, _ = apiClient.APIKeys()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return APIKeys", func() {
			Expect(returnedAPIKeys).To(Equal(mockAPIKeys))
		})
	})

	Describe("Create APIKeys", func() {
		var createdAPIKey *APIKey
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			createAPIKeyPayload := APIKeyCreatePayload{}
			copier.Copy(&createAPIKeyPayload, &mockAPIKey)

			expectedCreateRequest := APIKeyCreatePayloadWith{
				APIKeyCreatePayload: createAPIKeyPayload,
				OrganizationId:      organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/api-keys", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *APIKey) {
					*response = mockAPIKey
				})

			createdAPIKey, err = apiClient.APIKeyCreate(createAPIKeyPayload)
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

		It("Should return created APIKey", func() {
			Expect(*createdAPIKey).To(Equal(mockAPIKey))
		})
	})

	Describe("Delete APIKey", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/api-keys/" + mockAPIKey.Id)
			apiClient.APIKeyDelete(mockAPIKey.Id)
		})

		It("Should send DELETE request with APIKey id", func() {
			httpCall.Times(1)
		})
	})
})
