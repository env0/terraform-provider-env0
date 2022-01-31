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
				Description: "Max number of environments a single user can have in this project, 0 indicates no limit",
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
				Description: "Max number of environments in this project, 0 indicates no limit",
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
				Description: "Requires approval default value when creating a new environment in the project",
				Optional:    true,
			},
			"include_cost_estimation": {
				Type:        schema.TypeBool,
				Description: "Enable cost estimation for the project",
				Optional:    true,
			},
			"skip_apply_when_plan_is_empty": {
				Type:        schema.TypeBool,
				Description: "Skip apply when plan has no changes",
				Optional:    true,
			},
			"disable_destroy_environments": {
				Type:        schema.TypeBool,
				Description: "Disallow destroying environment in the project",
				Optional:    true,
			},
			"skip_redundant_deployments": {
				Type:        schema.TypeBool,
				Description: "Automatically skip queued deployments when a newer deployment is triggered",
				Optional:    true,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Description: "updated by",
				Computed:    true,
			},
			"run_pull_request_plan_default": {
				Type:        schema.TypeBool,
				Description: "Run Terraform Plan on Pull Requests for new environments targeting their branch default value",
				Optional:    true,
			},
			"continuous_deployment_default": {
				Type:        schema.TypeBool,
				Description: "Redeploy on every push to the git branch default value",
				Optional:    true,
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
	d.Set("skip_redundant_deployments", policy.SkipRedundantDeployments)
	d.Set("updated_by", policy.UpdatedBy)
	d.Set("run_pull_request_plan_default", policy.RunPullRequestPlanDefault)
	d.Set("continuous_deployment_default", policy.ContinuousDeploymentDefault)
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
		if !d.IsNewResource() {
			d.SetId("")
			return nil
		}
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
	if skipRedundantDeployments, ok := d.GetOk("skip_redundant_deployments"); ok {
		payload.SkipRedundantDeployments = skipRedundantDeployments.(bool)
	}
	if runPullRequestPlanDefault, ok := d.GetOk("run_pull_request_plan_default"); ok {
		payload.RunPullRequestPlanDefault = runPullRequestPlanDefault.(bool)
	}
	if continuousDeploymentDefault, ok := d.GetOk("continuous_deployment_default"); ok {
		payload.ContinuousDeploymentDefault = continuousDeploymentDefault.(bool)
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
