package api_test

import (
	. "github.com/env0/terraform-provider-env0/internal/api"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Variable", func() {
	var httpCall *gomock.Call
	var configVar ConfigurationVariable
	mockConfigurationVariable := &ConfigurationVariable{
		Id:             "id",
		Name:           "config-key",
		Value:          "config-value",
		IsSensitive:    false,
		Scope:          "GLOBAL",
		OrganizationId: organizationId,
	}

	Describe("Create", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post(
				"/configuration",
				map[string]interface{}{
					"name":           "testing_org_wide_var",
					"value":          "fake value",
					"isSensitive":    false,
					"scope":          ScopeGlobal,
					"type":           ConfigurationVariableTypeTerraform,
					"organizationId": organizationId,
				},
				&configVar,
			).Return(mockConfigurationVariable)

			configVar, _ = apiClient.ConfigurationVariableCreate(
				"testing_org_wide_var",
				"fake value",
				false,
				ScopeGlobal,
				"",
				ConfigurationVariableTypeTerraform,
				nil)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return created configuration variable", func() {
			Expect(configVar).To(Equal(mockConfigurationVariable))
		})
	})

	Describe("Delete", func() {
		JustBeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("configuration/" + mockConfigurationVariable.Id).Times(1)

			apiClient.ConfigurationVariableDelete(mockConfigurationVariable.Id)
		})

		It("Should call DELETE request with param", func() {
			httpCall.Times(1)
		})
	})
})
