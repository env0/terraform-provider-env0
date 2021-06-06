package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestUnitConfigurationVariableData(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)
	configurationVariable := client.ConfigurationVariable{
		Id:             "id0",
		Name:           "name0",
		ScopeId:        "scope0",
		Value:          "value0",
		OrganizationId: "organization0",
		UserId:         "user0",
		IsSensitive:    false,
		Scope:          client.ScopeEnvironment,
		Type:           client.ConfigurationVariableTypeEnvironment,
		Schema:         client.ConfigurationVariableSchema{Type: "string"},
	}

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
						resource.TestCheckResourceAttr(accessor, "type", strconv.Itoa(int(configurationVariable.Type))),
						//resource.TestCheckResourceAttr(accessor, "project_id", strconv.FormatBool(configurationVariable)),
						//resource.TestCheckResourceAttr(accessor, "template_id", strconv.FormatBool(configurationVariable.)),
						//resource.TestCheckResourceAttr(accessor, "environment_id", strconv.FormatBool(configurationVariable.)),
						//resource.TestCheckResourceAttr(accessor, "deployment_log_id", strconv.FormatBool(configurationVariable.)),
						resource.TestCheckResourceAttr(accessor, "value", configurationVariable.Value),
						resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
						resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(configurationVariable.IsSensitive)),
					),
				},
			},
		}
	}

	t.Run("By id", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(map[string]interface{}{"id": configurationVariable.Id}),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").AnyTimes().Return(configurationVariable, nil)
			})
	})
}
