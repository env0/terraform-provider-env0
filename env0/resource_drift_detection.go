package env0

import (
	"context"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDriftDetection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentDriftCreateOrUpdate,
		ReadContext:   resourceEnvironmentDriftRead,
		UpdateContext: resourceEnvironmentDriftCreateOrUpdate,
		DeleteContext: resourceEnvironmentDriftDelete,

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
		},
	}
}

func resourceEnvironmentDriftRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Id()

	drift, err := apiClient.EnvironmentDriftDetection(environmentId)

	if err != nil {
		return diag.Errorf("could not get environment drift detection: %v", err)
	}

	if drift.Enabled {
		d.Set("cron", drift.Cron)
	} else {
		d.SetId("")
	}
	return nil
}

func resourceEnvironmentDriftCreateOrUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	cron := d.Get("cron").(string)

	payload := EnvironmentSchedulingExpression{Cron: cron, Enabled: true}

	_, err := apiClient.EnvironmentUpdateDriftDetection(environmentId, payload)

	if err != nil {
		return diag.Errorf("could not create or update environment drift detection: %v", err)
	}

	d.SetId(environmentId)
	return nil
}

func resourceEnvironmentDriftDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Id()

	err := apiClient.EnvironmentStopDriftDetection(environmentId)

	if err != nil {
		return diag.Errorf("could not stop environment drift detection: %v", err)
	}

	return nil
}
