package env0

import (
	"context"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowTriggers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowTriggersCreateOrUpdate,
		ReadContext:   resourceWorkflowTriggersRead,
		UpdateContext: resourceWorkflowTriggersCreateOrUpdate,
		DeleteContext: resourceWorkflowTriggersDelete,

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
			},
		},
	}
}

func resourceWorkflowTriggersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Get("environment_id").(string)

	triggers, err := apiClient.WorkflowTrigger(environmentId)

	if err != nil {
		return diag.Errorf("could not get workflow triggers: %v", err)
	}

	d.Set("downstream_environment_ids", triggers)
	return nil
}

func resourceWorkflowTriggersCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)
	environmentId := d.Get("environment_id").(string)
	request := WorkflowTriggerCreatePayload{
		DownstreamEnvironmentIds: d.Get("downstream_environment_ids"),
	}
	triggers, err := apiClient.WorkflowTriggerCreate(environmentId, request)
	if err != nil {
		return diag.Errorf("could not Create or Update TeamProjectAssignment: %v", err)
	}

	d.SetId(environmentId)
	var downstreamIds []string
	for _, trigger := range triggers {

	}

	d.Set("downstream_environment_ids")
	return nil
}

func resourceWorkflowTriggersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	_, err := apiClient.WorkflowTriggerCreate(d.Id(), WorkflowTriggerCreatePayload{})
	if err != nil {
		return diag.Errorf("could not remove workflow triggers: %v", err)
	}

	return nil
}
