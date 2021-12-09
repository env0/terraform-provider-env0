package env0

import (
	"context"
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
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
				ForceNew:    true,
			},
			"workspace": {
				Type:        schema.TypeString,
				Description: "the terraform workspace of the environment",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "the revision the environment is to be run against",
				Optional:    true,
				Computed:    true,
			},
			"run_plan_on_pull_requests": {
				Type:        schema.TypeBool,
				Description: "should run terraform plan on pull requests creations",
				Optional:    true,
				Computed:    true,
			},
			"approve_plan_automatically": {
				Type:        schema.TypeBool,
				Description: "should deployments require manual approvals",
				Optional:    true,
				Computed:    true,
			},
			"deploy_on_push": {
				Type:        schema.TypeBool,
				Description: "should run terraform deploy on push events",
				Optional:    true,
				Computed:    true,
			},
			"auto_deploy_on_path_changes_only": {
				Type:        schema.TypeBool,
				Description: "redeploy only on path changes only",
				Optional:    true,
				Computed:    true,
			},
			"auto_deploy_by_custom_glob": {
				Type:         schema.TypeString,
				Description:  "redeploy on file filter pattern",
				RequiredWith: []string{"auto_deploy_on_path_changes_only"},
				Optional:     true,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Description: "id of the last deployment",
				Computed:    true,
			},
			"ttl": {
				Type:        schema.TypeString,
				Description: "the date the environment should be destroyed at (iso format)",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					utcPattern := `\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d`
					ttl := val.(string)
					matched, err := regexp.MatchString(utcPattern, ttl)
					if !matched || err != nil {
						errs = append(errs, fmt.Errorf("%q must be of iso format (for example: \"2021-12-13T10:00:00Z\"), got: %q", key, ttl))
					}
					return
				},
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

func setEnvironmentSchema(d *schema.ResourceData, environment client.Environment) {
	d.Set("id", environment.Id)
	d.Set("name", environment.Name)
	d.Set("project_id", environment.ProjectId)
	d.Set("workspace", environment.WorkspaceName)
	d.Set("auto_deploy_by_custom_glob", environment.AutoDeployByCustomGlob)
	d.Set("ttl", environment.LifespanEndAt)
	if environment.LatestDeploymentLog != (client.DeploymentLog{}) {
		d.Set("template_id", environment.LatestDeploymentLog.BlueprintId)
		d.Set("revision", environment.LatestDeploymentLog.BlueprintRevision)
	}
	if environment.PullRequestPlanDeployments != nil {
		d.Set("run_plan_on_pull_requests", *environment.PullRequestPlanDeployments)
	}
	if environment.RequiresApproval != nil {
		d.Set("approve_plan_automatically", !*environment.RequiresApproval)
	}
	if environment.ContinuousDeployment != nil {
		d.Set("deploy_on_push", *environment.ContinuousDeployment)
	}
	if environment.AutoDeployOnPathChangesOnly != nil {
		d.Set("auto_deploy_on_path_changes_only", *environment.AutoDeployOnPathChangesOnly)
	}
	//TODO: env\terraform variables
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := getCreatePayload(d)

	environment, err := apiClient.EnvironmentCreate(payload)
	if err != nil {
		return diag.Errorf("could not create environment: %v", err)
	}

	d.SetId(environment.Id)
	d.Set("deployment_id", environment.LatestDeploymentLogId)
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

	if shouldDeploy(d) {
		err := deploy(d, apiClient)
		if err != nil {
			return err
		}
	}

	// TODO: update TTL if needed, also consider not updating ttl if deploy happened (cause we update ttl there too)

	if shouldUpdate(d) {
		err := update(d, apiClient)
		if err != nil {
			return err
		}
	}

	if shouldUpdateTTL(d) {
		err := updateTTL(d, apiClient)
		if err != nil {
			return err
		}
	}

	return nil
}

func shouldDeploy(d *schema.ResourceData) bool {
	return d.HasChanges("revision", "configuration")
}

func shouldUpdate(d *schema.ResourceData) bool {
	return d.HasChanges("name", "approve_plan_automatically", "deploy_on_push", "run_plan_on_pull_requests", "auto_deploy_by_custom_glob")
}

func shouldUpdateTTL(d *schema.ResourceData) bool {
	return d.HasChange("ttl")
}

func deploy(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	deployPayload := getDeployPayload(d)
	deployResponse, err := apiClient.EnvironmentDeploy(d.Id(), deployPayload)
	if err != nil {
		return diag.Errorf("failed deploying environment: %v", err)
	}
	d.Set("deployment_id", deployResponse.Id)
	return nil
}

func update(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	payload := getUpdatePayload(d)
	_, err := apiClient.EnvironmentUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update environment: %v", err)
	}
	return nil
}

func updateTTL(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	ttl := d.Get("ttl").(string)
	payload := getTTl(ttl)
	_, err := apiClient.EnvironmentUpdateTTL(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update the environment's ttl: %v", err)
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
	if continuousDeployment, ok := d.GetOkExists("deploy_on_push"); ok {
		continuousDeployment := continuousDeployment.(bool)
		payload.ContinuousDeployment = &continuousDeployment
	}
	if requiresApproval, ok := d.GetOkExists("approve_plan_automatically"); ok {
		requiresApproval := !requiresApproval.(bool)
		payload.RequiresApproval = &requiresApproval
	}
	if pullRequestPlanDeployments, ok := d.GetOkExists("run_plan_on_pull_requests"); ok {
		pullRequestPlanDeployments := pullRequestPlanDeployments.(bool)
		payload.PullRequestPlanDeployments = &pullRequestPlanDeployments
	}
	if autoDeployOnPathChangesOnly, ok := d.GetOkExists("auto_deploy_on_path_changes_only"); ok {
		autoDeployOnPathChangesOnly := autoDeployOnPathChangesOnly.(bool)
		payload.AutoDeployOnPathChangesOnly = &autoDeployOnPathChangesOnly
	}
	if autoDeployByCustomGlob, ok := d.GetOk("auto_deploy_by_custom_glob"); ok {
		payload.AutoDeployByCustomGlob = autoDeployByCustomGlob.(string)
	}
	if ttl, ok := d.GetOk("ttl"); ok {
		ttlPayload := getTTl(ttl.(string))
		payload.TTL = &ttlPayload
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
	if requiresApproval, ok := d.GetOkExists("approve_plan_automatically"); ok {
		requiresApproval := requiresApproval.(bool)
		payload.RequiresApproval = &requiresApproval
	}
	if continuousDeployment, ok := d.GetOkExists("deploy_on_push"); ok {
		continuousDeployment := continuousDeployment.(bool)
		payload.ContinuousDeployment = &continuousDeployment
	}
	if pullRequestPlanDeployments, ok := d.GetOkExists("run_plan_on_pull_requests"); ok {
		pullRequestPlanDeployments := pullRequestPlanDeployments.(bool)
		payload.PullRequestPlanDeployments = &pullRequestPlanDeployments
	}
	if autoDeployOnPathChangesOnly, ok := d.GetOkExists("auto_deploy_on_path_changes_only"); ok {
		autoDeployOnPathChangesOnly := autoDeployOnPathChangesOnly.(bool)
		payload.AutoDeployOnPathChangesOnly = &autoDeployOnPathChangesOnly
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

	if configuration, ok := d.GetOk("configuration"); ok {
		configurationChanges := getConfigurationVariables(configuration.([]interface{}))
		payload.ConfigurationChanges = &configurationChanges
	}

	if userRequiresApproval, ok := d.GetOkExists("requires_approval"); ok {
		userRequiresApproval := userRequiresApproval.(bool)
		payload.UserRequiresApproval = &userRequiresApproval
	}

	return payload
}

func getTTl(date string) client.TTL {
	if date != "" {
		return client.TTL{
			Type:  client.TTLTypeDate,
			Value: date,
		}
	}
	return client.TTL{
		Type:  client.TTlTypeInfinite,
		Value: "",
	}
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
	varType := client.VariableTypes[variable["type"].(string)]

	configurationVariable := client.ConfigurationVariable{
		Name:  variable["name"].(string),
		Value: variable["value"].(string),
		Scope: client.ScopeDeployment,
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
