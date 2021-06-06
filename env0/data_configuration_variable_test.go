package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestUnitConfigurationVariableData(t *testing.T) {
	resourceType := "env0_organization"
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
		Scope:          "scope0",
		Type:           1,
		Schema:         "",
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, make(map[string]interface{})),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
					resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
					resource.TestCheckResourceAttr(accessor, "type", configurationVariable.CreatedBy),
					resource.TestCheckResourceAttr(accessor, "role", configurationVariable.Role),
					resource.TestCheckResourceAttr(accessor, "project_id", strconv.FormatBool(configurationVariable.IsSelfHosted)),
					resource.TestCheckResourceAttr(accessor, "template_id", strconv.FormatBool(configurationVariable.IsSelfHosted)),
					resource.TestCheckResourceAttr(accessor, "environment_id", strconv.FormatBool(configurationVariable.IsSelfHosted)),
					resource.TestCheckResourceAttr(accessor, "deployment_log_id", strconv.FormatBool(configurationVariable.IsSelfHosted)),
					resource.TestCheckResourceAttr(accessor, "value", strconv.FormatBool(configurationVariable.IsSelfHosted)),
					resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(configurationVariable.IsSensitive)),
					resource.TestCheckResourceAttr(accessor, "scope", strconv.FormatBool(configurationVariable.IsSelfHosted)),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Organization().AnyTimes().Return(configurationVariable, nil)
	})
}
