package env0

import (
	"context"
	"fmt"
	"strings"

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
				Description: fmt.Sprintf("template type (allowed values: %s)", strings.Join(allowedTemplateTypes, ", ")),
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
	var template client.Template
	var err diag.Diagnostics

	if name, ok := d.GetOk("name"); ok {
		if template, err = getTemplateByName(name, meta); err != nil {
			return err
		}
	} else if template, err = getTemplateById(d.Get("id").(string), meta); err != nil {
		return diag.Errorf("Could not query template: %v", err)
	}

	if err := writeResourceData(&template, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	templateReadRetryOnHelper("", d, "deploy", template.Retry.OnDeploy)
	templateReadRetryOnHelper("", d, "destroy", template.Retry.OnDestroy)

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
