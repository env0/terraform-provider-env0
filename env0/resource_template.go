package env0

import (
	"context"
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceTemplate() *schema.Resource {
	validateRetries := func(i interface{}, path cty.Path) diag.Diagnostics {
		retries := i.(int)
		if retries < 1 || retries > 3 {
			return diag.Errorf("Retries amount must be between 1 and 3")
		}

		return nil
	}

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
			"ssh_keys": {
				Type:        schema.TypeList,
				Description: "an array of references to 'data_ssh_key' to use when accessing git over ssh",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "a map of env0_ssh_key.id and env0_ssh_key.name for each project",
				},
			},
			"retries_on_deploy": {
				Type:             schema.TypeInt,
				Description:      "number of times to retry when deploying an environment based on this template",
				Optional:         true,
				ValidateDiagFunc: validateRetries,
			},
			"retry_on_deploy_only_when_matches_regex": {
				Type:         schema.TypeString,
				Description:  "if specified, will only retry (on deploy) if error matches specified regex",
				Optional:     true,
				RequiredWith: []string{"retries_on_deploy"},
			},
			"retries_on_destroy": {
				Type:             schema.TypeInt,
				Description:      "number of times to retry when destroying an environment based on this template",
				Optional:         true,
				ValidateDiagFunc: validateRetries,
			},
			"retry_on_destroy_only_when_matches_regex": {
				Type:         schema.TypeString,
				Description:  "if specified, will only retry (on destroy) if error matches specified regex",
				Optional:     true,
				RequiredWith: []string{"retries_on_destroy"},
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "The env0 application installation id on the relevant github repository",
				Optional:    true,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "The env0 application token id on the relevant gitlab account",
				Optional:    true,
			},
			"gitlab_project_id": {
				Type:        schema.TypeInt,
				Description: "The project id of the relevant repository",
				Optional:    true,
			},
			"terraform_version": {
				Type:        schema.TypeString,
				Description: "Terraform version to use",
				Optional:    true,
				Default:     "0.15.1",
			},
		},
	}
}

func templateCreatePayloadFromParameters(d *schema.ResourceData) (client.TemplateCreatePayload, diag.Diagnostics) {
	result := client.TemplateCreatePayload{
		Name:       d.Get("name").(string),
		Repository: d.Get("repository").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		result.Description = description.(string)
	}
	if githubInstallationId, ok := d.GetOk("github_installation_id"); ok {
		result.GithubInstallationId = githubInstallationId.(int)
	}
	if tokenId, ok := d.GetOk("token_id"); ok {
		result.TokenId = tokenId.(string)
	}
	if gitlabProjectId, ok := d.GetOk("gitlab_project_id"); ok {
		result.GitlabProjectId = gitlabProjectId.(int)
	}

	if result.GithubInstallationId != 0 && result.TokenId != "" {
		return client.TemplateCreatePayload{}, diag.Errorf("Cannot set token_id and github_installation_id for the same template")
	} else {
		result.IsGitLab = result.TokenId != ""
	}

	if path, ok := d.GetOk("path"); ok {
		result.Path = path.(string)
	}
	if revision, ok := d.GetOk("revision"); ok {
		result.Revision = revision.(string)
	}
	if type_, ok := d.GetOk("type"); ok {
		if type_ == string(client.TemplateTypeTerraform) {
			result.Type = client.TemplateTypeTerraform
		} else if type_ == string(client.TemplateTypeTerragrunt) {
			result.Type = client.TemplateTypeTerragrunt
		} else {
			return client.TemplateCreatePayload{}, diag.Errorf("'type' can either be 'terraform' or 'terragrunt': %s", type_)
		}
	}
	if sshKeys, ok := d.GetOk("ssh_keys"); ok {
		result.SshKeys = []client.TemplateSshKey{}
		for _, sshKey := range sshKeys.([]interface{}) {
			result.SshKeys = append(result.SshKeys, client.TemplateSshKey{
				Name: sshKey.(map[string]interface{})["name"].(string),
				Id:   sshKey.(map[string]interface{})["id"].(string)})
		}
	}
	onDeployRetries, hasRetriesOnDeploy := d.GetOk("retries_on_deploy")
	var onDeploy *client.TemplateRetryOn = nil
	var onDestroy *client.TemplateRetryOn = nil
	if hasRetriesOnDeploy {
		onDeploy = &client.TemplateRetryOn{
			Times: onDeployRetries.(int),
		}
	}
	if retryOnDeployOnlyIfMatchesRegex, ok := d.GetOk("retry_on_deploy_only_when_matches_regex"); ok {
		if !hasRetriesOnDeploy {
			return client.TemplateCreatePayload{}, diag.Errorf("may only specify 'retry_on_deploy_only_when_matches_regex'")
		}
		onDeploy.ErrorRegex = retryOnDeployOnlyIfMatchesRegex.(string)
	}

	onDestroyRetries, hasRetriesOnDestroy := d.GetOk("retries_on_destroy")
	if hasRetriesOnDestroy {
		onDestroy = &client.TemplateRetryOn{
			Times: onDestroyRetries.(int),
		}
	}
	if retryOnDestroyOnlyIfMatchesRegex, ok := d.GetOk("retry_on_destroy_only_when_matches_regex"); ok {
		if !hasRetriesOnDestroy {
			return client.TemplateCreatePayload{}, diag.Errorf("may only specify 'retry_on_destroy_only_when_matches_regex'")
		}
		onDestroy.ErrorRegex = retryOnDestroyOnlyIfMatchesRegex.(string)
	}

	if onDeploy != nil || onDestroy != nil {
		result.Retry = client.TemplateRetry{}
		if onDeploy != nil {
			result.Retry.OnDeploy = onDeploy
		}
		if onDestroy != nil {
			result.Retry.OnDestroy = onDestroy
		}
	}

	if terraformVersion, ok := d.GetOk("terraform_version"); ok {
		result.TerraformVersion = terraformVersion.(string)
	}
	return result, nil
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

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
	apiClient := meta.(client.ApiClientInterface)

	template, err := apiClient.Template(d.Id())
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}

	d.Set("name", template.Name)
	d.Set("description", template.Description)
	d.Set("github_installation_id", template.GithubInstallationId)
	d.Set("token_id", template.TokenId)
	d.Set("gitlab_project_id", template.GitlabProjectId)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("terraform_version", template.TerraformVersion)

	var rawSshKeys []map[string]string
	for _, sshKey := range template.SshKeys {
		rawSshKeys = append(rawSshKeys, map[string]string{"id": sshKey.Id, "name": sshKey.Name})
	}
	d.Set("ssh_keys", rawSshKeys)

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
	apiClient := meta.(client.ApiClientInterface)

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
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.TemplateDelete(id)
	if err != nil {
		return diag.Errorf("could not delete template: %v", err)
	}
	return nil
}

func resourceTemplateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving Template by id: ", id)
		_, getErr = getTemplateById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving Template by name: ", id)
		var template client.Template
		template, getErr = getTemplateByName(id, meta)
		d.SetId(template.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
