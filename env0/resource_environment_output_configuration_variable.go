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
	OutputName    string `json:"outputName"`
	ProjectId     string `json:"projectId"`
	EnvironmentId string `json:"environmentId,omitempty"`
	// TODO --- subenvironment...
}

func resourceEnvironmentOutputConfigurationVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentOutputConfigurationVariableCreate,
		ReadContext:   resourceEnvironmentOutputConfigurationVariableRead,
		UpdateContext: resourceConfigurationVariableUpdate,
		DeleteContext: resourceConfigurationVariableDelete,
		Description:   "for configuring environment output configuration variable: https://docs.env0.com/docs/environment-outputs",

		Importer: &schema.ResourceImporter{StateContext: resourceConfigurationVariableImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the variable",
				Required:    true,
			},
			"output_environment_id": {
				Type:         schema.TypeString,
				Description:  "the environment id of the output",
				Optional:     true,
				ExactlyOneOf: []string{"output_environment_id", "output_sub_environment_alias"},
			},
			"output_sub_environment_alias": {
				Type:         schema.TypeString,
				Description:  "TODO: the sub environment alias of the output",
				Optional:     true,
				ExactlyOneOf: []string{"output_environment_id", "output_sub_environment_alias"},
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

func serializeEnvironmentOutputConfigurationVariableValue(params *EnvironmentOutputConfigurationVariableParams, apiClient client.ApiClientInterface) (string, error) {
	environment, err := apiClient.Environment(params.OutputEnvironmentId)
	if err != nil {
		return "", fmt.Errorf("failed to get output environment details: %w", err)
	}

	value := EnvironmentOutputConfigurationVariableValue{
		OutputName:    params.OutputName,
		EnvironmentId: params.OutputEnvironmentId,
		ProjectId:     environment.ProjectId,
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal a value struct: %w", err)
	}

	return string(b), nil
}

func getEnvironmentOutputConfigurationVariableParams(d *schema.ResourceData, apiClient client.ApiClientInterface) (*EnvironmentOutputConfigurationVariableParams, error) {
	variable, err := apiClient.ConfigurationVariablesById(d.Id())
	if err != nil {
		return nil, err
	}

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

	// TODO scope... switch variable.Scope
	// TODO value...
	// TODO output values deserialize...

	return nil, nil
}

func getEnvironmentOutputConfigurationVariable(d *schema.ResourceData, apiClient client.ApiClientInterface) (*client.ConfigurationVariable, error) {
	var params EnvironmentOutputConfigurationVariableParams
	if err := readResourceData(&params, d); err != nil {
		return nil, fmt.Errorf("schema resource data deserialization failed: %v", err)
	}

	if params.Scope != string(client.ScopeProject) && params.IsReadOnly {
		return nil, errors.New("'is_read_only' can only be set to 'true' for the 'PROJECT' scope")
	}

	organizationId, err := apiClient.OrganizationId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization id: %w", err)
	}

	value, err := serializeEnvironmentOutputConfigurationVariableValue(&params, apiClient)
	if err != nil {
		return nil, err
	}

	variableType := client.ConfigurationVariableTypeEnvironment
	if params.Type == "terraform" {
		variableType = client.ConfigurationVariableTypeTerraform
	}

	variable := client.ConfigurationVariable{
		Schema: &client.ConfigurationVariableSchema{
			Type:   "string",
			Format: "ENVIRONMENT_OUTPUT",
		},
		IsSensitive:    boolPtr(false),
		IsRequired:     boolPtr(params.IsRequired),
		IsReadOnly:     boolPtr(params.IsReadOnly),
		ScopeId:        params.ScopeId,
		Scope:          client.Scope(params.Scope),
		OrganizationId: organizationId,
		Name:           params.Name,
		Description:    params.Description,
		Value:          value,
		Type:           &variableType,
	}

	return &variable, nil
}

func resourceEnvironmentOutputConfigurationVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	_, err := getEnvironmentOutputConfigurationVariable(d, apiClient)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	// TODO --- fix
	configurationVariable, err := apiClient.ConfigurationVariableCreate(client.ConfigurationVariableCreateParams{})
	if err != nil {
		return diag.Errorf("could not create environment output configuration variable: %v", err)
	}

	d.SetId(configurationVariable.Id)

	return nil
}

func resourceEnvironmentOutputConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	params, err := getEnvironmentOutputConfigurationVariableParams(d, apiClient)

	if err != nil {
		return ResourceGetFailure(ctx, "environment output configuration variable", d, err)
	}

	if err := writeResourceData(&params, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

/*
func resourceConfigurationVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	params, err := getConfigurationVariableCreateParams(d)
	if err != nil {
		return diag.Errorf(err.Error())
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
	if softDelete := d.Get("soft_delete"); softDelete.(bool) {
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
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", templateScope))
		} else {
			scopeName = strings.ToLower(fmt.Sprintf("%s_id", variable.Scope))
		}

		d.Set(scopeName, configurationParams.ScopeId)

		return []*schema.ResourceData{d}, nil
	}
}
*/
