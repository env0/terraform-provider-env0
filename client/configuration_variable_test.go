package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Variable", func() {
	mockConfigurationVariable := ConfigurationVariable{
		Id:             "idX",
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
})
