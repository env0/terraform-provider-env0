package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitConfigurationVariableResource_Create(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
					resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
					resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
				),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{configVar}, nil)
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
	})
}

func TestUnitConfigurationVariableResource_CreateWrongType(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
		Type:  6,
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
					"type":  configVar.Type,
				}),
				ExpectError: regexp.MustCompile(`(Error: 'type' can only receive either 'environment' or 'terraform')`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {

	})
}

func TestUnitConfigurationVariableResource_ReadWrongApiClientError(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				ExpectError: regexp.MustCompile(`(Error: could not get configurationVariable: error)`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{}, errors.New("error"))
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
	})
}

func TestUnitConfigurationVariableResource_ReadWrongNotFound(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				ExpectError: regexp.MustCompile(`(Error: variable .+ not found)`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{}, nil)
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
	})
}

func TestUnitConfigurationVariableResource_CreateApiClientError(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				ExpectError: regexp.MustCompile(`(could not create configurationVariable: error)`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
	})
}

func TestUnitConfigurationVariableResource_Update(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}
	newConfigVar := client.ConfigurationVariable{
		Id:    configVar.Id,
		Name:  configVar.Name,
		Value: "I want to be the config value",
	}

	updateTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
					resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
					resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  newConfigVar.Name,
					"value": newConfigVar.Value,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", newConfigVar.Id),
					resource.TestCheckResourceAttr(accessor, "name", newConfigVar.Name),
					resource.TestCheckResourceAttr(accessor, "value", newConfigVar.Value),
				),
			},
		},
	}

	runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		gomock.InOrder(
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil).Times(2),
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{newConfigVar}, nil),
		)
		mock.EXPECT().ConfigurationVariableUpdate(newConfigVar.Id, newConfigVar.Name, newConfigVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
	})
}

func TestUnitConfigurationVariableResource_UpdateWrongType(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:    "id0",
		Name:  "name0",
		Value: "Variable",
	}
	newConfigVar := client.ConfigurationVariable{
		Id:    configVar.Id,
		Name:  configVar.Name,
		Value: "I want to be the config value",
		Type:  6,
	}

	updateTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  configVar.Name,
					"value": configVar.Value,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
					resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
					resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name":  newConfigVar.Name,
					"value": newConfigVar.Value,
					"type":  newConfigVar.Type,
				}),
				ExpectError: regexp.MustCompile(`'type' can only receive either 'environment' or 'terraform'`),
			},
		},
	}

	runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil).Times(2)
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
	})
}
