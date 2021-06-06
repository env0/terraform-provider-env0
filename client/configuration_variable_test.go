package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Variable", func() {
	mockConfigurationVariable := ConfigurationVariable{
		Id:             "config-var-id-789",
		Name:           "configName",
		Value:          "configValue",
		OrganizationId: organizationId,
		IsSensitive:    true,
		Scope:          ScopeProject,
		Type:           ConfigurationVariableTypeEnvironment,
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
				"scope":          mockConfigurationVariable.Scope,
				"type":           mockConfigurationVariable.Type,
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
				mockConfigurationVariable.Scope,
				mockConfigurationVariable.ScopeId,
				mockConfigurationVariable.Type,
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
				"scope":          mockConfigurationVariable.Scope,
				"type":           mockConfigurationVariable.Type,
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
				mockConfigurationVariable.Scope,
				mockConfigurationVariable.ScopeId,
				mockConfigurationVariable.Type,
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

	Describe("ConfigurationVariables", func() {
		var returnedVariables []ConfigurationVariable
		mockVariables := []ConfigurationVariable{mockConfigurationVariable}
		expectedParams := map[string]string{"organizationId": organizationId}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration", expectedParams, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = mockVariables
				})
			returnedVariables, _ = apiClient.ConfigurationVariables(ScopeGlobal, "")
		})

		It("Should send GET request with expected params", func() {
			httpCall.Times(1)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should return variables", func() {
			Expect(returnedVariables).To(Equal(mockVariables))
		})

		DescribeTable("Different Scopes",
			func(scope string, expectedFieldName string) {
				scopeId := expectedFieldName + "-id"
				expectedParams := map[string]string{
					"organizationId":  organizationId,
					expectedFieldName: scopeId,
				}

				httpCall = mockHttpClient.EXPECT().
					Get("/configuration", expectedParams, gomock.Any()).
					Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
						*response = mockVariables
					})
				returnedVariables, _ = apiClient.ConfigurationVariables(Scope(scope), scopeId)
				httpCall.Times(1)
			},
			Entry("Template Scope", string(ScopeTemplate), "blueprintId"),
			Entry("Project Scope", string(ScopeProject), "projectId"),
			Entry("Environment Scope", string(ScopeEnvironment), "environmentId"),
			Entry("Project Scope", string(ScopeDeploymentLog), "deploymentLogId"),
		)
	})
})
