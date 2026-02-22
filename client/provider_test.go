package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Provider Client", func() {
	mockProvider := Provider{
		Id:          "id",
		Type:        "type",
		Description: "description",
	}

	Describe("Get Provider", func() {
		var returnedProvider *Provider

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/providers/"+mockProvider.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request any, response *Provider) {
					*response = mockProvider
				}).Times(1)
			returnedProvider, _ = apiClient.Provider(mockProvider.Id)
		})

		It("Should return provider", func() {
			Expect(*returnedProvider).To(Equal(mockProvider))
		})
	})

	Describe("Get All Providers", func() {
		var returnedProviders []Provider

		mockProviders := []Provider{mockProvider}

		BeforeEach(func() {
			mockOrganizationIdCall().Times(1)
			httpCall = mockHttpClient.EXPECT().
				Get("/providers", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]Provider) {
					*response = mockProviders
				}).Times(1)
			returnedProviders, _ = apiClient.Providers()
		})

		It("Should return providers", func() {
			Expect(returnedProviders).To(Equal(mockProviders))
		})
	})

	Describe("Create Provider", func() {
		var createdProvider *Provider

		BeforeEach(func() {
			mockOrganizationIdCall().Times(1)

			createProviderPayload := ProviderCreatePayload{
				Type:        mockProvider.Type,
				Description: mockProvider.Description,
			}

			expectedCreateRequest := struct {
				ProviderCreatePayload
				OrganizationId string `json:"organizationId"`
			}{
				createProviderPayload,
				organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/providers", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request any, response *Provider) {
					*response = mockProvider
				}).Times(1)

			createdProvider, _ = apiClient.ProviderCreate(createProviderPayload)
		})
		It("Should return created provider", func() {
			Expect(*createdProvider).To(Equal(mockProvider))
		})
	})

	Describe("Delete Provider", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/providers/"+mockProvider.Id, nil).Times(1)
			err = apiClient.ProviderDelete(mockProvider.Id)
		})

		It("Should send DELETE request with provider id", func() {})

		It("Should not return an error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Describe("Update Provider", func() {
		var updatedProvider *Provider

		updatedMockProvider := mockProvider
		updatedMockProvider.Description = "new-description"

		BeforeEach(func() {
			updateProviderPayload := ProviderUpdatePayload{Description: updatedMockProvider.Description}
			httpCall = mockHttpClient.EXPECT().
				Put("/providers/"+mockProvider.Id, updateProviderPayload, gomock.Any()).
				Do(func(path string, request any, response *Provider) {
					*response = updatedMockProvider
				}).Times(1)

			updatedProvider, _ = apiClient.ProviderUpdate(mockProvider.Id, updateProviderPayload)
		})

		It("Should return updated provider received from API", func() {
			Expect(*updatedProvider).To(Equal(updatedMockProvider))
		})
	})
})
