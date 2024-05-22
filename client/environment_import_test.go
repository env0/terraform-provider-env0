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
		Variables: []Variable{{
			Name:        "name",
			Value:       "value",
			IsSensitive: false,
			Type:        "string",
		},
		},
	}

	Describe("EnvironmentImportGet", func() {
		var result *EnvironmentImport

		BeforeEach(func() {
			mockHttpClient.EXPECT().
				Get("/staging-environments/"+mockEnvironmentImport.Id, nil, gomock.Any()).
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
				Variables: []Variable{{
					Name:        "name",
					Value:       "value",
					IsSensitive: false,
					Type:        "string",
				},
				}}

			expectedPayload := struct {
				OrganizationId string `json:"organizationId"`
				EnvironmentImportCreatePayload
			}{
				organizationId,
				payload,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/staging-environments", expectedPayload, gomock.Any()).
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
				Put("/staging-environments/"+mockEnvironmentImport.Id, expectedPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *EnvironmentImport) {
					*response = mockedResponse
				}).Times(1)

			result, _ = apiClient.EnvironmentImportUpdate(mockEnvironmentImport.Id, &payload)
		})

		It("Should return environment import received from API", func() {
			Expect(*result).To(Equal(mockedResponse))
		})
	})

})
