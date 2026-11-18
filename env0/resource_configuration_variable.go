package env0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const templateScope = "TEMPLATE"

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
			"sub_environment_alias": {
				Type:         schema.TypeString,
				Description:  "create the variable for a specific sub environment of a workflow. Requires template_id to be set. Omit to apply the variable to every sub environment of the workflow template.",
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"template_id"},
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
			"multi_select": {
				Type:          schema.TypeBool,
				Description:   "set to 'true' to create a multi select list variable. Requires 'enum' (the available options) and 'values' (the selected options). Cannot be used with 'value', 'format' or 'is_sensitive'",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"value", "format", "is_sensitive"},
			},
			"values": {
				Type:         schema.TypeList,
				Description:  "the selected options of a multi select list variable (requires 'multi_select'). Each must be one of 'enum', without duplicates. An empty list selects nothing",
				Optional:     true,
				RequiredWith: []string{"multi_select"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	if subEnvironmentAlias, ok := d.GetOk("sub_environment_alias"); ok {
		// The scope id is template_id:sub_environment_alias
		scope = client.ScopeSubEnvironment
		scopeId = d.Get("template_id").(string) + ":" + subEnvironmentAlias.(string)
	} else if templateId, ok := d.GetOk("template_id"); ok {
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

	var err error

	if params.EnumValues, err = getEnum(d); err != nil {
		return nil, err
	}

	if params.MultiSelect {
		if len(params.EnumValues) == 0 {
			return nil, errors.New("'enum' must be set when 'multi_select' is true")
		}

		if params.Value, err = getMultiSelectValue(d, params.EnumValues); err != nil {
			return nil, err
		}
	} else {
		if len(d.Get("values").([]any)) > 0 {
			return nil, errors.New("'values' can only be set when 'multi_select' is true")
		}

		if err := validateEnumValue(params.EnumValues, params.Value); err != nil {
			return nil, err
		}
	}

	if err := validateNilValue(params.IsReadOnly, params.IsRequired, params.Value); err != nil {
		return nil, err
	}

	return &params, nil
}

func resourceConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func getEnum(d *schema.ResourceData) ([]string, error) {
	var actualEnumValues []string

	if specified, ok := d.GetOk("enum"); ok {
		for i, enumValue := range specified.([]any) {
			if enumValue == nil {
				return nil, fmt.Errorf("an empty enum value is not allowed (at index %d)", i)
			}

			actualEnumValues = append(actualEnumValues, enumValue.(string))
		}
	}

	return actualEnumValues, nil
}

func validateEnumValue(enum []string, selectedValue string) error {
	if len(enum) > 0 && !slices.Contains(enum, selectedValue) {
		return fmt.Errorf("value - '%s' is not one of the enum options %v", selectedValue, enum)
	}

	return nil
}

func getMultiSelectValue(d *schema.ResourceData, enum []string) (string, error) {
	ivalues := d.Get("values").([]any)

	selected := make([]string, 0, len(ivalues))

	for _, ivalue := range ivalues {
		value, _ := ivalue.(string)
		selected = append(selected, value)
	}

	return encodeMultiSelectValue(enum, selected)
}

// A multi select value is always a JSON array string, never "" (ENG-3633).
func encodeMultiSelectValue(enum []string, selected []string) (string, error) {
	seen := map[string]bool{}

	for _, value := range selected {
		if err := validateEnumValue(enum, value); err != nil {
			return "", err
		}

		if seen[value] {
			return "", fmt.Errorf("value - '%s' is selected more than once", value)
		}

		seen[value] = true
	}

	encoded, err := json.Marshal(selected)
	if err != nil {
		return "", fmt.Errorf("could not encode the selected values: %w", err)
	}

	return string(encoded), nil
}

func parseMultiSelectValue(value string) ([]string, error) {
	selected := []string{}

	if value == "" {
		return selected, nil
	}

	if err := json.Unmarshal([]byte(value), &selected); err != nil {
		return nil, fmt.Errorf("multi select value '%s' is not a JSON array of strings: %w", value, err)
	}

	return selected, nil
}

func resourceConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if variable.Schema != nil && variable.Schema.Type == client.MultiSelectSchemaType {
		values, err := parseMultiSelectValue(variable.Value)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("multi_select", true)
		d.Set("enum", variable.Schema.Enum)
		d.Set("values", values)
		d.Set("value", "")

		return nil
	}

	d.Set("multi_select", false)
	d.Set("values", nil)

	if variable.IsSensitive == nil || !*variable.IsSensitive {
		d.Set("value", variable.Value)
	}

	return nil
}

func resourceConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceConfigurationVariableDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceConfigurationVariableImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	var configurationParams ConfigurationVariableParams

	inputData := d.Id()

	// soft delete isn't part of the configuration variable, so we need to set it
	d.Set("soft_delete", false)

	err := json.Unmarshal([]byte(inputData), &configurationParams)
	if err != nil {
		return nil, err
	}

	// We need this conversion since getConfigurationVariable query by the scope and in our BE we use blueprint as the scope name instead of template
	if string(configurationParams.Scope) == "TEMPLATE" {
		configurationParams.Scope = "BLUEPRINT"
	}

	variable, getErr := getConfigurationVariable(configurationParams, meta)
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	}

	if variable.Scope == client.ScopeSubEnvironment {
		// The scope id is template_id:sub_environment_alias
		parts := strings.Split(variable.ScopeId, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid sub environment scope id format: %s", variable.ScopeId)
		}

		templateId := parts[0]

		subEnvironmentAlias := parts[1]

		d.Set("template_id", templateId)

		d.Set("sub_environment_alias", subEnvironmentAlias)
	} else {
		var scopeName string

		if variable.Scope == client.ScopeTemplate {
			scopeName = strings.ToLower(templateScope + "_id")
		} else {
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", variable.Scope))
		}

		d.Set(scopeName, configurationParams.ScopeId)
	}

	d.SetId(variable.Id)

	return []*schema.ResourceData{d}, nil
}
