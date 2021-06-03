package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitConfigurationVariableResource(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	resourceFullName := resourceAccessor(resourceType, resourceName)
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
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", configVar.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", configVar.Name),
					resource.TestCheckResourceAttr(resourceFullName, "value", configVar.Value),
				),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ConfigurationVariableCreate(configVar.Name, configVar.Value, false, client.ScopeGlobal, "", client.ConfigurationVariableTypeEnvironment,
			nil).Times(1).Return(configVar, nil)
		mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").Times(1).Return([]client.ConfigurationVariable{configVar}, nil)
	})
}
