package env0

import (
	"context"
	"regexp"

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
			"filter": {
				Type:        schema.TypeString,
				Description: "add an optional regex filter, only names matching the filter will be added to 'names'. Note: The regex filter is Golang flavor",
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func dataTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var regex *regexp.Regexp

	var err error

	filter := d.Get("filter").(string)
	if filter != "" {
		regex, err = regexp.Compile(filter)
		if err != nil {
			return diag.Errorf("Invalid filter: %v", err)
		}
	}

	teams, err := apiClient.Teams()
	if err != nil {
		return diag.Errorf("Could not get teams: %v", err)
	}

	data := []string{}

	for _, team := range teams {
		if regex == nil || regex.MatchString(team.Name) {
			data = append(data, team.Name)
		}
	}

	d.Set("names", data)

	d.SetId("all_teams_names" + filter)

	return nil
}
