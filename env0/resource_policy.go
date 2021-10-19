package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyUpdate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyReset,

		Importer: nil,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "id  of the project",
				Required:    true,
				ForceNew:    true,
			},
			"number_of_environments": {
				Type:        schema.TypeInt,
				Description: "number of environments",
				Optional:    true,
			},
			"requires_approval_default": {
				Type:        schema.TypeBool,
				Description: "requires approval",
				Optional:    true,
			},
			"include_cost_estimation": {
				Type:        schema.TypeBool,
				Description: "include cost estimation",
				Optional:    true,
			},
			"skip_apply_when_plan_is_empty": {
				Type:        schema.TypeBool,
				Description: "skip apply when plan is empty",
				Optional:    true,
			},
			"disable_destroy_environments": {
				Type:        schema.TypeBool,
				Description: "disable destroy environments",
				Optional:    true,
			},
		},
	}
}

func setPolicySchema(d *schema.ResourceData, policy client.Policy) {
	d.Set("number_of_environments", policy.NumberOfEnvironments)
	d.Set("requires_approval_default", policy.RequiresApprovalDefault)
	d.Set("include_cost_estimation", policy.IncludeCostEstimation)
	d.Set("skip_apply_when_plan_is_empty", policy.SkipApplyWhenPlanIsEmpty)
	d.Set("disable_destroy_environments", policy.DisableDestroyEnvironments)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	policy, err := apiClient.Policy(d.Id())
	if err != nil {
		return diag.Errorf("could not get policy: %v", err)
	}

	setPolicySchema(d, policy)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.PolicyUpdatePayload{}
	if numberOfEnvironments, ok := d.GetOk("number_of_environments"); ok {
		payload.NumberOfEnvironments = numberOfEnvironments.(int)
	}
	if requiresApprovalDefault, ok := d.GetOk("requires_approval_default"); ok {
		payload.RequiresApprovalDefault = requiresApprovalDefault.(bool)
	}
	if includeCostEstimation, ok := d.GetOk("include_cost_estimation"); ok {
		payload.IncludeCostEstimation = includeCostEstimation.(bool)
	}
	if skipApplyWhenPlanIsEmpty, ok := d.GetOk("skip_apply_when_plan_is_empty"); ok {
		payload.SkipApplyWhenPlanIsEmpty = skipApplyWhenPlanIsEmpty.(bool)
	}
	if disableDestroyEnvironments, ok := d.GetOk("disable_destroy_environments"); ok {
		payload.DisableDestroyEnvironments = disableDestroyEnvironments.(bool)
	}

	_, err := apiClient.PolicyUpdate(payload)
	if err != nil {
		return diag.Errorf("could not create policy: %v", err)
	}

	return nil
}

func resourcePolicyReset(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.PolicyUpdatePayload{}

	_, err := apiClient.PolicyUpdate(payload)
	if err != nil {
		return diag.Errorf("could not delete policy: %v", err)
	}

	return nil
}
