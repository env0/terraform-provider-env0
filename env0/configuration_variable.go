package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getConfigurationVariablesFromSchema(configuration []interface{}) client.ConfigurationChanges {
	configurationChanges := client.ConfigurationChanges{}

	for _, variable := range configuration {
		configurationVariable := getConfigurationVariableFromSchema(variable.(map[string]interface{}))
		configurationChanges = append(configurationChanges, configurationVariable)
	}

	return configurationChanges
}

func getConfigurationVariableFromSchema(variable map[string]interface{}) client.ConfigurationVariable {
	varType, _ := client.GetConfigurationVariableType(variable["type"].(string))

	configurationVariable := client.ConfigurationVariable{
		Name:  variable["name"].(string),
		Value: variable["value"].(string),
		Scope: client.ScopeDeployment,
		Type:  &varType,
	}

	if variable["scope_id"] != nil {
		configurationVariable.ScopeId = variable["scope_id"].(string)
	}

	if variable["is_sensitive"] != nil {
		isSensitive := variable["is_sensitive"].(bool)
		configurationVariable.IsSensitive = &isSensitive
	}

	if variable["is_read_only"] != nil {
		isReadOnly := variable["is_read_only"].(bool)
		configurationVariable.IsReadOnly = &isReadOnly
	}

	if variable["is_required"] != nil {
		isRequired := variable["is_required"].(bool)
		configurationVariable.IsRequired = &isRequired
	}

	if variable["description"] != nil {
		configurationVariable.Description = variable["description"].(string)
	}

	if variable["regex"] != nil {
		configurationVariable.Regex = variable["regex"].(string)
	}

	configurationSchema := client.ConfigurationVariableSchema{
		Format: client.Format(variable["schema_format"].(string)),
		Enum:   nil,
		Type:   variable["schema_type"].(string),
	}

	if variable["schema_type"] != "" && len(variable["schema_enum"].([]interface{})) > 0 {
		enumOfAny := variable["schema_enum"].([]interface{})
		enum := make([]string, len(enumOfAny))

		for i := range enum {
			enum[i] = enumOfAny[i].(string)
		}

		configurationSchema.Type = variable["schema_type"].(string)
		configurationSchema.Enum = enum
	}

	configurationVariable.Schema = &configurationSchema

	return configurationVariable
}

func setEnvironmentConfigurationSchema(ctx context.Context, d *schema.ResourceData, configurationVariables []client.ConfigurationVariable) {
	ivariables, ok := d.GetOk("configuration")
	if !ok {
		return
	}

	if ivariables == nil {
		ivariables = make([]interface{}, 0)
	}

	variables := ivariables.([]interface{})

	newVariables := make([]interface{}, 0)

	// The goal is to maintain existing state order as much as possible. (The backend response order may vary from state).
	for _, ivariable := range variables {
		variable := ivariable.(map[string]interface{})
		variableName := variable["name"].(string)

		for _, configurationVariable := range configurationVariables {
			if configurationVariable.Name == variableName {
				newVariable := createVariable(&configurationVariable)

				if configurationVariable.IsSensitive != nil && *configurationVariable.IsSensitive {
					// To avoid drift for sensitive variables, don't override with the variable value received from API. Use the one in the schema instead.
					newVariable.(map[string]interface{})["value"] = variable["value"]
				}

				newVariables = append(newVariables, newVariable)

				break
			}
		}
	}

	// Check for drifts: add new configuration variables received from the backend.
	for _, configurationVariable := range configurationVariables {
		found := false

		for _, ivariable := range variables {
			variable := ivariable.(map[string]interface{})
			variableName := variable["name"].(string)

			if configurationVariable.Name == variableName {
				found = true

				break
			}
		}

		if !found {
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"configuration name": configurationVariable.Name})
			newVariables = append(newVariables, createVariable(&configurationVariable))
		}
	}

	if len(newVariables) > 0 {
		d.Set("configuration", newVariables)
	} else {
		d.Set("configuration", nil)
	}
}
