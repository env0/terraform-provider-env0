package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTeamImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the team",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description for the team",
				Optional:    true,
			},
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	team, err := apiClient.TeamCreate(payload)
	if err != nil {
		return diag.Errorf("could not create team: %v", err)
	}

	d.SetId(team.Id)

	return nil
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	team, err := apiClient.Team(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "team", d, err)
	}

	if err := writeResourceData(&team, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.TeamUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update team: %v", err)
	}

	return nil
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	err := apiClient.TeamDelete(d.Id())
	if err != nil {
		return diag.Errorf("could not delete team: %v", err)
	}

	return nil
}

func resourceTeamImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	var getErr diag.Diagnostics

	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		tflog.Info(ctx, "Resolving team by id", map[string]interface{}{"id": id})
		_, getErr = getTeamById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving team by name", map[string]interface{}{"name": id})

		var team client.Team

		team, getErr = getTeamByName(id, meta)
		d.SetId(team.Id)
	}

	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}

func getTeamByName(name string, meta interface{}) (client.Team, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)

	teams, err := apiClient.TeamsByName(name)
	if err != nil {
		return client.Team{}, diag.Errorf("Could not get teams: %v", err)
	}

	var teamsByName []client.Team

	for _, candidate := range teams {
		if candidate.Name == name {
			teamsByName = append(teamsByName, candidate)
		}
	}

	if len(teamsByName) > 1 {
		return client.Team{}, diag.Errorf("Found multiple teams for name: %s. Use ID instead or make sure team names are unique %v", name, teamsByName)
	}

	if len(teamsByName) == 0 {
		return client.Team{}, diag.Errorf("Could not find an env0 team with name %s", name)
	}

	return teamsByName[0], nil
}

func getTeamById(id interface{}, meta interface{}) (client.Team, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)

	team, err := apiClient.Team(id.(string))
	if err != nil {
		return client.Team{}, diag.Errorf("Could not get team: %v", err)
	}

	return team, nil
}
