package env0

import (
	"context"
	"slices"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "the environment's id",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "name of the environment",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"exclude_archived": {
				Type:          schema.TypeBool,
				Description:   "set to 'true' to exclude archived environments when getting an environment by name",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"id"},
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "project id of the environment",
				Computed:    true,
				Optional:    true,
			},
			"approve_plan_automatically": {
				Type:        schema.TypeBool,
				Description: "the default require approval of the environment",
				Computed:    true,
			},
			"run_plan_on_pull_requests": {
				Type:        schema.TypeBool,
				Description: "does pr plan enable",
				Computed:    true,
			},
			"auto_deploy_on_path_changes_only": {
				Type:        schema.TypeBool,
				Description: "does continuous deployment on file changes in path enable",
				Computed:    true,
			},
			"deploy_on_push": {
				Type:        schema.TypeBool,
				Description: "does continuous deployment is enabled",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "the status of the environment",
				Computed:    true,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Description: "the id of the latest deployment",
				Computed:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the template id the environment is to be created from",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "the last deployed revision",
				Computed:    true,
			},
			"output": {
				Type:        schema.TypeString,
				Description: "the deployment log output. Returns a json string. It can be either a map of key-value, or an array of (in case of Terragrunt run-all) of moduleName and a map of key-value. Note: if the deployment is still in progress returns 'null'",
				Computed:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "Bitbucket client key",
				Computed:    true,
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "Github installation id",
				Computed:    true,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "The token id used for repo integrations (Used by Gitlab or Azure DevOps)",
				Computed:    true,
			},
			"sub_environment_configuration": {
				Type:        schema.TypeList,
				Description: "the sub environments of the workflow environment. (Empty for non workflow environments)",
				Computed:    true,
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
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err diag.Diagnostics

	var environment client.Environment

	projectId := d.Get("project_id").(string)

	id, ok := d.GetOk("id")
	if ok {
		environment, err = getEnvironmentById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name := d.Get("name").(string)
		excludeArchived := d.Get("exclude_archived")

		environment, err = getEnvironmentByName(meta, name, projectId, excludeArchived.(bool))
		if err != nil {
			return err
		}
	}

	if err := setEnvironmentSchema(ctx, d, environment, client.ConfigurationChanges{}, nil); err != nil {
		return diag.FromErr(err)
	}

	// Set this explicitly because these are not set in setEnvironmentSchema due to possible project inheritance.

	if environment.AutoDeployOnPathChangesOnly != nil {
		d.Set("auto_deploy_on_path_changes_only", *environment.AutoDeployOnPathChangesOnly)
	}

	if environment.ContinuousDeployment != nil {
		d.Set("deploy_on_push", *environment.ContinuousDeployment)
	}

	if environment.PullRequestPlanDeployments != nil {
		d.Set("run_plan_on_pull_requests", *environment.PullRequestPlanDeployments)
	}

	if environment.RequiresApproval != nil {
		d.Set("approve_plan_automatically", !*environment.RequiresApproval)
	}

	if environment.LifespanEndAt != "" {
		d.Set("ttl", environment.LifespanEndAt)
	}

	subEnvironments := []any{}

	if environment.LatestDeploymentLog.WorkflowFile != nil {
		for alias, subenv := range environment.LatestDeploymentLog.WorkflowFile.Environments {
			subEnvironment := map[string]any{
				"id":    subenv.EnvironmentId,
				"alias": alias,
			}
			subEnvironments = append(subEnvironments, subEnvironment)
		}
	}

	slices.SortFunc(subEnvironments, func(a, b any) int {
		amap := a.(map[string]any)
		bmap := b.(map[string]any)

		aalias := amap["alias"].(string)
		balias := bmap["alias"].(string)

		return strings.Compare(aalias, balias)
	})

	d.Set("sub_environment_configuration", subEnvironments)

	templateId := environment.LatestDeploymentLog.BlueprintId

	template, err := getTemplateById(templateId, meta)
	if err != nil {
		return err
	}

	templateUpdater := struct {
		GithubInstallationId int    `tfschema:",omitempty"`
		TokenId              string `tfschema:",omitempty"`
		BitbucketClientKey   string `tfschema:",omitempty"`
	}{
		GithubInstallationId: template.GithubInstallationId,
		TokenId:              template.TokenId,
		BitbucketClientKey:   template.BitbucketClientKey,
	}

	if err := writeResourceData(&templateUpdater, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
