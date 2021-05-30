package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Variable", func() {
	mockConfigurationVariable := ConfigurationVariable{
		Id:             "config-var-id-789",
		Name:           "configName",
		Value:          "configValue",
		OrganizationId: organizationId,
		IsSensitive:    true,
		Scope:          "PROJECT",
		Type:           0,
		ScopeId:        "project-123",
		UserId:         "user|123",
	}

	Describe("ConfigurationVariableCreate", func() {
		var createdConfigurationVariable ConfigurationVariable

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			expectedCreateRequest := []map[string]interface{}{{
				"name":           mockConfigurationVariable.Name,
				"isSensitive":    mockConfigurationVariable.IsSensitive,
				"value":          mockConfigurationVariable.Value,
				"organizationId": organizationId,
				"scopeId":        mockConfigurationVariable.ScopeId,
				"scope":          Scope(mockConfigurationVariable.Scope),
				"type":           ConfigurationVariableType(mockConfigurationVariable.Type),
			}}

			httpCall = mockHttpClient.EXPECT().
				Post("configuration", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = []ConfigurationVariable{mockConfigurationVariable}
				})

			createdConfigurationVariable, _ = apiClient.ConfigurationVariableCreate(
				mockConfigurationVariable.Name,
				mockConfigurationVariable.Value,
				mockConfigurationVariable.IsSensitive,
				Scope(mockConfigurationVariable.Scope),
				mockConfigurationVariable.ScopeId,
				ConfigurationVariableType(mockConfigurationVariable.Type),
				nil,
			)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return created configuration variable", func() {
			Expect(createdConfigurationVariable).To(Equal(mockConfigurationVariable))
		})
	})

	Describe("ConfigurationVariableDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("configuration/" + mockConfigurationVariable.Id)
			apiClient.ConfigurationVariableDelete(mockConfigurationVariable.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})
	})

	Describe("ConfigurationVariableUpdate", func() {
		var updatedConfigurationVariable ConfigurationVariable

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			newName := "new-" + mockConfigurationVariable.Name
			newValue := "new-" + mockConfigurationVariable.Value

			expectedUpdateRequest := []map[string]interface{}{{
				"name":           newName,
				"value":          newValue,
				"id":             mockConfigurationVariable.Id,
				"isSensitive":    mockConfigurationVariable.IsSensitive,
				"organizationId": organizationId,
				"scopeId":        mockConfigurationVariable.ScopeId,
				"scope":          Scope(mockConfigurationVariable.Scope),
				"type":           ConfigurationVariableType(mockConfigurationVariable.Type),
			}}

			httpCall = mockHttpClient.EXPECT().
				Post("/configuration", expectedUpdateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = []ConfigurationVariable{mockConfigurationVariable}
				})

			updatedConfigurationVariable, _ = apiClient.ConfigurationVariableUpdate(
				mockConfigurationVariable.Id,
				newName,
				newValue,
				mockConfigurationVariable.IsSensitive,
				Scope(mockConfigurationVariable.Scope),
				mockConfigurationVariable.ScopeId,
				ConfigurationVariableType(mockConfigurationVariable.Type),
				nil,
			)
		})

		It("Should send POST request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should return configuration value received from API", func() {
			Expect(updatedConfigurationVariable).To(Equal(mockConfigurationVariable))
		})
	})
})
