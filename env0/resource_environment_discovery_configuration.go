package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironmentDiscoveryConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentDiscoveryConfigurationPut,
		ReadContext:   resourceEnvironmentDiscoveryConfigurationGet,
		UpdateContext: resourceEnvironmentDiscoveryConfigurationPut,
		DeleteContext: resourceEnvironmentDiscoveryConfigurationDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceEnvironmentDiscoveryConfigurationImport},

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
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "the repository to run discovery on",
				Optional:    true,
			},
			"repository_regex": {
				Type:        schema.TypeString,
				Description: "Regex to select repositories for discovery-file configuration (enables discoveryFileConfiguration mode)",
				Optional:    true,
			},
			"type": {
				Type:             schema.TypeString,
				Description:      "the infrastructure type use. Valid values: 'opentofu', 'terraform', 'terragrunt', 'workflow' (default: 'opentofu')",
				ValidateDiagFunc: NewStringInValidator([]string{client.OPENTOFU, client.TERRAFORM, client.TERRAGRUNT, client.WORKFLOW}),
				Optional:         true,
				Computed:         true,
			},
			"environment_placement": {
				Type:             schema.TypeString,
				Description:      "the environment placement strategy with the project.",
				ValidateDiagFunc: NewStringInValidator([]string{"existingSubProject", "topProject"}),
				Optional:         true,
				Computed:         true,
			},
			"workspace_naming": {
				Type:             schema.TypeString,
				Description:      "the Workspace namimg strategy.",
				ValidateDiagFunc: NewStringInValidator([]string{"default", "environmentName"}),
				Optional:         true,
				Computed:         true,
			},
			"auto_deploy_by_custom_glob": {
				Type:        schema.TypeString,
				Description: "If specified, deploy/plan on changes matching the given pattern (glob). Otherwise, deploy on template folder changes only",
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
				ValidateDiagFunc: NewRegexValidator(`^(?:[0-9]\.[0-9]{1,2}\.[0-9]{1,2})|RESOLVE_FROM_TERRAFORM_CODE|latest$`),
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
				Description:      "The binary to use with Terragrunt. Valid values: 'opentofu' and 'terraform' (default: 'opentofu')",
				ValidateDiagFunc: NewStringInValidator([]string{client.OPENTOFU, client.TERRAFORM}),
				Default:          client.OPENTOFU,
			},
			"is_terragrunt_run_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to 'true', execute terragrunt commands with 'run all' (default: false)",
				Default:     false,
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
			"vcs_connection_id": {
				Type:        schema.TypeString,
				Description: "the VCS connection id to be used",
				Optional:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "bitbucket client",
				Optional:    true,
			},
			"gitlab_project_id": {
				Type:        schema.TypeInt,
				Description: "gitlab project id (deprecated)",
				Optional:    true,
				Deprecated:  "project id is now auto-fetched from the repository URL",
			},
			"is_azure_devops": {
				Type:         schema.TypeBool,
				Optional:     true,
				Description:  "set to true if azure devops is used",
				Default:      false,
				RequiredWith: []string{"token_id"},
			},
			"is_bitbucket_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true if Bitbucket Server is used",
				Default:     false,
			},
			"is_github_enterprise": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true if GitHub Enterprise is used",
				Default:     false,
			},
			"is_gitlab_enterprise": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true if GitLab Enterprise is used",
				Default:     false,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "a token id to be used with 'gitlab' or 'azure_devops'",
				Optional:    true,
			},
			"root_path": {
				Type:        schema.TypeString,
				Description: "start files Glob matching from this folder only",
				Optional:    true,
			},
			"create_new_environments_from_pull_requests": {
				Type:        schema.TypeBool,
				Description: "create new environments from pull requests (default: false)",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func discoveryReadSshKeyHelper(putPayload *client.EnvironmentDiscoveryPutPayload, d *schema.ResourceData) {
	sshKeyId := d.Get("ssh_key_id").(string)
	if sshKeyId != "" {
		sshKeyName := d.Get("ssh_key_name").(string)
		putPayload.SshKeys = append(putPayload.SshKeys, client.TemplateSshKey{
			Id:   sshKeyId,
			Name: sshKeyName,
		})
	}
}

func discoveryWriteSshKeyHelper(getPayload *client.EnvironmentDiscoveryPayload, d *schema.ResourceData) {
	var sshKey client.TemplateSshKey

	if len(getPayload.SshKeys) > 0 {
		sshKey = getPayload.SshKeys[0]
	}

	d.Set("ssh_key_id", sshKey.Id)
	d.Set("ssh_key_name", sshKey.Name)
}

func discoveryValidatePutPayload(putPayload *client.EnvironmentDiscoveryPutPayload) error {
	// If discovery-file configuration is used, skip specific validations
	if putPayload.DiscoveryFileConfiguration != nil && putPayload.DiscoveryFileConfiguration.RepositoryRegex != "" {
		return nil
	}

	if putPayload.Repository == "" {
		return errors.New("'repository' not set")
	}

	if putPayload.GlobPattern == "" {
		return errors.New("'glob_pattern' not set")
	}

	opentofuVersionSet := putPayload.OpentofuVersion != ""
	terraformVersionSet := putPayload.TerraformVersion != ""
	terragruntVersionSet := putPayload.TerragruntVersion != ""

	switch putPayload.Type {
	case client.OPENTOFU:
		if !opentofuVersionSet {
			return errors.New("'opentofu_version' not set")
		}
	case client.TERRAFORM:
		if !terraformVersionSet {
			return errors.New("'terraform_version' not set")
		}
	case client.TERRAGRUNT:
		if !terragruntVersionSet {
			return errors.New("'terragrunt_version' not set")
		}

		if putPayload.TerragruntTfBinary == client.OPENTOFU && !opentofuVersionSet {
			return errors.New("'terragrunt_tf_binary' is set to 'opentofu', but 'opentofu_version' not set")
		}

		if putPayload.TerragruntTfBinary == client.TERRAFORM && !terraformVersionSet {
			return errors.New("'terragrunt_tf_binary' is set to 'terraform', but 'terraform_version' not set")
		}
	case client.WORKFLOW:
	default:
		return fmt.Errorf("unhandled type %s", putPayload.Type)
	}

	return nil
}

func resourceEnvironmentDiscoveryConfigurationPut(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var putPayload client.EnvironmentDiscoveryPutPayload
	if err := readResourceData(&putPayload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if v, ok := d.GetOk("repository_regex"); ok {
		repoRegex := v.(string)
		if repoRegex != "" {
			if gv, exists := d.GetOk("glob_pattern"); exists && gv.(string) != "" {
				return diag.Errorf("'glob_pattern' cannot be set when 'repository_regex' is provided")
			}

			if rv, exists := d.GetOk("repository"); exists && rv.(string) != "" {
				return diag.Errorf("'repository' cannot be set when 'repository_regex' is provided")
			}

			putPayload = client.EnvironmentDiscoveryPutPayload{
				DiscoveryFileConfiguration: &client.DiscoveryFileConfiguration{RepositoryRegex: repoRegex},
			}
		}
	}

	// Provider-side defaults for non discovery-file mode (fields are Optional+Computed with no schema defaults)
	if putPayload.DiscoveryFileConfiguration == nil {
		if putPayload.Type == "" {
			putPayload.Type = client.OPENTOFU
		}

		if putPayload.EnvironmentPlacement == "" {
			putPayload.EnvironmentPlacement = "topProject"
		}

		if putPayload.WorkspaceNaming == "" {
			putPayload.WorkspaceNaming = "default"
		}
	}

	if err := putPayload.Invalidate(); err != nil {
		return diag.Errorf("invalid environment discovery payload: %v", err)
	}

	if putPayload.DiscoveryFileConfiguration == nil {
		discoveryReadSshKeyHelper(&putPayload, d)
		templateCreatePayloadRetryOnHelper("", d, "deploy", &putPayload.Retry.OnDeploy)
		templateCreatePayloadRetryOnHelper("", d, "destroy", &putPayload.Retry.OnDestroy)
	}

	if err := discoveryValidatePutPayload(&putPayload); err != nil {
		return diag.Errorf("validation error: %s", err.Error())
	}

	if putPayload.Type != client.TERRAGRUNT {
		putPayload.TerragruntTfBinary = ""
	}

	res, err := apiClient.PutEnvironmentDiscovery(d.Get("project_id").(string), &putPayload)
	if err != nil {
		return diag.Errorf("enable/update environment discovery configuration request failed: %s", err.Error())
	}

	// If vcs_connection_id is set in config, ignore github_installation_id from backend to avoid drift
	if _, ok := d.GetOk("vcs_connection_id"); ok {
		res.GithubInstallationId = 0
	}

	if err := setResourceEnvironmentDiscoveryConfiguration(d, res); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Id)

	return nil
}

func resourceEnvironmentDiscoveryConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	if err := apiClient.DeleteEnvironmentDiscovery(projectId); err != nil {
		return diag.Errorf("delete environment discovery configuration request failed: %s", err.Error())
	}

	return nil
}

func setResourceEnvironmentDiscoveryConfiguration(d *schema.ResourceData, getPayload *client.EnvironmentDiscoveryPayload) error {
	if err := writeResourceData(getPayload, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	discoveryWriteSshKeyHelper(getPayload, d)

	templateReadRetryOnHelper("", d, "deploy", getPayload.Retry.OnDeploy)
	templateReadRetryOnHelper("", d, "destroy", getPayload.Retry.OnDestroy)

	if getPayload.DiscoveryFileConfiguration != nil {
		_ = d.Set("repository_regex", getPayload.DiscoveryFileConfiguration.RepositoryRegex)
	}

	if getPayload.DiscoveryFileConfiguration == nil {
		// Apply defaults only when not using discovery-file configuration
		typeVal := getPayload.Type
		if typeVal == "" {
			typeVal = client.OPENTOFU
		}
		_ = d.Set("type", typeVal)

		envPlacement := getPayload.EnvironmentPlacement
		if envPlacement == "" {
			envPlacement = "topProject"
		}
		_ = d.Set("environment_placement", envPlacement)

		workspaceNaming := getPayload.WorkspaceNaming
		if workspaceNaming == "" {
			workspaceNaming = "default"
		}
		_ = d.Set("workspace_naming", workspaceNaming)
	}

	return nil
}

func resourceEnvironmentDiscoveryConfigurationGet(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	getPayload, err := apiClient.GetEnvironmentDiscovery(projectId)
	if err != nil {
		return ResourceGetFailure(ctx, "environment_discovery_configuration", d, err)
	}

	_, vcsConnectionIdOk := d.GetOk("vcs_connection_id")
	if vcsConnectionIdOk {
		getPayload.GithubInstallationId = 0
	}

	if err := setResourceEnvironmentDiscoveryConfiguration(d, getPayload); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEnvironmentDiscoveryConfigurationImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Id()

	getPayload, err := apiClient.GetEnvironmentDiscovery(projectId)
	if err != nil {
		return nil, err
	}

	if err := setResourceEnvironmentDiscoveryConfiguration(d, getPayload); err != nil {
		return nil, err
	}

	d.Set("project_id", projectId)

	if _, ok := d.GetOk("terragrunt_tf_binary"); !ok {
		d.Set("terragrunt_tf_binary", client.OPENTOFU)
	}

	return []*schema.ResourceData{d}, nil
}
