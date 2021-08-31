package env0

import (
	"context"
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
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
	payload := client.TeamCreatePayload{
		Name: d.Get("name").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		payload.Description = description.(string)
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
		return diag.Errorf("could not get team: %v", err)
	}

	d.Set("name", team.Name)
	d.Set("description", team.Description)

	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.TeamUpdatePayload{
		Name: d.Get("name").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		payload.Description = description.(string)
	}

	_, err := apiClient.TeamUpdate(d.Id(), payload)
	if err != nil {
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
		log.Println("[INFO] Resolving Team by id: ", id)
		_, getErr = getTeamById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving Team by name: ", id)
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

func getTeamByName(name interface{}, meta interface{}) (client.Team, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	teams, err := apiClient.Teams()
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
		return client.Team{}, diag.Errorf("Found multiple Teams for name: %s. Use ID instead or make sure Team names are unique %v", name, teamsByName)
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
