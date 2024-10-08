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
		Description:   "Note: do not use with an environment resource that has it's configuration variables defined in it's 'configuration' field (see env0_environment_resource -> configuration)",

		Importer: &schema.ResourceImporter{StateContext: resourceConfigurationVariableImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the configuration variable",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "a description of the variables",
				Optional:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "value for the configuration variable",
				Optional:    true,
				Sensitive:   true,
				Default:     "",
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
				ForceNew:      true,
			},
			"template_id": {
				Type:          schema.TypeString,
				Description:   "create the variable under this template, not globally",
				Optional:      true,
				ConflictsWith: []string{"project_id", "environment_id"},
				ForceNew:      true,
			},
			"environment_id": {
				Type:          schema.TypeString,
				Description:   "create the variable under this environment, not globally. Make sure to 'ignore changes' on environment.configuration to prevent drifts",
				Optional:      true,
				ConflictsWith: []string{"template_id", "project_id", "is_required", "is_read_only"},
				ForceNew:      true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "default 'environment'. set to 'terraform' to create a terraform variable",
				Optional:    true,
				Default:     "environment",
				ForceNew:    true,
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
			"format": {
				Type:         schema.TypeString,
				Description:  "specifies the format of the configuration value ('HCL' or 'JSON'). If none is specified, 'JSON' and 'HCL' values will be considered to be a 'string' (text) type",
				Default:      "",
				Optional:     true,
				ValidateFunc: ValidateConfigurationPropertySchema,
			},
			"is_read_only": {
				Type:          schema.TypeBool,
				Description:   "the value of this variable cannot be edited by lower scopes",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"environment_id"},
			},
			"is_required": {
				Type:          schema.TypeBool,
				Description:   "the value of this variable must be set by lower scopes",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"environment_id"},
			},
			"regex": {
				Type:        schema.TypeString,
				Description: "the value of this variable must match provided regular expression (enforced only in env0 UI)",
				Optional:    true,
			},
			"soft_delete": {
				Type:        schema.TypeBool,
				Description: "soft delete the configuration variable, once removed from the configuration it won't be deleted from env0",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

const templateScope = "TEMPLATE"

func validateNilValue(isReadOnly bool, isRequired bool, value string) error {
	if isReadOnly && isRequired && value == "" {
		return errors.New("'value' cannot be empty when 'is_read_only' and 'is_required' are true ")
	}
	return nil
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
	if environmentId, ok := d.GetOk("environment_id"); ok {
		scope = client.ScopeEnvironment
		scopeId = environmentId.(string)
	}

	return scope, scopeId
}

func getConfigurationVariableCreateParams(d *schema.ResourceData) (*client.ConfigurationVariableCreateParams, error) {
	scope, scopeId := whichScope(d)
	params := client.ConfigurationVariableCreateParams{Scope: scope, ScopeId: scopeId}
	if err := readResourceData(&params, d); err != nil {
		return nil, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	if err := validateNilValue(params.IsReadOnly, params.IsRequired, params.Value); err != nil {
		return nil, err
	}

	var err error
	if params.EnumValues, err = getEnum(d, params.Value); err != nil {
		return nil, err
	}

	return &params, nil
}

func resourceConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	params, err := getConfigurationVariableCreateParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(client.ApiClientInterface)

	configurationVariable, err := apiClient.ConfigurationVariableCreate(*params)
	if err != nil {
		return diag.Errorf("could not create configurationVariable: %v", err)
	}

	d.SetId(configurationVariable.Id)

	return nil
}

func getEnum(d *schema.ResourceData, selectedValue string) ([]string, error) {
	var enumValues []interface{}
	var actualEnumValues []string
	if specified, ok := d.GetOk("enum"); ok {
		enumValues = specified.([]interface{})
		valueExists := false

		for i, enumValue := range enumValues {
			if enumValue == nil {
				return nil, fmt.Errorf("an empty enum value is not allowed (at index %d)", i)
			}

			actualEnumValues = append(actualEnumValues, enumValue.(string))

			if enumValue == selectedValue {
				valueExists = true
			}
		}
		if !valueExists {
			return nil, fmt.Errorf("value - '%s' is not one of the enum options %v", selectedValue, actualEnumValues)
		}
	}
	return actualEnumValues, nil
}

func resourceConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	variable, err := apiClient.ConfigurationVariablesById(id)

	if err != nil {
		return ResourceGetFailure(ctx, "configuration variable", d, err)
	}

	d.Set("type", "environment")

	if err := writeResourceData(&variable, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if variable.IsSensitive == nil || !*variable.IsSensitive {
		d.Set("value", variable.Value)
	}

	return nil
}

func resourceConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	params, err := getConfigurationVariableCreateParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	if _, err := apiClient.ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{Id: id, CommonParams: *params}); err != nil {
		return diag.Errorf("could not update configurationVariable: %v", err)
	}

	return nil
}

func resourceConfigurationVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// don't delete if soft delete is set
	if softDelete := d.Get("soft_delete"); softDelete != nil && softDelete.(bool) {
		return nil
	}

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

	// soft delete isn't part of the configuration variable, so we need to set it
	d.Set("soft_delete", false)

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

		var scopeName string

		if variable.Scope == client.ScopeTemplate {
			scopeName = strings.ToLower(templateScope + "_id")
		} else {
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", variable.Scope))
		}

		d.Set(scopeName, configurationParams.ScopeId)

		return []*schema.ResourceData{d}, nil
	}
}
