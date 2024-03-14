package env0

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
				Required:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "the repository to run discovery on",
				Required:    true,
			},
			"type": {
				Type:             schema.TypeString,
				Description:      "the infrastructure type use. Valid values: 'opentofu', 'terraform', 'terragrunt', 'workflow' (default: 'opentofu')",
				Default:          "opentofu",
				ValidateDiagFunc: NewStringInValidator([]string{"opentofu", "terraform", "terragrunt", "workflow"}),
				Optional:         true,
			},
			"environment_placement": {
				Type:             schema.TypeString,
				Description:      "the environment placement strategy with the project (default: 'topProject')",
				Default:          "topProject",
				ValidateDiagFunc: NewStringInValidator([]string{"existingSubProject", "topProject"}),
				Optional:         true,
			},
			"workspace_naming": {
				Type:             schema.TypeString,
				Description:      "the Workspace namimg strategy (default: 'default')",
				Default:          "default",
				ValidateDiagFunc: NewStringInValidator([]string{"default", "environmentName"}),
				Optional:         true,
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
				Description:      "The binary to use with Terragrunt. Valid values: 'opentofu' and 'terraform' (default: 'opentofu')",
				ValidateDiagFunc: NewStringInValidator([]string{"opentofu", "terraform"}),
				Default:          "opentofu",
			},
			"is_terragrunt_run_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to 'true', execute terragrunt commands with 'run all'",
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
				Default:      false,
				RequiredWith: []string{"token_id"},
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "a token id to be used with 'gitlab' or 'azure_devops'",
				Optional:    true,
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
	opentofuVersionSet := putPayload.OpentofuVersion != ""
	terraformVersionSet := putPayload.TerraformVersion != ""
	terragruntVersionSet := putPayload.TerragruntVersion != ""

	switch putPayload.Type {
	case "opentofu":
		if !opentofuVersionSet {
			return errors.New("'opentofu_version' not set")
		}
	case "terraform":
		if !terraformVersionSet {
			return errors.New("'terraform_version' not set")
		}
	case "terragrunt":
		if !terragruntVersionSet {
			return errors.New("'terragrunt_version' not set")
		}

		if putPayload.TerragruntTfBinary == "opentofu" && !opentofuVersionSet {
			return errors.New("'terragrunt_tf_binary' is set to 'opentofu', but 'opentofu_version' not set")
		}

		if putPayload.TerragruntTfBinary == "terraform" && !terraformVersionSet {
			return errors.New("'terragrunt_tf_binary' is set to 'terraform', but 'terraform_version' not set")
		}
	case "workflow":
	default:
		return fmt.Errorf("unhandled type %s", putPayload.Type)
	}

	vcsCounter := 0
	vcsEnabledAttributes := []string{}

	if putPayload.GithubInstallationId != 0 {
		vcsCounter++
		vcsEnabledAttributes = append(vcsEnabledAttributes, "github_installation_id")
	}

	if putPayload.BitbucketClientKey != "" {
		vcsCounter++
		vcsEnabledAttributes = append(vcsEnabledAttributes, "bitbucket_client_key")
	}

	if putPayload.GitlabProjectId != 0 {
		vcsCounter++
		vcsEnabledAttributes = append(vcsEnabledAttributes, "gitlab_project_id")
	}

	if putPayload.IsAzureDevops {
		vcsCounter++
		vcsEnabledAttributes = append(vcsEnabledAttributes, "is_azure_devops")
	}

	if vcsCounter == 0 {
		return errors.New("must set exactly one vcs, none were configured: github_installation_id, bitbucket_client_key, gitlab_project_id, or is_azure_devops")
	}

	if vcsCounter > 1 {
		return fmt.Errorf("must set exactly one vcs, but more were configured: %s", strings.Join(vcsEnabledAttributes, ", "))
	}

	return nil
}

func resourceEnvironmentDiscoveryConfigurationPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var putPayload client.EnvironmentDiscoveryPutPayload
	if err := readResourceData(&putPayload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	discoveryReadSshKeyHelper(&putPayload, d)

	templateCreatePayloadRetryOnHelper("", d, "deploy", &putPayload.Retry.OnDeploy)
	templateCreatePayloadRetryOnHelper("", d, "destroy", &putPayload.Retry.OnDestroy)

	if err := discoveryValidatePutPayload(&putPayload); err != nil {
		return diag.Errorf("validation error: %s", err.Error())
	}

	if putPayload.Type != "terragrunt" {
		// Remove the default terragrunt_tf_binary if terragrunt isn't used.
		putPayload.TerragruntTfBinary = ""
	}

	res, err := apiClient.EnableUpdateEnvironmentDiscovery(d.Get("project_id").(string), &putPayload)
	if err != nil {
		return diag.Errorf("enable/update environment discovery configuration request failed: %s", err.Error())
	}

	d.SetId(res.Id)

	return nil
}

func resourceEnvironmentDiscoveryConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	if err := apiClient.DeleteEnvironmentDiscovery(projectId); err != nil {
		return diag.Errorf("delete environment discovery configuration request failed: %s", err.Error())
	}

	return nil
}

func setResourceEnvironmentDiscoveryConfiguration(d *schema.ResourceData, getPayload *client.EnvironmentDiscoveryPayload) error {
	if err := writeResourceData(getPayload, d); err != nil {
		return fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	discoveryWriteSshKeyHelper(getPayload, d)

	templateReadRetryOnHelper("", d, "deploy", getPayload.Retry.OnDeploy)
	templateReadRetryOnHelper("", d, "destroy", getPayload.Retry.OnDestroy)

	return nil
}

func resourceEnvironmentDiscoveryConfigurationGet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	getPayload, err := apiClient.GetEnvironmentDiscovery(projectId)
	if err != nil {
		return ResourceGetFailure(ctx, "environment_discovery_configuration", d, err)
	}

	if err := setResourceEnvironmentDiscoveryConfiguration(d, getPayload); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEnvironmentDiscoveryConfigurationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
		d.Set("terragrunt_tf_binary", "opentofu")
	}

	return []*schema.ResourceData{d}, nil
}
