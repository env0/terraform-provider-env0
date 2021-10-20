package env0

import (
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
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
		Description:    "desc0",
		ScopeId:        "scope0",
		Value:          "value0",
		OrganizationId: "organization0",
		UserId:         "user0",
		IsSensitive:    false,
		Scope:          client.ScopeEnvironment,
		Type:           client.ConfigurationVariableTypeEnvironment,
		Schema:         client.ConfigurationVariableSchema{Type: "string"},
	}

	checkResources := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
		resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
		resource.TestCheckResourceAttr(accessor, "description", configurationVariable.Description),
		resource.TestCheckResourceAttr(accessor, "type", "environment"),
		resource.TestCheckResourceAttr(accessor, "value", configurationVariable.Value),
		resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
		resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(configurationVariable.IsSensitive)),
	)

	t.Run("ScopeGlobal", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"id": configurationVariable.Id}),
						Check:  checkResources,
					},
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})
	t.Run("ScopeGlobal Enum", func(t *testing.T) {

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
			Schema:         client.ConfigurationVariableSchema{Type: "string", Enum: []string{"a", "b"}},
		}

		checkResources := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
			resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
			resource.TestCheckResourceAttr(accessor, "type", "environment"),
			resource.TestCheckResourceAttr(accessor, "value", configurationVariable.Value),
			resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
			resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(configurationVariable.IsSensitive)),
			resource.TestCheckResourceAttr(accessor, "enum.0", "a"),
			resource.TestCheckResourceAttr(accessor, "enum.1", "b"),
		)

		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"id": configurationVariable.Id}),
						Check:  checkResources,
					},
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})

	t.Run("ScopeTemplate", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"id": configurationVariable.Id, "template_id": "template_id"}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeTemplate, "template_id").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})

	t.Run("configuration variable not exists in the server", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "invalid"}),
						ExpectError: regexp.MustCompile("Could not query variables"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").AnyTimes().
					Return(nil, errors.New("not found"))
			})
	})

	t.Run("configuration variable not match to name", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "invalid"}),
						ExpectError: regexp.MustCompile("Could not find variable"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariables(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})
}
