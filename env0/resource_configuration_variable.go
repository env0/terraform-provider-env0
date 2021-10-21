package env0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConfigurationVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationVariableCreate,
		ReadContext:   resourceConfigurationVariableRead,
		UpdateContext: resourceConfigurationVariableUpdate,
		DeleteContext: resourceConfigurationVariableDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceConfigurationVariableImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the configuration variable",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "value for the configuration variable",
				Required:    true,
				Sensitive:   true,
			},
			"is_sensitive": {
				Type:        schema.TypeBool,
				Description: "is the variable sensitive, defaults to false",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"project_id": {
				Type:          schema.TypeString,
				Description:   "create the variable under this project, not globally",
				Optional:      true,
				ConflictsWith: []string{"template_id", "environment_id"},
			},
			"template_id": {
				Type:          schema.TypeString,
				Description:   "create the variable under this template, not globally",
				Optional:      true,
				ConflictsWith: []string{"project_id", "environment_id"},
			},
			"environment_id": {
				Type:          schema.TypeString,
				Description:   "create the variable under this environment, not globally",
				Optional:      true,
				ConflictsWith: []string{"template_id", "project_id"},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "default 'environment'. set to 'terraform' to create a terraform variable",
				Optional:    true,
				Default:     "environment",
			},
			"enum": {
				Type:        schema.TypeList,
				Description: "limit possible values to values from this list",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "name to give the configuration variable",
				},
			},
		},
	}
}

func whichScope(d *schema.ResourceData) (client.Scope, string) {
	scope := client.ScopeGlobal
	scopeId := ""
	if projectId, ok := d.GetOk("project_id"); ok {
		scope = client.ScopeProject
		scopeId = projectId.(string)
	}
	if templateId, ok := d.GetOk("template_id"); ok {
		scope = client.ScopeTemplate
		scopeId = templateId.(string)
	}
	if templateId, ok := d.GetOk("blueprint_id"); ok {
		scope = client.ScopeTemplate
		scopeId = templateId.(string)
	}
	if environmentId, ok := d.GetOk("environment_id"); ok {
		scope = client.ScopeEnvironment
		scopeId = environmentId.(string)
	}

	return scope, scopeId
}

func resourceConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	scope, scopeId := whichScope(d)
	name := d.Get("name").(string)
	value := d.Get("value").(string)
	isSensitive := d.Get("is_sensitive").(bool)
	typeAsString := d.Get("type").(string)
	var type_ client.ConfigurationVariableType
	switch typeAsString {
	case "environment":
		type_ = client.ConfigurationVariableTypeEnvironment
	case "terraform":
		type_ = client.ConfigurationVariableTypeTerraform
	default:
		return diag.Errorf("'type' can only receive either 'environment' or 'terraform': %s", typeAsString)
	}
	actualEnumValues, getEnumErr := getEnum(d, value)
	if getEnumErr != nil {
		return getEnumErr
	}

	configurationVariable, err := apiClient.ConfigurationVariableCreate(name, value, isSensitive, scope, scopeId, type_, actualEnumValues)
	if err != nil {
		return diag.Errorf("could not create configurationVariable: %v", err)
	}

	d.SetId(configurationVariable.Id)

	return nil
}

func getEnum(d *schema.ResourceData, selectedValue string) ([]string, diag.Diagnostics) {
	var enumValues []interface{}
	var actualEnumValues []string
	if specified, ok := d.GetOk("enum"); ok {
		enumValues = specified.([]interface{})
		valueExists := false
		for _, enumValue := range enumValues {
			actualEnumValues = append(actualEnumValues, enumValue.(string))
			if enumValue == selectedValue {
				valueExists = true
			}
		}
		if !valueExists {
			return nil, diag.Errorf("value - '%s' is not one of the enum options %v", selectedValue, actualEnumValues)
		}
	}
	return actualEnumValues, nil
}

func resourceConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	scope, scopeId := whichScope(d)
	variables, err := apiClient.ConfigurationVariables(scope, scopeId)
	if err != nil {
		return diag.Errorf("could not get configurationVariable: %v", err)
	}
	for _, variable := range variables {
		if variable.Id == id {
			d.Set("name", variable.Name)
			d.Set("value", variable.Value)
			d.Set("is_sensitive", variable.IsSensitive)
			if variable.Type == client.ConfigurationVariableTypeTerraform {
				d.Set("type", "terraform")
			} else {
				d.Set("type", "environment")
			}
			if len(variable.Schema.Enum) > 0 {
				d.Set("enum", variable.Schema.Enum)
			}
			return nil
		}
	}
	return diag.Errorf("variable %s not found (under this scope): %v", id, err)
}

func resourceConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	scope, scopeId := whichScope(d)
	name := d.Get("name").(string)
	value := d.Get("value").(string)
	isSensitive := d.Get("is_sensitive").(bool)
	typeAsString := d.Get("type").(string)
	var type_ client.ConfigurationVariableType
	switch typeAsString {
	case "environment":
		type_ = client.ConfigurationVariableTypeEnvironment
	case "terraform":
		type_ = client.ConfigurationVariableTypeTerraform
	default:
		return diag.Errorf("'type' can only receive either 'environment' or 'terraform': %s", typeAsString)
	}
	actualEnumValues, getEnumErr := getEnum(d, value)
	if getEnumErr != nil {
		return getEnumErr
	}
	_, err := apiClient.ConfigurationVariableUpdate(id, name, value, isSensitive, scope, scopeId, type_, actualEnumValues)
	if err != nil {
		return diag.Errorf("could not update configurationVariable: %v", err)
	}

	return nil
}

func resourceConfigurationVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.ConfigurationVariableDelete(id)
	if err != nil {
		return diag.Errorf("could not delete configurationVariable: %v", err)
	}
	return nil
}

func resourceConfigurationVariableImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var configurationParams ConfigurationVariableParams
	inputData := d.Id()
	err := json.Unmarshal([]byte(inputData), &configurationParams)
	// We need this conversion since getConfigurationVariable query by the scope and in our BE we use blueprint as the scope name instead of template
	if string(configurationParams.Scope) == "TEMPLATE" {
		configurationParams.Scope = "BLUEPRINT"
	}
	if err != nil {
		return nil, err
	}
	variable, getErr := getConfigurationVariable(configurationParams, meta)
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		d.SetId(variable.Id)
		scopeName := strings.ToLower(fmt.Sprintf("%s_id", variable.Scope))

		d.Set(scopeName, configurationParams.ScopeId)

		return []*schema.ResourceData{d}, nil
	}
}
