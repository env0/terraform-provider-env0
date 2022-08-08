package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkflowTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataWorkflowTriggersRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the source environment",
				Required:    true,
			},
			"downstream_environment_ids": {
				Type:        schema.TypeList,
				Description: "environments to trigger",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "id of the downstream environments to trigger",
				},
			},
		},
	}
}

func dataWorkflowTriggersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentId := d.Get("environment_id").(string)

	triggers, err := apiClient.WorkflowTrigger(environmentId)

	if err != nil {
		return diag.Errorf("could not get workflow triggers: %v", err)
	}

	d.SetId(environmentId)
	var triggerIds []string
	for _, value := range triggers {
		triggerIds = append(triggerIds, value.Id)
	}

	d.Set(`downstream_environment_ids`, triggerIds)

	return nil
}
