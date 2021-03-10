package env0tfprovider

import (
	"context"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
			"type": {
				Type:        schema.TypeString,
				Description: "'terraform' or 'environment'",
				Computed:    true,
			},
		},
	}
}

func dataConfigurationVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	scope := env0apiclient.ScopeGlobal
	scopeId := ""
	if projectId, ok := d.GetOk("project_id"); ok {
		scope = env0apiclient.ScopeProject
		scopeId = projectId.(string)
	}
	if templateId, ok := d.GetOk("template_id"); ok {
		scope = env0apiclient.ScopeTemplate
		scopeId = templateId.(string)
	}
	if environmentId, ok := d.GetOk("environment_id"); ok {
		scope = env0apiclient.ScopeEnvironment
		scopeId = environmentId.(string)
	}
	if deploymentLogId, ok := d.GetOk("deployment_log_id"); ok {
		scope = env0apiclient.ScopeDeploymentLog
		scopeId = deploymentLogId.(string)
	}
	variables, err := apiClient.ConfigurationVariables(scope, scopeId)
	if err != nil {
		return diag.Errorf("Could not query variables: %v", err)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")
	var variable env0apiclient.ConfigurationVariable
	for _, candidate := range variables {
		if idOk && candidate.Id == id.(string) {
			variable = candidate
			break
		}
		if nameOk && candidate.Name == name.(string) {
			variable = candidate
			break
		}
	}

	d.SetId(variable.Id)
	d.Set("name", variable.Name)
	d.Set("value", variable.Value)
	d.Set("is_sensitive", variable.IsSensitive)
	d.Set("scope", variable.Scope)
	if variable.Type == int64(env0apiclient.ConfigurationVariableTypeEnvironment) {
		d.Set("type", "environment")
	} else {
		d.Set("type", "terraform")
	}

	return nil
}
