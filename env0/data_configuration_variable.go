package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigurationVariableParams struct {
	Scope             client.Scope `tfschema:"-"`
	ScopeId           string       `tfschema:"-"`
	Id                string
	Name              string
	ConfigurationType string `json:"-" tfschema:"type"`
}

func dataConfigurationVariable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataConfigurationVariableRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the configuration variable",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "a description of the variable",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "'terraform' or 'environment'. If specified as an argument, limits searching by variable name only to variables of this type.",
				Optional:    true,
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the configuration variable",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Description:   "search for the variable under this project, not globally",
				Optional:      true,
				ConflictsWith: []string{"template_id", "environment_id", "deployment_log_id"},
			},
			"template_id": {
				Type:          schema.TypeString,
				Description:   "search for the variable under this template, not globally",
				Optional:      true,
				ConflictsWith: []string{"project_id", "environment_id", "deployment_log_id"},
			},
			"environment_id": {
				Type:          schema.TypeString,
				Description:   "search for the variable under this environment, not globally",
				Optional:      true,
				ConflictsWith: []string{"template_id", "project_id", "deployment_log_id"},
			},
			"deployment_log_id": {
				Type:          schema.TypeString,
				Description:   "search for the variable under this deployment log, not globally",
				Optional:      true,
				ConflictsWith: []string{"template_id", "environment_id", "project_id"},
			},
			"value": {
				Type:        schema.TypeString,
				Description: "value stored in the variable",
				Computed:    true,
				Sensitive:   true,
			},
			"is_sensitive": {
				Type:        schema.TypeBool,
				Description: "is the variable defined as sensitive",
				Computed:    true,
			},
			"scope": {
				Type:        schema.TypeString,
				Description: "scope of the variable",
				Computed:    true,
			},
			"enum": {
				Type:        schema.TypeList,
				Description: "possible values of this variable",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the configuration variable option",
				},
			},
			"format": {
				Type:        schema.TypeString,
				Description: "specifies the format of the configuration value (HCL/JSON)",
				Computed:    true,
			},
			"is_read_only": {
				Type:        schema.TypeBool,
				Description: "specifies if the value of this variable cannot be edited by lower scopes",
				Computed:    true,
				Optional:    true,
			},
			"is_required": {
				Type:        schema.TypeBool,
				Description: "specifies if the value of this variable must be set by lower scopes",
				Computed:    true,
				Optional:    true,
			},
			"regex": {
				Type:        schema.TypeString,
				Description: "specifies a regular expression to validate variable value (enforced only in env0 UI)",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	scope, scopeId := getScopeAndId(d)

	params := ConfigurationVariableParams{Scope: scope, ScopeId: scopeId}
	if err := readResourceData(&params, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	variable, err := getConfigurationVariable(params, meta)
	if err != nil {
		return err
	}

	if err := writeResourceData(&variable, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	d.Set("enum", variable.Schema.Enum)

	if variable.Schema.Format != client.Text {
		d.Set("format", string(variable.Schema.Format))
	}

	return nil
}

func getScopeAndId(d *schema.ResourceData) (client.Scope, string) {
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

	if deploymentLogId, ok := d.GetOk("deployment_log_id"); ok {
		scope = client.ScopeDeploymentLog
		scopeId = deploymentLogId.(string)
	}

	return scope, scopeId
}

func getConfigurationVariable(params ConfigurationVariableParams, meta interface{}) (client.ConfigurationVariable, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	id, idOk := params.Id, params.Id != ""

	if idOk {
		variable, err := apiClient.ConfigurationVariablesById(id)
		if err != nil {
			return client.ConfigurationVariable{}, diag.Errorf("Could not query variable: %v", err)
		}
		return variable, nil
	}

	variables, err := apiClient.ConfigurationVariablesByScope(params.Scope, params.ScopeId)

	if err != nil {
		return client.ConfigurationVariable{}, diag.Errorf("Could not query variables: %v", err)
	}

	name, nameOk := params.Name, params.Name != ""
	typeString, ok := params.ConfigurationType, params.ConfigurationType != ""
	type_ := -1
	if ok {
		if !nameOk {
			return client.ConfigurationVariable{}, diag.Errorf("Specify 'type' only when searching configuration variables by 'name' (not by 'id')")
		}
		switch typeString {
		case "environment":
			type_ = int(client.ConfigurationVariableTypeEnvironment)
		case "terraform":
			type_ = int(client.ConfigurationVariableTypeTerraform)
		default:
			return client.ConfigurationVariable{}, diag.Errorf("Invalid value for 'type': %s. can be either 'environment' or 'terraform'", typeString)
		}
	}
	var variable client.ConfigurationVariable
	for _, candidate := range variables {
		if idOk && candidate.Id == id {
			variable = candidate
			break
		}
		if nameOk && candidate.Name == name {
			if type_ != -1 {
				if int(*candidate.Type) != type_ {
					continue
				}
			}
			variable = candidate
			break
		}
	}
	if variable.Id == "" {
		return client.ConfigurationVariable{}, diag.Errorf("Could not find variable")
	}
	return variable, nil
}
