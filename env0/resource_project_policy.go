package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/client"
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
				Type:             schema.TypeString,
				Description:      "id of the project",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
			"number_of_environments": {
				Type:             schema.TypeInt,
				Description:      "Max number of environments a single user can have in this project, `null` indicates no limit",
				Optional:         true,
				ValidateDiagFunc: NewGreaterThanValidator(0),
			},
			"number_of_environments_total": {
				Type:             schema.TypeInt,
				Description:      "Max number of environments in this project, `null` indicates no limit",
				Optional:         true,
				ValidateDiagFunc: NewGreaterThanValidator(0),
			},
			"requires_approval_default": {
				Type:        schema.TypeBool,
				Description: "Requires approval default value when creating a new environment in the project",
				Default:     true,
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

func setPolicySchema(d *schema.ResourceData, policy client.Policy) error {
	return writeResourceData(&policy, d)
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

	if err := setPolicySchema(d, policy); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	d.SetId(projectId)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()
	d.SetId(projectId)

	payload := client.PolicyUpdatePayload{}

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
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
