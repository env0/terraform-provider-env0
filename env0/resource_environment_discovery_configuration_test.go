package env0

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentDiscoveryConfigurationResource(t *testing.T) {
	resourceType := "env0_environment_discovery_configuration"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName

	accessor := resourceAccessor(resourceType, resourceName)

	projectId := "pid"
	id := "id"

	t.Run("default (opentofu and terragrunt)", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "opentofu",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			OpentofuVersion:      "1.6.2",
			GithubInstallationId: 12345,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			OpentofuVersion:      putPayload.OpentofuVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
		}

		updatePutPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**/**",
			Repository:           "https://re.po",
			Type:                 "terragrunt",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			OpentofuVersion:      "1.6.3",
			TerragruntVersion:    "0.63.0",
			GithubInstallationId: 3213,
			TerragruntTfBinary:   "opentofu",
		}

		updateGetPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          updatePutPayload.GlobPattern,
			Repository:           updatePutPayload.Repository,
			Type:                 updatePutPayload.Type,
			EnvironmentPlacement: updatePutPayload.EnvironmentPlacement,
			WorkspaceNaming:      updatePutPayload.WorkspaceNaming,
			OpentofuVersion:      updatePutPayload.OpentofuVersion,
			TerragruntVersion:    updatePutPayload.TerragruntVersion,
			GithubInstallationId: updatePutPayload.GithubInstallationId,
			TerragruntTfBinary:   "opentofu",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"opentofu_version":       putPayload.OpentofuVersion,
						"github_installation_id": putPayload.GithubInstallationId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "opentofu_version", putPayload.OpentofuVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           updatePutPayload.GlobPattern,
						"repository":             updatePutPayload.Repository,
						"opentofu_version":       updatePutPayload.OpentofuVersion,
						"github_installation_id": updatePutPayload.GithubInstallationId,
						"type":                   updatePutPayload.Type,
						"terragrunt_version":     updatePutPayload.TerragruntVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", updatePutPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", updatePutPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", updatePutPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", updatePutPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", updatePutPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "opentofu_version", updatePutPayload.OpentofuVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(updatePutPayload.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "terragrunt_version", updatePutPayload.TerragruntVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(2).Return(&getPayload, nil),
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &updatePutPayload).Times(1).Return(&updateGetPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&updateGetPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("terraform", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "terraform",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			TerraformVersion:     "1.6.2",
			GithubInstallationId: 12345,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			TerraformVersion:     putPayload.TerraformVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"type":                   putPayload.Type,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"terraform_version":      putPayload.TerraformVersion,
						"github_installation_id": putPayload.GithubInstallationId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "terraform_version", putPayload.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("workflow", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "workflow",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			GithubInstallationId: 12345,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			GithubInstallationId: putPayload.GithubInstallationId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"type":                   putPayload.Type,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"github_installation_id": putPayload.GithubInstallationId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("bitbucket & terragrunt & terraform", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "terragrunt",
			TerragruntTfBinary:   "terraform",
			TerragruntVersion:    "0.65.0",
			TerraformVersion:     "1.5.0",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			BitbucketClientKey:   "key",
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			TerragruntTfBinary:   putPayload.TerragruntTfBinary,
			TerragruntVersion:    putPayload.TerragruntVersion,
			TerraformVersion:     putPayload.TerraformVersion,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			BitbucketClientKey:   putPayload.BitbucketClientKey,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":           projectId,
						"type":                 putPayload.Type,
						"glob_pattern":         putPayload.GlobPattern,
						"repository":           putPayload.Repository,
						"terragrunt_tf_binary": putPayload.TerragruntTfBinary,
						"terragrunt_version":   putPayload.TerragruntVersion,
						"terraform_version":    putPayload.TerraformVersion,
						"bitbucket_client_key": putPayload.BitbucketClientKey,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "bitbucket_client_key", putPayload.BitbucketClientKey),
						resource.TestCheckResourceAttr(accessor, "terragrunt_tf_binary", putPayload.TerragruntTfBinary),
						resource.TestCheckResourceAttr(accessor, "terragrunt_version", putPayload.TerragruntVersion),
						resource.TestCheckResourceAttr(accessor, "terraform_version", putPayload.TerraformVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("terraform + gitlab", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "terraform",
			TerraformVersion:     "1.7.8",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			GitlabProjectId:      12345,
			TokenId:              "abcdefg",
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			TerraformVersion:     putPayload.TerraformVersion,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			GitlabProjectId:      putPayload.GitlabProjectId,
			TokenId:              putPayload.TokenId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":        projectId,
						"type":              putPayload.Type,
						"glob_pattern":      putPayload.GlobPattern,
						"repository":        putPayload.Repository,
						"gitlab_project_id": putPayload.GitlabProjectId,
						"token_id":          putPayload.TokenId,
						"terraform_version": putPayload.TerraformVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "gitlab_project_id", strconv.Itoa(putPayload.GitlabProjectId)),
						resource.TestCheckResourceAttr(accessor, "token_id", putPayload.TokenId),
						resource.TestCheckResourceAttr(accessor, "terraform_version", putPayload.TerraformVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("azure devops", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "workflow",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			TokenId:              "12345",
			IsAzureDevops:        true,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			TokenId:              putPayload.TokenId,
			IsAzureDevops:        putPayload.IsAzureDevops,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":      projectId,
						"type":            putPayload.Type,
						"glob_pattern":    putPayload.GlobPattern,
						"repository":      putPayload.Repository,
						"token_id":        putPayload.TokenId,
						"is_azure_devops": putPayload.IsAzureDevops,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "token_id", putPayload.TokenId),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", "true"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("sshkey", func(t *testing.T) {
		sshKeyId := "sshi"
		sshKeyName := "sshn"
		sshKeys := []client.TemplateSshKey{
			{
				Id:   sshKeyId,
				Name: sshKeyName,
			},
		}

		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "terraform",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			TerraformVersion:     "1.6.2",
			GithubInstallationId: 12345,
			SshKeys:              sshKeys,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			TerraformVersion:     putPayload.TerraformVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
			SshKeys:              putPayload.SshKeys,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"type":                   putPayload.Type,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"terraform_version":      putPayload.TerraformVersion,
						"github_installation_id": putPayload.GithubInstallationId,
						"ssh_key_name":           sshKeyName,
						"ssh_key_id":             sshKeyId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "terraform_version", putPayload.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "ssh_key_name", sshKeyName),
						resource.TestCheckResourceAttr(accessor, "ssh_key_id", sshKeyId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("retry", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "terraform",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			TerraformVersion:     "1.6.2",
			GithubInstallationId: 12345,
			Retry: client.TemplateRetry{
				OnDeploy: &client.TemplateRetryOn{
					Times:      3,
					ErrorRegex: "abc",
				},
				OnDestroy: &client.TemplateRetryOn{
					Times: 1,
				},
			},
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			TerraformVersion:     putPayload.TerraformVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
			Retry:                putPayload.Retry,
		}

		updatePutPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**/**",
			Repository:           "https://re.po",
			Type:                 "terraform",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			TerraformVersion:     "1.6.2",
			GithubInstallationId: 12345,
			Retry: client.TemplateRetry{
				OnDestroy: &client.TemplateRetryOn{
					Times: 1,
				},
			},
		}

		updateGetPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          updatePutPayload.GlobPattern,
			Repository:           updatePutPayload.Repository,
			Type:                 updatePutPayload.Type,
			EnvironmentPlacement: updatePutPayload.EnvironmentPlacement,
			WorkspaceNaming:      updatePutPayload.WorkspaceNaming,
			TerraformVersion:     updatePutPayload.TerraformVersion,
			GithubInstallationId: updatePutPayload.GithubInstallationId,
			Retry:                updatePutPayload.Retry,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"type":                   putPayload.Type,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"terraform_version":      putPayload.TerraformVersion,
						"github_installation_id": putPayload.GithubInstallationId,
						"retries_on_deploy":      3,
						"retry_on_deploy_only_when_matches_regex": "abc",
						"retries_on_destroy":                      1,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "terraform_version", putPayload.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "retries_on_deploy", "3"),
						resource.TestCheckResourceAttr(accessor, "retry_on_deploy_only_when_matches_regex", "abc"),
						resource.TestCheckResourceAttr(accessor, "retries_on_destroy", "1"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"type":                   updatePutPayload.Type,
						"glob_pattern":           updatePutPayload.GlobPattern,
						"repository":             updatePutPayload.Repository,
						"terraform_version":      updatePutPayload.TerraformVersion,
						"github_installation_id": updatePutPayload.GithubInstallationId,
						"retries_on_destroy":     1,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", updatePutPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", updatePutPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", updatePutPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", updatePutPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", updatePutPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "terraform_version", updatePutPayload.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(updatePutPayload.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "retries_on_destroy", "1"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(2).Return(&getPayload, nil),
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &updatePutPayload).Times(1).Return(&updateGetPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&updateGetPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("error: default (opentofu) with no version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"github_installation_id": 1234,
					}),
					ExpectError: regexp.MustCompile("'opentofu_version' not set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: terraform with no version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"github_installation_id": 1234,
						"type":                   "terraform",
					}),
					ExpectError: regexp.MustCompile("'terraform_version' not set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: terragrunt with no version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"github_installation_id": 1234,
						"type":                   "terragrunt",
					}),
					ExpectError: regexp.MustCompile("'terragrunt_version' not set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: opentofu (with terragrunt) version not set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"github_installation_id": 1234,
						"type":                   "terragrunt",
						"terragrunt_version":     "0.65.1",
					}),
					ExpectError: regexp.MustCompile("'terragrunt_tf_binary' is set to 'opentofu', but 'opentofu_version' not set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: opentofu (with terragrunt) version not set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"github_installation_id": 1234,
						"type":                   "terragrunt",
						"terragrunt_version":     "0.65.1",
						"terragrunt_tf_binary":   "terraform",
					}),
					ExpectError: regexp.MustCompile("'terragrunt_tf_binary' is set to 'terraform', but 'terraform_version' not set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: no vcs set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":   projectId,
						"glob_pattern": "**",
						"repository":   "https://re.po",
						"type":         "workflow",
					}),
					ExpectError: regexp.MustCompile("must set exactly one vcs, none were configured"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("error: more than one vcs set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           "**",
						"repository":             "https://re.po",
						"type":                   "workflow",
						"github_installation_id": 1234,
						"gitlab_project_id":      5678,
						"token_id":               "1345",
					}),
					ExpectError: regexp.MustCompile("must set exactly one vcs, but more were configured: github_installation_id, gitlab_project_id"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("import", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "opentofu",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			OpentofuVersion:      "1.6.2",
			GithubInstallationId: 12345,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			OpentofuVersion:      putPayload.OpentofuVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":             projectId,
						"glob_pattern":           putPayload.GlobPattern,
						"repository":             putPayload.Repository,
						"opentofu_version":       putPayload.OpentofuVersion,
						"github_installation_id": putPayload.GithubInstallationId,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     projectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(3).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})
}
