package client_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Configuration Variable", func() {
	isSensitive := true
	varType := ConfigurationVariableTypeEnvironment
	schema := ConfigurationVariableSchema{
		Type:   "string",
		Format: HCL,
	}
	isReadOnly := true
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
		IsReadOnly:     &isReadOnly,
		IsRequired:     &isRequired,
		Regex:          "regex",
	}

	mockTemplateConfigurationVariable := ConfigurationVariable{
		Id:             "config-var-id-1111",
		Name:           "ignore",
		Description:    "ignore",
		Value:          "ignore",
		OrganizationId: organizationId,
		Scope:          ScopeTemplate,
		ScopeId:        "scope-id",
	}

	Describe("ConfigurationVariable", func() {
		Describe("Schema", func() {
			It("On schema type is free text, enum should be nil", func() {
				var parsedPayload ConfigurationVariable
				_ = json.Unmarshal([]byte(`{"schema": {"type": "string"}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeNil())
			})

			It("On schema type is dropdown, enum should be present", func() {
				var parsedPayload ConfigurationVariable
				_ = json.Unmarshal([]byte(`{"schema": {"type": "string", "enum": ["hello"]}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeEquivalentTo([]string{"hello"}))
			})
		})

		Describe("Enums", func() {
			It("Should convert enums correctly", func() {
				var parsedPayload ConfigurationVariable
				_ = json.Unmarshal([]byte(`{"scope":"PROJECT", "type": 1}`), &parsedPayload)
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
				"isReadonly":     *mockConfig.IsReadOnly,
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
					IsReadOnly:  *mockConfig.IsReadOnly,
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
			httpCall = mockHttpClient.EXPECT().Delete("configuration/"+mockConfigurationVariable.Id, nil)
			_ = apiClient.ConfigurationVariableDelete(mockConfigurationVariable.Id)
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
				"isReadonly": *mockConfigurationVariable.IsReadOnly,
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
						IsReadOnly:  *mockConfigurationVariable.IsReadOnly,
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
		scopeId := mockTemplateConfigurationVariable.ScopeId

		var returnedVariables []ConfigurationVariable
		mockVariables := []ConfigurationVariable{mockTemplateConfigurationVariable}
		expectedParams := map[string]string{"organizationId": organizationId, "blueprintId": scopeId}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/configuration", expectedParams, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = mockVariables
				})
			returnedVariables, _ = apiClient.ConfigurationVariablesByScope(ScopeTemplate, scopeId)
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
	})
})

func TestConfigurationVariableMarshelling(t *testing.T) {
	str := "this is a string"

	variable := ConfigurationVariable{
		Value: "a",
		Overwrites: &ConfigurationVariableOverwrites{
			Value: str,
		},
	}

	b, err := json.Marshal(&variable)
	if assert.NoError(t, err) {
		assert.False(t, strings.Contains(string(b), str))
	}

	type ConfigurationVariableDummy ConfigurationVariable

	dummy := ConfigurationVariableDummy(variable)

	b, err = json.Marshal(&dummy)
	if assert.NoError(t, err) {
		assert.True(t, strings.Contains(string(b), str))
	}

	var variable2 ConfigurationVariable

	err = json.Unmarshal(b, &variable2)

	if assert.NoError(t, err) && assert.NotNil(t, variable2.Overwrites) {
		assert.Equal(t, str, variable2.Overwrites.Value)
		assert.Equal(t, variable, variable2)
	}
}
