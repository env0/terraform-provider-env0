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
				Description: "The environment's name",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id of the environment",
				Required:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the template id the environment is to be created from",
				Required:    true,
			},
		},
	}
}

// TODO: make it a const
var VariableTypes = map[string]client.ConfigurationVariableType{
	"terraform":   client.ConfigurationVariableTypeTerraform,
	"environment": client.ConfigurationVariableTypeEnvironment,
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

	// TODO: deploy if needed

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

	payload.DeployRequest = &client.DeployRequest{}

	if blueprintId, ok := d.GetOk("template_id"); ok {
		payload.DeployRequest.BlueprintId = blueprintId.(string)
	}

	if blueprintRepository, ok := d.GetOk("repository"); ok {
		payload.DeployRequest.BlueprintRepository = blueprintRepository.(string)
	}

	if blueprintRevision, ok := d.GetOk("revision"); ok {
		payload.DeployRequest.BlueprintRevision = blueprintRevision.(string)
	}

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

func getDeployPayload(d *schema.ResourceData) (client.DeployRequest, diag.Diagnostics) {
	payload := client.DeployRequest{}

	if templateId, ok := d.GetOk("templateId"); ok {
		payload.BlueprintId = templateId.(string)
	}

	if revision, ok := d.GetOk("revision"); ok {
		payload.BlueprintRevision = revision.(string)
	}

	if repository, ok := d.GetOk("repository"); ok {
		payload.BlueprintRepository = repository.(string)
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		configurationChanges, err := getConfigurationVariables(configuration.([]interface{}))
		if err != nil {
			return client.DeployRequest{}, err
		}
		payload.ConfigurationChanges = &configurationChanges
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		payload.TTL = &client.TTL{
			Type:  ttl.(map[string]interface{})["type"].(string),
			Value: ttl.(map[string]interface{})["value"].(string),
		}
	}

	if envName, ok := d.GetOk("name"); ok {
		payload.EnvName = envName.(string)
	}

	if userRequiresApproval, ok := d.GetOk("requires_approval"); ok {
		payload.UserRequiresApproval = userRequiresApproval.(bool)
	}

	return payload, nil
}

func getConfigurationVariables(configuration []interface{}) (client.ConfigurationChanges, diag.Diagnostics) {
	configurationChanges := client.ConfigurationChanges{}
	for _, variable := range configuration {
		configurationVariable, err := getConfigurationVariableForEnvironment(variable.(map[string]interface{}))
		if err != nil {
			return client.ConfigurationChanges{}, err
		}
		configurationChanges = append(configurationChanges, configurationVariable)
	}
	return configurationChanges, nil
}

func getConfigurationVariableForEnvironment(variable map[string]interface{}) (client.ConfigurationVariable, diag.Diagnostics) {
	configurationVariable := client.ConfigurationVariable{}

	if variable["name"] == nil {
		return client.ConfigurationVariable{}, diag.Errorf("failed reading configuration variables. name, value, scope and type are required")
	}
	configurationVariable.Name = variable["name"].(string)

	if variable["value"] == nil {
		return client.ConfigurationVariable{}, diag.Errorf("failed reading configuration variables. name, value, scope and type are required")
	}
	configurationVariable.Value = variable["value"].(string)

	if variable["scope"] == nil {
		return client.ConfigurationVariable{}, diag.Errorf("failed reading configuration variables. name, value, scope and type are required")
	}
	configurationVariable.Scope = variable["scope"].(client.Scope)

	if variable["type"] == nil {
		return client.ConfigurationVariable{}, diag.Errorf("failed reading configuration variables. name, value, scope and type are required")
	}
	configurationVariable.Type = VariableTypes[variable["type"].(string)]

	if variable["scope_id"] != nil {
		configurationVariable.ScopeId = variable["scope_id"].(string)
	}

	if variable["organization_id"] != nil {
		configurationVariable.OrganizationId = variable["organization_id"].(string)
	}

	if variable["user_id"] != nil {
		configurationVariable.UserId = variable["user_id"].(string)
	}

	if variable["is_sensitive"] != nil {
		configurationVariable.IsSensitive = variable["is_sensitive"].(bool)
	}

	if variable["description"] != nil {
		configurationVariable.Description = variable["description"].(string)
	}

	if variable["schema"] != nil {
		configurationVariable.Schema = getConfigurationVariableSchema(variable["schema"].(map[string]interface{}))
	}

	return configurationVariable, nil
}

func getConfigurationVariableSchema(schema map[string]interface{}) client.ConfigurationVariableSchema {
	return client.ConfigurationVariableSchema{
		Type: schema["type"].(string),
		Enum: schema["enum"].([]string),
	}
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
