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

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "the environment's id",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "the environment's name",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id of the environment",
				Required:    true,
				ForceNew:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the template id the environment is to be created from",
				Required:    true,
			},
			"workspace": {
				Type:        schema.TypeString,
				Description: "the terraform workspace of the environment",
				Optional:    true,
				ForceNew:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "the revision the environment is to be run against",
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "the repository the environment should use",
				Optional:    true,
			},
			"run_plan_on_pull_requests": {
				Type:        schema.TypeBool,
				Description: "should run terraform plan on pull requests creations",
				Optional:    true,
				Default:     false,
			},
			"approve_plan_automatically": {
				Type:        schema.TypeBool,
				Description: "should deployments require manual approvals ( defaults to true )",
				Optional:    true,
				Default:     true,
			},
			"deploy_on_push": {
				Type:        schema.TypeBool,
				Description: "should run terraform deploy on push events",
				Optional:    true,
				Default:     false,
			},
			"auto_deploy_by_custom_glob": {
				Type: schema.TypeBool,
				// TODO: description
				Description: "should deploy by custom glob",
				Optional:    true,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Description: "id of the last deployment",
				Computed:    true,
			},
			"configuration": {
				Type:        schema.TypeList,
				Description: "terraform and environment variables for the environment",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "variable name",
							Required:    true,
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Description: "variable value",
							Required:    true,
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Description: "variable type (allowed values are: terraform, environment)",
							Default:     "environment",
							Optional:    true,
						},
						"scope": &schema.Schema{
							Type:        schema.TypeString,
							Description: "variable scope ( allowed values are: GLOBAL, BLUEPRINT, PROJECT, ENVIRONMENT, DEPLOYMENT )",
							Optional:    true,
						},
						"scope_id": &schema.Schema{
							Type:        schema.TypeString,
							Description: "the scope's id",
							Optional:    true,
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Description: "description for the variable",
							Optional:    true,
						},
						"is_sensitive": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "should the variable value be hidden",
							Optional:    true,
						},
						"schema_type": &schema.Schema{
							Type:        schema.TypeString,
							Description: "the type the variable must be of",
							Optional:    true,
						},
						"schema_enum": &schema.Schema{
							Type:        schema.TypeList,
							Description: "a list of possible variable values",
							Optional:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "name to give the configuration variable",
							},
						},
					},
				},
			},
		},
	}
}

var VariableTypes = map[string]client.ConfigurationVariableType{
	"terraform":   client.ConfigurationVariableTypeTerraform,
	"environment": client.ConfigurationVariableTypeEnvironment,
}

func setEnvironmentSchema(d *schema.ResourceData, environment client.Environment) {
	d.Set("id", environment.Id)
	d.Set("name", environment.Name)
	d.Set("project_id", environment.ProjectId)
	d.Set("template_id", environment.LatestDeploymentLog.BlueprintId)
	d.Set("workspace", environment.WorkspaceName)
	d.Set("revision", environment.LatestDeploymentLog.BlueprintRevision)
	d.Set("repository", environment.LatestDeploymentLog.BlueprintRepository)
	d.Set("run_plan_on_pull_requests", environment.PullRequestPlanDeployments)
	d.Set("approve_plan_automatically", !environment.RequiresApproval)
	d.Set("deploy_on_push", environment.ContinuousDeployment)
	d.Set("deployment_id", environment.LatestDeploymentLogId)
	d.Set("auto_deploy_by_custom_glob", environment.AutoDeployByCustomGlob)
	//TODO: TTL and env\terraform variables
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

	if d.HasChanges("template_id", "revision", "repository", "configuration") {
		deployPayload := getDeployPayload(d)
		deployResponse, err := apiClient.EnvironmentDeploy(d.Id(), deployPayload)
		if err != nil {
			return diag.Errorf("failed deploying environment: %v", err)
		}
		d.Set("deployment_id", deployResponse.Id)
	}

	// TODO: update TTL if needed, also consider not updating ttl if deploy happened (cause we update ttl there too)

	if d.HasChanges("name", "approve_plan_automatically", "deploy_on_push", "run_plan_on_pull_requests", "auto_deploy_by_custom_glob") {
		payload := getUpdatePayload(d)
		_, err := apiClient.EnvironmentUpdate(d.Id(), payload)
		if err != nil {
			return diag.Errorf("could not update environment: %v", err)
		}
	}

	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	_, err := apiClient.EnvironmentDestroy(d.Id())
	if err != nil {
		return diag.Errorf("could not delete environment: %v", err)
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

	deployPayload := getDeployPayload(d)

	payload.DeployRequest = &deployPayload

	return payload
}

func getUpdatePayload(d *schema.ResourceData) client.EnvironmentUpdate {
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

func getDeployPayload(d *schema.ResourceData) client.DeployRequest {
	payload := client.DeployRequest{}

	if templateId, ok := d.GetOk("template_id"); ok {
		payload.BlueprintId = templateId.(string)
	}

	if revision, ok := d.GetOk("revision"); ok {
		payload.BlueprintRevision = revision.(string)
	}

	if repository, ok := d.GetOk("repository"); ok {
		payload.BlueprintRepository = repository.(string)
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		configurationChanges := getConfigurationVariables(configuration.([]interface{}))
		payload.ConfigurationChanges = &configurationChanges
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		payload.TTL = &client.TTL{
			Type:  ttl.(map[string]interface{})["type"].(string),
			Value: ttl.(map[string]interface{})["value"].(string),
		}
	}

	if userRequiresApproval, ok := d.GetOk("requires_approval"); ok {
		payload.UserRequiresApproval = userRequiresApproval.(bool)
	}

	return payload
}

func getConfigurationVariables(configuration []interface{}) client.ConfigurationChanges {
	configurationChanges := client.ConfigurationChanges{}
	for _, variable := range configuration {
		configurationVariable := getConfigurationVariableForEnvironment(variable.(map[string]interface{}))
		configurationChanges = append(configurationChanges, configurationVariable)
	}
	return configurationChanges
}

func getConfigurationVariableForEnvironment(variable map[string]interface{}) client.ConfigurationVariable {
	varType := VariableTypes[variable["type"].(string)]

	configurationVariable := client.ConfigurationVariable{
		Name:  variable["name"].(string),
		Value: variable["value"].(string),
		Scope: client.Scope(variable["scope"].(string)),
		Type:  &varType,
	}

	if variable["scope_id"] != nil {
		configurationVariable.ScopeId = variable["scope_id"].(string)
	}

	if variable["is_sensitive"] != nil {
		isSensitive := variable["is_sensitive"].(bool)
		configurationVariable.IsSensitive = &isSensitive
	}

	if variable["description"] != nil {
		configurationVariable.Description = variable["description"].(string)
	}

	if variable["schema_type"] != nil && variable["schema_enum"] != nil {
		enumOfAny := variable["schema_enum"].([]interface{})
		enum := make([]string, len(enumOfAny))
		for i := range enum {
			enum[i] = enumOfAny[i].(string)
		}
		schema := client.ConfigurationVariableSchema{
			Type: variable["schema_type"].(string),
			Enum: enum,
		}
		configurationVariable.Schema = &schema
	}

	return configurationVariable
}

func getEnvironmentByName(name interface{}, meta interface{}) (client.Environment, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	environments, err := apiClient.Environments()
	if err != nil {
		return client.Environment{}, diag.Errorf("Could not get Environment: %v", err)
	}

	var environmentsByName []client.Environment
	for _, candidate := range environments {
		if candidate.Name == name {
			environmentsByName = append(environmentsByName, candidate)
		}
	}

	if len(environmentsByName) > 1 {
		return client.Environment{}, diag.Errorf("Found multiple environments for name: %s. Use ID instead or make sure environment names are unique %v", name, environmentsByName)
	}

	if len(environmentsByName) == 0 {
		return client.Environment{}, diag.Errorf("Could not find an env0 environment with name %s", name)
	}

	return environmentsByName[0], nil
}

func getEnvironmentById(environmentId string, meta interface{}) (client.Environment, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	environment, err := apiClient.Environment(environmentId)
	if err != nil {
		return client.Environment{}, diag.Errorf("Could not find environment: %v", err)
	}
	return environment, nil
}
