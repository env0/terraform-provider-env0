package env0

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SubEnvironment struct {
	Id                       string
	Alias                    string
	Revision                 string
	Workflow                 string
	Workspace                string
	Configuration            client.ConfigurationChanges `tfschema:"-"`
	ApprovePlanAutomatically bool
}

func getSubEnvironments(d *schema.ResourceData) ([]SubEnvironment, error) {
	isubEnvironments, ok := d.GetOk("sub_environment_configuration")
	if !ok {
		return nil, nil
	}

	numberOfSubEnvironments := len(isubEnvironments.([]interface{}))

	var subEnvironments []SubEnvironment

	for i := 0; i < numberOfSubEnvironments; i++ {
		prefix := fmt.Sprintf("sub_environment_configuration.%d", i)

		var subEnvironment SubEnvironment

		if err := readResourceDataEx(prefix, &subEnvironment, d); err != nil {
			return nil, err
		}

		configurationPrefix := prefix + ".configuration"
		if configuration, ok := d.GetOk(configurationPrefix); ok {
			subEnvironment.Configuration = getConfigurationVariablesFromSchema(configuration.([]interface{}))

			for i := range subEnvironment.Configuration {
				subEnvironment.Configuration[i].Scope = client.ScopeEnvironment
			}
		}

		subEnvironments = append(subEnvironments, subEnvironment)
	}

	return subEnvironments, nil
}

func isTemplateless(d *schema.ResourceData) bool {
	_, ok := d.GetOk("without_template_settings.0")
	return ok
}

func resourceEnvironment() *schema.Resource {
	configurationSchema := &schema.Resource{
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
				Default:     false,
			},
			"schema_type": {
				Type:        schema.TypeString,
				Description: "the type the variable",
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
				Description:  "the variable format",
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
	}

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
			},
			"template_id": {
				Type:         schema.TypeString,
				Description:  "the template id the environment is to be created from.\nImportant note: the template must first be assigned to the same project as the environment (project_id). Use 'env0_template_project_assignment' to assign the template to the project. In addition, be sure to leverage 'depends_on' if applicable.\nImportant note: After the environment is created, this field cannot be modified.",
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
				Type:          schema.TypeString,
				Description:   "the revision the environment is to be run against. Please note that changing this attribute will require environment redeploy",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"without_template_settings"},
			},
			"run_plan_on_pull_requests": {
				Type:        schema.TypeBool,
				Description: "should run terraform plan on pull requests creations.\nIf true must specify one of the following - 'github_installation_id' if using GitHub, 'gitlab_project_id' and 'token_id' if using GitLab, or 'bitbucket_client_key' if using BitBucket.",
				Optional:    true,
			},
			"approve_plan_automatically": {
				Type:        schema.TypeBool,
				Description: "should deployments require manual approvals",
				Optional:    true,
			},
			"deploy_on_push": {
				Type:        schema.TypeBool,
				Description: "should run terraform deploy on push events.\nIf true must specify one of the following - 'github_installation_id' if using GitHub, 'gitlab_project_id' and 'token_id' if using GitLab, or 'bitbucket_client_key' if using BitBucket.",
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
			"prevent_auto_deploy": {
				Type:        schema.TypeBool,
				Description: "use this flag to prevent auto deploy on environment creation",
				Optional:    true,
			},
			"terragrunt_working_directory": {
				Type:        schema.TypeString,
				Description: "The working directory path to be used by a Terragrunt template. If left empty '/' is used. Note: modifying this field destroys the current environment and creates a new one",
				ForceNew:    true,
				Optional:    true,
			},
			"vcs_commands_alias": {
				Type:        schema.TypeString,
				Description: "set an alias for this environment in favor of running VCS commands using PR comments against it. Additional details: https://docs.env0.com/docs/plan-and-apply-from-pr-comments",
				Optional:    true,
			},
			"vcs_pr_comments_enabled": {
				Type:        schema.TypeBool,
				Description: "set to 'true' to enable running VCS PR plan/apply commands using PR comments. This can be set to 'true' (enabled) without setting alias in 'vcs_commands_alias'. Additional details: https://docs.env0.com/docs/plan-and-apply-from-pr-comments#configuration",
				Optional:    true,
			},
			"is_inactive": {
				Type:        schema.TypeBool,
				Description: "If 'true', it marks the environment as inactive. It can be re-activated by setting it to 'false' or removing this field. Note: it's not allowed to create an inactive environment",
				Default:     false,
				Optional:    true,
			},
			"configuration": {
				Type:        schema.TypeList,
				Description: "terraform and environment variables for the environment. Note: do not use with 'env0_configuration_variable' resource",
				Optional:    true,
				Elem:        configurationSchema,
			},
			"without_template_settings": {
				Type:         schema.TypeList,
				Description:  "settings for creating an environment without a template",
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				ExactlyOneOf: []string{"without_template_settings", "template_id"},
				Elem: &schema.Resource{
					Schema: getTemplateSchema("without_template_settings.0."),
				},
			},
			"sub_environment_configuration": {
				Type:        schema.TypeList,
				Description: "the subenvironments for a workflow environment. Template type must be 'workflow'. Must match the configuration as defined in 'env0.workflow.yml'",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "the id of the sub environment",
							Computed:    true,
						},
						"alias": {
							Type:        schema.TypeString,
							Description: "sub environment alias name",
							Required:    true,
						},
						"revision": {
							Type:        schema.TypeString,
							Description: "sub environment revision",
							Required:    true,
						},
						"workspace": {
							Type:        schema.TypeString,
							Description: "sub environment workspace (overrides the configurtion in the yml file)",
							Optional:    true,
						},
						"configuration": {
							Type:        schema.TypeList,
							Description: "sub environment configuration variables. Note: do not use with 'env0_configuration_variable' resource",
							Optional:    true,
							Elem:        configurationSchema,
						},
						"approve_plan_automatically": {
							Type:        schema.TypeBool,
							Description: "when 'true' (default) plans are approved automatically, otherwise ('false') deployment require manual approval",
							Optional:    true,
							Default:     true,
						},
					},
				},
			},
			"drift_detection_cron": {
				Type:             schema.TypeString,
				Description:      "cron expression for scheduled drift detection of the environment (cannot be used with resource_drift_detection resource)",
				Optional:         true,
				ValidateDiagFunc: ValidateCronExpression,
			},
			"is_remote_apply_enabled": {
				Type:        schema.TypeBool,
				Description: "enables remote apply when set to true (defaults to false). Can only be enabled when is_remote_backend and approve_plan_automatically are enabled",
				Optional:    true,
				Default:     false,
			},
			"removal_strategy": {
				Type:             schema.TypeString,
				Description:      "by default when removing an environment, it gets destroyed. Setting this value to 'mark_as_archived' will force the environment to be archived instead of tying to destroy it ('Mark as inactive' in the UI)",
				Optional:         true,
				Default:          "destroy",
				ValidateDiagFunc: NewStringInValidator([]string{"destroy", "mark_as_archived"}),
			},
			"k8s_namespace": {
				Type:        schema.TypeString,
				Description: "kubernetes (or helm) namespace to be used. If modified deletes current environment and creates a new one",
				Optional:    true,
				ForceNew:    true,
			},
			"variable_sets": {
				Type:        schema.TypeList,
				Description: "a list of IDs of variable sets to assign to this environment. Note: must not be used with 'env0_variable_set_assignment'",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "variable set id",
				},
			},
		},
		CustomizeDiff: customdiff.ValidateChange("template_id", func(ctx context.Context, oldValue, newValue, meta interface{}) error {
			if oldValue != "" && oldValue != newValue {
				return errors.New("template_id may not be modified, create a new environment instead")
			}

			return nil
		}),
	}
}

func setEnvironmentSchema(ctx context.Context, d *schema.ResourceData, environment client.Environment, configurationVariables client.ConfigurationChanges, variableSetsIds []string) error {
	if err := writeResourceData(&environment, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if _, exists := d.GetOkExists("vcs_pr_comments_enabled"); !exists {
		// VcsPrCommentsEnabled may have been "forced" to be 'true', ignore any drifts if not explicitly configured in the environment resource.
	} else {
		d.Set("vcs_pr_comments_enabled", environment.VcsPrCommentsEnabled)
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

	if environment.LatestDeploymentLog.WorkflowFile != nil && len(environment.LatestDeploymentLog.WorkflowFile.Environments) > 0 {
		iSubEnvironments, ok := d.GetOk("sub_environment_configuration")

		if ok {
			var newSubEnvironments []interface{}

			for i, iSubEnvironment := range iSubEnvironments.([]interface{}) {
				subEnvironment := iSubEnvironment.(map[string]interface{})

				alias := d.Get(fmt.Sprintf("sub_environment_configuration.%d.alias", i)).(string)

				workkflowSubEnvironment, ok := environment.LatestDeploymentLog.WorkflowFile.Environments[alias]
				if ok {
					subEnvironment["id"] = workkflowSubEnvironment.EnvironmentId
				}

				newSubEnvironments = append(newSubEnvironments, subEnvironment)
			}

			d.Set("sub_environment_configuration", newSubEnvironments)
		}
	}

	setEnvironmentConfigurationSchema(ctx, d, configurationVariables)

	if d.Get("variable_sets") != nil {
		// To avoid drifts keep the schema order as much as possible.
		variableSetsFromSchema := getEnvironmentVariableSetIdsFromSchema(d)
		sortedVariablesSet := []string{}

		for _, schemav := range variableSetsFromSchema {
			for _, newv := range variableSetsIds {
				if schemav == newv {
					sortedVariablesSet = append(sortedVariablesSet, schemav)
					break
				}
			}
		}

		for _, newv := range variableSetsIds {
			found := false
			for _, sortedv := range sortedVariablesSet {
				if newv == sortedv {
					found = true
					break
				}
			}

			if !found {
				sortedVariablesSet = append(sortedVariablesSet, newv)
			}
		}

		d.Set("variable_sets", sortedVariablesSet)
	}

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

func setEnvironmentConfigurationSchema(ctx context.Context, d *schema.ResourceData, configurationVariables []client.ConfigurationVariable) {
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
				newVariable := createVariable(&configurationVariable)

				if configurationVariable.IsSensitive != nil && *configurationVariable.IsSensitive {
					// To avoid drift for sensitive variables, don't override with the variable value received from API. Use the one in the schema instead.
					newVariable.(map[string]interface{})["value"] = variable["value"]
				}

				newVariables = append(newVariables, newVariable)

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
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"configuration name": configurationVariable.Name})
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

	if d.Get("is_inactive").(bool) {
		return diag.Errorf("cannot create an inactive environment (remove 'is_inactive' or set it to 'false')")
	}

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
		// Note: the blueprint id field of the environment is returned only during creation of a template without environment.
		// Afterward, it will be omitted from future response.
		// setEnvironmentSchema() (several lines below) sets the blueprint id in the resource (under "without_template_settings.0.id").
		environment, err = apiClient.EnvironmentCreateWithoutTemplate(payload)
	}

	if err != nil {
		return diag.Errorf("could not create environment: %v", err)
	}

	environmentConfigurationVariables := client.ConfigurationChanges{}
	if environmentPayload.ConfigurationChanges != nil {
		environmentConfigurationVariables = *environmentPayload.ConfigurationChanges
	}

	d.SetId(environment.Id)
	d.Set("deployment_id", environment.LatestDeploymentLogId)
	if environment.AutoDeployOnPathChangesOnly != nil {
		d.Set("auto_deploy_on_path_changes_only", *environment.AutoDeployOnPathChangesOnly)
	}

	var environmentVariableSetIds []string
	if environmentPayload.ConfigurationSetChanges != nil {
		environmentVariableSetIds = environmentPayload.ConfigurationSetChanges.Assign
	}

	if err := setEnvironmentSchema(ctx, d, environment, environmentConfigurationVariables, environmentVariableSetIds); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getEnvironmentVariableSetIdsFromApi(d *schema.ResourceData, apiClient client.ApiClientInterface) ([]string, error) {
	environmentVariableSets, err := apiClient.ConfigurationSetsAssignments("ENVIRONMENT", d.Id())
	if err != nil {
		return nil, err
	}

	var environmentVariableSetIds []string
	for _, variableSet := range environmentVariableSets {
		if variableSet.AssignmentScope == "environment" {
			environmentVariableSetIds = append(environmentVariableSetIds, variableSet.Id)
		}
	}

	return environmentVariableSetIds, nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environment, err := apiClient.Environment(d.Id())
	if err != nil {
		return diag.Errorf("could not get environment: %v", err)
	}

	scope := client.ScopeEnvironment
	if _, ok := d.GetOk("sub_environment_configuration"); ok {
		scope = client.ScopeWorkflow
	}

	environmentConfigurationVariables, err := apiClient.ConfigurationVariablesByScope(scope, environment.Id)
	if err != nil {
		return diag.Errorf("could not get environment configuration variables: %v", err)
	}

	environmentVariableSetIds, err := getEnvironmentVariableSetIdsFromApi(d, apiClient)
	if err != nil {
		return diag.Errorf("could not get environment variable sets: %v", err)
	}

	if err := setEnvironmentSchema(ctx, d, environment, environmentConfigurationVariables, environmentVariableSetIds); err != nil {
		return diag.FromErr(err)
	}

	if isTemplateless(d) {
		// environment with no template.
		templateId := d.Get("without_template_settings.0.id").(string)
		template, err := apiClient.Template(templateId)
		if err != nil {
			return diag.Errorf("could not get template: %v", err)
		}

		if err := templateRead("without_template_settings", template, d, false); err != nil {
			return diag.Errorf("schema resource data serialization failed: %v", err)
		}
	}

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if d.HasChange("project_id") {
		newProjectId := d.Get("project_id").(string)
		if err := apiClient.EnvironmentMove(d.Id(), newProjectId); err != nil {
			return diag.Errorf("failed to move environment to project id '%s': %s", newProjectId, err)
		}
	}

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

	if shouldUpdateDriftDetection(d) {
		if err := updateDriftDetection(d, apiClient); err != nil {
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
	if _, ok := d.GetOk("without_template_settings.0"); ok {
		if d.HasChange("without_template_settings.0.revision") {
			return true
		}
	}

	return d.HasChanges("revision", "configuration", "sub_environment_configuration", "variable_sets")
}

func shouldUpdate(d *schema.ResourceData) bool {
	if d.HasChanges("name", "approve_plan_automatically", "deploy_on_push", "run_plan_on_pull_requests", "auto_deploy_by_custom_glob", "auto_deploy_on_path_changes_only", "vcs_commands_alias", "is_remote_backend", "is_inactive", "is_remote_apply_enabled", "vcs_pr_comments_enabled") {
		return true
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("vcs_pr_comments_enabled"); exists {
		// if this field is set to 'false' will return that there is change each time.
		// this is because the terraform SDK is unable to detecred changes between 'unset' and 'false' (sdk limitation).
		if !val.(bool) {
			return true
		}
	}

	return false
}

func shouldUpdateTTL(d *schema.ResourceData) bool {
	return d.HasChange("ttl")
}

func shouldUpdateDriftDetection(d *schema.ResourceData) bool {
	return d.HasChange("drift_detection_cron")
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

func updateDriftDetection(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	drift_detection_cron, ok := d.GetOk("drift_detection_cron")
	if !ok || drift_detection_cron.(string) == "" {
		if err := apiClient.EnvironmentStopDriftDetection(d.Id()); err != nil {
			return diag.Errorf("could not stop drift detection: %v", err)
		}
	} else {
		if _, err := apiClient.EnvironmentUpdateDriftDetection(d.Id(), client.EnvironmentSchedulingExpression{
			Enabled: true,
			Cron:    drift_detection_cron.(string),
		}); err != nil {
			return diag.Errorf("could not update drift detection: %v", err)
		}
	}

	return nil
}

func deploy(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	deployPayload, err := getDeployPayload(d, apiClient, true)
	if err != nil {
		return diag.FromErr(err)
	}

	subEnvironments, err := getSubEnvironments(d)
	if err != nil {
		return diag.Errorf("failed to extract subenvrionments from resourcedata: %v", err)
	}

	if len(subEnvironments) > 0 {
		deployPayload.SubEnvironments = make(map[string]client.SubEnvironment)

		for i, subEnvironment := range subEnvironments {
			configuration := d.Get(fmt.Sprintf("sub_environment_configuration.%d.configuration", i)).([]interface{})
			configurationChanges := getConfigurationVariablesFromSchema(configuration)
			configurationChanges, err = getUpdateConfigurationVariables(configurationChanges, subEnvironment.Id, client.ScopeEnvironment, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}

			for i := range configurationChanges {
				configurationChanges[i].Scope = client.ScopeEnvironment
			}

			deployPayload.SubEnvironments[subEnvironment.Alias] = client.SubEnvironment{
				Revision:             subEnvironment.Revision,
				Workspace:            subEnvironment.Workspace,
				ConfigurationChanges: configurationChanges,
				UserRequiresApproval: !subEnvironment.ApprovePlanAutomatically,
			}
		}
	}

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

func getEnvironmentVariableSetIdsFromSchema(d *schema.ResourceData) []string {
	var variableSets []string

	if ivariableSets, ok := d.GetOk("variable_sets"); ok {
		for _, ivariableSet := range ivariableSets.([]interface{}) {
			variableSets = append(variableSets, ivariableSet.(string))
		}
	}

	return variableSets
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	markAsArchived := d.Get("removal_strategy").(string) == "mark_as_archived"

	if markAsArchived {
		if err := apiClient.EnvironmentMarkAsArchived(d.Id()); err != nil {
			return diag.Errorf("could not archive the environment: %v", err)
		}

		return nil
	}

	canDestroy := d.Get("force_destroy")

	if canDestroy != true {
		return diag.Errorf(`must enable "force_destroy" safeguard in order to destroy`)
	}

	_, err := apiClient.EnvironmentDestroy(d.Id())
	if err != nil {
		if frerr, ok := err.(*http.FailedResponseError); ok && frerr.BadRequest() {
			tflog.Warn(ctx, "Could not delete environment. Already deleted?", map[string]interface{}{"id": d.Id(), "error": frerr.Error()})
			return nil
		}
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
	if val, exists := d.GetOkExists("vcs_pr_comments_enabled"); exists {
		payload.VcsPrCommentsEnabled = boolPtr(val.(bool))
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
	if val, exists := d.GetOkExists("prevent_auto_deploy"); exists {
		payload.PreventAutoDeploy = boolPtr(val.(bool))
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

	if err := assertEnvironment(d); err != nil {
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

	if drift_detection_cron, ok := d.GetOk("drift_detection_cron"); ok && drift_detection_cron.(string) != "" {
		payload.DriftDetectionRequest = &client.DriftDetectionRequest{
			Enabled: true,
			Cron:    drift_detection_cron.(string),
		}
	}

	variableSets := getEnvironmentVariableSetIdsFromSchema(d)
	if len(variableSets) > 0 {
		payload.ConfigurationSetChanges = &client.ConfigurationSetChanges{
			Assign: variableSets,
		}
	}

	deployPayload, err := getDeployPayload(d, apiClient, false)
	if err != nil {
		return client.EnvironmentCreate{}, diag.FromErr(err)
	}

	subEnvironments, err := getSubEnvironments(d)
	if err != nil {
		return client.EnvironmentCreate{}, diag.Errorf("failed to extract subenvrionments from resourcedata: %v", err)
	}

	if len(subEnvironments) > 0 {
		payload.Type = "workflow"

		deployPayload.SubEnvironments = make(map[string]client.SubEnvironment)

		for _, subEnvironment := range subEnvironments {
			deployPayload.SubEnvironments[subEnvironment.Alias] = client.SubEnvironment{
				Revision:             subEnvironment.Revision,
				ConfigurationChanges: subEnvironment.Configuration,
				Workspace:            subEnvironment.Workspace,
				UserRequiresApproval: !subEnvironment.ApprovePlanAutomatically,
			}
		}
	}

	payload.DeployRequest = &deployPayload

	return payload, nil
}

func assertEnvironment(d *schema.ResourceData) diag.Diagnostics {
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

	isRemoteApplyEnabled := d.Get("is_remote_apply_enabled").(bool)
	isRemotedBackend := d.Get("is_remote_backend").(bool)
	approvePlanAutomatically := d.Get("approve_plan_automatically").(bool)
	if isRemoteApplyEnabled && (!isRemotedBackend || !approvePlanAutomatically) {
		return diag.Errorf("cannot set is_remote_apply_enabled when approve_plan_automatically or is_remote_backend are disabled")
	}

	return nil
}

func getUpdatePayload(d *schema.ResourceData) (client.EnvironmentUpdate, diag.Diagnostics) {
	var payload client.EnvironmentUpdate

	if err := readResourceData(&payload, d); err != nil {
		return client.EnvironmentUpdate{}, diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	// Because the terraform SDK is unable to detecred changes between 'unset' and 'false' (sdk limitation): always set a value here (even if there's no change).
	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("vcs_pr_comments_enabled"); exists {
		payload.VcsPrCommentsEnabled = boolPtr(val.(bool))
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

	if d.HasChange("is_inactive") {
		payload.IsArchived = boolPtr(d.Get("is_inactive").(bool))
	}

	if err := assertEnvironment(d); err != nil {
		return client.EnvironmentUpdate{}, err
	}

	return payload, nil
}

func getEnvironmentConfigurationSetChanges(d *schema.ResourceData, apiClient client.ApiClientInterface) (*client.ConfigurationSetChanges, error) {
	variableSetsFromSchema := getEnvironmentVariableSetIdsFromSchema(d)
	variableSetFromApi, err := getEnvironmentVariableSetIdsFromApi(d, apiClient)
	if err != nil {
		return nil, err
	}

	var assignVariableSets []string
	var unassignVariableSets []string

	for _, sv := range variableSetsFromSchema {
		found := false

		for _, av := range variableSetFromApi {
			if sv == av {
				found = true
				break
			}
		}

		if !found {
			assignVariableSets = append(assignVariableSets, sv)
		}
	}

	for _, av := range variableSetFromApi {
		found := false

		for _, sv := range variableSetsFromSchema {
			if sv == av {
				found = true
				break
			}
		}

		if !found {
			unassignVariableSets = append(unassignVariableSets, av)
		}
	}

	if assignVariableSets == nil && unassignVariableSets == nil {
		return nil, ErrNoChanges
	}

	return &client.ConfigurationSetChanges{
		Assign:   assignVariableSets,
		Unassign: unassignVariableSets,
	}, nil
}

func getDeployPayload(d *schema.ResourceData, apiClient client.ApiClientInterface, isRedeploy bool) (client.DeployRequest, error) {
	payload := client.DeployRequest{}
	var err error

	if isTemplateless(d) {
		if templateId, ok := d.GetOk("without_template_settings.0.id"); ok {
			payload.BlueprintId = templateId.(string)
		}
	} else {
		payload.BlueprintId = d.Get("template_id").(string)
	}

	if revision, ok := d.GetOk("revision"); ok {
		payload.BlueprintRevision = revision.(string)
	}

	// For 'Workflows', the 'root' environment should never require a user approval.
	if _, ok := d.GetOk("sub_environment_configuration"); ok {
		payload.UserRequiresApproval = boolPtr(false)
	}

	if isRedeploy {
		if revision, ok := d.GetOk("without_template_settings.0.revision"); ok {
			payload.BlueprintRevision = revision.(string)
		}

		if configuration, ok := d.GetOk("configuration"); ok && isRedeploy {
			configurationChanges := getConfigurationVariablesFromSchema(configuration.([]interface{}))
			scope := client.ScopeEnvironment
			if _, ok := d.GetOk("sub_environment_configuration"); ok {
				scope = client.ScopeWorkflow
			}
			configurationChanges, err = getUpdateConfigurationVariables(configurationChanges, d.Get("id").(string), scope, apiClient)
			if err != nil {
				return client.DeployRequest{}, err
			}
			payload.ConfigurationChanges = &configurationChanges
		}

		payload.ConfigurationSetChanges, err = getEnvironmentConfigurationSetChanges(d, apiClient)
		if err != nil && !errors.Is(err, ErrNoChanges) {
			return client.DeployRequest{}, err
		}
	}

	//lint:ignore SA1019 reason: https://github.com/hashicorp/terraform-plugin-sdk/issues/817
	if val, exists := d.GetOkExists("approve_plan_automatically"); exists {
		payload.UserRequiresApproval = boolPtr(!val.(bool))
	}

	return payload, nil
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

func getUpdateConfigurationVariables(configurationChanges client.ConfigurationChanges, environmentId string, scope client.Scope, apiClient client.ApiClientInterface) (client.ConfigurationChanges, error) {
	existVariables, err := apiClient.ConfigurationVariablesByScope(scope, environmentId)
	if err != nil {
		return client.ConfigurationChanges{}, fmt.Errorf("could not get environment configuration variables: %w", err)
	}
	configurationChanges = linkToExistConfigurationVariables(configurationChanges, existVariables)
	configurationChanges = deleteUnusedConfigurationVariables(configurationChanges, existVariables)

	return configurationChanges, nil
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
	varType, _ := client.GetConfigurationVariableType(variable["type"].(string))

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

func getEnvironmentByName(meta interface{}, name string, projectId string, excludeArchived bool) (client.Environment, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	environmentsByName, err := apiClient.EnvironmentsByName(name)
	if err != nil {
		return client.Environment{}, diag.Errorf("Could not get Environment: %v", err)
	}

	filteredEnvironments := []client.Environment{}
	for _, candidate := range environmentsByName {
		if excludeArchived && candidate.IsArchived != nil && *candidate.IsArchived {
			continue
		}

		if projectId != "" && candidate.ProjectId != projectId {
			continue
		}

		filteredEnvironments = append(filteredEnvironments, candidate)
	}

	if len(filteredEnvironments) > 1 {
		return client.Environment{}, diag.Errorf("Found multiple environments for name: %s. Use ID instead or make sure environment names are unique %v", name, environmentsByName)
	}

	if len(filteredEnvironments) == 0 {
		return client.Environment{}, diag.Errorf("Could not find an env0 environment with name %s", name)
	}

	return filteredEnvironments[0], nil
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
		tflog.Info(ctx, "Resolving environment by id", map[string]interface{}{"id": id})
		environment, getErr = getEnvironmentById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving environment by name", map[string]interface{}{"name": id})

		environment, getErr = getEnvironmentByName(meta, id, "", false)
	}

	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	}

	apiClient := meta.(client.ApiClientInterface)

	d.SetId(environment.Id)

	scope := client.ScopeEnvironment
	if environment.LatestDeploymentLog.WorkflowFile != nil && len(environment.LatestDeploymentLog.WorkflowFile.Environments) > 0 {
		scope = client.ScopeWorkflow
	}

	environmentConfigurationVariables, err := apiClient.ConfigurationVariablesByScope(scope, environment.Id)
	if err != nil {
		return nil, fmt.Errorf("could not get environment configuration variables: %v", err)
	}

	environmentVariableSetIds, err := getEnvironmentVariableSetIdsFromApi(d, apiClient)
	if err != nil {
		return nil, fmt.Errorf("could not get environment variable sets: %v", err)
	}

	d.Set("deployment_id", environment.LatestDeploymentLogId)

	if environment.IsSingleUseBlueprint {
		templateId := environment.BlueprintId
		if templateId == "" {
			templateId = environment.LatestDeploymentLog.BlueprintId
		}

		template, err := apiClient.Template(templateId)
		if err != nil {
			return nil, fmt.Errorf("failed to get template with id %s: %w", templateId, err)
		}

		if err := templateRead("without_template_settings", template, d, true); err != nil {
			return nil, fmt.Errorf("failed to write template to schema: %w", err)
		}
	}

	if err := setEnvironmentSchema(ctx, d, environment, environmentConfigurationVariables, environmentVariableSetIds); err != nil {
		return nil, err
	}

	if environment.IsRemoteBackend != nil {
		d.Set("is_remote_backend", *environment.IsRemoteBackend)
	}

	if environment.AutoDeployOnPathChangesOnly != nil {
		d.Set("auto_deploy_on_path_changes_only", *environment.AutoDeployOnPathChangesOnly)
	}

	d.Set("is_inactive", false) // default is false.
	if environment.IsArchived != nil {
		d.Set("is_inactive", *environment.IsArchived)
	}

	d.Set("force_destroy", false)
	d.Set("removal_strategy", "destroy")

	d.Set("vcs_pr_comments_enabled", environment.VcsCommandsAlias != "" || environment.VcsPrCommentsEnabled)

	return []*schema.ResourceData{d}, nil
}
