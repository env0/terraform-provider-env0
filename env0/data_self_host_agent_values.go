package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSelfHostedAgentValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSelfHostedAgentValuesRead,

		Schema: map[string]*schema.Schema{
			"agent_key": {
				Type:        schema.TypeString,
				Description: "the agent key",
				Required:    true,
			},
			"values": {
				Type:        schema.TypeString,
				Description: "Self hosted agent helm values. The values can be passed to a helm release resource (https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release)",
				Computed:    true,
			},
		},
	}
}

func dataSelfHostedAgentValuesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agentKey := d.Get("agent_key").(string)

	values, err := apiClient.AgentValues(agentKey)
	if err != nil {
		return diag.Errorf("could not get agent values: %v", err)
	}

	d.Set("values", values)
	d.SetId(agentKey)

	return nil
}
