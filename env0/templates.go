package env0

import (
	"fmt"
	"sort"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TemplateType int

const (
	TemplateTypeSingle = 1
	TemplateTypeShared = 2
)

var allowedTemplateTypes = []string{
	"terraform",
	"terragrunt",
	"pulumi",
	"k8s",
	"workflow",
	"cloudformation",
}

func getTemplateSchema(templateType TemplateType) map[string]*schema.Schema {
	/*
	 *	VCS Constraints:
	 *		GitHub - githubInstallationId
	 *		GitLab -  gitlabProjectId and tokenId
	 *		Bitbucket - bitbucketClientKey
	 *		GH Enterprise - isGitHubEnterprise
	 *		GL Enterprise - isGitLabEnterprise
	 *		BB Server - isBitbucketServer
	 *		Other - tokenId (optional field)
	 */

	allVCSAttributes := []string{
		"token_id",
		"gitlab_project_id",
		"github_installation_id",
		"bitbucket_client_key",
		"is_gitlab_enterprise",
		"is_bitbucket_server",
		"is_github_enterprise",
	}

	appendByType := func(strs []string, str string) []string {
		if templateType == TemplateTypeShared {
			return append(strs, str)
		} else {
			return append(strs, "template.0."+str)
		}
	}

	allVCSAttributesBut := func(strs ...string) []string {
		sort.Strings(strs)
		butAttrs := []string{}

		for _, attr := range allVCSAttributes {
			if sort.SearchStrings(strs, attr) >= len(strs) {
				butAttrs = appendByType(butAttrs, attr)
			}
		}

		return butAttrs
	}

	requiredWith := func(strs ...string) []string {
		attrs := []string{}
		for _, str := range strs {
			attrs = appendByType(attrs, str)
		}

		return attrs
	}

	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "name to give the template",
			Required:    templateType == TemplateTypeShared,
			Computed:    templateType == TemplateTypeSingle,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "description for the template",
			Optional:    true,
			Computed:    templateType == TemplateTypeSingle,
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
			ConflictsWith: allVCSAttributesBut("github_installation_id"),
		},
		"token_id": {
			Type:          schema.TypeString,
			Description:   "the token id used for private git repos or for integration with GitLab, you can get this value by using a data resource of an existing Gitlab template or contact our support team",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("token_id", "gitlab_project_id"),
		},
		"gitlab_project_id": {
			Type:          schema.TypeInt,
			Description:   "the project id of the relevant repository",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("token_id", "gitlab_project_id"),
			RequiredWith:  requiredWith("token_id"),
		},
		"terraform_version": {
			Type:             schema.TypeString,
			Description:      "the Terraform version to use (example: 0.15.1). Setting to `RESOLVE_FROM_TERRAFORM_CODE` defaults to the version of `terraform.required_version` during run-time (resolve from terraform code).",
			Optional:         true,
			ValidateDiagFunc: NewRegexValidator(`^(?:[0-9]\.[0-9]{1,2}\.[0-9]{1,2})|RESOLVE_FROM_TERRAFORM_CODE$`),
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
			ConflictsWith: allVCSAttributesBut("is_gitlab_enterprise"),
		},
		"bitbucket_client_key": {
			Type:          schema.TypeString,
			Description:   "the bitbucket client key used for integration",
			Optional:      true,
			ConflictsWith: allVCSAttributesBut("bitbucket_client_key"),
		},
		"is_bitbucket_server": {
			Type:          schema.TypeBool,
			Description:   "true if this template uses bitbucket server repository",
			Optional:      true,
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_bitbucket_server"),
		},
		"is_github_enterprise": {
			Type:          schema.TypeBool,
			Description:   "true if this template uses github enterprise repository",
			Optional:      true,
			Default:       "false",
			ConflictsWith: allVCSAttributesBut("is_github_enterprise"),
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
	}

	if templateType == TemplateTypeSingle {
		s["id"] = &schema.Schema{
			Type:        schema.TypeString,
			Description: "id of the template",
			Computed:    true,
		}
	}

	return s
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
	if prefix != "" {
		tokenIdKey = prefix + "." + tokenIdKey
	}
	if tokenId, ok := d.GetOk(tokenIdKey); ok {
		payload.IsGitLab = tokenId != ""
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
	if prefix != "" {
		templates := []client.Template{template}
		if err := writeResourceDataSlice(templates, prefix, d); err != nil {
			return fmt.Errorf("schema resource data serialization failed: %v", err)
		}
	} else if err := writeResourceData(&template, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	templateReadRetryOnHelper(prefix, d, "deploy", template.Retry.OnDeploy)
	templateReadRetryOnHelper(prefix, d, "destroy", template.Retry.OnDestroy)

	return nil
}

// Helpers function for templateRead.
func templateReadRetryOnHelper(prefix string, d *schema.ResourceData, retryType string, retryOn *client.TemplateRetryOn) {
	if prefix != "" {
		values := d.Get(prefix)
		valuesSlice := values.([]interface{})
		valueMap := valuesSlice[0].(map[string]interface{})
		if retryOn != nil {
			valueMap["retries_on_"+retryType] = retryOn.Times
			valueMap["retry_on_"+retryType+"_only_when_matches_regex"] = retryOn.ErrorRegex
		} else {
			valueMap["retries_on_"+retryType] = 0
			valueMap["retry_on_"+retryType+"_only_when_matches_regex"] = ""
		}

		d.Set(prefix, values)

		return
	}

	if retryOn != nil {
		d.Set(prefix+"retries_on_"+retryType, retryOn.Times)
		d.Set(prefix+"retry_on_"+retryType+"_only_when_matches_regex", retryOn.ErrorRegex)
	} else {
		d.Set(prefix+"retries_on_"+retryType, 0)
		d.Set(prefix+"retry_on_"+retryType+"_only_when_matches_regex", "")
	}
}
