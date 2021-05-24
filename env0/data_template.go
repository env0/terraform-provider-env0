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
			"ssh_keys": {
				Type:        schema.TypeList,
				Description: "which ssh keys are used for accessing git over ssh",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "a map of env0_ssh_key.id and env0_ssh_key.name for each project",
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
		},
	}
}

func dataTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	name, nameSpecified := d.GetOk("name")
	var template client.Template
	var err error
	if nameSpecified {
		templates, err := apiClient.Templates()
		if err != nil {
			return diag.Errorf("Could not query templates: %v", err)
		}
		for _, candidate := range templates {
			if candidate.Name == name {
				template = candidate
			}
		}
		if template.Name == "" {
			return diag.Errorf("Could not find an env0 template with name %s", name)
		}
	} else {
		template, err = apiClient.Template(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Could not query template: %v", err)
		}
	}

	d.SetId(template.Id)
	d.Set("name", template.Name)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("project_ids", template.ProjectIds)
	d.Set("ssh_keys", template.SshKeys)
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

	//TODO: sshkeys
	return nil
}
