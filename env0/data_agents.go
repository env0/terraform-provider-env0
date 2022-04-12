package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAgents() *schema.Resource {
	return &schema.Resource{
		ReadContext: agentsRead,

		Schema: map[string]*schema.Schema{
			"agents": {
				Type:        schema.TypeList,
				Description: "list of organization agents",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_key": {
							Type:        schema.TypeString,
							Description: "agent key identifier",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func agentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agents, err := apiClient.Agents()
	if err != nil {
		return diag.Errorf("could not read agents: %v", err)
	}

	data := struct {
		Agents []client.Agent
	}{agents}

	if err := writeResourceData(&data, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	// Not really needed. But required by Terraform SDK - https://github.com/hashicorp/terraform-plugin-sdk/issues/541
	d.SetId("1")

	return nil
}
