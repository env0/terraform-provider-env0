package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Environment Import Client", func() {
	mockEnvironmentImport := EnvironmentImport{
		Id:         "id",
		Name:       "name",
		IacType:    "tofu",
		IacVersion: "1.0",
		Workspace:  "workspace",
	}

	Describe("EnvironmentImportGet", func() {
		var result *EnvironmentImport

		BeforeEach(func() {
			mockHttpClient.EXPECT().
				Get("/environment-imports/"+mockEnvironmentImport.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *EnvironmentImport) {
					*response = mockEnvironmentImport
				}).Times(1)
			result, _ = apiClient.EnvironmentImportGet(mockEnvironmentImport.Id)
		})

		It("Should return the environment import", func() {
			Expect(*result).To(Equal(mockEnvironmentImport))
		})
	})

	Describe("EnvironmentImportCreate", func() {
		var result *EnvironmentImport
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)

			payload := EnvironmentImportCreatePayload{
				Name:       "name",
				IacType:    "tofu",
				IacVersion: "1.0",
				Workspace:  "workspace",
			}

			expectedPayload := struct {
				OrganizationId string `json:"organizationId"`
				EnvironmentImportCreatePayload
			}{
				organizationId,
				payload,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/environment-imports", expectedPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *EnvironmentImport) {
					*response = mockEnvironmentImport
				}).Times(1)

			result, err = apiClient.EnvironmentImportCreate(&payload)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return the created environment import", func() {
			Expect(*result).To(Equal(mockEnvironmentImport))
		})
	})

	Describe("EnvironmentImportUpdate", func() {
		var result *EnvironmentImport
		var mockedResponse EnvironmentImport
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)
			payload := EnvironmentImportUpdatePayload{
				Name: "new name",
			}

			mockedResponse = mockEnvironmentImport
			mockedResponse.Name = payload.Name

			expectedPayload := struct {
				OrganizationId string `json:"organizationId"`
				EnvironmentImportUpdatePayload
			}{
				organizationId,
				payload,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/environment-imports/"+mockEnvironmentImport.Id, expectedPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *EnvironmentImport) {
					*response = mockedResponse
				}).Times(1)

			result, _ = apiClient.EnvironmentImportUpdate(mockEnvironmentImport.Id, &payload)
		})

		It("Should return environment import received from API", func() {
			Expect(*result).To(Equal(mockedResponse))
		})
	})

	Describe("EnvironmentImportDelete", func() {
		var err error

		BeforeEach(func() {
			mockHttpClient.EXPECT().Delete("/environment-imports/"+mockEnvironmentImport.Id, nil).Times(1)
			err = apiClient.EnvironmentImportDelete(mockEnvironmentImport.Id)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
