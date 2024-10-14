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
				Description:  "The name of the team",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "ID of the team",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Textual description of the team",
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

	if err := writeResourceData(&team, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
