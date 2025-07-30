package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPolicyRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id of the policy",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
			},
			"number_of_environments": {
				Type:        schema.TypeInt,
				Description: "Max number of environments a single user can have in this project, 0 indicates no limit",
				Computed:    true,
			},
			"number_of_environments_total": {
				Type:        schema.TypeInt,
				Description: "Max number of environments in this project, 0 indicates no limit",
				Computed:    true,
			},
			"requires_approval_default": {
				Type:        schema.TypeBool,
				Description: "Requires approval default value when creating a new environment in the project",
				Computed:    true,
			},
			"include_cost_estimation": {
				Type:        schema.TypeBool,
				Description: "Enable cost estimation for the project",
				Computed:    true,
			},
			"skip_apply_when_plan_is_empty": {
				Type:        schema.TypeBool,
				Description: "Skip apply when plan has no changes",
				Computed:    true,
			},
			"disable_destroy_environments": {
				Type:        schema.TypeBool,
				Description: "Disallow destroying environment in the project",
				Computed:    true,
			},
			"skip_redundant_deployments": {
				Type:        schema.TypeBool,
				Description: "Automatically skip queued deployments when a newer deployment is triggered",
				Computed:    true,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Description: "updated by",
				Computed:    true,
			},
			"run_pull_request_plan_default": {
				Type:        schema.TypeBool,
				Description: "Run Terraform Plan on Pull Requests for new environments targeting their branch default value",
				Computed:    true,
			},
			"continuous_deployment_default": {
				Type:        schema.TypeBool,
				Description: "Redeploy on every push to the git branch default value",
				Computed:    true,
			},
			"max_ttl": {
				Type:        schema.TypeString,
				Description: "the maximum environment time-to-live allowed on deploy time",
				Computed:    true,
			},
			"default_ttl": {
				Type:        schema.TypeString,
				Description: "the default environment time-to-live allowed on deploy time",
				Computed:    true,
			},
			"force_remote_backend": {
				Type:        schema.TypeBool,
				Description: "force env0 remote backend",
				Computed:    true,
			},
			"auto_drift_remediation": {
				Type:        schema.TypeString,
				Description: "Auto drift remediation strategy (DISABLED, CODE_TO_CLOUD, CLOUD_TO_CODE, SMART_REMEDIATION)",
				Computed:    true,
			},
		},
	}
}

func dataPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err diag.Diagnostics

	var policy client.Policy

	projectId, ok := d.GetOk("project_id")
	if ok {
		policy, err = getPolicyByProjectId(projectId.(string), meta)
		if err != nil {
			return err
		}
	}

	if err := writeResourceData(&policy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	d.SetId(policy.Id)

	return nil
}

func getPolicyByProjectId(projectId string, meta any) (client.Policy, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)

	policy, err := apiClient.Policy(projectId)
	if err != nil {
		return client.Policy{}, diag.Errorf("Could not query policy: %v", err)
	}

	return policy, nil
}
