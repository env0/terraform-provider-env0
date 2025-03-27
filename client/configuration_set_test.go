package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Configuration Set", func() {
	id := "id12345"
	projectId := "projectId123"

	mockConfigurationSet := ConfigurationSet{
		Id:          id,
		Name:        "name",
		Description: "description",
	}

	var configurationSet *ConfigurationSet

	Describe("create organization configuration set", func() {
		BeforeEach(func() {
			mockOrganizationIdCall().Times(1)

			createPayload := CreateConfigurationSetPayload{
				Name:        "name1",
				Description: "des1",
				Scope:       "organization",
			}

			createPayloadWithScopeId := CreateConfigurationSetPayload{
				Name:        "name1",
				Description: "des1",
				Scope:       "organization",
				ScopeId:     organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/configuration-sets", &createPayloadWithScopeId, gomock.Any()).
				Do(func(path string, request any, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSetCreate(&createPayload)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("create project configuration set", func() {
		BeforeEach(func() {
			createPayload := CreateConfigurationSetPayload{
				Name:        "name1",
				Description: "des1",
				Scope:       "project",
				ScopeId:     projectId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/configuration-sets", &createPayload, gomock.Any()).
				Do(func(path string, request any, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSetCreate(&createPayload)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("update configuration set", func() {
		BeforeEach(func() {
			updatePayload := UpdateConfigurationSetPayload{
				Name:        "name2",
				Description: "des2",
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/configuration-sets/"+id, &updatePayload, gomock.Any()).
				Do(func(path string, request any, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSetUpdate(id, &updatePayload)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("get configuration set by id", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets/"+id, nil, gomock.Any()).
				Do(func(path string, request any, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSet(id)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("delete configuration set", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Delete("/configuration-sets/"+id, nil).
				Do(func(path string, request any) {}).
				Times(1)

			err = apiClient.ConfigurationSetDelete(id)
		})

		It("Should call delete once", func() {})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("get configuration variables by set id", func() {
		mockVariables := []ConfigurationVariable{
			{
				ScopeId: "a",
				Value:   "b",
				Scope:   "c",
				Id:      "d",
			},
		}

		var variables []ConfigurationVariable

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration", map[string]string{
					"setId":          id,
					"organizationId": organizationId,
				}, gomock.Any()).
				Do(func(path string, request any, response *[]ConfigurationVariable) {
					*response = mockVariables
				}).Times(1)

			variables, _ = apiClient.ConfigurationVariablesBySetId(id)
		})

		It("Should return configuration variables", func() {
			Expect(variables).To(Equal(mockVariables))
		})
	})

	Describe("get configuration variables by set project id", func() {
		mockVariables := []ConfigurationSet{
			{
				Id:              "id",
				Name:            "name",
				CreationScopeId: "create_scope_id",
			},
		}

		var variables []ConfigurationSet

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets", map[string]string{
					"scopeId": mockVariables[0].CreationScopeId,
					"scope":   "project",
				}, gomock.Any()).
				Do(func(path string, request any, response *[]ConfigurationSet) {
					*response = mockVariables
				}).Times(1)

			variables, _ = apiClient.ConfigurationSets("PROJECT", mockVariables[0].CreationScopeId)
		})

		It("Should return configuration sets", func() {
			Expect(variables).To(Equal(mockVariables))
		})
	})

	Describe("get configuration variables by set organization id", func() {
		mockVariables := []ConfigurationSet{
			{
				Id:              "id",
				Name:            "name",
				CreationScopeId: "create_scope_id",
			},
		}

		var variables []ConfigurationSet

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets", map[string]string{
					"scopeId": mockVariables[0].CreationScopeId,
					"scope":   "organization",
				}, gomock.Any()).
				Do(func(path string, request any, response *[]ConfigurationSet) {
					*response = mockVariables
				}).Times(1)

			variables, _ = apiClient.ConfigurationSets("ORGANIZATION", mockVariables[0].CreationScopeId)
		})

		It("Should return configuration sets", func() {
			Expect(variables).To(Equal(mockVariables))
		})
	})
})
