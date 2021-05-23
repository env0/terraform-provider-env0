package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		ReadContext:   resourceTemplateRead,
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceTemplateDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTemplateImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the template",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description for the template",
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "git repository for the template source code",
				Required:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "terraform / terragrunt file folder inside source code",
				Optional:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "source code revision (branch / tag) to use",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "'terraform' or 'terragrunt'",
				Optional:    true,
				Default:     "terraform",
			},
			"project_ids": {
				Type:        schema.TypeList,
				Description: "which projects may access this template (id of project)",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "env0_project.id for each project",
				},
			},
			"ssh_key_names": {
				Type:        schema.TypeList,
				Description: "names of env0 defined ssh keys to use when accessing git over ssh",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "env0_ssh_key.name for each ssh key",
				},
			},
			"retries_on_deploy": {
				Type:        schema.TypeInt,
				Description: "number of times to retry when deploying an environment based on this template",
				Optional:    true,
			},
			"retry_on_deploy_only_when_matches_regex": {
				Type:        schema.TypeString,
				Description: "if specified, will only retry (on deploy) if error matches specified regex",
				Optional:    true,
			},
			"retries_on_destroy": {
				Type:        schema.TypeInt,
				Description: "number of times to retry when destroying an environment based on this template",
				Optional:    true,
			},
			"retry_on_destroy_only_when_matches_regex": {
				Type:        schema.TypeString,
				Description: "if specified, will only retry (on destroy) if error matches specified regex",
				Optional:    true,
			},
		},
	}
}

func templateCreatePayloadFromParameters(d *schema.ResourceData) (env0apiclient.TemplateCreatePayload, diag.Diagnostics) {
	result := env0apiclient.TemplateCreatePayload{
		Name:       d.Get("name").(string),
		Repository: d.Get("repository").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		result.Description = description.(string)
	}
	if path, ok := d.GetOk("path"); ok {
		result.Path = path.(string)
	}
	if revision, ok := d.GetOk("revision"); ok {
		result.Revision = revision.(string)
	}
	if type_, ok := d.GetOk("type"); ok {
		if type_ == string(env0apiclient.TemplateTypeTerraform) {
			result.Type = env0apiclient.TemplateTypeTerraform
		} else if type_ == string(env0apiclient.TemplateTypeTerragrunt) {
			result.Type = env0apiclient.TemplateTypeTerragrunt
		} else {
			return env0apiclient.TemplateCreatePayload{}, diag.Errorf("'type' can either be 'terraform' or 'terragrunt': %s", type_)
		}
	}
	if projectIds, ok := d.GetOk("project_ids"); ok {
		result.ProjectIds = []string{}
		for _, projectId := range projectIds.([]interface{}) {
			result.ProjectIds = append(result.ProjectIds, projectId.(string))
		}
	}
	if sshKeyNames, ok := d.GetOk("ssh_key_names"); ok {
		result.SshKeys = []env0apiclient.TemplateSshKey{}
		for _, sshKeyName := range sshKeyNames.([]interface{}) {
			result.SshKeys = append(result.SshKeys, env0apiclient.TemplateSshKey{Name: sshKeyName.(string)})
		}
	}
	onDeployRetries, hasRetriesOnDeploy := d.GetOk("retries_on_deploy")
	if hasRetriesOnDeploy {
		if result.Retry == nil {
			result.Retry = &env0apiclient.TemplateRetry{}
		}
		result.Retry.OnDeploy = &env0apiclient.TemplateRetryOn{
			Times: onDeployRetries.(int),
		}
	}
	if retryOnDeployOnlyIfMatchesRegex, ok := d.GetOk("retry_on_deploy_only_if_matches_regex"); ok {
		if !hasRetriesOnDeploy {
			return env0apiclient.TemplateCreatePayload{}, diag.Errorf("may only specify 'retry_on_deploy_only_if_matches_regex'")
		}
		result.Retry.OnDeploy.ErrorRegex = retryOnDeployOnlyIfMatchesRegex.(string)
	}

	onDestroyRetries, hasRetriesOnDestroy := d.GetOk("retries_on_destroy")
	if hasRetriesOnDestroy {
		if result.Retry == nil {
			result.Retry = &env0apiclient.TemplateRetry{}
		}
		result.Retry.OnDestroy = &env0apiclient.TemplateRetryOn{
			Times: onDestroyRetries.(int),
		}
	}
	if retryOnDestroyOnlyIfMatchesRegex, ok := d.GetOk("retry_on_destroy_only_if_matches_regex"); ok {
		if !hasRetriesOnDestroy {
			return env0apiclient.TemplateCreatePayload{}, diag.Errorf("may only specify 'retry_on_destroy_only_if_matches_regex'")
		}
		result.Retry.OnDestroy.ErrorRegex = retryOnDestroyOnlyIfMatchesRegex.(string)
	}
	return result, nil
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	request, problem := templateCreatePayloadFromParameters(d)
	if problem != nil {
		return problem
	}
	template, err := apiClient.TemplateCreate(request)
	if err != nil {
		return diag.Errorf("could not create template: %v", err)
	}

	d.SetId(template.Id)

	return nil
}

func resourceTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	template, err := apiClient.Template(d.Id())
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}

	d.Set("name", template.Name)
	d.Set("description", template.Description)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("project_ids", template.ProjectIds)
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

	return nil
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	request, problem := templateCreatePayloadFromParameters(d)
	if problem != nil {
		return problem
	}
	_, err := apiClient.TemplateUpdate(d.Id(), request)
	if err != nil {
		return diag.Errorf("could not update template: %v", err)
	}

	return nil
}

func resourceTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	id := d.Id()
	err := apiClient.TemplateDelete(id)
	if err != nil {
		return diag.Errorf("could not delete template: %v", err)
	}
	return nil
}

func resourceTemplateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, errors.New("Not implemented")
	// apiClient := meta.(*env0apiclient.ApiClient)

	// id := d.Id()
	// template, err := apiClient.Template(id)
	// if err != nil {
	// 	return nil, err
	// }

	// d.Set("name", template.Name)

	// return []*schema.ResourceData{d}, nil
}
