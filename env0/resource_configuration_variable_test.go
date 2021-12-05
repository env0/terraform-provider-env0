package env0

import (
	"errors"
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitConfigurationVariableResource(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:          "id0",
		Name:        "name0",
		Description: "desc0",
		Value:       "Variable",
	}
	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"name":        configVar.Name,
		"description": configVar.Description,
		"value":       configVar.Value,
	})

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
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{configVar}, nil)
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
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				configVar.Schema.Enum, configVar.Description).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{configVar}, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})
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
					ExpectError: regexp.MustCompile(`(Error: could not get configurationVariable: error)`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{}, errors.New("error"))
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Read not found", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(Error: variable .+ not found)`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{}, nil)
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
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
		})
	})

	t.Run("Update", func(t *testing.T) {
		newConfigVar := client.ConfigurationVariable{
			Id:          configVar.Id,
			Name:        configVar.Name,
			Description: configVar.Description,
			Value:       "I want to be the config value",
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        configVar.Name,
						"description": configVar.Description,
						"value":       configVar.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        newConfigVar.Name,
						"description": newConfigVar.Description,
						"value":       newConfigVar.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", newConfigVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", newConfigVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", newConfigVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", newConfigVar.Value),
					),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(configVar, nil)
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil).Times(2),
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{newConfigVar}, nil),
			)
			mock.EXPECT().ConfigurationVariableUpdate(newConfigVar.Id, newConfigVar.Name, newConfigVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, newConfigVar.Description).Times(1).Return(configVar, nil)
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
						"name":        configVar.Name,
						"description": configVar.Description,
						"value":       configVar.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
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
			mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
				nil, configVar.Description).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil).Times(2)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})
}
