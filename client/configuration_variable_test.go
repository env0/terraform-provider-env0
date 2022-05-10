package client_test

import (
	"encoding/json"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Variable", func() {
	isSensitive := true
	varType := ConfigurationVariableTypeEnvironment
	schema := ConfigurationVariableSchema{
		Type:   "string",
		Format: HCL,
	}
	isReadonly := true
	isRequired := true

	mockConfigurationVariable := ConfigurationVariable{
		Id:             "config-var-id-789",
		Name:           "configName",
		Description:    "configDescription",
		Value:          "configValue",
		OrganizationId: organizationId,
		IsSensitive:    &isSensitive,
		Scope:          ScopeProject,
		Type:           &varType,
		ScopeId:        "project-123",
		UserId:         "user|123",
		Schema:         &schema,
		IsReadonly:     &isReadonly,
		IsRequired:     &isRequired,
		Regex:          "regex",
	}

	mockGlobalConfigurationVariable := ConfigurationVariable{
		Id:             "config-var-id-789",
		Name:           "configName",
		Description:    "configDescription",
		Value:          "configValue",
		OrganizationId: organizationId,
		IsSensitive:    &isSensitive,
		Scope:          ScopeGlobal,
		Type:           &varType,
		ScopeId:        "project-123",
		UserId:         "user|123",
		Schema:         &schema,
		IsReadonly:     &isReadonly,
		IsRequired:     &isRequired,
		Regex:          "regex",
	}

	Describe("ConfigurationVariable", func() {
		Describe("Schema", func() {
			It("On schema type is free text, enum should be nil", func() {
				var parsedPayload ConfigurationVariable
				json.Unmarshal([]byte(`{"schema": {"type": "string"}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeNil())
			})

			It("On schema type is dropdown, enum should be present", func() {
				var parsedPayload ConfigurationVariable
				json.Unmarshal([]byte(`{"schema": {"type": "string", "enum": ["hello"]}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeEquivalentTo([]string{"hello"}))
			})
		})

		Describe("Enums", func() {
			It("Should convert enums correctly", func() {
				var parsedPayload ConfigurationVariable
				json.Unmarshal([]byte(`{"scope":"PROJECT", "type": 1}`), &parsedPayload)
				Expect(parsedPayload.Scope).Should(Equal(ScopeProject))
				Expect(*parsedPayload.Type).Should(Equal(ConfigurationVariableTypeTerraform))
			})
		})
	})

	Describe("ConfigurationVariablesById", func() {
		id := "configurationId"
		var found ConfigurationVariable
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration/"+id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *ConfigurationVariable) {
					*response = mockConfigurationVariable
				})

			found, _ = apiClient.ConfigurationVariablesById(id)
		})

		It("Should return variable", func() {
			Expect(found).To(Equal(mockConfigurationVariable))
		})
	})

	Describe("ConfigurationVariableCreate", func() {
		var createdConfigurationVariable ConfigurationVariable

		var GetExpectedRequest = func(mockConfig ConfigurationVariable) []map[string]interface{} {
			schema := map[string]interface{}{
				"type":   mockConfig.Schema.Type,
				"format": mockConfig.Schema.Format,
			}

			if mockConfig.Schema.Format == Text {
				delete(schema, "format")
			}

			request := []map[string]interface{}{{
				"name":           mockConfig.Name,
				"description":    mockConfig.Description,
				"isSensitive":    *mockConfig.IsSensitive,
				"value":          mockConfig.Value,
				"organizationId": organizationId,
				"scopeId":        mockConfig.ScopeId,
				"scope":          mockConfig.Scope,
				"type":           *mockConfig.Type,
				"schema":         schema,
				"isReadonly":     *mockConfig.IsReadonly,
				"isRequired":     *mockConfig.IsRequired,
				"regex":          mockConfig.Regex,
			}}
			return request
		}

		var SetCreateRequestExpectation = func(mockConfig ConfigurationVariable) {
			expectedCreateRequest := GetExpectedRequest(mockConfig)
			httpCall = mockHttpClient.EXPECT().
				Post("configuration", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = []ConfigurationVariable{mockConfig}
				})
		}

		var DoCreateRequest = func(mockConfig ConfigurationVariable) {
			createdConfigurationVariable, _ = apiClient.ConfigurationVariableCreate(
				ConfigurationVariableCreateParams{
					Name:        mockConfig.Name,
					Value:       mockConfig.Value,
					Description: mockConfig.Description,
					IsSensitive: *mockConfig.IsSensitive,
					Scope:       mockConfig.Scope,
					ScopeId:     mockConfig.ScopeId,
					Type:        *mockConfig.Type,
					EnumValues:  nil,
					Format:      mockConfig.Schema.Format,
					IsReadonly:  *mockConfig.IsReadonly,
					IsRequired:  *mockConfig.IsRequired,
					Regex:       mockConfig.Regex,
				},
			)
		}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			SetCreateRequestExpectation(mockConfigurationVariable)
			DoCreateRequest(mockConfigurationVariable)
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

		DescribeTable("Create with different schema format", func(schemaFormat Format) {
			var mockWithFormat = ConfigurationVariable{}
			copier.Copy(&mockWithFormat, &mockConfigurationVariable)
			mockWithFormat.Schema.Format = schemaFormat
			SetCreateRequestExpectation(mockWithFormat)

			DoCreateRequest(mockWithFormat)

			httpCall.Times(1)
			Expect(createdConfigurationVariable).To(Equal(mockWithFormat))
		},
			Entry("Text Format", Text),
			Entry("JSON Format", JSON),
			Entry("HCL Format", HCL),
		)
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
			newDescription := "new-" + mockConfigurationVariable.Description
			newValue := "new-" + mockConfigurationVariable.Value

			expectedUpdateRequest := []map[string]interface{}{{
				"name":           newName,
				"description":    newDescription,
				"value":          newValue,
				"id":             mockConfigurationVariable.Id,
				"isSensitive":    *mockConfigurationVariable.IsSensitive,
				"organizationId": organizationId,
				"scopeId":        mockConfigurationVariable.ScopeId,
				"scope":          mockConfigurationVariable.Scope,
				"type":           *mockConfigurationVariable.Type,
				"schema": map[string]interface{}{
					"type": mockConfigurationVariable.Schema.Type,
				},
				"isReadonly": *mockConfigurationVariable.IsReadonly,
				"isRequired": *mockConfigurationVariable.IsRequired,
				"regex":      mockConfigurationVariable.Regex,
			}}

			httpCall = mockHttpClient.EXPECT().
				Post("/configuration", expectedUpdateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = []ConfigurationVariable{mockConfigurationVariable}
				})

			updatedConfigurationVariable, _ = apiClient.ConfigurationVariableUpdate(
				ConfigurationVariableUpdateParams{
					Id: mockConfigurationVariable.Id,
					CommonParams: ConfigurationVariableCreateParams{
						Name:        newName,
						Value:       newValue,
						Description: newDescription,
						IsSensitive: *mockConfigurationVariable.IsSensitive,
						Scope:       mockConfigurationVariable.Scope,
						ScopeId:     mockConfigurationVariable.ScopeId,
						Type:        *mockConfigurationVariable.Type,
						EnumValues:  nil,
						IsReadonly:  *mockConfigurationVariable.IsReadonly,
						IsRequired:  *mockConfigurationVariable.IsRequired,
						Regex:       mockConfigurationVariable.Regex,
					},
				},
			)
		})

		It("Should send POST request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should return configuration value received from API", func() {
			Expect(updatedConfigurationVariable).To(Equal(mockConfigurationVariable))
		})
	})

	Describe("ConfigurationVariablesByScope", func() {
		var returnedVariables []ConfigurationVariable
		mockVariables := []ConfigurationVariable{mockConfigurationVariable, mockGlobalConfigurationVariable}
		expectedParams := map[string]string{"organizationId": organizationId}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration", expectedParams, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = mockVariables
				})
			returnedVariables, _ = apiClient.ConfigurationVariablesByScope(ScopeGlobal, "")
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
				returnedVariables, _ = apiClient.ConfigurationVariablesByScope(Scope(scope), scopeId)
				Expect(returnedVariables).To((Equal([]ConfigurationVariable{mockConfigurationVariable})))
				httpCall.Times(1)
			},
			Entry("Template Scope", string(ScopeTemplate), "blueprintId"),
			Entry("Project Scope", string(ScopeProject), "projectId"),
			Entry("Environment Scope", string(ScopeEnvironment), "environmentId"),
			Entry("Project Scope", string(ScopeDeploymentLog), "deploymentLogId"),
		)
	})
})
