package env0

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyReset,

		Importer: &schema.ResourceImporter{StateContext: resourcePolicyImport},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id of the policy",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id  of the project",
				Required:    true,
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					if strings.TrimSpace(i.(string)) == "" {
						return diag.Errorf("project id must not be empty")
					}
					return nil
				},
			},
			"number_of_environments": {
				Type:        schema.TypeInt,
				Description: "number of environments per project, 0 indicates no limit",
				Optional:    true,
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					n := i.(int)
					if n < 1 {
						return diag.Errorf("Number of environments must be greater than zero")
					}
					return nil
				},
			},
			"number_of_environments_total": {
				Type:        schema.TypeInt,
				Description: "number of environments total, 0 indicates no limit",
				Optional:    true,
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					n := i.(int)
					if n < 1 {
						return diag.Errorf("Number of total environments must be greater that zero")
					}
					return nil
				},
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
				Description: "Disallow destroying environment in the project",
				Optional:    true,
			},
			"skip_redundant_deployments": {
				Type:        schema.TypeBool,
				Description: "skip redundant deployments",
				Optional:    true,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Description: "updated by",
				Computed:    true,
			},
		},
	}
}

func setPolicySchema(d *schema.ResourceData, policy client.Policy) {
	d.Set("id", policy.Id)
	d.Set("project_id", policy.ProjectId)
	d.Set("number_of_environments", policy.NumberOfEnvironments)
	d.Set("number_of_environments_total", policy.NumberOfEnvironmentsTotal)
	d.Set("requires_approval_default", policy.RequiresApprovalDefault)
	d.Set("include_cost_estimation", policy.IncludeCostEstimation)
	d.Set("skip_apply_when_plan_is_empty", policy.SkipApplyWhenPlanIsEmpty)
	d.Set("disable_destroy_environments", policy.DisableDestroyEnvironments)
	d.Set("skip_redundant_deployments", policy.SkipRedundantDepolyments)
	d.Set("updated_by", policy.UpdatedBy)
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	d.SetId(projectId)
	return resourcePolicyUpdate(ctx, d, meta)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()

	policy, err := apiClient.Policy(projectId)
	if err != nil {
		return diag.Errorf("could not get policy: %v", err)
	}

	d.SetId(projectId)
	setPolicySchema(d, policy)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()
	d.SetId(projectId)
	log.Printf("[INFO] updating policy for project: %s\n", projectId)

	payload := client.PolicyUpdatePayload{}
	if projectId, ok := d.GetOk("project_id"); ok {
		payload.ProjectId = projectId.(string)
	}
	if numberOfEnvironments, ok := d.GetOk("number_of_environments"); ok {
		payload.NumberOfEnvironments = numberOfEnvironments.(int)
		// return diag.Errorf("number of environments: %d", payload.NumberOfEnvironments)
	}
	if numberOfEnvironmentsTotal, ok := d.GetOk("number_of_environments_total"); ok {
		payload.NumberOfEnvironmentsTotal = numberOfEnvironmentsTotal.(int)
		// return diag.Errorf("number of environments total: %d", payload.NumberOfEnvironmentsTotal)
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
	if skipRedundantDeployments, ok := d.GetOk("skip_redundant_depolyments"); ok {
		payload.SkipRedundantDepolyments = skipRedundantDeployments.(bool)
	}

	_, err := apiClient.PolicyUpdate(payload)
	if err != nil {
		return diag.Errorf("could not update policy: %v", err)
	}
	return nil
}

func resourcePolicyReset(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()

	payload := client.PolicyUpdatePayload{
		ProjectId:               projectId,
		RequiresApprovalDefault: true,
	}

	policy, err := apiClient.PolicyUpdate(payload)
	if err != nil {
		return diag.Errorf("could not reset policy: %v", err)
	}

	d.SetId(projectId)
	setPolicySchema(d, policy)

	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectId := d.Id()

	policy, err := getPolicyByProjectId(projectId, meta)
	if err != nil {
		return nil, errors.New(err[0].Summary)
	}

	d.SetId(projectId)
	setPolicySchema(d, policy)

	return []*schema.ResourceData{d}, nil
}
