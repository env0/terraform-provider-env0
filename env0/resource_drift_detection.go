package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DriftRemediationDisabled    = "DISABLED"
	DriftRemediationCodeToCloud = "CODE_TO_CLOUD"
)

func resourceDriftDetection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentDriftCreateOrUpdate,
		ReadContext:   resourceEnvironmentDriftRead,
		UpdateContext: resourceEnvironmentDriftCreateOrUpdate,
		DeleteContext: resourceEnvironmentDriftDelete,

		Description: "note: instead of using this resource, setting drift detection can be configured directly through the environment resource",

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "The environment's id",
				Required:    true,
				ForceNew:    true,
			},
			"cron": {
				Type:             schema.TypeString,
				Description:      "Cron expression for scheduled drift detection of the environment",
				Required:         true,
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
		},
	}
}

func resourceEnvironmentDriftRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Id()

	drift, err := apiClient.EnvironmentDriftDetection(environmentId)

	if err != nil {
		return diag.Errorf("could not get environment drift detection: %v", err)
	}

	if !drift.Enabled {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")

		return nil
	}

	if err := writeResourceData(&drift, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceEnvironmentDriftCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	cron := d.Get("cron").(string)
	autoDriftRemediation := d.Get("auto_drift_remediation").(string)

	payload := client.EnvironmentSchedulingExpression{
		Cron:                 cron,
		Enabled:              true,
		AutoDriftRemediation: autoDriftRemediation,
	}

	if _, err := apiClient.EnvironmentUpdateDriftDetection(environmentId, payload); err != nil {
		return diag.Errorf("could not create or update environment drift detection: %v", err)
	}

	d.SetId(environmentId)

	return nil
}

func resourceEnvironmentDriftDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.EnvironmentStopDriftDetection(d.Id()); err != nil {
		return diag.Errorf("could not stop environment drift detection: %v", err)
	}

	return nil
}
