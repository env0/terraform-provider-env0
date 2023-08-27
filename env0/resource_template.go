package env0

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allowedTemplateTypes = []string{
	"terraform",
	"terragrunt",
	"pulumi",
	"k8s",
	"workflow",
	"cloudformation",
	"helm",
}

func getTemplateSchema(prefix string) map[string]*schema.Schema {
	var allVCSAttributes = []string{
		"token_id",
		"gitlab_project_id",
		"github_installation_id",
		"bitbucket_client_key",
		"is_gitlab_enterprise",
		"is_bitbucket_server",
		"is_github_enterprise",
		"is_azure_devops",
		"helm_chart_name",
		"is_helm_repository",
		"path",
	}

	allVCSAttributesBut := func(strs ...string) []string {
		butAttrs := []string{}

		for _, attr := range allVCSAttributes {
			var found bool
			for _, str := range strs {
				if str == attr {
					found = true
					break
				}
			}

			if !found {
				if prefix != "" {
					attr = prefix + attr
				}
				butAttrs = append(butAttrs, attr)
			}
		}

		return butAttrs
	}

	requiredWith := func(strs ...string) []string {
		ret := []string{}

		for _, str := range strs {
			if prefix != "" {
				str = prefix + str
			}
			ret = append(ret, str)
		}

		return ret
	}

	templateSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the template",
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "description for the template",
			Optional:    true,
		},
		"repository": {
			Type:        schema.TypeString,
			Description: "git repository url for the template source code",
			Required:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "terraform / terragrunt file folder inside source code",
			Optional:    true,
		},
		"type": {
			Type:             schema.TypeString,
			Description:      fmt.Sprintf("template type (allowed values: %s)", strings.Join(allowedTemplateTypes, ", ")),
			Optional:         true,
			Default:          "terraform",
			ValidateDiagFunc: NewStringInValidator(allowedTemplateTypes),
		},
		"revision": {
			Type:        schema.TypeString,
			Description: "source code revision (branch / tag) to use",
			Optional:    true,
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
			ValidateDiagFunc: ValidateRetries,
		},
		"retry_on_deploy_only_when_matches_regex": {
			Type:         schema.TypeString,
			Description:  "if specified, will only retry (on deploy) if error matches specified regex",
			Optional:     true,
			RequiredWith: requiredWith("retries_on_deploy"),
		},
		"retries_on_destroy": {
			Type:             schema.TypeInt,
			Description:      "number of times to retry when destroying an environment based on this template",
			Optional:         true,
			ValidateDiagFunc: ValidateRetries,
		},
		"retry_on_destroy_only_when_matches_regex": {
			Type:         schema.TypeString,
			Description:  "if specified, will only retry (on destroy) if error matches specified regex",
			Optional:     true,
			RequiredWith: requiredWith("retries_on_destroy"),
		},
		"github_installation_id": {
			Type:          schema.TypeInt,
			Description:   "the env0 application installation id on the relevant github repository",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("github_installation_id", "path"),
		},
		"token_id": {
			Type:          schema.TypeString,
			Description:   "the git token id to be used",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("token_id", "gitlab_project_id", "is_azure_devops", "path"),
		},
		"gitlab_project_id": {
			Type:          schema.TypeInt,
			Description:   "the project id of the relevant repository",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("token_id", "gitlab_project_id", "path"),
			RequiredWith:  requiredWith("token_id"),
		},
		"terraform_version": {
			Type:             schema.TypeString,
			Description:      "the Terraform version to use (example: 0.15.1). Setting to `RESOLVE_FROM_TERRAFORM_CODE` defaults to the version of `terraform.required_version` during run-time (resolve from terraform code). Setting to `latest`, the version used will be the most recent one available for Terraform.",
			Optional:         true,
			ValidateDiagFunc: NewRegexValidator(`^(?:[0-9]\.[0-9]{1,2}\.[0-9]{1,2})|RESOLVE_FROM_TERRAFORM_CODE|latest$`),
			Default:          "0.15.1",
		},
		"terragrunt_version": {
			Type:             schema.TypeString,
			Description:      "the Terragrunt version to use (example: 0.36.5)",
			ValidateDiagFunc: NewRegexValidator(`^[0-9]\.[0-9]{1,2}\.[0-9]{1,2}$`),
			Optional:         true,
		},
		"is_gitlab_enterprise": {
			Type:          schema.TypeBool,
			Description:   "true if this template uses gitlab enterprise repository",
			Optional:      true,
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_gitlab_enterprise", "path"),
		},
		"bitbucket_client_key": {
			Type:          schema.TypeString,
			Description:   "the bitbucket client key used for integration",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("bitbucket_client_key", "path"),
		},
		"is_bitbucket_server": {
			Type:          schema.TypeBool,
			Description:   "true if this template uses bitbucket server repository",
			Optional:      true,
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_bitbucket_server", "path"),
		},
		"is_github_enterprise": {
			Type:          schema.TypeBool,
			Description:   "true if this template uses github enterprise repository",
			Optional:      true,
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_github_enterprise", "path"),
		},
		"file_name": {
			Type:        schema.TypeString,
			Description: "the cloudformation file name. Required if the template type is cloudformation",
			Optional:    true,
		},
		"is_terragrunt_run_all": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: `true if this template should execute run-all commands on multiple modules (check https://terragrunt.gruntwork.io/docs/features/execute-terraform-commands-on-multiple-modules-at-once/#the-run-all-command for additional details). Can only be true with "terragrunt" template type and terragrunt version 0.28.1 and above`,
			Default:     "false",
		},
		"is_azure_devops": {
			Type:          schema.TypeBool,
			Optional:      true,
			Description:   "true if this template integrates with azure dev ops",
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_azure_devops", "token_id", "path"),
			RequiredWith:  requiredWith("token_id"),
		},
		"helm_chart_name": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "the helm chart name. Required if is_helm_repository is set to 'true'",
			ConflictsWith: allVCSAttributesBut("helm_chart_name", "is_helm_repository"),
		},
		"is_helm_repository": {
			Type:          schema.TypeBool,
			Optional:      true,
			Description:   "true if this template integrates with a helm repository",
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("helm_chart_name", "is_helm_repository"),
			RequiredWith:  requiredWith("helm_chart_name"),
		},
	}

	if prefix == "" {
		templateSchema["name"] = &schema.Schema{
			Type:        schema.TypeString,
			Description: "name to give the template",
			Required:    true,
		}
	}

	return templateSchema
}

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		ReadContext:   resourceTemplateRead,
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceTemplateDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTemplateImport},

		Schema: getTemplateSchema(""),
	}
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request, problem := templateCreatePayloadFromParameters("", d)
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
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	if err := templateRead("", template, d); err != nil {
		return diag.Errorf("%v", err)
	}

	return nil
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request, problem := templateCreatePayloadFromParameters("", d)
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
		tflog.Info(ctx, "Resolving template by id", map[string]interface{}{"id": id})
		_, getErr = getTemplateById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving template by name", map[string]interface{}{"name": id})
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

func templateCreatePayloadRetryOnHelper(prefix string, d *schema.ResourceData, retryType string, retryOnPtr **client.TemplateRetryOn) {
	if prefix != "" {
		prefix += "."
	}

	retries, hasRetries := d.GetOk(prefix + "retries_on_" + retryType)
	if hasRetries {
		retryOn := &client.TemplateRetryOn{
			Times: retries.(int),
		}
		if retryIfMatchesRegex, ok := d.GetOk(prefix + "retry_on_" + retryType + "_only_when_matches_regex"); ok {
			retryOn.ErrorRegex = retryIfMatchesRegex.(string)
		}

		*retryOnPtr = retryOn
	}
}

func templateCreatePayloadFromParameters(prefix string, d *schema.ResourceData) (client.TemplateCreatePayload, diag.Diagnostics) {
	var payload client.TemplateCreatePayload
	if err := readResourceDataEx(prefix, &payload, d); err != nil {
		return payload, diag.Errorf("schema resource data serialization failed: %v", err)
	}

	tokenIdKey := "token_id"
	isAzureDevOpsKey := "is_azure_devops"

	if prefix != "" {
		tokenIdKey = prefix + "." + tokenIdKey
		isAzureDevOpsKey = prefix + "." + isAzureDevOpsKey
	}

	// IsGitLab is implicitly assumed to be true if tokenId is non-empty. Unless AzureDevOps is explicitly used.
	if tokenId, ok := d.GetOk(tokenIdKey); ok {
		if isAzureDevOps := d.Get(isAzureDevOpsKey); !isAzureDevOps.(bool) {
			payload.IsGitLab = tokenId != ""
		}
	}

	templateCreatePayloadRetryOnHelper(prefix, d, "deploy", &payload.Retry.OnDeploy)
	templateCreatePayloadRetryOnHelper(prefix, d, "destroy", &payload.Retry.OnDestroy)

	if err := payload.Validate(); err != nil {
		return payload, diag.Errorf(err.Error())
	}

	return payload, nil
}

// Reads template and writes to the resource data.
func templateRead(prefix string, template client.Template, d *schema.ResourceData) error {
	pathPrefix := "path"
	if prefix != "" {
		pathPrefix = prefix + ".0.path"
	}

	path, pathOk := d.GetOk(pathPrefix)

	if err := writeResourceDataEx(prefix, &template, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	// https://github.com/env0/terraform-provider-env0/issues/699 - backend removes the "/".
	if pathOk && path.(string) == "/"+template.Path {
		d.Set(pathPrefix, path.(string))
	}

	templateReadRetryOnHelper(prefix, d, "deploy", template.Retry.OnDeploy)
	templateReadRetryOnHelper(prefix, d, "destroy", template.Retry.OnDestroy)

	return nil
}

// Helpers function for templateRead.
func templateReadRetryOnHelper(prefix string, d *schema.ResourceData, retryType string, retryOn *client.TemplateRetryOn) {
	if prefix != "" {
		value := d.Get(prefix + ".0").(map[string]interface{})
		if retryOn != nil {
			value["retries_on_"+retryType] = retryOn.Times
			value["retry_on_"+retryType+"_only_when_matches_regex"] = retryOn.ErrorRegex
		} else {
			value["retries_on_"+retryType] = 0
			value["retry_on_"+retryType+"_only_when_matches_regex"] = ""
		}
		d.Set(prefix, []interface{}{value})
	} else {
		if retryOn != nil {
			d.Set("retries_on_"+retryType, retryOn.Times)
			d.Set("retry_on_"+retryType+"_only_when_matches_regex", retryOn.ErrorRegex)
		} else {
			d.Set("retries_on_"+retryType, 0)
			d.Set("retry_on_"+retryType+"_only_when_matches_regex", "")
		}
	}
}
