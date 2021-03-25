package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/env0/terraform-provider-env0/api"
)

var _ = Describe("Organization wide configuration variable", func() {
	var configurationVariable ConfigurationVariable
	var configurationVariableErr error

	JustBeforeEach(func() {
		configurationVariable = ConfigurationVariable{}
		configurationVariable, configurationVariableErr = apiClient.ConfigurationVariableCreate(
			"testing_org_wide_var",
			"fake value",
			false,
			ScopeGlobal,
			"",
			ConfigurationVariableTypeTerraform,
			nil)
	})

	Specify("Configuration variable creation should succeed", func() {
		Expect(configurationVariableErr).To(BeNil())
		Expect(configurationVariable.Name).To(Equal("testing_org_wide_var"))
	})

	AfterEach(func() {
		if configurationVariable.Id == "" {
			return
		}
		err := apiClient.ConfigurationVariableDelete(configurationVariable.Id)
		Expect(err).To(BeNil())
	})

	When("Fetching all the global configuration variables", func() {
		var configurationVariables []ConfigurationVariable
		var configurationVariablesErr error

		JustBeforeEach(func() {
			configurationVariables, configurationVariablesErr = apiClient.ConfigurationVariables(ScopeGlobal, "")
		})

		Specify("Fetch should succeed", func() {
			Expect(configurationVariablesErr).To(BeNil())
			Expect(len(configurationVariables)).ToNot(BeZero())
		})

		When("Seaching for created configuration variable in fetched variable list", func() {
			var found ConfigurationVariable

			JustBeforeEach(func() {
				found = ConfigurationVariable{}
				for _, candidate := range configurationVariables {
					if candidate.Id == configurationVariable.Id {
						found = candidate
					}
				}
			})

			It("Should have found the created variable", func() {
				Expect(found.Name).To(Equal("testing_org_wide_var"))
			})

			It("Shoud not be sensitive", func() {
				Expect(found.IsSensitive).To(BeFalse())
			})
		})
	})
})
