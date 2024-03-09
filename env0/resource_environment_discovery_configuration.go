package env0

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceEnvironmentDiscoveryConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentDriftCreateOrUpdate,
		ReadContext:   resourceEnvironmentDriftRead,
		UpdateContext: resourceEnvironmentDriftCreateOrUpdate,
		DeleteContext: resourceEnvironmentDriftDelete,

		Description: "See https://docs.env0.com/docs/environment-discovery for additional details",

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "the project id",
				Required:    true,
				ForceNew:    true,
			},
			"glob_pattern": {
				Type:        schema.TypeString,
				Description: "the environments glob pattern. Any match to this pattern will result in an Environment creation and plan",
				Required:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "the repository to run discovery on",
				Required:    true,
			},
			"type": {
				Type:             schema.TypeString,
				Description:      "the infrastructure type use. Valid values: 'opentofu', 'terraform', 'terragrunt', 'workflow'",
				Required:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"opentofu", "terraform", "terragrunt", "workflow"}),
			},
			"environment_placement": {
				Type:             schema.TypeString,
				Description:      "the environment placement strategy with the project (default: 'topProject')",
				Required:         true,
				Default:          "topProject",
				ValidateDiagFunc: NewStringInValidator([]string{"existingSubProject", "topProject"}),
			},
			"workspace_naming": {
				Type:             schema.TypeString,
				Description:      "the Workspace namimg strategy (default: 'default')",
				Optional:         true,
				Default:          "default",
				ValidateDiagFunc: NewStringInValidator([]string{"default", "environmentName"}),
			},
			"auto_deploy_by_custom_glob": {
				Type:        schema.TypeString,
				Description: "if configured, deploy/plan on changes matching the given pattern (glob). Otherwise, deploy/plan on changes in directories matching the main glob pattern",
				Optional:    true,
			},
			"terraform_version": {
				Type:             schema.TypeString,
				Description:      "the Terraform version to use (example: 1.7.4). Setting to `RESOLVE_FROM_TERRAFORM_CODE` defaults to the version of `terraform.required_version` during run-time (resolve from terraform code). Setting to `latest`, the version used will be the most recent one available for Terraform.",
				Optional:         true,
				ValidateDiagFunc: NewRegexValidator(`^(?:[0-9]\.[0-9]{1,2}\.[0-9]{1,2})|RESOLVE_FROM_TERRAFORM_CODE|latest$`),
			},
			"opentofu_version": {
				Type:             schema.TypeString,
				Description:      "the Opentofu version to use (example: 1.6.1). Setting to `latest`, the version used will be the most recent one available for OpenTofu.",
				Optional:         true,
				ValidateDiagFunc: NewRegexValidator(`^(?:[0-9]\.[0-9]{1,2}\.[0-9]{1,2})|latest$`),
			},
			"terragrunt_version": {
				Type:             schema.TypeString,
				Description:      "the Terragrunt version to use (example: 0.52.0)",
				ValidateDiagFunc: NewRegexValidator(`^[0-9]\.[0-9]{1,2}\.[0-9]{1,2}$`),
				Optional:         true,
			},
			"terragrunt_tf_binary": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The binary to use with Terragrunt. Valid values: 'opentofu' and 'terraform'",
				ValidateDiagFunc: NewStringInValidator([]string{"opentofu", "terraform"}),
				RequiredWith:     []string{"terragrunt_version"},
			},
			"is_terragrunt_run_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to 'true', execute terragrunt commands with 'run all'",
				Default:     "false",
			},
			"ssh_key_id": {
				Type:         schema.TypeString,
				Description:  "The ssh key id that will be available during deployment",
				Optional:     true,
				RequiredWith: []string{"ssh_key_name"},
			},
			"ssh_key_name": {
				Type:         schema.TypeString,
				Description:  "The ssh key name that will be available during deployment",
				Optional:     true,
				RequiredWith: []string{"ssh_key_id"},
			},
			"retries_on_deploy": {
				Type:             schema.TypeInt,
				Description:      "number of times to retry when deploy fails (between 1 and 3)",
				Optional:         true,
				ValidateDiagFunc: ValidateRetries,
			},
			"retry_on_deploy_only_when_matches_regex": {
				Type:         schema.TypeString,
				Description:  "retry (on deploy) if error matches the specified regex",
				Optional:     true,
				RequiredWith: []string{"retries_on_deploy"},
			},
			"retries_on_destroy": {
				Type:             schema.TypeInt,
				Description:      "number of times to retry when destroy fails (between 1 and 3)",
				Optional:         true,
				ValidateDiagFunc: ValidateRetries,
			},
			"retry_on_destroy_only_when_matches_regex": {
				Type:         schema.TypeString,
				Description:  "retry (on destroy) if error matches the specified regex",
				Optional:     true,
				RequiredWith: []string{"retries_on_destroy"},
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "github repository id",
				Optional:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "bitbucket client",
				Optional:    true,
			},
			"gitlab_project_id": {
				Type:         schema.TypeInt,
				Description:  "gitlab project id",
				Optional:     true,
				RequiredWith: []string{"token_id"},
			},
			"is_azure_devops": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "set to true if azure devops is used",
				Default:      "false",
				RequiredWith: []string{"token_id"},
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "a token id to be used with 'gitlab' or 'azure_devops'",
				Optional:    true,
			},
			// TODO test - and conflicts with...
		},
	}
}
