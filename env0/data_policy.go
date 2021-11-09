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
				Description: "number of environments per project",
				Computed:    true,
			},
			"number_of_environments_total": {
				Type:        schema.TypeInt,
				Description: "number of environments total",
				Computed:    true,
			},
			"requires_approval_default": {
				Type:        schema.TypeBool,
				Description: "requires approval",
				Computed:    true,
			},
			"include_cost_estimation": {
				Type:        schema.TypeBool,
				Description: "include cost estimation",
				Computed:    true,
			},
			"skip_apply_when_plan_is_empty": {
				Type:        schema.TypeBool,
				Description: "skip apply when plan is empty",
				Computed:    true,
			},
			"disable_destroy_environments": {
				Type:        schema.TypeBool,
				Description: "disable destroy environments",
				Computed:    true,
			},
			"skip_redundant_deployments": {
				Type:        schema.TypeBool,
				Description: "skip redundant deployments",
				Computed:    true,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Description: "updated by",
				Computed:    true,
			},
		},
	}
}

func dataPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var policy client.Policy

	projectId, ok := d.GetOk("project_id")
	if ok {
		policy, err = getPolicyByProjectId(projectId.(string), meta)
		if err != nil {
			return err
		}
	}
	d.SetId(policy.Id)
	setPolicySchema(d, policy)
	return nil
}

func getPolicyByProjectId(projectId string, meta interface{}) (client.Policy, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	policy, err := apiClient.Policy(projectId)
	if err != nil {
		return client.Policy{}, diag.Errorf("Could not query policy: %v", err)
	}
	return policy, nil
}
