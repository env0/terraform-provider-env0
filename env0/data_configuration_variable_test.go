package env0

import (
	"errors"

	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitConfigurationVariableData(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)
	isSensitive := false
	isReadonly := true
	isRequired := false
	variableType := client.ConfigurationVariableTypeEnvironment
	configurationVariable := client.ConfigurationVariable{
		Id:             "id0",
		Name:           "name0",
		Description:    "desc0",
		ScopeId:        "scope0",
		Value:          "value0",
		OrganizationId: "organization0",
		UserId:         "user0",
		IsSensitive:    &isSensitive,
		Scope:          client.ScopeEnvironment,
		Type:           &variableType,
		Schema:         &client.ConfigurationVariableSchema{Type: "string", Format: client.HCL},
		IsReadOnly:     &isReadonly,
		IsRequired:     &isRequired,
		Regex:          "regex",
	}

	checkResources := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
		resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
		resource.TestCheckResourceAttr(accessor, "description", configurationVariable.Description),
		resource.TestCheckResourceAttr(accessor, "type", "environment"),
		resource.TestCheckResourceAttr(accessor, "value", configurationVariable.Value),
		resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
		resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(*configurationVariable.IsSensitive)),
		resource.TestCheckResourceAttr(accessor, "format", string(configurationVariable.Schema.Format)),
		resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configurationVariable.IsReadOnly)),
		resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configurationVariable.IsRequired)),
		resource.TestCheckResourceAttr(accessor, "regex", "regex"),
	)

	t.Run("ScopeGlobal", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"id": configurationVariable.Id}),
						Check:  checkResources,
					},
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).AnyTimes().
					Return(configurationVariable, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})
	t.Run("ScopeGlobal Enum", func(t *testing.T) {
		isSensitive := false
		variableType := client.ConfigurationVariableTypeEnvironment
		configurationVariable := client.ConfigurationVariable{
			Id:             "id0",
			Name:           "name0",
			ScopeId:        "scope0",
			Value:          "value0",
			OrganizationId: "organization0",
			UserId:         "user0",
			IsSensitive:    &isSensitive,
			Scope:          client.ScopeEnvironment,
			Type:           &variableType,
			Schema:         &client.ConfigurationVariableSchema{Type: "string", Enum: []string{"a", "b"}},
		}

		checkResources := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(accessor, "id", configurationVariable.Id),
			resource.TestCheckResourceAttr(accessor, "name", configurationVariable.Name),
			resource.TestCheckResourceAttr(accessor, "type", "environment"),
			resource.TestCheckResourceAttr(accessor, "value", configurationVariable.Value),
			resource.TestCheckResourceAttr(accessor, "scope", string(configurationVariable.Scope)),
			resource.TestCheckResourceAttr(accessor, "is_sensitive", strconv.FormatBool(*configurationVariable.IsSensitive)),
			resource.TestCheckResourceAttr(accessor, "enum.0", "a"),
			resource.TestCheckResourceAttr(accessor, "enum.1", "b"),
		)

		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"id": configurationVariable.Id}),
						Check:  checkResources,
					},
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesById(configurationVariable.Id).AnyTimes().
					Return(configurationVariable, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})

	t.Run("ScopeTemplate", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"template_id": "template_id", "name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeTemplate, "template_id").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})

	t.Run("ScopeEnvironment", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"environment_id": configurationVariable.Id, "name": configurationVariable.Name}),
						Check:  checkResources,
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, configurationVariable.Id).AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})

	t.Run("configuration variable not exists in the server", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{"name": "invalid"}),
						ExpectError: regexp.MustCompile("Could not query variables"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeGlobal, "").AnyTimes().
					Return(nil, errors.New("not found"))
			})
	})

	t.Run("configuration variable not match to name", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{"name": "invalid"}),
						ExpectError: regexp.MustCompile("Could not find variable"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeGlobal, "").AnyTimes().
					Return([]client.ConfigurationVariable{configurationVariable}, nil)
			})
	})
}
