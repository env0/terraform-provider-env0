package env0

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentOutputConfigurationVariableResource(t *testing.T) {
	resourceType := "env0_environment_output_configuration_variable"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName

	accessor := resourceAccessor(resourceType, resourceName)

	value := EnvironmentOutputConfigurationVariableValue{
		OutputName:    "output_name",
		EnvironmentId: "output_environment_id",
	}

	valueBytes, _ := json.Marshal(&value)
	valueStr := string(valueBytes)

	updatedValue := EnvironmentOutputConfigurationVariableValue{
		OutputName:    "output_name2",
		EnvironmentId: "output_environment_id2",
	}

	updatedValueBytes, _ := json.Marshal(&updatedValue)
	updatedValueStr := string(updatedValueBytes)

	configurationVariable := client.ConfigurationVariable{
		Id:          "id0",
		Name:        "name0",
		Description: "desc0",
		Value:       valueStr,
		IsReadOnly:  boolPtr(true),
		ScopeId:     "scope_id",
		Scope:       "PROJECT",
		IsSensitive: boolPtr(false),
		Type:        (*client.ConfigurationVariableType)(intPtr(1)),
		IsRequired:  boolPtr(true),
	}

	updatedConfigurationVariable := configurationVariable
	updatedConfigurationVariable.Description = "desc1"
	updatedConfigurationVariable.Value = updatedValueStr

	t.Run("create and update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  configurationVariable.Name,
						"description":           configurationVariable.Description,
						"is_read_only":          strconv.FormatBool(*configurationVariable.IsReadOnly),
						"is_required":           strconv.FormatBool(*configurationVariable.IsRequired),
						"output_environment_id": value.EnvironmentId,
						"output_name":           value.OutputName,
						"scope":                 configurationVariable.Scope,
						"scope_id":              configurationVariable.ScopeId,
						"type":                  "terraform",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationVariable.Description),
						resource.TestCheckResourceAttr(accessor, "output_environment_id", value.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "output_name", value.OutputName),
						resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
						resource.TestCheckResourceAttr(accessor, "scope_id", configurationVariable.ScopeId),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configurationVariable.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configurationVariable.IsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  configurationVariable.Name,
						"description":           updatedConfigurationVariable.Description,
						"is_read_only":          strconv.FormatBool(*configurationVariable.IsReadOnly),
						"is_required":           strconv.FormatBool(*configurationVariable.IsRequired),
						"output_environment_id": updatedValue.EnvironmentId,
						"output_name":           updatedValue.OutputName,
						"scope":                 configurationVariable.Scope,
						"scope_id":              configurationVariable.ScopeId,
						"type":                  "terraform",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedConfigurationVariable.Description),
						resource.TestCheckResourceAttr(accessor, "output_environment_id", updatedValue.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "output_name", updatedValue.OutputName),
						resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
						resource.TestCheckResourceAttr(accessor, "scope_id", configurationVariable.ScopeId),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configurationVariable.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configurationVariable.IsRequired)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariableCreate(client.ConfigurationVariableCreateParams{
					Name:        configurationVariable.Name,
					Value:       valueStr,
					Scope:       configurationVariable.Scope,
					ScopeId:     configurationVariable.ScopeId,
					Type:        *configurationVariable.Type,
					Description: configurationVariable.Description,
					Format:      client.ENVIRONMENT_OUTPUT,
					IsReadOnly:  *configurationVariable.IsReadOnly,
					IsRequired:  *configurationVariable.IsRequired,
				}).Times(1).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).Times(2).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{CommonParams: client.ConfigurationVariableCreateParams{
					Name:        configurationVariable.Name,
					Value:       updatedValueStr,
					Scope:       configurationVariable.Scope,
					ScopeId:     configurationVariable.ScopeId,
					Type:        *configurationVariable.Type,
					Description: updatedConfigurationVariable.Description,
					Format:      client.ENVIRONMENT_OUTPUT,
					IsReadOnly:  *configurationVariable.IsReadOnly,
					IsRequired:  *configurationVariable.IsRequired,
				}, Id: configurationVariable.Id}).Times(1).Return(updatedConfigurationVariable, nil),
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).Times(1).Return(updatedConfigurationVariable, nil),
				mock.EXPECT().ConfigurationVariableDelete(configurationVariable.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("create read_only in non project scope", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  configurationVariable.Name,
						"description":           configurationVariable.Description,
						"is_read_only":          strconv.FormatBool(*configurationVariable.IsReadOnly),
						"is_required":           strconv.FormatBool(*configurationVariable.IsRequired),
						"output_environment_id": value.EnvironmentId,
						"output_name":           value.OutputName,
						"scope":                 "ENVIRONMENT",
						"scope_id":              configurationVariable.ScopeId,
						"type":                  "terraform",
					}),
					ExpectError: regexp.MustCompile(`'is_read_only' can only be set to 'true' for the 'PROJECT' scope`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	importById := fmt.Sprintf(`{  "Scope": "%s",  "ScopeId": "%s", "Id": "%s"}`, configurationVariable.Scope, configurationVariable.ScopeId, configurationVariable.Id)
	importByName := fmt.Sprintf(`{  "Scope": "%s",  "ScopeId": "%s", "Name": "%s"}`, configurationVariable.Scope, configurationVariable.ScopeId, configurationVariable.Name)

	t.Run("import by id", func(t *testing.T) {
		createTestCaseForImport := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  configurationVariable.Name,
						"description":           configurationVariable.Description,
						"is_read_only":          strconv.FormatBool(*configurationVariable.IsReadOnly),
						"is_required":           strconv.FormatBool(*configurationVariable.IsRequired),
						"output_environment_id": value.EnvironmentId,
						"output_name":           value.OutputName,
						"scope":                 configurationVariable.Scope,
						"scope_id":              configurationVariable.ScopeId,
						"type":                  "terraform",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     importById,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, createTestCaseForImport, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariableCreate(client.ConfigurationVariableCreateParams{
					Name:        configurationVariable.Name,
					Value:       valueStr,
					Scope:       configurationVariable.Scope,
					ScopeId:     configurationVariable.ScopeId,
					Type:        *configurationVariable.Type,
					Description: configurationVariable.Description,
					Format:      client.ENVIRONMENT_OUTPUT,
					IsReadOnly:  *configurationVariable.IsReadOnly,
					IsRequired:  *configurationVariable.IsRequired,
				}).Times(1).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).Times(3).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariableDelete(configurationVariable.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		createTestCaseForImport := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  configurationVariable.Name,
						"description":           configurationVariable.Description,
						"is_read_only":          strconv.FormatBool(*configurationVariable.IsReadOnly),
						"is_required":           strconv.FormatBool(*configurationVariable.IsRequired),
						"output_environment_id": value.EnvironmentId,
						"output_name":           value.OutputName,
						"scope":                 configurationVariable.Scope,
						"scope_id":              configurationVariable.ScopeId,
						"type":                  "terraform",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     importByName,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, createTestCaseForImport, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariableCreate(client.ConfigurationVariableCreateParams{
					Name:        configurationVariable.Name,
					Value:       valueStr,
					Scope:       configurationVariable.Scope,
					ScopeId:     configurationVariable.ScopeId,
					Type:        *configurationVariable.Type,
					Description: configurationVariable.Description,
					Format:      client.ENVIRONMENT_OUTPUT,
					IsReadOnly:  *configurationVariable.IsReadOnly,
					IsRequired:  *configurationVariable.IsRequired,
				}).Times(1).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).Times(1).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariablesByScope(configurationVariable.Scope, configurationVariable.ScopeId).Times(1).Return([]client.ConfigurationVariable{configurationVariable}, nil),
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).Times(1).Return(configurationVariable, nil),
				mock.EXPECT().ConfigurationVariableDelete(configurationVariable.Id).Times(1).Return(nil),
			)
		})
	})
}
