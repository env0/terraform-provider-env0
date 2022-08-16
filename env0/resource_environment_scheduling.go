package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironmentScheduling() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentSchedulingCreateOrUpdate,
		ReadContext:   resourceEnvironmentSchedulingRead,
		UpdateContext: resourceEnvironmentSchedulingCreateOrUpdate,
		DeleteContext: resourceEnvironmentSchedulingDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "The environment's id",
				Required:    true,
				ForceNew:    true,
			},
			"destroy_cron": {
				Type:             schema.TypeString,
				Description:      "Cron expression for scheduled destroy of the environment. Destroy and Deploy cron expressions must not be the same.",
				AtLeastOneOf:     []string{"destroy_cron", "deploy_cron"},
				Optional:         true,
				ValidateDiagFunc: ValidateCronExpression,
			},
			"deploy_cron": {
				Type:             schema.TypeString,
				Description:      "Cron expression for scheduled deploy of the environment. Destroy and Deploy cron expressions must not be the same.",
				AtLeastOneOf:     []string{"destroy_cron", "deploy_cron"},
				Optional:         true,
				ValidateDiagFunc: ValidateCronExpression,
			},
		},
	}
}

func resourceEnvironmentSchedulingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Id()

	environmentScheduling, err := apiClient.EnvironmentScheduling(environmentId)

	if err != nil {
		return diag.Errorf("could not get environment scheduling: %v", err)
	}

	if err := writeResourceData(&environmentScheduling, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceEnvironmentSchedulingCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)

	var payload client.EnvironmentScheduling

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.EnvironmentSchedulingUpdate(environmentId, payload); err != nil {
		return diag.Errorf("could not create or update environment scheduling: %v", err)
	}

	d.SetId(environmentId)

	return nil
}

func resourceEnvironmentSchedulingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.EnvironmentSchedulingDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete environment scheduling: %v", err)
	}

	return nil
}
