package env0

import (
	"context"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the team",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the team",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "textual description of the team",
				Computed:    true,
			},
		},
	}
}

func dataTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var team client.Team

	id, ok := d.GetOk("id")
	if ok {
		team, err = getTeamById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("Either 'name' or 'id' must be specified")
		}
		team, err = getTeamByName(name.(string), meta)
		if err != nil {
			return err
		}
	}

	d.SetId(team.Id)
	setTeamSchema(d, team)

	return nil
}
