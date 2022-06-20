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

func templateCreatePayloadRetryOnHelper(d *schema.ResourceData, retryType string, retryOnPtr **client.TemplateRetryOn) {
	retries, hasRetries := d.GetOk("retries_on_" + retryType)
	if hasRetries {
		retryOn := &client.TemplateRetryOn{
			Times: retries.(int),
		}
		if retryIfMatchesRegex, ok := d.GetOk("retry_on_" + retryType + "_only_when_matches_regex"); ok {
			retryOn.ErrorRegex = retryIfMatchesRegex.(string)
		}

		*retryOnPtr = retryOn
	}
}

func templateCreatePayloadFromParameters(d *schema.ResourceData) (client.TemplateCreatePayload, diag.Diagnostics) {
	var payload client.TemplateCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return payload, diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if tokenId, ok := d.GetOk("token_id"); ok {
		payload.IsGitLab = tokenId != ""
	}

	templateCreatePayloadRetryOnHelper(d, "deploy", &payload.Retry.OnDeploy)
	templateCreatePayloadRetryOnHelper(d, "destroy", &payload.Retry.OnDestroy)

	return payload, nil
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
	d.Set("is_github_enterprise", template.IsGithubEnterprise)

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
