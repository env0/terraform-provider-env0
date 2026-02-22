package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowTriggers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowTriggersCreateOrUpdate,
		ReadContext:   resourceWorkflowTriggersRead,
		UpdateContext: resourceWorkflowTriggersCreateOrUpdate,
		DeleteContext: resourceWorkflowTriggersDelete,
		Description:   "cannot be used with env0_workflow_trigger",

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the source environment",
				Required:    true,
			},
			"downstream_environment_ids": {
				Type:        schema.TypeList,
				Description: "environments to trigger",
				Required:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "id of the downstream environments to trigger",
				},
				MinItems: 1,
			},
		},
	}
}

func resourceWorkflowTriggersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)

	triggers, err := apiClient.WorkflowTrigger(environmentId)
	if err != nil {
		return diag.Errorf("could not get workflow triggers: %v", err)
	}

	triggerIds := []string{}
	for _, value := range triggers {
		triggerIds = append(triggerIds, value.Id)
	}

	d.Set("downstream_environment_ids", triggerIds)

	return nil
}

func resourceWorkflowTriggersCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	environmentId := d.Get("environment_id").(string)
	rawDownstreamIds := d.Get("downstream_environment_ids").([]any)

	requestDownstreamIds := []string{}

	for _, rawId := range rawDownstreamIds {
		requestDownstreamIds = append(requestDownstreamIds, rawId.(string))
	}

	request := client.WorkflowTriggerUpsertPayload{
		DownstreamEnvironmentIds: requestDownstreamIds,
	}

	triggers, err := apiClient.WorkflowTriggerUpsert(environmentId, request)
	if err != nil {
		return diag.Errorf("could not Create or Update workflow triggers: %v", err)
	}

	downstreamIds := []string{}
	for _, trigger := range triggers {
		downstreamIds = append(downstreamIds, trigger.Id)
	}

	d.SetId(environmentId)
	d.Set("downstream_environment_ids", downstreamIds)

	return nil
}

func resourceWorkflowTriggersDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	_, err := apiClient.WorkflowTriggerUpsert(d.Id(), client.WorkflowTriggerUpsertPayload{DownstreamEnvironmentIds: []string{}})
	if err != nil {
		return diag.Errorf("could not remove workflow triggers: %v", err)
	}

	return nil
}
