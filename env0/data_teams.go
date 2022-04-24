package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Description: "list of all teams (by name)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the team name",
				},
			},
		},
	}
}

func dataTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	teams, err := apiClient.Teams()
	if err != nil {
		return diag.Errorf("Could not get teams: %v", err)
	}

	data := []string{}

	for _, team := range teams {
		data = append(data, team.Name)
	}

	d.Set("names", data)

	// Not really needed. But required by Terraform SDK - https://github.com/hashicorp/terraform-plugin-sdk/issues/541
	d.SetId("1")

	return nil
}
