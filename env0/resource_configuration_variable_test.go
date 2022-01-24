package env0

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitConfigurationVariableResource(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	isReadonly := true
	isRequired := false
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:          "id0",
		Name:        "name0",
		Description: "desc0",
		Value:       "Variable",
		IsReadonly:  &isReadonly,
		IsRequired:  &isRequired,
	}
	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"name":         configVar.Name,
		"description":  configVar.Description,
		"value":        configVar.Value,
		"is_read_only": strconv.FormatBool(*configVar.IsReadonly),
		"is_required":  strconv.FormatBool(*configVar.IsRequired),
	})

	configurationVariableCreateParams := client.ConfigurationVariableCreateParams{
		Name:        configVar.Name,
		Value:       configVar.Value,
		IsSensitive: false,
		Scope:       client.ScopeGlobal,
		ScopeId:     "",
		Type:        client.ConfigurationVariableTypeEnvironment,
		EnumValues:  nil,
		Description: configVar.Description,
		Format:      client.Text,
		IsRequired:  *configVar.IsRequired,
		IsReadonly:  *configVar.IsReadonly,
	}
	t.Run("Create", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadonly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create Enum", func(t *testing.T) {
		schema := client.ConfigurationVariableSchema{
			Type: "string",
			Enum: []string{"Variable", "a"},
		}
		configVar := client.ConfigurationVariable{
			Id:          "id0",
			Name:        "name0",
			Description: "desc0",
			Value:       "Variable",
			Schema:      &schema,
		}
		stepConfig := fmt.Sprintf(`
	resource "%s" "test" {
		name = "%s"
		description = "%s"
		value= "%s"
		enum = ["%s","%s"]
	}`, resourceType, configVar.Name, configVar.Description, configVar.Value, configVar.Schema.Enum[0], configVar.Schema.Enum[1])

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "enum.0", configVar.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "enum.1", configVar.Schema.Enum[1]),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(
				client.ConfigurationVariableCreateParams{
					Name:        configVar.Name,
					Value:       configVar.Value,
					IsSensitive: false,
					Scope:       client.ScopeGlobal,
					ScopeId:     "",
					Type:        client.ConfigurationVariableTypeEnvironment,
					EnumValues:  configVar.Schema.Enum,
					Description: configVar.Description,
					Format:      client.Text,
				}).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	for _, format := range []client.Format{client.HCL, client.JSON} {
		t.Run("Create "+string(format)+" Variable", func(t *testing.T) {

			expectedVariable := `{
A = "A"
B = "B"
C = "C"
}
`

			schema := client.ConfigurationVariableSchema{
				Type:   "string",
				Format: format,
			}
			configVar := client.ConfigurationVariable{
				Id:          "id0",
				Name:        "name0",
				Description: "desc0",
				Value:       expectedVariable,
				Schema:      &schema,
			}
			terraformDirective := `<<EOT
{
%{ for key, value in var.map ~}
${key} = "${value}"
%{ endfor ~}
}
EOT`
			stepConfig := fmt.Sprintf(`
variable "map" {
		description = "a mapped variable"
		type        = map(string)
		default = %s
	}


resource "%s" "test" {
		name = "%s"
		description = "%s"
		value = %s
		format = "%s"
}`, expectedVariable, resourceType, configVar.Name, configVar.Description, terraformDirective, string(format))

			createTestCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: stepConfig,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
							resource.TestCheckResourceAttr(accessor, "format", string(format)),
						),
					},
				},
			}

			runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariableCreate(
					client.ConfigurationVariableCreateParams{
						Name:        configVar.Name,
						Value:       expectedVariable,
						IsSensitive: false,
						Scope:       client.ScopeGlobal,
						ScopeId:     "",
						Type:        client.ConfigurationVariableTypeEnvironment,
						EnumValues:  configVar.Schema.Enum,
						Description: configVar.Description,
						Format:      format,
					}).Times(1).Return(configVar, nil)
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
				mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
			})
		})
	}

	t.Run("Create Enum with wrong value", func(t *testing.T) {
		stepConfig := fmt.Sprintf(`
	resource "%s" "test" {
		name = "%s"
		description = "%s"
		value= "%s"
		enum = ["a","b"]
	}`, resourceType, configVar.Name, configVar.Description, configVar.Value)
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(fmt.Sprintf("value - '%s' is not one of the enum options", configVar.Value)),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create with wrong type", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  configVar.Name,
						"value": configVar.Value,
						"type":  6,
					}),
					ExpectError: regexp.MustCompile(`(Error: 'type' can only receive either 'environment' or 'terraform')`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {

		})
	})

	t.Run("Read with wrong api error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not get configurationVariable: error"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(
				configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create api client error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(could not create configurationVariable: error)`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
		})
	})

	t.Run("Update", func(t *testing.T) {
		newIsReadonly := false
		newIsRequired := true
		newConfigVar := client.ConfigurationVariable{
			Id:          configVar.Id,
			Name:        configVar.Name,
			Description: configVar.Description,
			Value:       "I want to be the config value",
			IsReadonly:  &newIsReadonly,
			IsRequired:  &newIsRequired,
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(*configVar.IsReadonly),
						"is_required":  strconv.FormatBool(*configVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadonly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         newConfigVar.Name,
						"description":  newConfigVar.Description,
						"value":        newConfigVar.Value,
						"format":       client.HCL,
						"is_read_only": strconv.FormatBool(*newConfigVar.IsReadonly),
						"is_required":  strconv.FormatBool(*newConfigVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", newConfigVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", newConfigVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", newConfigVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", newConfigVar.Value),
						resource.TestCheckResourceAttr(accessor, "format", string(client.HCL)),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*newConfigVar.IsReadonly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*newConfigVar.IsRequired)),
					),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			createParams := configurationVariableCreateParams
			updateParams := createParams
			updateParams.Name = newConfigVar.Name
			updateParams.Value = newConfigVar.Value
			updateParams.Description = newConfigVar.Description
			updateParams.Format = client.HCL
			updateParams.IsReadonly = *newConfigVar.IsReadonly
			updateParams.IsRequired = *newConfigVar.IsRequired

			mock.EXPECT().ConfigurationVariableCreate(createParams).Times(1).Return(configVar, nil)
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(2).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Return(newConfigVar, nil),
			)
			mock.EXPECT().ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{CommonParams: updateParams, Id: newConfigVar.Id}).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Update with wrong type", func(t *testing.T) {
		wrongType := client.ConfigurationVariableType(6)
		newConfigVar := client.ConfigurationVariable{
			Id:    configVar.Id,
			Name:  configVar.Name,
			Value: "I want to be the config value",
			Type:  &wrongType,
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(*configVar.IsReadonly),
						"is_required":  strconv.FormatBool(*configVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadonly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        newConfigVar.Name,
						"description": newConfigVar.Description,
						"value":       newConfigVar.Value,
						"type":        newConfigVar.Type,
					}),
					ExpectError: regexp.MustCompile(`'type' can only receive either 'environment' or 'terraform'`),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(2).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})
}
