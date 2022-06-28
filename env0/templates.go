package env0

import (
	"fmt"
	"sort"
	"strings"

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
		ret := []string{}

		for _, attr := range allVCSAttributes {
			if sort.SearchStrings(strs, attr) >= len(strs) {
				ret = appendByType(ret, attr)
			}
		}

		return ret
	}

	requiredWith := func(strs ...string) []string {
		ret := []string{}
		for _, str := range strs {
			ret = appendByType(ret, str)
		}

		return ret
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
		"revision": {
			Type:        schema.TypeString,
			Description: "source code revision (branch / tag) to use",
			Optional:    true,
		},
		"type": {
			Type:             schema.TypeString,
			Description:      fmt.Sprintf("template type (allowed values: %s)", strings.Join(allowedTemplateTypes, ", ")),
			Optional:         true,
			Default:          "terraform",
			ValidateDiagFunc: NewStringInValidator(allowedTemplateTypes),
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
