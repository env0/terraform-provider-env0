package env0

import (
	"context"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTeamImport},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "the environment's id",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The environment's name",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the environment's project id",
				Required:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the environment's template id",
				Required:    true,
			},
		},
	}
}

//func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	apiClient := meta.(client.ApiClientInterface)
//	payload := client.TeamCreatePayload{
//		Name: d.Get("name").(string),
//	}
//	if description, ok := d.GetOk("description"); ok {
//		payload.Description = description.(string)
//	}
//
//	team, err := apiClient.TeamCreate(payload)
//	if err != nil {
//		return diag.Errorf("could not create team: %v", err)
//	}
//
//	d.SetId(team.Id)
//
//	return nil
//}

func setEnvironmentSchema(d *schema.ResourceData, environment client.Environment) {
	d.Set("name", environment.Name)
	d.Set("project_id", environment.ProjectId)
	d.Set("template_id", environment.TemplateId)
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request := client.EnvironmentCreatePayload{
		Name:      d.Get("name").(string),
		ProjectId: d.Get("project_id").(string),
		DeployRequest: client.DeployRequest{
			BlueprintId: d.Get("template_id").(string),
		},
	}

	environment, err := apiClient.EnvironmentCreate(request)
	if err != nil {
		return diag.Errorf("could not create environment: %v", err)
	}

	d.SetId(environment.Id)
	setEnvironmentSchema(d, environment)

	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environment, err := apiClient.Environment(d.Id())
	if err != nil {
		return diag.Errorf("could not get environment: %v", err)
	}

	setEnvironmentSchema(d, environment)

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.EnvironmentUpdatePayload{
		Name: d.Get("name").(string),
	}

	_, err := apiClient.EnvironmentUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update environment: %v", err)
	}

	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//apiClient := meta.(client.ApiClientInterface)

	//err := apiClient.EnvironmentDelete(d.Id())
	//if err != nil {
	//	return diag.Errorf("could not delete team: %v", err)
	//}
	return nil
}

//
//func resourceTeamImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
//	id := d.Id()
//	var getErr diag.Diagnostics
//	_, uuidErr := uuid.Parse(id)
//	if uuidErr == nil {
//		log.Println("[INFO] Resolving team by id: ", id)
//		_, getErr = getTeamById(id, meta)
//	} else {
//		log.Println("[DEBUG] ID is not a valid env0 id ", id)
//		log.Println("[INFO] Resolving team by name: ", id)
//		var team client.Team
//		team, getErr = getTeamByName(id, meta)
//		d.SetId(team.Id)
//	}
//	if getErr != nil {
//		return nil, errors.New(getErr[0].Summary)
//	} else {
//		return []*schema.ResourceData{d}, nil
//	}
//}
//
//func getTeamByName(name interface{}, meta interface{}) (client.Team, diag.Diagnostics) {
//	apiClient := meta.(client.ApiClientInterface)
//	teams, err := apiClient.Teams()
//	if err != nil {
//		return client.Team{}, diag.Errorf("Could not get teams: %v", err)
//	}
//
//	var teamsByName []client.Team
//	for _, candidate := range teams {
//		if candidate.Name == name {
//			teamsByName = append(teamsByName, candidate)
//		}
//	}
//
//	if len(teamsByName) > 1 {
//		return client.Team{}, diag.Errorf("Found multiple teams for name: %s. Use ID instead or make sure team names are unique %v", name, teamsByName)
//	}
//
//	if len(teamsByName) == 0 {
//		return client.Team{}, diag.Errorf("Could not find an env0 team with name %s", name)
//	}
//
//	return teamsByName[0], nil
//}
//
//func getTeamById(id interface{}, meta interface{}) (client.Team, diag.Diagnostics) {
//	apiClient := meta.(client.ApiClientInterface)
//	team, err := apiClient.Team(id.(string))
//	if err != nil {
//		return client.Team{}, diag.Errorf("Could not get team: %v", err)
//	}
//	return team, nil
//}
