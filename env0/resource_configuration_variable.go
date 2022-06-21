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
				Description:  "specifies the format of the configuration value (HCL/JSON)",
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

func resourceConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope, scopeId := whichScope(d)
	params := client.ConfigurationVariableCreateParams{Scope: scope, ScopeId: scopeId}
	if err := readResourceData(&params, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := validateNilValue(params.IsReadOnly, params.IsRequired, params.Value); err != nil {
		return diag.Errorf(err.Error())
	}

	actualEnumValues, getEnumErr := getEnum(d, params.Value)
	if getEnumErr != nil {
		return getEnumErr
	}

	params.EnumValues = actualEnumValues

	apiClient := meta.(client.ApiClientInterface)

	configurationVariable, err := apiClient.ConfigurationVariableCreate(params)
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
		for i, enumValue := range enumValues {
			if enumValue == nil {
				return nil, diag.Errorf("an empty enum value is not allowed (at index %d)", i)
			}
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
	variable, err := apiClient.ConfigurationVariablesById(id)

	if err != nil {
		return ResourceGetFailure("configuration variable", d, err)
	}

	d.Set("name", variable.Name)
	d.Set("description", variable.Description)
	d.Set("value", variable.Value)
	d.Set("is_sensitive", variable.IsSensitive)
	d.Set("is_read_only", variable.IsReadOnly)
	d.Set("is_required", variable.IsRequired)
	d.Set("regex", variable.Regex)
	if variable.Type != nil && *variable.Type == client.ConfigurationVariableTypeTerraform {
		d.Set("type", "terraform")
	} else {
		d.Set("type", "environment")
	}
	if variable.Schema != nil {
		if len(variable.Schema.Enum) > 0 {
			d.Set("enum", variable.Schema.Enum)
		}

		if variable.Schema.Format != "" {
			d.Set("format", variable.Schema.Format)
		}
	}

	return nil
}

func resourceConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	scope, scopeId := whichScope(d)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	value := d.Get("value").(string)
	isSensitive := d.Get("is_sensitive").(bool)
	typeAsString := d.Get("type").(string)
	format := client.Format(d.Get("format").(string))
	isReadOnly := d.Get("is_read_only").(bool)
	isRequired := d.Get("is_required").(bool)
	regex := d.Get("regex").(string)

	if err := validateNilValue(isReadOnly, isRequired, value); err != nil {
		return diag.Errorf(err.Error())
	}

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
	_, err := apiClient.ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{Id: id, CommonParams: client.ConfigurationVariableCreateParams{
		Name:        name,
		Value:       value,
		IsSensitive: isSensitive,
		Scope:       scope,
		ScopeId:     scopeId,
		Type:        type_,
		EnumValues:  actualEnumValues,
		Description: description,
		Format:      format,
		IsReadOnly:  isReadOnly,
		IsRequired:  isRequired,
		Regex:       regex,
	}})
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

		var scopeName string

		if variable.Scope == client.ScopeTemplate {
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", templateScope))
		} else {
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", variable.Scope))
		}

		d.Set(scopeName, configurationParams.ScopeId)

		return []*schema.ResourceData{d}, nil
	}
}
