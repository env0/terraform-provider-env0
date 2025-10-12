package env0

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TemplateType string

const (
	CustomFlow     TemplateType = "custom-flow"
	ApprovalPolicy TemplateType = "approval-policy"
)

func getConfigurationTemplateSchema(templateType TemplateType) map[string]*schema.Schema {
	var text string

	switch templateType {
	case CustomFlow:
		text = "custom flow"
	case ApprovalPolicy:
		text = "approval policy"
	}

	s := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "id of the " + text,
			Computed:    true,
		},
		"repository": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("repository url for the %s source code", text),
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
			Type:        schema.TypeInt,
			Description: "the project id of the relevant repository (deprecated)",
			Deprecated:  "project id is now auto-fetched from the repository URL",
			Optional:    true,
		},
		"github_installation_id": {
			Type:          schema.TypeInt,
			Description:   "the env0 application installation id on the relevant github repository",
			Optional:      true,
			ConflictsWith: []string{"vcs_connection_id"},
		},
		"vcs_connection_id": {
			Type:          schema.TypeString,
			Description:   "the VCS connection id to be used",
			Optional:      true,
			ConflictsWith: []string{"github_installation_id"},
		},
		"bitbucket_client_key": {
			Type:        schema.TypeString,
			Description: "the bitbucket client key used for integration",
			Optional:    true,
		},
		"is_bitbucket_server": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses bitbucket server repository", text),
			Optional:    true,
			Default:     false,
		},
		"is_gitlab_enterprise": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses gitlab enterprise repository", text),
			Optional:    true,
			Default:     false,
		},
		"is_github_enterprise": {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("true if this %s uses github enterprise repository", text),
			Optional:    true,
			Default:     false,
		},
		"is_gitlab": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: fmt.Sprintf("true if this %s integrates with gitlab repository", text),
			Default:     false,
		},
		"is_azure_devops": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: fmt.Sprintf("true if this %s integrates with azure dev ops repository", text),
			Default:     false,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name for the " + text,
			Required:    true,
		},
	}

	return s
}
