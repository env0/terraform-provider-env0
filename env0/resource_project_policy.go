package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,

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
				ForceNew:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
			"number_of_environments": {
				Type:             schema.TypeInt,
				Description:      "Max number of environments a single user can have in this project.\nOmitting removes the restriction.",
				Optional:         true,
				ValidateDiagFunc: NewGreaterThanValidator(0),
			},
			"number_of_environments_total": {
				Type:             schema.TypeInt,
				Description:      "Max number of environments in this project.\nOmitting removes the restriction.",
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
			"max_ttl": {
				Type:             schema.TypeString,
				Description:      "the maximum environment time-to-live allowed on deploy time. Format is <number>-<M/w/d/h> (Examples: 12-h, 3-d, 1-w, 1-M). Default value is 'inherit' which inherits the organization policy. must be equal or longer than default_ttl",
				Optional:         true,
				Default:          "inherit",
				ValidateDiagFunc: ValidateTtl,
			},
			"default_ttl": {
				Type:             schema.TypeString,
				Description:      "the default environment time-to-live allowed on deploy time. Format is <number>-<M/w/d/h> (Examples: 12-h, 3-d, 1-w, 1-M). Default value is 'inherit' which inherits the organization policy. must be equal or shorter than max_ttl",
				Optional:         true,
				Default:          "inherit",
				ValidateDiagFunc: ValidateTtl,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := resourcePolicyUpdate(ctx, d, meta); err != nil {
		return err
	}

	projectId := d.Get("project_id").(string)
	d.SetId(projectId)

	return nil
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()

	policy, err := apiClient.Policy(projectId)
	if err != nil {
		return diag.Errorf("could not get policy: %v", err)
	}

	if err := writeResourceData(&policy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	d.SetId(projectId)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.PolicyUpdatePayload{}

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	// Validate if one is "inherit", the other must be too.
	if (payload.MaxTtl == "inherit" || payload.DefaultTtl == "inherit") && payload.MaxTtl != payload.DefaultTtl {
		return diag.Errorf("max_ttl and default_ttl must both inherit organization settings or override them")
	}

	if err := validateTtl(&payload.DefaultTtl, &payload.MaxTtl); err != nil {
		return diag.FromErr(err)
	}

	if payload.DefaultTtl == "Infinite" {
		payload.DefaultTtl = ""
	}

	if payload.MaxTtl == "Infinite" {
		payload.MaxTtl = ""
	}

	if _, err := apiClient.PolicyUpdate(payload); err != nil {
		return diag.Errorf("could not update policy: %v", err)
	}

	return nil
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err := writeResourceData(&policy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectId := d.Id()

	policy, err := getPolicyByProjectId(projectId, meta)
	if err != nil {
		return nil, errors.New(err[0].Summary)
	}

	d.SetId(projectId)
	if err := writeResourceData(&policy, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
