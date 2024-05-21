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
		Id: id,
	}

	mockConfigurationSetList := []ConfigurationSet{mockConfigurationSet}

	var configurationSet *ConfigurationSet
	var configurationSetList []ConfigurationSet

	Describe("create organization configuration set", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)

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
				Do(func(path string, request interface{}, response *ConfigurationSet) {
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
				Do(func(path string, request interface{}, response *ConfigurationSet) {
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
				Do(func(path string, request interface{}, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSetUpdate(id, &updatePayload)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("list organization scope configuration sets", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets", map[string]string{
					"scope":               "organization",
					"scopeId":             organizationId,
					"includeHigherScopes": "false",
				}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationSet) {
					*response = mockConfigurationSetList
				}).Times(1)

			configurationSetList, _ = apiClient.ConfigurationSets("organization", "")
		})

		It("Should return configuration set list", func() {
			Expect(configurationSetList).To(Equal(mockConfigurationSetList))
		})
	})

	Describe("list project scope configuration sets", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets", map[string]string{
					"scope":               "project",
					"scopeId":             projectId,
					"includeHigherScopes": "false",
				}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationSet) {
					*response = mockConfigurationSetList
				}).Times(1)

			configurationSetList, _ = apiClient.ConfigurationSets("project", projectId)
		})

		It("Should return configuration set list", func() {
			Expect(configurationSetList).To(Equal(mockConfigurationSetList))
		})
	})

	Describe("get configuration set by id", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets/"+id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *ConfigurationSet) {
					*response = mockConfigurationSet
				}).Times(1)

			configurationSet, _ = apiClient.ConfigurationSet(id)
		})

		It("Should return configuration set", func() {
			Expect(*configurationSet).To(Equal(mockConfigurationSet))
		})
	})

	Describe("delete configuration set", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Delete("/configuration-sets/"+id, nil).
				Do(func(path string, request interface{}) {}).
				Times(1)

			apiClient.ConfigurationSetDelete(id)
		})

		It("Should call delete once", func() {})
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
			mockOrganizationIdCall(organizationId).Times(1)

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration", map[string]string{
					"setId":          id,
					"organizationId": organizationId,
				}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = mockVariables
				}).Times(1)

			variables, _ = apiClient.ConfigurationVariablesBySetId(id)
		})

		It("Should return configuration variables", func() {
			Expect(variables).To(Equal(mockVariables))
		})
	})
})
