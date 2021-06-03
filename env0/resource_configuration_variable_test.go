package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitConfigurationVariableResourceCreate(t *testing.T) {
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
				Config: resourceConfigCreate(resourceType, resourceName, map[string]string{
					"name":  configVar.Name,
					"value": configVar.Value,
				}, make(map[string]int), make(map[string]bool)),
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
func TestUnitConfigurationVariableResourceUpdate(t *testing.T) {
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
				Config: resourceConfigCreate(resourceType, resourceName, map[string]string{
					"name":  configVar.Name,
					"value": configVar.Value,
				}, make(map[string]int), make(map[string]bool)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
					resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
					resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]string{
					"name":  newConfigVar.Name,
					"value": newConfigVar.Value,
				}, make(map[string]int), make(map[string]bool)),
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
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil),
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{configVar}, nil),
			mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Return([]client.ConfigurationVariable{newConfigVar}, nil),
		)
		mock.EXPECT().ConfigurationVariableUpdate(newConfigVar.Id, newConfigVar.Name, newConfigVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)

	})
}
