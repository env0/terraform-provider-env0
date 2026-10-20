package client_test

import (
	"encoding/json"

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

	Describe("Get All ApiKeys with an ip allowlist", func() {
		var returnedApiKeys []ApiKey

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/api-keys", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]ApiKey) {
					Expect(json.Unmarshal([]byte(`[{"id":"id","ipAllowlist":[{"type":"subnet","network":"0.0.0.0","prefix":0,"family":"ipv4"}]}]`), response)).To(Succeed())
				})
			returnedApiKeys, _ = apiClient.ApiKeys()
		})

		It("Should decode the ip allowlist", func() {
			prefix := 0
			Expect(returnedApiKeys[0].IpAllowlist).To(Equal([]IpAllowlistRule{
				{Type: "subnet", Network: "0.0.0.0", Prefix: &prefix, Family: "ipv4"},
			}))
		})
	})

	Describe("Update ApiKey ip allowlist", func() {
		var (
			updatedApiKey *ApiKey
			err           error
		)

		prefix := 0
		rules := []IpAllowlistRule{
			{Type: "address", Address: "198.51.100.7", Family: "ipv4"},
			{Type: "subnet", Network: "0.0.0.0", Prefix: &prefix, Family: "ipv4"},
		}
		mockApiKeyWithAllowlist := mockApiKey
		mockApiKeyWithAllowlist.IpAllowlist = rules

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Put("/api-keys/"+mockApiKey.Id+"/ip-allowlist", ApiKeyUpdateIpAllowlistPayload{IpAllowlist: rules}, gomock.Any()).
				Do(func(path string, request any, response *ApiKey) {
					*response = mockApiKeyWithAllowlist
				})

			updatedApiKey, err = apiClient.ApiKeyUpdateIpAllowlist(mockApiKey.Id, rules)
		})

		It("Should send PUT request with the rules", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return the updated ApiKey", func() {
			Expect(*updatedApiKey).To(Equal(mockApiKeyWithAllowlist))
		})

		// The server rejects a field that doesn't belong to the rule's type, and a missing prefix.
		It("Should send only each rule's own fields, and prefix 0", func() {
			body, _ := json.Marshal(ApiKeyUpdateIpAllowlistPayload{IpAllowlist: rules})
			Expect(string(body)).To(Equal(`{"ipAllowlist":[{"type":"address","address":"198.51.100.7","family":"ipv4"},{"type":"subnet","network":"0.0.0.0","prefix":0,"family":"ipv4"}]}`))
		})

		It("Should send an empty list as []", func() {
			body, _ := json.Marshal(ApiKeyUpdateIpAllowlistPayload{IpAllowlist: []IpAllowlistRule{}})
			Expect(string(body)).To(Equal(`{"ipAllowlist":[]}`))
		})
	})

	Describe("Create ApiKeys", func() {
		var (
			createdApiKey *ApiKey
			err           error
		)

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
		var (
			returnedOidcSub string
			err             error
		)

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
