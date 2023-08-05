package env0

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TemplateType string

const (
	CustomFlow     TemplateType = "custom flow"
	ApprovalPolicy TemplateType = "approval policy"
)

func getTemplate(templateType TemplateType) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("id of the %s", templateType),
			Computed:    true,
		},
		"repository": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("repository url for the %s source code", templateType),
			Required:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "terraform / terragrunt file folder inside source code. Should be the full path including the .yaml/.yml file",
			Optional:    true,
		},
		"revision": {
			Type:        schema.TypeString,
			Description: "source code revision (branch / tag) to use",
			Optional:    true,
		},
		"token_id": {
			Type:        schema.TypeString,
			Description: "the git token id to be used",
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
		"gitlab_project_id": {
			Type:         schema.TypeInt,
			Description:  "the project id of the relevant repository",
			Optional:     true,
			RequiredWith: []string{"token_id"},
		},
		"github_installation_id": {
			Type:        schema.TypeInt,
			Description: "the env0 application installation id on the relevant github repository",
			Optional:    true,
		},
		"bitbucket_client_key": {
			Type:        schema.TypeString,
			Description: "the bitbucket client key used for integration",
			Optional:    true,
		},
		"is_bitbucket_server": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses bitbucket server repository", templateType),
			Optional:    true,
			Default:     false,
		},
		"is_gitlab_enterprise": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses gitlab enterprise repository", templateType),
			Optional:    true,
			Default:     false,
		},
		"is_github_enterprise": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses github enterprise repository", templateType),
			Optional:    true,
			Default:     false,
		},
		"is_gitlab": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: fmt.Sprintf("true if this %s integrates with gitlab repository", templateType),
			Default:     false,
		},
		"is_azure_devops": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: fmt.Sprintf("true if this %s integrates with azure dev ops repository", templateType),
			Default:     false,
		},
	}

	if templateType == CustomFlow {
		s["name"] = &schema.Schema{
			Type:        schema.TypeString,
			Description: "name for the custom flow. note: for the UI to render the custom-flow please use `project-<project.id>`",
			Required:    true,
		}
	}

	if templateType == ApprovalPolicy {
		s["name"] = &schema.Schema{
			Type:        schema.TypeString,
			Description: "name for the approval policy. The name must be in the following format `approval-policy-{scope}-{scopeId}` (E.g. `approval-policy-PROJECT-<project.id>`) - see examples",
			Required:    true,
			ForceNew:    true,
		}
	}

	return s
}
