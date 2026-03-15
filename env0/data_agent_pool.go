package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAgentPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAgentPoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the agent pool",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the agent pool",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "the agent pool description",
				Computed:    true,
			},
			"agent_key": {
				Type:        schema.TypeString,
				Description: "the agent key identifier. Deprecated: use id instead",
				Computed:    true,
				Deprecated:  "use id instead",
			},
			"logs": {
				Type:        schema.TypeList,
				Description: "self-hosted logs configuration (dynamo self-hosted)",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Description: "the AWS account id for self-hosted logs",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "the AWS region for self-hosted logs",
							Computed:    true,
						},
						"external_id": {
							Type:        schema.TypeString,
							Description: "the external id for assuming the role",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataAgentPoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var agentPool *client.AgentPool
	var err error

	id, ok := d.GetOk("id")
	if ok {
		agentPool, err = apiClient.AgentPool(id.(string))
		if err != nil {
			return diag.Errorf("could not read agent pool: %v", err)
		}
	} else {
		agentPool, err = getAgentPoolByName(d.Get("name").(string), meta)
		if err != nil {
			return diag.Errorf("could not read agent pool: %v", err)
		}
	}

	d.SetId(agentPool.Id)
	setAgentPoolSchema(d, agentPool)

	return nil
}
