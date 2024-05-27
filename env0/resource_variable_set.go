package env0

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var variableSetVariableSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "variable name",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "variable value for 'hcl', 'json', or 'text' format",
			Optional:    true,
		},
		"dropdown_values": {
			Type:        schema.TypeList,
			Description: "a list of variable values for 'dropdown' format",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"type": {
			Type:             schema.TypeString,
			Description:      "variable type: terraform or environment (defaults to 'environment')",
			Default:          "environment",
			Optional:         true,
			ValidateDiagFunc: NewStringInValidator([]string{"environment", "terraform"}),
		},
		"is_sensitive": {
			Type:        schema.TypeBool,
			Description: "is the value sensitive (defaults to 'false'). Note: 'dropdown' value format cannot be senstive.",
			Optional:    true,
			Default:     false,
		},
		"format": {
			Type:             schema.TypeString,
			Description:      "the value format: 'text' (free text), 'dropdown' (dropdown list), 'hcl', 'json'. Note: 'hcl' and 'json' can only be used in terraform variables.",
			Optional:         true,
			Default:          "text",
			ValidateDiagFunc: NewStringInValidator([]string{"text", "dropdown", "hcl", "json"}),
		},
	},
}

func resourceVariableSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableSetCreate,
		ReadContext:   resourceVariableSetRead,
		UpdateContext: resourceSshKeyUpdate, // TODO - update
		DeleteContext: resourceVariableSetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the variable set",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "the description of the variable set",
				Optional:    true,
			},
			"scope": {
				Type:             schema.TypeString,
				Description:      "the scope of the variable set: 'organization', or 'project' (defaults to 'organization')",
				Optional:         true,
				Default:          "organization",
				ValidateDiagFunc: NewStringInValidator([]string{"organization", "project"}),
				ForceNew:         true,
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the scope id (e.g. project id). Note: not required for organization scope.",
				Optional:    true,
				ForceNew:    true,
			},
			"variable": {
				Type:        schema.TypeList,
				Description: "terraform or environment variable",
				Optional:    true,
				Elem:        variableSetVariableSchema,
			},
		},
	}
}

func getVariableFromSchema(d map[string]interface{}) (*client.ConfigurationVariable, error) {
	var res client.ConfigurationVariable

	res.Scope = "SET"
	res.Name = d["name"].(string)

	isSensitive, ok := d["is_senstive"].(bool)
	if !ok {
		isSensitive = false
	}
	res.IsSensitive = boolPtr(isSensitive)

	variableType := d["type"].(string)
	if variableType == "terraform" {
		res.Type = (*client.ConfigurationVariableType)(intPtr(1))
	} else {
		res.Type = (*client.ConfigurationVariableType)(intPtr(0))
	}

	value, ok := d["value"].(string)
	if !ok {
		value = ""
	} else {
		res.Value = value
	}

	switch format := d["format"].(string); format {
	case "text":
		if len(value) == 0 {
			return nil, fmt.Errorf("free text variable %s must have a value: ", res.Name)
		}
		res.Schema = &client.ConfigurationVariableSchema{
			Type: "string",
		}
	case "hcl":
		if len(value) == 0 {
			return nil, fmt.Errorf("HCL variable %s must have a value: ", res.Name)
		}
		res.Schema = &client.ConfigurationVariableSchema{
			Format: "HCL",
		}
	case "json":
		if len(value) == 0 {
			return nil, fmt.Errorf("JSON variable %s must have a value: ", res.Name)
		}
		res.Schema = &client.ConfigurationVariableSchema{
			Format: "JSON",
		}
		// validate JSON.
		var js json.RawMessage
		if err := json.Unmarshal([]byte(value), &js); err != nil {
			return nil, fmt.Errorf("JSON variable %s is not a valid json value: %w", res.Name, err)
		}
	case "dropdown":
		ivalues, ok := d["dropdown_values"].([]interface{})
		if !ok || len(ivalues) == 0 {
			return nil, fmt.Errorf("dropdown variables %s must have dropdown_values", res.Name)
		}

		var values []string
		for _, ivalue := range ivalues {
			values = append(values, ivalue.(string))
		}

		res.Value = ivalues[0].(string)
		res.Schema = &client.ConfigurationVariableSchema{
			Type: "string",
			Enum: values,
		}
	}

	return &res, nil
}

func getSchemaFromVariables(variables []client.ConfigurationVariable) (interface{}, error) {
	res := make([]interface{}, 0)

	for _, variable := range variables {
		ivariable := make(map[string]interface{})
		res = append(res, ivariable)

		ivariable["name"] = variable.Name
		if len(variable.Description) > 0 {
			ivariable["description"] = variable.Description
		}

		if variable.Type == nil || *variable.Type == client.ConfigurationVariableTypeEnvironment {
			ivariable["type"] = "environment"
		} else {
			ivariable["type"] = "terraform"
		}

		if variable.IsSensitive == nil || !*variable.IsSensitive {
			ivariable["is_sensitive"] = false
		} else {
			ivariable["is_sensitive"] = true
		}

		if variable.Schema.Type == "string" {
			if len(variable.Schema.Enum) > 0 {
				ivariable["format"] = "dropdown"
				ivalues := make([]interface{}, 0)
				ivariable["dropdown_values"] = ivalues
				for _, value := range variable.Schema.Enum {
					ivalues = append(ivalues, value)
				}
			} else {
				ivariable["format"] = "text"
				ivariable["value"] = variable.Value
			}
		} else if variable.Schema.Format == "HCL" {
			ivariable["format"] = "hcl"
			ivariable["value"] = variable.Value
		} else if variable.Schema.Format == "JSON" {
			ivariable["format"] = "json"
			ivariable["value"] = variable.Value
		} else {
			return nil, fmt.Errorf("unhandled variable use-case: %s", variable.Name)
		}
	}

	return res, nil
}

func getVariablesFromSchema(d *schema.ResourceData, organizationId string) ([]client.ConfigurationVariable, error) {
	res := []client.ConfigurationVariable{}

	ivariables, ok := d.GetOk("variable")
	if !ok {
		return res, nil
	}

	for _, ivariable := range ivariables.([]interface{}) {
		variable, err := getVariableFromSchema(ivariable.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		variable.OrganizationId = organizationId
		res = append(res, *variable)
	}

	return res, nil
}

func resourceVariableSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error

	apiClient := meta.(client.ApiClientInterface)

	organizationId, err := apiClient.OrganizationId()
	if err != nil {
		return diag.Errorf("failed to get organization id")
	}

	var payload client.CreateConfigurationSetPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if payload.Scope != "organization" && payload.ScopeId == "" {
		return diag.Errorf("scope_id must be configured for the scope '%s'", payload.Scope)
	}

	if payload.ConfigurationProperties, err = getVariablesFromSchema(d, organizationId); err != nil {
		return diag.Errorf("failed to get variables from schema: %v", err)
	}

	configurationSet, err := apiClient.ConfigurationSetCreate(&payload)
	if err != nil {
		return diag.Errorf("failed to create a variable set: %v", err)
	}

	d.SetId(configurationSet.Id)

	return nil
}

func resourceVariableSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	if err := apiClient.ConfigurationSetDelete(id); err != nil {
		return diag.Errorf("failed to delete a variable set: %v", err)
	}

	return nil
}

type mergedVariables struct {
	currentVariables []client.ConfigurationVariable
	newVariables     []client.ConfigurationVariable
	deletedVariables []client.ConfigurationVariable
}

func mergeVariables(schema []client.ConfigurationVariable, api []client.ConfigurationVariable) *mergedVariables {
	var res mergedVariables

	// To avoid false drifts, keep the order of the 'currentVariables' list similiar to the schema as much as possible.
	for _, svariable := range schema {
		found := false

		for _, avariable := range api {
			if svariable.Name == avariable.Name {
				found = true
				if avariable.IsSensitive != nil && *avariable.IsSensitive {
					// Senstive - to avoid drift use the value from the schema
					avariable.Value = svariable.Value
				}
				res.currentVariables = append(res.currentVariables, avariable)
				break
			}

			if !found {
				// found a variable in the schema but not in the api - this is a new variable.
				res.newVariables = append(res.newVariables, svariable)
			}
		}
	}

	for _, avariable := range api {
		found := false

		for _, svariable := range schema {
			if svariable.Name == avariable.Name {
				found = true
				break
			}
		}

		if !found {
			// found a variable in the api but not in the schema - this is a deleted variable.
			res.deletedVariables = append(res.deletedVariables, avariable)
		}
	}

	return &res
}

func resourceVariableSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	configurationSet, err := apiClient.ConfigurationSet(id)
	if err != nil {
		return ResourceGetFailure(ctx, "variable_set", d, err)
	}

	if err := writeResourceData(configurationSet, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	variablesFromApi, err := apiClient.ConfigurationVariablesBySetId(id)
	if err != nil {
		return diag.Errorf("failed to get variables from the variables set: %v", err)
	}

	variablesFromSchema, err := getVariablesFromSchema(d, "")
	if err != nil {
		return diag.Errorf("failed to get variables from schema: %v", err)
	}

	mergedVariables := mergeVariables(variablesFromSchema, variablesFromApi)

	// for "READ" the source of truth is the variables from the API - existing + deleted.
	variables := append(mergedVariables.currentVariables, mergedVariables.deletedVariables...)

	ivariables, err := getSchemaFromVariables(variables)
	if err != nil {
		return diag.Errorf("failed to get schema from variables: %v", err)
	}

	if err := d.Set("variable", ivariables); err != nil {
		return diag.Errorf("failed to set variables in schema: %v", err)
	}

	return nil
}
