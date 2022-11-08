package env0

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isTemplateless(d *schema.ResourceData) bool {
	_, ok := d.GetOk("without_template_settings.0")
	return ok
}

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceEnvironmentImport},

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
				Type:         schema.TypeString,
				Description:  "the template id the environment is to be created from.\nImportant note: the template must first be assigned to the same project as the environment (project_id). Use 'env0_template_project_assignment' to assign the template to the project. In addition, be sure to leverage 'depends_on' if applicable.",
				Optional:     true,
				ExactlyOneOf: []string{"without_template_settings", "template_id"},
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
			},
			"approve_plan_automatically": {
				Type:        schema.TypeBool,
				Description: "should deployments require manual approvals",
				Optional:    true,
			},
			"deploy_on_push": {
				Type:        schema.TypeBool,
				Description: "should run terraform deploy on push events",
				Optional:    true,
			},
			"auto_deploy_on_path_changes_only": {
				Type:        schema.TypeBool,
				Description: "redeploy only on path changes only",
				Optional:    true,
			},
			"auto_deploy_by_custom_glob": {
				Type:         schema.TypeString,
				Description:  "redeploy on file filter pattern.\nWhen used 'auto_deploy_on_path_changes_only' must be configured to true and 'deploy_on_push' or 'run_plan_on_pull_requests' must be configured to true.",
				Optional:     true,
				RequiredWith: []string{"auto_deploy_on_path_changes_only"},
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Description: "id of the last deployment",
				Computed:    true,
			},
			"output": {
				Type:        schema.TypeString,
				Description: "the deployment log output. Returns a json string. It can be either a map of key-value, or an array of (in case of Terragrunt run-all) of moduleName and a map of key-value. Note: if the deployment is still in progress returns 'null'",
				Computed:    true,
				Optional:    true,
			},
			"ttl": {
				Type:        schema.TypeString,
				Description: "the date the environment should be destroyed at (iso format). omitting this attribute will result in infinite ttl.",
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
			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "Destroy safeguard. Must be enabled before delete/destroy",
				Optional:    true,
			},
			"is_remote_backend": {
				Type:        schema.TypeBool,
				Description: "should use remote backend",
				Optional:    true,
			},
			"terragrunt_working_directory": {
				Type:        schema.TypeString,
				Description: "The working directory path to be used by a Terragrunt template. If left empty '/' is used.",
				Optional:    true,
			},
			"vcs_commands_alias": {
				Type:        schema.TypeString,
				Description: "set an alias for this environment in favor of running VCS commands using PR comments against it. Additional details: https://docs.env0.com/docs/plan-and-apply-from-pr-comments",
				Optional:    true,
			},
			"configuration": {
				Type:        schema.TypeList,
				Description: "terraform and environment variables for the environment",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "variable name",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "variable value",
							Required:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "variable type (allowed values are: terraform, environment)",
							Default:     "environment",
							Optional:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								value := val.(string)
								if value != "environment" && value != "terraform" {
									errs = append(errs, fmt.Errorf("%q can be either \"environment\" or \"terraform\", got: %q", key, value))
								}
								return
							},
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description for the variable",
							Optional:    true,
						},
						"is_sensitive": {
							Type:        schema.TypeBool,
							Description: "should the variable value be hidden",
							Optional:    true,
						},
						"schema_type": {
							Type:        schema.TypeString,
							Description: "the type the variable must be of",
							Optional:    true,
							Default:     "string",
						},
						"schema_enum": {
							Type:        schema.TypeList,
							Description: "a list of possible variable values",
							Optional:    true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "name to give the configuration variable",
							},
						},
						"schema_format": {
							Type:         schema.TypeString,
							Description:  "the variable format:",
							Default:      "",
							Optional:     true,
							ValidateFunc: ValidateConfigurationPropertySchema,
						},
						"is_read_only": {
							Type:        schema.TypeBool,
							Description: "is the variable read only",
							Optional:    true,
							Default:     false,
						},
						"is_required": {
							Type:        schema.TypeBool,
							Description: "is the variable required",
							Optional:    true,
							Default:     false,
						},
						"regex": {
							Type:        schema.TypeString,
							Description: "the value of this variable must match provided regular expression (enforced only in env0 UI)",
							Optional:    true,
						},
					},
				},
			},
			"without_template_settings": {
				Type:         schema.TypeList,
				Description:  "settings for creating an environment without a template. Is not imported when running the import command",
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				ExactlyOneOf: []string{"without_template_settings", "template_id"},
				Elem: &schema.Resource{
					Schema: getTemplateSchema("without_template_settings.0."),
				},
			},
		},

		CustomizeDiff: customdiff.ForceNewIf("template_id", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			// For templateless: any changes in template_id, no need to do anything (template id can't change).
			// This is done due to historical bugs/issues.
			if _, ok := d.GetOk("without_template_settings.0"); ok {
				return false
			}
			return true
		}),
	}
}

func setEnvironmentSchema(d *schema.ResourceData, environment client.Environment, configurationVariables client.ConfigurationChanges) error {
	if err := writeResourceData(&environment, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	if !isTemplateless(d) {
		if environment.LatestDeploymentLog.BlueprintId != "" {
			d.Set("template_id", environment.LatestDeploymentLog.BlueprintId)
			d.Set("revision", environment.LatestDeploymentLog.BlueprintRevision)
		}
	} else if environment.BlueprintId != "" || environment.LatestDeploymentLog.BlueprintId != "" {
		settings := d.Get("without_template_settings").([]interface{})
		elem := settings[0].(map[string]interface{})

		if environment.BlueprintId != "" {
			elem["id"] = environment.BlueprintId
		} else {
			elem["id"] = environment.LatestDeploymentLog.BlueprintId
		}

		d.Set("without_template_settings", settings)
	}

	if len(environment.LatestDeploymentLog.Output) == 0 {
		d.Set("output", "null")
	} else {
		d.Set("output", string(environment.LatestDeploymentLog.Output))
	}

	if environment.RequiresApproval != nil {
		d.Set("approve_plan_automatically", !*environment.RequiresApproval)
	}

	setEnvironmentConfigurationSchema(d, configurationVariables)

	return nil
}

func createVariable(configurationVariable *client.ConfigurationVariable) interface{} {
	variable := make(map[string]interface{})

	variable["name"] = configurationVariable.Name
	variable["value"] = configurationVariable.Value

	if configurationVariable.Type == nil || *configurationVariable.Type == 0 {
		variable["type"] = "environment"
	} else {
		variable["type"] = "terraform"
	}

	if configurationVariable.Description != "" {
		variable["description"] = configurationVariable.Description
	}

	if configurationVariable.Regex != "" {
		variable["regex"] = configurationVariable.Regex
	}

	if configurationVariable.IsSensitive != nil {
		variable["is_sensitive"] = configurationVariable.IsSensitive
	}

	if configurationVariable.IsReadOnly != nil {
		variable["is_read_only"] = configurationVariable.IsReadOnly
	}

	if configurationVariable.IsRequired != nil {
		variable["is_required"] = configurationVariable.IsRequired
	}

	if configurationVariable.Schema != nil {
		variable["schema_type"] = configurationVariable.Schema.Type
		variable["schema_enum"] = configurationVariable.Schema.Enum
		variable["schema_format"] = configurationVariable.Schema.Format
	}

	return variable
}

func setEnvironmentConfigurationSchema(d *schema.ResourceData, configurationVariables []client.ConfigurationVariable) {
	ivariables, ok := d.GetOk("configuration")
	if !ok {
		return
	}

	if ivariables == nil {
		ivariables = make([]interface{}, 0)
	}

	variables := ivariables.([]interface{})

	newVariables := make([]interface{}, 0)

	// The goal is to maintain existing state order as much as possible. (The backend response order may vary from state).
	for _, ivariable := range variables {
		variable := ivariable.(map[string]interface{})
		variableName := variable["name"].(string)

		for _, configurationVariable := range configurationVariables {
			if configurationVariable.Name == variableName {
				newVariables = append(newVariables, createVariable(&configurationVariable))
				break
			}
		}
	}

	// Check for drifts: add new configuration variables received from the backend.
	for _, configurationVariable := range configurationVariables {
		found := false

		for _, ivariable := range variables {
			variable := ivariable.(map[string]interface{})
			variableName := variable["name"].(string)
			if configurationVariable.Name == variableName {
				found = true
				break
			}
		}

		if !found {
			log.Printf("[WARN] Drift Detected for configuration: %s", configurationVariable.Name)
			newVariables = append(newVariables, createVariable(&configurationVariable))
		}
	}

	if len(newVariables) > 0 {
		d.Set("configuration", newVariables)
	} else {
		d.Set("configuration", nil)
	}
}

// Validate that the template is assigned to the "project_id".
func validateTemplateProjectAssignment(d *schema.ResourceData, apiClient client.ApiClientInterface) error {
	projectId := d.Get("project_id").(string)
	templateId := d.Get("template_id").(string)

	template, err := apiClient.Template(templateId)
	if err != nil {
		return fmt.Errorf("could not get template: %v", err)
	}

	if projectId != template.ProjectId && !stringInSlice(projectId, template.ProjectIds) {
		return errors.New("could not create environment: template is not assigned to project")
	}

	return nil
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentPayload, createEnvPayloadErr := getCreatePayload(d, apiClient)
	if createEnvPayloadErr != nil {
		return createEnvPayloadErr
	}

	var environment client.Environment
	var err error

	if !isTemplateless(d) {
		if err := validateTemplateProjectAssignment(d, apiClient); err != nil {
			return diag.Errorf("%v\n", err)
		}
		environment, err = apiClient.EnvironmentCreate(environmentPayload)
	} else {
		templatePayload, createTemPayloadErr := templateCreatePayloadFromParameters("without_template_settings.0", d)
		if createTemPayloadErr != nil {
			return createTemPayloadErr
		}
		payload := client.EnvironmentCreateWithoutTemplate{
			EnvironmentCreate: environmentPayload,
			TemplateCreate:    templatePayload,
		}
		// Note: the blueprint id field of the environment is returned only during creation of a template without envrionment.
		// Afterward, it will be omitted from future response.
		// setEnvironmentSchema() (several lines below) sets the blueprint id in the resource (under "without_template_settings.0.id").
		environment, err = apiClient.EnvironmentCreateWithoutTemplate(payload)
	}
	if err != nil {
		return diag.Errorf("could not create environment: %v", err)
	}
	environmentConfigurationVariables := client.ConfigurationChanges{}
	if environmentPayload.DeployRequest.ConfigurationChanges != nil {
		environmentConfigurationVariables = *environmentPayload.DeployRequest.ConfigurationChanges
	}
	d.SetId(environment.Id)
	d.Set("deployment_id", environment.LatestDeploymentLogId)
	d.Set("auto_deploy_on_path_changes_only", environment.AutoDeployOnPathChangesOnly)
	setEnvironmentSchema(d, environment, environmentConfigurationVariables)

	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environment, err := apiClient.Environment(d.Id())
	if err != nil {
		return diag.Errorf("could not get environment: %v", err)
	}

	environmentConfigurationVariables, err := apiClient.ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id)
	if err != nil {
		return diag.Errorf("could not get environment configuration variables: %v", err)
	}

	setEnvironmentSchema(d, environment, environmentConfigurationVariables)

	if isTemplateless(d) {
		// envrionment with no template.
		templateId := d.Get("without_template_settings.0.id").(string)
		template, err := apiClient.Template(templateId)
		if err != nil {
			return diag.Errorf("could not get template: %v", err)
		}
		if err := templateRead("without_template_settings", template, d); err != nil {
			return diag.Errorf("schema resource data serialization failed: %v", err)
		}
	}

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if shouldUpdate(d) {
		if err := update(d, apiClient); err != nil {
			return err
		}
	}

	if shouldUpdateTTL(d) {

		if err := updateTTL(d, apiClient); err != nil {
			return err
		}
	}

	if shouldUpdateTemplate(d) {
		if err := updateTemplate(d, apiClient); err != nil {
			return err
		}
	}

	if shouldDeploy(d) {
		if err := deploy(d, apiClient); err != nil {
			return err
		}
	}

	return nil
}

func shouldUpdateTemplate(d *schema.ResourceData) bool {
	return isTemplateless(d) && d.HasChange("without_template_settings.0")
}

func shouldDeploy(d *schema.ResourceData) bool {
	return d.HasChanges("revision", "configuration")
}

func shouldUpdate(d *schema.ResourceData) bool {
	return d.HasChanges("name", "approve_plan_automatically", "deploy_on_push", "run_plan_on_pull_requests", "auto_deploy_by_custom_glob", "auto_deploy_on_path_changes_only", "terragrunt_working_directory", "vcs_commands_alias", "is_remote_backend")
}

func shouldUpdateTTL(d *schema.ResourceData) bool {
	return d.HasChange("ttl")
}

func updateTemplate(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	payload, problem := templateCreatePayloadFromParameters("without_template_settings.0", d)
	if problem != nil {
		return problem
	}

	templateId := d.Get("without_template_settings.0.id").(string)

	if _, err := apiClient.TemplateUpdate(templateId, payload); err != nil {
		return diag.Errorf("could not update template: %v", err)
	}

	return nil
}

func deploy(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	deployPayload := getDeployPayload(d, apiClient, true)
	deployResponse, err := apiClient.EnvironmentDeploy(d.Id(), deployPayload)
	if err != nil {
		return diag.Errorf("failed deploying environment: %v", err)
	}
	d.Set("deployment_id", deployResponse.Id)
	return nil
}

func update(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	payload, updateEnvPayloadErr := getUpdatePayload(d)

	if updateEnvPayloadErr != nil {
		return diag.Errorf("%v", updateEnvPayloadErr)
	}

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
	canDestroy := d.Get("force_destroy")

	if canDestroy != true {
		return diag.Errorf(`must enable "force_destroy" safeguard in order to destroy`)
	}

	apiClient := meta.(client.ApiClientInterface)

	_, err := apiClient.EnvironmentDestroy(d.Id())
	if err != nil {
		return diag.Errorf("could not delete environment: %v", err)
	}
	return nil
}

func getCreatePayload(d *schema.ResourceData, apiClient client.ApiClientInterface) (client.EnvironmentCreate, diag.Diagnostics) {
	var payload client.EnvironmentCreate

	if err := readResourceData(&payload, d); err != nil {
		return client.EnvironmentCreate{}, diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("deploy_on_push"); exists {
		payload.ContinuousDeployment = boolPtr(val.(bool))
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("run_plan_on_pull_requests"); exists {
		payload.PullRequestPlanDeployments = boolPtr(val.(bool))
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("auto_deploy_on_path_changes_only"); exists {
		payload.AutoDeployOnPathChangesOnly = boolPtr(val.(bool))
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("approve_plan_automatically"); exists {
		payload.RequiresApproval = boolPtr(!val.(bool))
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("is_remote_backend"); exists {
		payload.IsRemoteBackend = boolPtr(val.(bool))
	}

	if err := assertDeploymentTriggers(d); err != nil {
		return client.EnvironmentCreate{}, err
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		configurationChanges := getConfigurationVariablesFromSchema(configuration.([]interface{}))
		payload.ConfigurationChanges = &configurationChanges
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		ttlPayload := getTTl(ttl.(string))
		payload.TTL = &ttlPayload
	}

	deployPayload := getDeployPayload(d, apiClient, false)
	payload.DeployRequest = &deployPayload

	return payload, nil
}

func assertDeploymentTriggers(d *schema.ResourceData) diag.Diagnostics {
	continuousDeployment := d.Get("deploy_on_push").(bool)
	pullRequestPlanDeployments := d.Get("run_plan_on_pull_requests").(bool)
	autoDeployOnPathChangesOnly := d.Get("auto_deploy_on_path_changes_only").(bool)
	autoDeployByCustomGlob := d.Get("auto_deploy_by_custom_glob").(string)

	if autoDeployByCustomGlob != "" {
		if !continuousDeployment && !pullRequestPlanDeployments {
			return diag.Errorf("run_plan_on_pull_requests or deploy_on_push must be enabled for auto_deploy_by_custom_glob")
		}
		if !autoDeployOnPathChangesOnly {
			return diag.Errorf("cannot set auto_deploy_by_custom_glob when auto_deploy_on_path_changes_only is disabled")
		}
	}

	return nil
}

func getUpdatePayload(d *schema.ResourceData) (client.EnvironmentUpdate, diag.Diagnostics) {
	var payload client.EnvironmentUpdate

	if err := readResourceData(&payload, d); err != nil {
		return client.EnvironmentUpdate{}, diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if d.HasChange("approve_plan_automatically") {
		payload.RequiresApproval = boolPtr(!d.Get("approve_plan_automatically").(bool))
	}

	if d.HasChange("deploy_on_push") {
		payload.ContinuousDeployment = boolPtr(d.Get("deploy_on_push").(bool))
	}

	if d.HasChange("run_plan_on_pull_requests") {
		payload.PullRequestPlanDeployments = boolPtr(d.Get("run_plan_on_pull_requests").(bool))
	}

	if d.HasChange("auto_deploy_on_path_changes_only") {
		payload.AutoDeployOnPathChangesOnly = boolPtr(d.Get("auto_deploy_on_path_changes_only").(bool))
	}

	if d.HasChange("is_remote_backend") {
		payload.IsRemoteBackend = boolPtr(d.Get("is_remote_backend").(bool))
	}

	if err := assertDeploymentTriggers(d); err != nil {
		return client.EnvironmentUpdate{}, err
	}

	return payload, nil
}

func getDeployPayload(d *schema.ResourceData, apiClient client.ApiClientInterface, isRedeploy bool) client.DeployRequest {
	payload := client.DeployRequest{}

	if isTemplateless(d) {
		templateId, ok := d.GetOk("without_template_settings.0.id")
		if ok {
			payload.BlueprintId = templateId.(string)
		}
	} else {
		payload.BlueprintId = d.Get("template_id").(string)
	}

	if revision, ok := d.GetOk("revision"); ok {
		payload.BlueprintRevision = revision.(string)
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		configurationChanges := getConfigurationVariablesFromSchema(configuration.([]interface{}))
		if isRedeploy {
			configurationChanges = getUpdateConfigurationVariables(configurationChanges, d.Get("id").(string), apiClient)
		}
		payload.ConfigurationChanges = &configurationChanges
	}

	if userRequiresApproval, ok := d.GetOk("requires_approval"); ok {
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

func getUpdateConfigurationVariables(configurationChanges client.ConfigurationChanges, environmentId string, apiClient client.ApiClientInterface) client.ConfigurationChanges {
	existVariables, err := apiClient.ConfigurationVariablesByScope(client.ScopeEnvironment, environmentId)
	if err != nil {
		diag.Errorf("could not get environment configuration variables: %v", err)
	}
	configurationChanges = linkToExistConfigurationVariables(configurationChanges, existVariables)
	configurationChanges = deleteUnusedConfigurationVariables(configurationChanges, existVariables)
	return configurationChanges
}

func getConfigurationVariablesFromSchema(configuration []interface{}) client.ConfigurationChanges {
	configurationChanges := client.ConfigurationChanges{}
	for _, variable := range configuration {
		configurationVariable := getConfigurationVariableFromSchema(variable.(map[string]interface{}))
		configurationChanges = append(configurationChanges, configurationVariable)
	}

	return configurationChanges
}

func deleteUnusedConfigurationVariables(configurationChanges client.ConfigurationChanges, existVariables client.ConfigurationChanges) client.ConfigurationChanges {
	for _, existVariable := range existVariables {
		if isExist, _ := isVariableExist(configurationChanges, existVariable); !isExist {
			toDelete := true
			existVariable.ToDelete = &toDelete
			configurationChanges = append(configurationChanges, existVariable)
		}
	}
	return configurationChanges
}

func linkToExistConfigurationVariables(configurationChanges client.ConfigurationChanges, existVariables client.ConfigurationChanges) client.ConfigurationChanges {
	updateConfigurationChanges := client.ConfigurationChanges{}
	for _, change := range configurationChanges {
		if isExist, existVariable := isVariableExist(existVariables, change); isExist {
			change.Id = existVariable.Id
		}
		updateConfigurationChanges = append(updateConfigurationChanges, change)
	}
	return updateConfigurationChanges
}

func isVariableExist(variables client.ConfigurationChanges, search client.ConfigurationVariable) (bool, client.ConfigurationVariable) {
	for _, variable := range variables {
		if variable.Name == search.Name && typeEqual(variable, search) {
			return true, variable
		}
	}
	return false, client.ConfigurationVariable{}
}

func typeEqual(variable client.ConfigurationVariable, search client.ConfigurationVariable) bool {
	return *variable.Type == *search.Type ||
		variable.Type == nil && *search.Type == client.ConfigurationVariableTypeEnvironment ||
		search.Type == nil && *variable.Type == client.ConfigurationVariableTypeEnvironment
}

func getConfigurationVariableFromSchema(variable map[string]interface{}) client.ConfigurationVariable {
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

	if variable["is_read_only"] != nil {
		isReadOnly := variable["is_read_only"].(bool)
		configurationVariable.IsReadOnly = &isReadOnly
	}

	if variable["is_required"] != nil {
		isRequired := variable["is_required"].(bool)
		configurationVariable.IsRequired = &isRequired
	}

	if variable["description"] != nil {
		configurationVariable.Description = variable["description"].(string)
	}

	if variable["regex"] != nil {
		configurationVariable.Regex = variable["regex"].(string)
	}

	configurationSchema := client.ConfigurationVariableSchema{
		Format: client.Format(variable["schema_format"].(string)),
		Enum:   nil,
		Type:   variable["schema_type"].(string),
	}

	if variable["schema_type"] != "" && len(variable["schema_enum"].([]interface{})) > 0 {
		enumOfAny := variable["schema_enum"].([]interface{})
		enum := make([]string, len(enumOfAny))
		for i := range enum {
			enum[i] = enumOfAny[i].(string)
		}
		configurationSchema.Type = variable["schema_type"].(string)
		configurationSchema.Enum = enum
	}

	configurationVariable.Schema = &configurationSchema
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

func resourceEnvironmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	var environment client.Environment
	_, err := uuid.Parse(id)
	if err == nil {
		log.Println("[INFO] Resolving Environment by id: ", id)
		environment, getErr = getEnvironmentById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving Environment by name: ", id)

		environment, getErr = getEnvironmentByName(id, meta)
	}
	apiClient := meta.(client.ApiClientInterface)
	d.SetId(environment.Id)
	environmentConfigurationVariables, err := apiClient.ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id)
	if err != nil {
		return nil, fmt.Errorf("could not get environment configuration variables: %v", err)
	}

	d.Set("deployment_id", environment.LatestDeploymentLogId)
	setEnvironmentSchema(d, environment, environmentConfigurationVariables)

	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
