package env0

import (
	"context"
	"errors"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
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

		Schema: getTemplateSchema(TemplateTypeShared),
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
		result.GitlabProjectId = d.Get("gitlab_project_id").(int)
		result.IsGitLab = result.TokenId != ""
	}
	if isGitlabEnterprise, ok := d.GetOk("is_gitlab_enterprise"); ok {
		result.IsGitlabEnterprise = isGitlabEnterprise.(bool)
	}
	if isBitbucketServer, ok := d.GetOk("is_bitbucket_server"); ok {
		result.IsBitbucketServer = isBitbucketServer.(bool)
	}
	if isGitHubEnterprise, ok := d.GetOk("is_github_enterprise"); ok {
		result.IsGitHubEnterprise = isGitHubEnterprise.(bool)
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
	if terragruntVersion, ok := d.GetOk("terragrunt_version"); ok {
		result.TerragruntVersion = terragruntVersion.(string)
	}
	if bitbucketClientKey, ok := d.GetOk("bitbucket_client_key"); ok {
		result.BitbucketClientKey = bitbucketClientKey.(string)
	}

	if terragruntVersion, ok := d.GetOk("terragrunt_version"); ok {
		result.TerragruntVersion = terragruntVersion.(string)
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

	if template.IsDeleted && !d.IsNewResource() {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", template.Name)
	d.Set("description", template.Description)
	if template.GithubInstallationId != 0 {
		d.Set("github_installation_id", template.GithubInstallationId)
	}
	if template.TokenId != "" {
		d.Set("token_id", template.TokenId)
	}
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("terraform_version", template.TerraformVersion)
	// 'gitlab_project_id' should not be set because it doesn't exist on 'template'

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

	if template.BitbucketClientKey != "" {
		d.Set("bitbucket_client_key", template.BitbucketClientKey)
	}

	d.Set("is_gitlab_enterprise", template.IsGitlabEnterprise)
	d.Set("is_bitbucket_server", template.IsBitbucketServer)
	d.Set("is_github_enterprise", template.IsGitHubEnterprise)

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
