package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Gpg Token Client", func() {
	mockGpgKey := GpgKey{
		Id:      "id",
		Name:    "name",
		KeyId:   "keyId",
		Content: "content",
	}

	Describe("Get All Gpg Keys", func() {
		var returnedGpgKeys []GpgKey
		mockGpgKeys := []GpgKey{mockGpgKey}

		BeforeEach(func() {
			mockOrganizationIdCall()
			mockHttpClient.EXPECT().
				Get("/gpg-keys", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]GpgKey) {
					*response = mockGpgKeys
				})
			returnedGpgKeys, _ = apiClient.GpgKeys()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should return GpgKeys", func() {
			Expect(returnedGpgKeys).To(Equal(mockGpgKeys))
		})
	})

	Describe("Delete Gpg Key", func() {
		var err error

		BeforeEach(func() {
			mockHttpClient.EXPECT().Delete("/gpg-keys/"+mockGpgKey.Id, nil)
			err = apiClient.GpgKeyDelete(mockGpgKey.Id)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Create GpgKey", func() {
		var createdGpgKey *GpgKey
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall()

			payload := struct {
				OrganizationId string `json:"organizationId"`
				GpgKeyCreatePayload
			}{
				organizationId,
				GpgKeyCreatePayload{
					Name:    mockGpgKey.Name,
					KeyId:   mockGpgKey.KeyId,
					Content: mockGpgKey.Content,
				},
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/gpg-keys", payload, gomock.Any()).
				Do(func(path string, request any, response *GpgKey) {
					*response = mockGpgKey
				})

			createdGpgKey, err = apiClient.GpgKeyCreate(&GpgKeyCreatePayload{
				Name:    mockGpgKey.Name,
				KeyId:   mockGpgKey.KeyId,
				Content: mockGpgKey.Content,
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created GpgToken", func() {
			Expect(*createdGpgKey).To(Equal(mockGpgKey))
		})
	})
})
