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

func setEnvironmentSchema(d *schema.ResourceData, environment client.Environment) {
	d.Set("name", environment.Name)
	d.Set("project_id", environment.ProjectId)
	d.Set("template_id", environment.TemplateId)
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := getCreatePayload(d)

	environment, err := apiClient.EnvironmentCreate(payload)
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

	payload := getUpdatePayload(d)

	// TODO: deploy if needed?

	_, err := apiClient.EnvironmentUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update environment: %v", err)
	}

	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	_, err := apiClient.EnvironmentDestroy(d.Id())
	if err != nil {
		return diag.Errorf("could not delete team: %v", err)
	}
	return nil
}

func getCreatePayload(d *schema.ResourceData) client.EnvironmentCreate {
	payload := client.EnvironmentCreate{}

	if name, ok := d.GetOk("name"); ok {
		payload.Name = name.(string)
	}

	if projectId, ok := d.GetOk("project_id"); ok {
		payload.ProjectId = projectId.(string)
	}

	payload.DeployRequest = &client.DeployRequest{}

	if blueprintId, ok := d.GetOk("template_id"); ok {
		payload.DeployRequest.BlueprintId = blueprintId.(string)
	}

	if blueprintRepository, ok := d.GetOk("repository"); ok {
		payload.DeployRequest.BlueprintRepository = blueprintRepository.(string)
	}

	if blueprintRevision, ok := d.GetOk("revision"); ok {
		payload.DeployRequest.BlueprintRepository = blueprintRevision.(string)
	}

	return payload
}

func getUpdatePayload(d *schema.ResourceData) client.EnvironmentUpdate {
	// TODO: check if not filling them make these null or false
	payload := client.EnvironmentUpdate{}

	if name, ok := d.GetOk("name"); ok {
		payload.Name = name.(string)
	}
	if requiresApproval, ok := d.GetOk("requires_approval"); ok {
		payload.RequiresApproval = requiresApproval.(bool)
	}
	if isArchived, ok := d.GetOk("is_archived"); ok {
		payload.IsArchived = isArchived.(bool)
	}
	if continuousDeployment, ok := d.GetOk("redeploy_on_push"); ok {
		payload.ContinuousDeployment = continuousDeployment.(bool)
	}
	if pullRequestPlanDeployments, ok := d.GetOk("pr_plan_on_pull_request"); ok {
		payload.PullRequestPlanDeployments = pullRequestPlanDeployments.(bool)
	}
	if autoDeployOnPathChangesOnly, ok := d.GetOk("auto_deploy_on_path_change_only"); ok {
		payload.AutoDeployOnPathChangesOnly = autoDeployOnPathChangesOnly.(bool)
	}
	if autoDeployByCustomGlob, ok := d.GetOk("auto_deploy_by_custom_glob"); ok {
		payload.AutoDeployByCustomGlob = autoDeployByCustomGlob.(string)
	}

	return payload
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
