package env0

import (
	"context"

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
		},
	}
}

func dataEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	setEnvironmentSchema(ctx, d, environment, client.ConfigurationChanges{})

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
