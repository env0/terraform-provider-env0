package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowTriggerCreate,
		ReadContext:   resourceWorkflowTriggerRead,
		DeleteContext: resourceWorkflowTriggerDelete,
		Description:   "cannot be used with env0_workflow_triggers",

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the source environment",
				Required:    true,
				ForceNew:    true,
			},
			"downstream_environment_id": {
				Type:        schema.TypeString,
				Description: "environment to trigger",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceWorkflowTriggerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	downstreamEnvironmentId := d.Get("downstream_environment_id").(string)

	triggers, err := apiClient.WorkflowTrigger(environmentId)
	if err != nil {
		return diag.Errorf("could not get workflow triggers: %v", err)
	}

	for _, trigger := range triggers {
		if trigger.Id == downstreamEnvironmentId {
			return nil
		}
	}

	tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
	d.SetId("")

	return nil
}

func resourceWorkflowTriggerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	downstreamEnvironmentId := d.Get("downstream_environment_id").(string)

	payload := client.WorkflowTriggerEnvironments{DownstreamEnvironmentIds: []string{downstreamEnvironmentId}}

	if err := apiClient.SubscribeWorkflowTrigger(environmentId, payload); err != nil {
		return diag.Errorf("failed to subscribe a workflow trigger: %v", err)
	}

	d.SetId(environmentId + "_" + downstreamEnvironmentId)

	return nil
}

func resourceWorkflowTriggerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	downstreamEnvironmentId := d.Get("downstream_environment_id").(string)

	payload := client.WorkflowTriggerEnvironments{DownstreamEnvironmentIds: []string{downstreamEnvironmentId}}

	if err := apiClient.UnsubscribeWorkflowTrigger(environmentId, payload); err != nil {
		return diag.Errorf("failed to unsubscribe a workflow trigger: %v", err)
	}

	return nil
}
