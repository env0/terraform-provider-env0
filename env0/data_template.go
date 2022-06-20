package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the template",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the template",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "template source code repository url",
				Computed:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "terraform / terrgrunt folder inside source code repository",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "source code revision (branch / tag) to use",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "'terraform' or 'terragrunt'",
				Computed:    true,
			},
			"project_ids": {
				Type:        schema.TypeList,
				Description: "which projects may access this template (id of project)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "env0_project.id for each project",
				},
			},
			"retries_on_deploy": {
				Type:        schema.TypeInt,
				Description: "number of times to retry when deploying an environment based on this template",
				Computed:    true,
			},
			"retry_on_deploy_only_when_matches_regex": {
				Type:        schema.TypeString,
				Description: "if specified, will only retry (on deploy) if error matches specified regex",
				Computed:    true,
			},
			"retries_on_destroy": {
				Type:        schema.TypeInt,
				Description: "number of times to retry when destroying an environment based on this template",
				Computed:    true,
			},
			"retry_on_destroy_only_when_matches_regex": {
				Type:        schema.TypeString,
				Description: "if specified, will only retry (on destroy) if error matches specified regex",
				Computed:    true,
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "The env0 application installation id on the relevant github repository",
				Optional:    true,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "The token id used for private git repos or for integration with GitLab",
				Optional:    true,
			},
			"gitlab_project_id": {
				Type:        schema.TypeInt,
				Description: "The project id of the relevant repository",
				Optional:    true,
			},
			"terraform_version": {
				Type:        schema.TypeString,
				Description: "terraform version to use",
				Computed:    true,
			},
			"terragrunt_version": {
				Type:        schema.TypeString,
				Description: "terragrunt version to use",
				Computed:    true,
				Optional:    true,
			},
			"is_gitlab_enterprise": {
				Type:        schema.TypeBool,
				Description: "Does this template use gitlab enterprise repository?",
				Optional:    true,
				Computed:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "the bitbucket client key used for integration",
				Optional:    true,
				Computed:    true,
			},
			"is_bitbucket_server": {
				Type:        schema.TypeBool,
				Description: "true if this template uses bitbucket server repository",
				Optional:    true,
				Computed:    true,
			},
			"is_github_enterprise": {
				Type:        schema.TypeBool,
				Description: "true if this template uses github enterprise repository",
				Optional:    true,
				Computed:    true,
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Description: "an array of references to 'data_ssh_key' to use when accessing git over ssh",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ssh key id",
							Required:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "ssh key name",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func dataTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name, nameSpecified := d.GetOk("name")
	var template client.Template
	var err diag.Diagnostics
	if nameSpecified {
		template, err = getTemplateByName(name, meta)
		if err != nil {
			return err
		}
	} else {
		template, err = getTemplateById(d.Get("id").(string), meta)
		if err != nil {
			return diag.Errorf("Could not query template: %v", err)
		}
	}

	// TODO: use writeResourceData instead.
	d.SetId(template.Id)
	d.Set("name", template.Name)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("project_ids", template.ProjectIds)
	d.Set("terraform_version", template.TerraformVersion)
	if template.TerragruntVersion != "" {
		d.Set("terragrunt_version", template.TerragruntVersion)
	}
	if template.Retry.OnDeploy != nil {
		d.Set("retries_on_deploy", template.Retry.OnDeploy.Times)
		d.Set("retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex)
	} else {
		d.Set("retries_on_deploy", 0)
		d.Set("retry_on_deploy_only_when_matches_regex", "")
	}
	if template.Retry.OnDestroy != nil {
		d.Set("retries_on_destroy", template.Retry.OnDestroy.Times)
		d.Set("retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex)
	} else {
		d.Set("retries_on_destroy", 0)
		d.Set("retry_on_destroy_only_when_matches_regex", "")
	}

	if template.GithubInstallationId != 0 {
		d.Set("github_installation_id", template.GithubInstallationId)
	}

	if template.TokenId != "" {
		d.Set("token_id", template.TokenId)
	}

	if template.BitbucketClientKey != "" {
		d.Set("bitbucket_client_key", template.BitbucketClientKey)
	}

	if template.GitlabProjectId != 0 {
		d.Set("gitlab_project_id", template.GitlabProjectId)
	}

	if template.IsGitlabEnterprise {
		d.Set("is_gitlab_enterprise", template.IsGitlabEnterprise)
	}

	if template.IsBitbucketServer {
		d.Set("is_bitbucket_server", template.IsBitbucketServer)
	}

	if template.IsGithubEnterprise {
		d.Set("is_github_enterprise", template.IsGithubEnterprise)
	}

	var sshKeys []interface{}

	for _, sshKey := range template.SshKeys {
		newSshKey := make(map[string]interface{})
		newSshKey["id"] = sshKey.Id
		newSshKey["name"] = sshKey.Name
		sshKeys = append(sshKeys, newSshKey)
	}

	if sshKeys != nil {
		d.Set("ssh_keys", sshKeys)
	}

	return nil
}

func getTemplateByName(name interface{}, meta interface{}) (client.Template, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	templates, err := apiClient.Templates()
	if err != nil {
		return client.Template{}, diag.Errorf("Could not query templates: %v", err)
	}

	var templatesByName []client.Template
	for _, candidate := range templates {
		if candidate.Name == name && !candidate.IsDeleted {
			templatesByName = append(templatesByName, candidate)
		}
	}

	if len(templatesByName) > 1 {
		return client.Template{}, diag.Errorf("Found multiple Templates for name: %s. Use ID instead or make sure Template names are unique %v", name, templatesByName)
	}

	if len(templatesByName) == 0 {
		return client.Template{}, diag.Errorf("Could not find an env0 template with name %s", name)
	}

	return templatesByName[0], nil
}

func getTemplateById(id interface{}, meta interface{}) (client.Template, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	template, err := apiClient.Template(id.(string))
	if err != nil {
		return client.Template{}, diag.Errorf("Could not query template: %v", err)
	}
	return template, nil
}
