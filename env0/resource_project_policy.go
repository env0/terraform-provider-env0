package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const INFINITE = "Infinite"
const INHERIT = "inherit"

func resourceProjectPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectPolicyCreate,
		ReadContext:   resourceProjectPolicyRead,
		UpdateContext: resourceProjectPolicyUpdate,
		DeleteContext: resourceProjectPolicyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceProjectPolicyImport},

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
				Default:          INHERIT,
				ValidateDiagFunc: ValidateTtl,
			},
			"default_ttl": {
				Type:             schema.TypeString,
				Description:      "the default environment time-to-live allowed on deploy time. Format is <number>-<M/w/d/h> (Examples: 12-h, 3-d, 1-w, 1-M). Default value is 'inherit' which inherits the organization policy. must be equal or shorter than max_ttl",
				Optional:         true,
				Default:          INHERIT,
				ValidateDiagFunc: ValidateTtl,
			},
			"force_remote_backend": {
				Type:        schema.TypeBool,
				Description: "if 'true' all environments created in this project will be forced to use env0 remote backend. Default is 'false'",
				Optional:    true,
				Default:     false,
			},
			"drift_detection_cron": {
				Type:             schema.TypeString,
				Description:      "default cron expression for new environments",
				Optional:         true,
				ValidateDiagFunc: ValidateCronExpression,
			},
			"auto_drift_remediation": {
				Type:        schema.TypeString,
				Description: "Auto drift remediation setting (DISABLED or CODE_TO_CLOUD). Defaults to DISABLED",
				Optional:    true,
				Default:     DriftRemediationDisabled,
				ValidateDiagFunc: NewStringInValidator([]string{
					DriftRemediationDisabled,
					DriftRemediationCodeToCloud,
				}),
			},
			"vcs_pr_comments_enabled_default": {
				Type:        schema.TypeBool,
				Description: "if 'true' all environments created in this project will be created with an 'enabled' running VCS PR plan/apply commands using PR comments. Default is 'false'",
				Optional:    true,
				Default:     false,
			},
			"outputs_as_inputs_enabled": {
				Type:        schema.TypeBool,
				Description: "if 'true' enables 'environment outputs'. Default is 'false'",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceProjectPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := resourceProjectPolicyUpdate(ctx, d, meta); err != nil {
		return err
	}

	projectId := d.Get("project_id").(string)
	d.SetId(projectId)

	return nil
}

func resourceProjectPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceProjectPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.PolicyUpdatePayload{}

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	// Validate if one is "inherit", the other must be too.
	if (payload.MaxTtl == INHERIT || payload.DefaultTtl == INHERIT) && payload.MaxTtl != payload.DefaultTtl {
		return diag.Errorf("max_ttl and default_ttl must both inherit organization settings or override them")
	}

	if err := validateTtl(&payload.DefaultTtl, &payload.MaxTtl); err != nil {
		return diag.FromErr(err)
	}

	if payload.DefaultTtl == INFINITE {
		payload.DefaultTtl = ""
	}

	if payload.MaxTtl == INFINITE {
		payload.MaxTtl = ""
	}

	if payload.DriftDetectionCron != "" {
		payload.DriftDetectionEnabled = true
	}

	if _, err := apiClient.PolicyUpdate(payload); err != nil {
		return diag.Errorf("could not update policy: %v", err)
	}

	return nil
}

func resourceProjectPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceProjectPolicyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	projectId := d.Id()

	policy, err := getPolicyByProjectId(projectId, meta)
	if err != nil {
		return nil, errors.New(err[0].Summary)
	}

	d.SetId(projectId)

	if err := writeResourceData(&policy, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
