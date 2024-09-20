package env0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EnvironmentOutputConfigurationVariableParams struct {
	Name                      string
	OutputEnvironmentId       string
	OutputSubEnvironmentAlias string
	OutputName                string
	Scope                     string
	ScopeId                   string
	Description               string
	Type                      string
	IsReadOnly                bool
	IsRequired                bool
}

type EnvironmentOutputConfigurationVariableValue struct {
	OutputName          string `json:"outputName"`
	EnvironmentId       string `json:"environmentId,omitempty"`
	SubEnvironmentAlias string `json:"subEnvironmentAlias,omitempty"`
}

func resourceEnvironmentOutputConfigurationVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentOutputConfigurationVariableCreate,
		ReadContext:   resourceEnvironmentOutputConfigurationVariableRead,
		UpdateContext: resourceEnvironmentOutputConfigurationVariableUpdate,
		DeleteContext: resourceConfigurationVariableDelete,
		Description:   "for additional details check: https://docs.env0.com/docs/environment-outputs",

		Importer: &schema.ResourceImporter{StateContext: resourceEnvironmentOutputConfigurationVariableImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the variable",
				Required:    true,
			},
			"output_environment_id": {
				Type:        schema.TypeString,
				Description: "the environment id of the output",
				Required:    true,
			},
			"output_name": {
				Type:        schema.TypeString,
				Description: "the name of the output value",
				Required:    true,
			},
			"scope": {
				Type:             schema.TypeString,
				Description:      "the type of resource to assign to. Valid values: 'PROJECT', 'ENVIRONMENT', 'WORKFLOW', and 'DEPLOYMENT'. Default value: 'ENVIRONMENT'",
				Optional:         true,
				Default:          "ENVIRONMENT",
				ForceNew:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"PROJECT", "ENVIRONMENT", "WORKFLOW", "DEPLOYMENT"}),
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the id of the of the resource to assign to (E.g. the environment id)",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "a description of the variable",
				Optional:    true,
			},
			"type": {
				Type:             schema.TypeString,
				Description:      "defaults to 'environment'. Set to 'terraform' to create a terraform output variable",
				Optional:         true,
				Default:          "environment",
				ForceNew:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"environment", "terraform"}),
			},
			"is_read_only": {
				Type:        schema.TypeBool,
				Description: "set to 'true' if the value of this variable cannot be edited in lower scopes (applicable only to 'PROJECT' scope)",
				Optional:    true,
				Default:     false,
			},
			"is_required": {
				Type:        schema.TypeBool,
				Description: "set to 'true' if the value of this variable is required in lower scopes",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func serializeEnvironmentOutputConfigurationVariableValue(params *EnvironmentOutputConfigurationVariableParams) (string, error) {
	value := EnvironmentOutputConfigurationVariableValue{
		OutputName:    params.OutputName,
		EnvironmentId: params.OutputEnvironmentId,
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal a value struct: %w", err)
	}

	return string(b), nil
}

func deserializeEnvironmentOutputConfigurationVariableValue(valueStr string) (*EnvironmentOutputConfigurationVariableValue, error) {
	var value EnvironmentOutputConfigurationVariableValue

	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value string: %w", err)
	}

	if value.OutputName == "" {
		return nil, errors.New("after unmarshal 'outputName' is empty")
	}

	if value.EnvironmentId == "" && value.SubEnvironmentAlias == "" {
		return nil, errors.New("after unmarshal both 'environmentId' and 'subEnvironmentAlias' are empty")
	}

	return &value, nil
}

func getEnvironmentOutputConfigurationVariableParamsFromVariable(d *schema.ResourceData, variable *client.ConfigurationVariable) (*EnvironmentOutputConfigurationVariableParams, error) {
	var params EnvironmentOutputConfigurationVariableParams
	if err := readResourceData(&params, d); err != nil {
		return nil, fmt.Errorf("schema resource data deserialization failed: %v", err)
	}

	params.Name = variable.Name
	params.Description = variable.Description

	if variable.IsReadOnly != nil {
		params.IsReadOnly = *variable.IsReadOnly
	} else {
		params.IsReadOnly = false
	}

	if variable.IsRequired != nil {
		params.IsRequired = *variable.IsRequired
	} else {
		params.IsRequired = false
	}

	if variable.Type == nil || *variable.Type == client.ConfigurationVariableTypeEnvironment {
		params.Type = "environment"
	} else {
		params.Type = "terraform"
	}

	params.ScopeId = variable.ScopeId

	switch scope := variable.Scope; scope {
	case client.ScopeEnvironment, client.ScopeDeployment, client.ScopeWorkflow, client.ScopeProject:
		params.Scope = string(scope)
	default:
		return nil, fmt.Errorf("invalid scope %s", scope)
	}

	value, err := deserializeEnvironmentOutputConfigurationVariableValue(variable.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	params.OutputEnvironmentId = value.EnvironmentId
	params.OutputSubEnvironmentAlias = value.SubEnvironmentAlias
	params.OutputName = value.OutputName

	return &params, nil
}

func getEnvironmentOutputCreateParams(d *schema.ResourceData) (*client.ConfigurationVariableCreateParams, error) {
	var params EnvironmentOutputConfigurationVariableParams
	if err := readResourceData(&params, d); err != nil {
		return nil, fmt.Errorf("schema resource data deserialization failed: %v", err)
	}

	if params.Scope != string(client.ScopeProject) && params.IsReadOnly {
		return nil, errors.New("'is_read_only' can only be set to 'true' for the 'PROJECT' scope")
	}

	value, err := serializeEnvironmentOutputConfigurationVariableValue(&params)
	if err != nil {
		return nil, err
	}

	variableType := client.ConfigurationVariableTypeEnvironment
	if params.Type == "terraform" {
		variableType = client.ConfigurationVariableTypeTerraform
	}

	createParams := client.ConfigurationVariableCreateParams{
		Format:      client.ENVIRONMENT_OUTPUT,
		IsRequired:  params.IsRequired,
		IsReadOnly:  params.IsReadOnly,
		ScopeId:     params.ScopeId,
		Scope:       client.Scope(params.Scope),
		Name:        params.Name,
		Description: params.Description,
		Value:       value,
		Type:        variableType,
	}

	return &createParams, nil
}

func resourceEnvironmentOutputConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	createParams, err := getEnvironmentOutputCreateParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	configurationVariable, err := apiClient.ConfigurationVariableCreate(*createParams)
	if err != nil {
		return diag.Errorf("could not create environment output configuration variable: %v", err)
	}

	d.SetId(configurationVariable.Id)

	return nil
}

func resourceEnvironmentOutputConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	variable, err := apiClient.ConfigurationVariablesById(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "environment output configuration variable", d, err)
	}

	params, err := getEnvironmentOutputConfigurationVariableParamsFromVariable(d, &variable)
	if err != nil {
		return diag.Errorf("failed to get params from configuration variable: %v", err)
	}

	if err := writeResourceData(params, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceEnvironmentOutputConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	createParams, err := getEnvironmentOutputCreateParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	if _, err := apiClient.ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{Id: id, CommonParams: *createParams}); err != nil {
		return diag.Errorf("could not update environment output configuration variable: %v", err)
	}

	return nil
}

func resourceEnvironmentOutputConfigurationVariableImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	}

	d.SetId(variable.Id)

	params, err := getEnvironmentOutputConfigurationVariableParamsFromVariable(d, &variable)
	if err != nil {
		return nil, fmt.Errorf("failed to get params from configuration variable: %w", err)
	}

	if err := writeResourceData(params, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
