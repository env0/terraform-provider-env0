package env0

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitTemplateResource(t *testing.T) {
	const resourceType = "env0_template"

	const resourceName = "test"

	const defaultVersion = "0.15.1"

	const testType = client.OPENTOFU

	const defaultOpentofuVersion = "1.6.0"

	var resourceFullName = resourceAccessor(resourceType, resourceName)

	gleeTemplate := client.Template{
		Id:          "id0",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:               "terragrunt",
		IsGitlabEnterprise: true,
		TerraformVersion:   "0.12.24",
		TerragruntVersion:  "0.35.1",
		TerragruntTfBinary: "terraform",
	}
	gleeUpdatedTemplate := client.Template{
		Id:          gleeTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:               "terragrunt",
		TerragruntVersion:  "0.35.1",
		IsGitlabEnterprise: true,
		TerraformVersion:   "0.15.1",
		IsTerragruntRunAll: true,
		TerragruntTfBinary: "terraform",
	}
	gitlabTemplate := client.Template{
		Id:          "id0-gitlab",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:             "terraform",
		TokenId:          "1",
		TerraformVersion: "0.12.24",
		TokenName:        "token_name",
		IsGitlab:         true,
	}
	gitlabUpdatedTemplate := client.Template{
		Id:          gitlabTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:               "terragrunt",
		TerragruntVersion:  "0.26.1",
		TokenId:            "2",
		TerraformVersion:   "0.15.1",
		TokenName:          "token_name2",
		TerragruntTfBinary: "terraform",
	}
	githubTemplate := client.Template{
		Id:          "id0",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:                 "terraform",
		GithubInstallationId: 1,
		TerraformVersion:     "0.12.24",
	}
	githubUpdatedTemplate := client.Template{
		Id:          githubTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:                 "terragrunt",
		TerragruntVersion:    "0.35.1",
		GithubInstallationId: 2,
		TerraformVersion:     "0.15.1",
		IsTerragruntRunAll:   true,
		TerragruntTfBinary:   "terraform",
	}
	bitbucketTemplate := client.Template{
		Id:          "id0",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:               "terraform",
		TerraformVersion:   "0.12.24",
		BitbucketClientKey: "clientkey",
	}
	bitbucketUpdatedTemplate := client.Template{
		Id:          bitbucketTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:               "terragrunt",
		BitbucketClientKey: "clientkey2",
		TerragruntVersion:  "0.35.1",
		TerraformVersion:   "0.15.1",
		TerragruntTfBinary: "terraform",
	}
	gheeTemplate := client.Template{
		Id:          "id0",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:               "terraform",
		TerraformVersion:   "0.12.24",
		IsGithubEnterprise: true,
	}
	gheeUpdatedTemplate := client.Template{
		Id:          gheeTemplate.Id,
		Name:        "template1",
		Description: "description1",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:               "terraform",
		TerraformVersion:   "0.12.24",
		IsGithubEnterprise: true,
	}
	bitbucketServerTemplate := client.Template{
		Id:          "id011",
		Name:        "template011",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:              "terraform",
		TerraformVersion:  "0.12.24",
		IsBitbucketServer: true,
	}
	bitbucketServerUpdatedTemplate := client.Template{
		Id:          bitbucketServerTemplate.Id,
		Name:        "template222",
		Description: "description1",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:              "terraform",
		TerraformVersion:  "0.12.24",
		IsBitbucketServer: true,
	}
	cloudformationTemplate := client.Template{
		Id:          "id0",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:             "cloudformation",
		FileName:         "cool.yaml",
		TerraformVersion: "0.15.1",
	}
	cloudformationUpdatedTemplate := client.Template{
		Id:          gleeTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:             "cloudformation",
		FileName:         "stack.yaml",
		TerraformVersion: "0.15.1",
	}
	azureDevOpsTemplate := client.Template{
		Id:          "id0-azure-dev-ops",
		Name:        "template0",
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:             "terraform",
		TokenId:          "1",
		TerraformVersion: "latest",
		IsAzureDevOps:    true,
	}
	azureDevOpsUpdatedTemplate := client.Template{
		Id:          azureDevOpsTemplate.Id,
		Name:        "new-name",
		Description: "new-description",
		Repository:  "env0/repo-new",
		Path:        "path/zero/new",
		Revision:    "branch-zero-new",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "NewForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "NewForDestroy.*",
			},
		},
		Type:               "terragrunt",
		TerragruntVersion:  "0.26.1",
		TokenId:            "2",
		TerraformVersion:   "0.15.1",
		TerragruntTfBinary: "terraform",
		IsAzureDevOps:      true,
	}

	helmTemplate := client.Template{
		Id:               "helmTemplate",
		Name:             "template0",
		Description:      "description0",
		Repository:       "env0/repo",
		Type:             "helm",
		HelmChartName:    "chart1",
		IsHelmRepository: true,
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		TerraformVersion: "0.12.24",
	}
	helmUpdatedTemplate := client.Template{
		Id:               helmTemplate.Id,
		Name:             "new-name",
		Description:      "new-description",
		Repository:       "env0/repo-new",
		Type:             "helm",
		HelmChartName:    "chart1",
		IsHelmRepository: true,
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		TerraformVersion: "0.12.24",
	}

	opentofuTemplate := client.Template{
		Id:               "opentofu",
		Name:             "template0",
		Description:      "description0",
		Repository:       "env0/repo",
		Type:             "opentofu",
		OpentofuVersion:  "1.6.0-alpha",
		TerraformVersion: "0.15.1",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
	}
	opentofuUpdatedTemplate := client.Template{
		Id:               opentofuTemplate.Id,
		Name:             "new-name",
		Description:      "new-description",
		Repository:       "env0/repo-new",
		Type:             "opentofu",
		OpentofuVersion:  "RESOLVE_FROM_CODE",
		TerraformVersion: "0.15.1",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
	}

	ansibleTemplate := client.Template{
		Id:             "ansible",
		Name:           "template0",
		Description:    "description0",
		Repository:     "env0/repo",
		Type:           "ansible",
		AnsibleVersion: "3.5.6",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
	}
	ansibleUpdatedTemplate := client.Template{
		Id:             ansibleTemplate.Id,
		Name:           "new-name",
		Description:    "new-description",
		Repository:     "env0/repo-new",
		Type:           "ansible",
		AnsibleVersion: "latest",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
	}

	fullTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
		templateAsDictionary := map[string]interface{}{
			"name":       template.Name,
			"repository": template.Repository,
		}

		if template.Type != "" {
			templateAsDictionary["type"] = template.Type
		} else {
			templateAsDictionary["type"] = string(testType)
		}

		if template.Description != "" {
			templateAsDictionary["description"] = template.Description
		}

		if template.Revision != "" {
			templateAsDictionary["revision"] = template.Revision
		}

		if template.Path != "" {
			templateAsDictionary["path"] = template.Path
		}

		if template.Retry != (client.TemplateRetry{}) && template.Retry.OnDeploy != nil {
			templateAsDictionary["retries_on_deploy"] = template.Retry.OnDeploy.Times
			if template.Retry.OnDeploy.ErrorRegex != "" {
				templateAsDictionary["retry_on_deploy_only_when_matches_regex"] = template.Retry.OnDeploy.ErrorRegex
			}
		}

		if template.Retry != (client.TemplateRetry{}) && template.Retry.OnDestroy != nil {
			templateAsDictionary["retries_on_destroy"] = template.Retry.OnDestroy.Times
			if template.Retry.OnDestroy.ErrorRegex != "" {
				templateAsDictionary["retry_on_destroy_only_when_matches_regex"] = template.Retry.OnDestroy.ErrorRegex
			}
		}

		if template.TerraformVersion != "" {
			templateAsDictionary["terraform_version"] = template.TerraformVersion
		}

		if template.OpentofuVersion != "" {
			templateAsDictionary["opentofu_version"] = template.OpentofuVersion
		}

		if template.TokenId != "" {
			templateAsDictionary["token_id"] = template.TokenId
		}

		if template.TokenName != "" {
			templateAsDictionary["token_name"] = template.TokenName
		}

		if template.GithubInstallationId != 0 {
			templateAsDictionary["github_installation_id"] = template.GithubInstallationId
		}

		if template.IsGitlabEnterprise != false {
			templateAsDictionary["is_gitlab_enterprise"] = template.IsGitlabEnterprise
		}

		if template.BitbucketClientKey != "" {
			templateAsDictionary["bitbucket_client_key"] = template.BitbucketClientKey
		}

		if template.IsGithubEnterprise != false {
			templateAsDictionary["is_github_enterprise"] = template.IsGithubEnterprise
		}

		if template.IsBitbucketServer != false {
			templateAsDictionary["is_bitbucket_server"] = template.IsBitbucketServer
		}

		if template.FileName != "" {
			templateAsDictionary["file_name"] = template.FileName
		}

		if template.TerragruntVersion != "" {
			templateAsDictionary["terragrunt_version"] = template.TerragruntVersion
		}

		if template.IsTerragruntRunAll {
			templateAsDictionary["is_terragrunt_run_all"] = true
		}

		if template.IsAzureDevOps {
			templateAsDictionary["is_azure_devops"] = true
		}

		if template.IsHelmRepository {
			templateAsDictionary["is_helm_repository"] = true
		}

		if template.HelmChartName != "" {
			templateAsDictionary["helm_chart_name"] = template.HelmChartName
		}

		if template.IsGitlab {
			templateAsDictionary["is_gitlab"] = true
		}

		if template.AnsibleVersion != "" {
			templateAsDictionary["ansible_version"] = template.AnsibleVersion
		}

		if template.TerragruntTfBinary != "" {
			templateAsDictionary["terragrunt_tf_binary"] = template.TerragruntTfBinary
		}

		return resourceConfigCreate(resourceType, resourceName, templateAsDictionary)
	}
	fullTemplateResourceCheck := func(resourceFullName string, template client.Template) resource.TestCheckFunc {
		tokenIdAssertion := resource.TestCheckResourceAttr(resourceFullName, "token_id", template.TokenId)
		if template.TokenId == "" {
			tokenIdAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "token_id")
		}

		helmChartNameAssertion := resource.TestCheckResourceAttr(resourceFullName, "helm_chart_name", template.HelmChartName)
		if template.HelmChartName == "" {
			helmChartNameAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "helm_chart_name")
		}

		filenameAssertion := resource.TestCheckResourceAttr(resourceFullName, "file_name", template.FileName)
		if template.FileName == "" {
			filenameAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "file_name")
		}

		terragruntVersionAssertion := resource.TestCheckResourceAttr(resourceFullName, "terragrunt_version", template.TerragruntVersion)
		if template.TerragruntVersion == "" {
			terragruntVersionAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "terragrunt_version")
		}

		opentofuVersionAssertion := resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", template.OpentofuVersion)
		if template.OpentofuVersion == "" {
			opentofuVersionAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "opentofu_version")
		}

		githubInstallationIdAssertion := resource.TestCheckResourceAttr(resourceFullName, "github_installation_id", strconv.Itoa(template.GithubInstallationId))
		if template.GithubInstallationId == 0 {
			githubInstallationIdAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "github_installation_id")
		}

		pathAssertion := resource.TestCheckResourceAttr(resourceFullName, "path", template.Path)
		if template.Path == "" {
			pathAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "path")
		}

		tokenNameAssertion := resource.TestCheckResourceAttr(resourceFullName, "token_name", template.TokenName)
		if template.TokenName == "" {
			tokenNameAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "token_name")
		}

		ansibleVersionAssertion := resource.TestCheckResourceAttr(resourceFullName, "ansible_version", template.AnsibleVersion)
		if template.AnsibleVersion == "" {
			ansibleVersionAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "ansible_version")
		}

		terraformVersionAssertion := resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion)
		if template.TerraformVersion == "" {
			terraformVersionAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "terraform_version")
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
			resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
			resource.TestCheckResourceAttr(resourceFullName, "description", template.Description),
			resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
			resource.TestCheckResourceAttr(resourceFullName, "type", template.Type),
			resource.TestCheckResourceAttr(resourceFullName, "retries_on_deploy", strconv.Itoa(template.Retry.OnDeploy.Times)),
			resource.TestCheckResourceAttr(resourceFullName, "retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex),
			resource.TestCheckResourceAttr(resourceFullName, "retries_on_destroy", strconv.Itoa(template.Retry.OnDestroy.Times)),
			resource.TestCheckResourceAttr(resourceFullName, "retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex),
			resource.TestCheckResourceAttr(resourceFullName, "is_gitlab_enterprise", strconv.FormatBool(template.IsGitlabEnterprise)),
			tokenIdAssertion,
			filenameAssertion,
			terragruntVersionAssertion,
			githubInstallationIdAssertion,
			helmChartNameAssertion,
			pathAssertion,
			opentofuVersionAssertion,
			terraformVersionAssertion,
			resource.TestCheckResourceAttr(resourceFullName, "is_terragrunt_run_all", strconv.FormatBool(template.IsTerragruntRunAll)),
			resource.TestCheckResourceAttr(resourceFullName, "is_azure_devops", strconv.FormatBool(template.IsAzureDevOps)),
			resource.TestCheckResourceAttr(resourceFullName, "is_helm_repository", strconv.FormatBool(template.IsHelmRepository)),
			tokenNameAssertion,
			resource.TestCheckResourceAttr(resourceFullName, "is_gitlab", strconv.FormatBool(template.IsGitlab)),
			ansibleVersionAssertion,
		)
	}

	var templateUseCases = []struct {
		vcs             string
		template        client.Template
		updatedTemplate client.Template
	}{
		{"GitLab EE", gleeTemplate, gleeUpdatedTemplate},
		{"GitLab", gitlabTemplate, gitlabUpdatedTemplate},
		{"GitHub", githubTemplate, githubUpdatedTemplate},
		{"Bitbucket", bitbucketTemplate, bitbucketUpdatedTemplate},
		{"GitHub EE", gheeTemplate, gheeUpdatedTemplate},
		{"Bitbucket Server", bitbucketServerTemplate, bitbucketServerUpdatedTemplate},
		{"Cloudformation", cloudformationTemplate, cloudformationUpdatedTemplate},
		{"Azure DevOps", azureDevOpsTemplate, azureDevOpsUpdatedTemplate},
		{"Helm Chart", helmTemplate, helmUpdatedTemplate},
		{"Opentofu", opentofuTemplate, opentofuUpdatedTemplate},
		{"Ansible", ansibleTemplate, ansibleUpdatedTemplate},
	}

	for _, templateUseCase := range templateUseCases {
		t.Run("Full "+templateUseCase.vcs+" template (without SSH keys)", func(t *testing.T) {
			templateCreatePayload := client.TemplateCreatePayload{
				Name:                 templateUseCase.template.Name,
				Repository:           templateUseCase.template.Repository,
				Description:          templateUseCase.template.Description,
				GithubInstallationId: templateUseCase.template.GithubInstallationId,
				IsGitlabEnterprise:   templateUseCase.template.IsGitlabEnterprise,
				IsGitlab:             templateUseCase.template.IsGitlab,
				TokenId:              templateUseCase.template.TokenId,
				Path:                 templateUseCase.template.Path,
				Revision:             templateUseCase.template.Revision,
				Type:                 templateUseCase.template.Type,
				Retry:                templateUseCase.template.Retry,
				TerraformVersion:     templateUseCase.template.TerraformVersion,
				BitbucketClientKey:   templateUseCase.template.BitbucketClientKey,
				IsGithubEnterprise:   templateUseCase.template.IsGithubEnterprise,
				IsBitbucketServer:    templateUseCase.template.IsBitbucketServer,
				FileName:             templateUseCase.template.FileName,
				TerragruntVersion:    templateUseCase.template.TerragruntVersion,
				IsTerragruntRunAll:   templateUseCase.template.IsTerragruntRunAll,
				IsAzureDevOps:        templateUseCase.template.IsAzureDevOps,
				IsHelmRepository:     templateUseCase.template.IsHelmRepository,
				HelmChartName:        templateUseCase.template.HelmChartName,
				OpentofuVersion:      templateUseCase.template.OpentofuVersion,
				TokenName:            templateUseCase.template.TokenName,
				AnsibleVersion:       templateUseCase.template.AnsibleVersion,
			}

			updateTemplateCreateTemplate := client.TemplateCreatePayload{
				Name:                 templateUseCase.updatedTemplate.Name,
				Repository:           templateUseCase.updatedTemplate.Repository,
				Description:          templateUseCase.updatedTemplate.Description,
				GithubInstallationId: templateUseCase.updatedTemplate.GithubInstallationId,
				IsGitlabEnterprise:   templateUseCase.updatedTemplate.IsGitlabEnterprise,
				IsGitlab:             templateUseCase.updatedTemplate.IsGitlab,
				TokenId:              templateUseCase.updatedTemplate.TokenId,
				Path:                 templateUseCase.updatedTemplate.Path,
				Revision:             templateUseCase.updatedTemplate.Revision,
				Type:                 templateUseCase.updatedTemplate.Type,
				Retry:                templateUseCase.updatedTemplate.Retry,
				TerraformVersion:     templateUseCase.updatedTemplate.TerraformVersion,
				BitbucketClientKey:   templateUseCase.updatedTemplate.BitbucketClientKey,
				IsGithubEnterprise:   templateUseCase.updatedTemplate.IsGithubEnterprise,
				IsBitbucketServer:    templateUseCase.updatedTemplate.IsBitbucketServer,
				FileName:             templateUseCase.updatedTemplate.FileName,
				TerragruntVersion:    templateUseCase.updatedTemplate.TerragruntVersion,
				IsTerragruntRunAll:   templateUseCase.updatedTemplate.IsTerragruntRunAll,
				IsAzureDevOps:        templateUseCase.updatedTemplate.IsAzureDevOps,
				IsHelmRepository:     templateUseCase.updatedTemplate.IsHelmRepository,
				HelmChartName:        templateUseCase.updatedTemplate.HelmChartName,
				OpentofuVersion:      templateUseCase.updatedTemplate.OpentofuVersion,
				TokenName:            templateUseCase.updatedTemplate.TokenName,
				AnsibleVersion:       templateUseCase.updatedTemplate.AnsibleVersion,
			}

			if templateUseCase.template.Type == "" {
				templateCreatePayload.Type = string(testType)
				updateTemplateCreateTemplate.Type = string(testType)
			}

			if templateUseCase.template.Type == "terragrunt" {
				templateCreatePayload.TerragruntTfBinary = templateUseCase.template.TerragruntTfBinary
			}

			if templateUseCase.updatedTemplate.Type == "terragrunt" {
				updateTemplateCreateTemplate.TerragruntTfBinary = templateUseCase.updatedTemplate.TerragruntTfBinary
			}

			if templateUseCase.template.Type != "terraform" && templateUseCase.template.Type != "terragrunt" {
				templateCreatePayload.TerraformVersion = ""
				updateTemplateCreateTemplate.TerraformVersion = ""
			}

			if templateUseCase.vcs == "Cloudformation" {
				templateCreatePayload.Type = "cloudformation"
				updateTemplateCreateTemplate.Type = "cloudformation"
			}

			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: fullTemplateResourceConfig(resourceType, resourceName, templateUseCase.template),
						Check:  fullTemplateResourceCheck(resourceFullName, templateUseCase.template),
					},
					{
						Config: fullTemplateResourceConfig(resourceType, resourceName, templateUseCase.updatedTemplate),
						Check:  fullTemplateResourceCheck(resourceFullName, templateUseCase.updatedTemplate),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Template(templateUseCase.template.Id).Times(2).Return(templateUseCase.template, nil),        // 1 after create, 1 before update
					mock.EXPECT().Template(templateUseCase.template.Id).Times(1).Return(templateUseCase.updatedTemplate, nil), // 1 after update
				)
				mock.EXPECT().TemplateCreate(templateCreatePayload).Times(1).Return(templateUseCase.template, nil)
				mock.EXPECT().TemplateUpdate(templateUseCase.updatedTemplate.Id, updateTemplateCreateTemplate).Times(1).Return(templateUseCase.updatedTemplate, nil)
				mock.EXPECT().TemplateDelete(templateUseCase.updatedTemplate.Id).Times(1).Return(nil)
			})
		})
	}

	t.Run("Basic template", func(t *testing.T) {
		template := client.Template{
			Id:         "id0",
			Name:       "template0",
			Repository: "env0/repo",
		}

		templateWithDefaults := client.Template{
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}

		basicTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":              template.Name,
				"repository":        template.Repository,
				"type":              string(testType),
				"terraform_version": defaultVersion,
				"opentofu_version":  defaultOpentofuVersion,
			})
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: basicTemplateResourceConfig(resourceType, resourceName, template),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", string(testType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
						resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", defaultOpentofuVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(template.Id).AnyTimes().Return(templateWithDefaults, nil)
			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:            template.Name,
				Repository:      template.Repository,
				Type:            testType,
				OpentofuVersion: defaultOpentofuVersion,
			}).Times(1).Return(template, nil)
			mock.EXPECT().TemplateDelete(template.Id).Times(1).Return(nil)
		})
	})

	t.Run("Invalid type", func(t *testing.T) {
		template := client.Template{
			Name:       "template0",
			Repository: "env0/repo",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":       template.Name,
						"repository": template.Repository,
						"type":       "gruntyform",
					}),
					ExpectError: regexp.MustCompile(`must be one of`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("SSH keys", func(t *testing.T) {
		initialSshKey1 := client.TemplateSshKey{
			Id:   "ssh-key-id-1",
			Name: "ssh-key-name-1",
		}
		initialSshKey2 := client.TemplateSshKey{
			Id:   "ssh-key-id-2",
			Name: "ssh-key-name-2",
		}
		updatedSshKey1 := client.TemplateSshKey{
			Id:   "updated-ssh-key-id-1",
			Name: "updated-ssh-key-name-1",
		}
		updatedSshKey2 := client.TemplateSshKey{
			Id:   "updated-ssh-key-id-2",
			Name: "updated-ssh-key-name-2",
		}

		template := client.Template{
			Id:               "id0",
			Name:             "template0",
			Repository:       "env0/repo",
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			SshKeys:          []client.TemplateSshKey{initialSshKey1, initialSshKey2},
			OpentofuVersion:  defaultOpentofuVersion,
		}

		updatedTemplate := client.Template{
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			SshKeys:          []client.TemplateSshKey{updatedSshKey1, updatedSshKey2},
			OpentofuVersion:  defaultOpentofuVersion,
		}

		sshKeyTemplateResourceConfig := func(name string, repository string, sshKey1 client.TemplateSshKey, sshKey2 client.TemplateSshKey) string {
			return fmt.Sprintf(`
	resource "env0_template" "test" {
		name = "%s"
		repository = "%s"
		terraform_version = "%s"
		type = "%s"
		opentofu_version = "%s"
		ssh_keys = [{
			id   = "%s"
			name = "%s"
			}, {
			id   = "%s"
			name = "%s"
		}]
	}`, name, repository, defaultVersion, string(testType), defaultOpentofuVersion, sshKey1.Id, sshKey1.Name, sshKey2.Id, sshKey2.Name)
		}

		sshTemplateResourceCheck := func(resourceFullName string, template client.Template, sshKey1 client.TemplateSshKey, sshKey2 client.TemplateSshKey) resource.TestCheckFunc {
			return resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
				resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
				resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
				resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.0.id", sshKey1.Id),
				resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.0.name", sshKey1.Name),
				resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.1.id", sshKey2.Id),
				resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.1.name", sshKey2.Name),
			)
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: sshKeyTemplateResourceConfig(template.Name, template.Repository, initialSshKey1, initialSshKey2),
					Check:  sshTemplateResourceCheck(resourceFullName, template, initialSshKey1, initialSshKey2),
				},
				{
					Config: sshKeyTemplateResourceConfig(template.Name, template.Repository, updatedSshKey1, updatedSshKey2),
					Check:  sshTemplateResourceCheck(resourceFullName, updatedTemplate, updatedSshKey1, updatedSshKey2),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(template.Id).Times(2).Return(template, nil),        // 1 after create, 1 before update
				mock.EXPECT().Template(template.Id).Times(1).Return(updatedTemplate, nil), // 1 after update
			)
			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:            template.Name,
				Repository:      template.Repository,
				Type:            testType,
				SshKeys:         template.SshKeys,
				OpentofuVersion: defaultOpentofuVersion,
			}).Times(1).Return(template, nil)
			mock.EXPECT().TemplateUpdate(updatedTemplate.Id, client.TemplateCreatePayload{
				Name:            updatedTemplate.Name,
				Repository:      updatedTemplate.Repository,
				Type:            testType,
				SshKeys:         updatedTemplate.SshKeys,
				OpentofuVersion: defaultOpentofuVersion,
			}).Times(1).Return(updatedTemplate, nil)
			mock.EXPECT().TemplateDelete(template.Id).Times(1).Return(nil)
		})
	})

	t.Run("Invalid retry times field", func(t *testing.T) {
		testMatrix := map[string][]int{
			"retries_on_deploy":  {-1, 0, 4, 5},
			"retries_on_destroy": {-1, 0, 4, 5},
		}

		var testCases []resource.TestCase

		for attribute, amounts := range testMatrix {
			for _, amount := range amounts {
				testCases = append(testCases, resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", attribute: amount}),
							ExpectError: regexp.MustCompile("retries amount must be between 1 and 3"),
						},
					},
				})
			}
		}

		for _, testCase := range testCases {
			tc := testCase

			t.Run("Invalid retry times field", func(t *testing.T) {
				runUnitTest(t, tc, func(mockFunc *client.MockApiClientInterface) {})
			})
		}
	})

	t.Run("Invalid retry regex field", func(t *testing.T) {
		testMatrix := map[string]string{
			"retries_on_deploy":  "retry_on_deploy_only_when_matches_regex",
			"retries_on_destroy": "retry_on_destroy_only_when_matches_regex",
		}

		var testCases []resource.TestCase
		for timesAttribute, regexAttribute := range testMatrix {
			testCases = append(testCases, resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", regexAttribute: "bla"}),
						ExpectError: regexp.MustCompile(fmt.Sprintf("`%s,%s`\\s+must\\s+be\\s+specified", timesAttribute, regexAttribute)),
					},
				},
			})
		}

		for _, testCase := range testCases {
			tc := testCase

			t.Run("Invalid retry regex field", func(t *testing.T) {
				runUnitTest(t, tc, func(mockFunc *client.MockApiClientInterface) {})
			})
		}
	})

	var mixedUsecases = []struct {
		firstVcs  string
		secondVcs string
		tfObject  map[string]interface{}
		exception string
	}{
		{"GitLab", "GitHub", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "github_installation_id": 1, "token_id": "2"}, "\"github_installation_id\": conflicts with token_id"},
		{"GitLab", "GitLab EE", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "token_id": "2", "is_gitlab_enterprise": "true"}, "\"is_gitlab_enterprise\": conflicts with token_id"},
		{"GitHub", "GitLab EE", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "github_installation_id": 1, "is_gitlab_enterprise": "true"}, "\"github_installation_id\": conflicts with is_gitlab_enterprise"},
		{"GitHub", "Bitbucket", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "github_installation_id": 1, "bitbucket_client_key": "3"}, "\"bitbucket_client_key\": conflicts with github_installation_id"},
		{"GitLab", "Bitbucket", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "token_id": "2", "bitbucket_client_key": "3"}, "\"bitbucket_client_key\": conflicts with token_id"},
		{"GitLab EE", "GitHub EE", map[string]interface{}{"name": "test", "type": testType, "repository": "env0/test", "is_gitlab_enterprise": "true", "is_github_enterprise": "true"}, "\"is_github_enterprise\": conflicts with is_gitlab_enterprise"},
	}

	for _, mixUseCase := range mixedUsecases {
		t.Run("Mixed "+mixUseCase.firstVcs+" and "+mixUseCase.secondVcs+" template", func(t *testing.T) {
			var testCases []resource.TestCase

			testCases = append(testCases, resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      resourceConfigCreate(resourceType, resourceName, mixUseCase.tfObject),
						ExpectError: regexp.MustCompile(mixUseCase.exception),
					},
				},
			})

			for _, testCase := range testCases {
				runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
			}
		})
	}

	t.Run("detect drift", func(t *testing.T) {
		template := client.Template{
			Id:         "id0",
			Name:       "template0",
			Repository: "env0/repo",
			Type:       string(testType),
		}

		updateTemplate := client.Template{
			Id:         "id0-update",
			Name:       "template0-update",
			Repository: "env0/repo-update",
			Type:       string(testType),
		}

		templateWithDefaults := client.Template{
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}
		templateWithDefaultsUpdate := client.Template{
			Id:               updateTemplate.Id,
			Name:             updateTemplate.Name,
			Repository:       updateTemplate.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}

		templateWithDrift := client.Template{
			IsDeleted:        true,
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}

		basicTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":              template.Name,
				"type":              string(template.Type),
				"repository":        template.Repository,
				"terraform_version": defaultVersion,
				"opentofu_version":  defaultOpentofuVersion,
			})
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: basicTemplateResourceConfig(resourceType, resourceName, template),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", string(testType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
						resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", defaultOpentofuVersion),
					),
				},
				{
					Config: basicTemplateResourceConfig(resourceType, resourceName, updateTemplate),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", updateTemplate.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", updateTemplate.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", updateTemplate.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", string(testType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
						resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", defaultOpentofuVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:            template.Name,
				Repository:      template.Repository,
				Type:            testType,
				OpentofuVersion: defaultOpentofuVersion,
			}).Times(1).Return(template, nil)

			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:            updateTemplate.Name,
				Repository:      updateTemplate.Repository,
				Type:            testType,
				OpentofuVersion: defaultOpentofuVersion,
			}).Times(1).Return(updateTemplate, nil)

			gomock.InOrder(
				mock.EXPECT().Template(template.Id).Times(1).Return(templateWithDefaults, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(templateWithDrift, nil),
				mock.EXPECT().Template(updateTemplate.Id).Times(1).Return(templateWithDefaultsUpdate, nil),
			)
			mock.EXPECT().TemplateDelete(updateTemplate.Id).Times(1).Return(nil)
		})
	})

	// https://github.com/env0/terraform-provider-env0/issues/699
	t.Run("path drift", func(t *testing.T) {
		pathTemplate := client.Template{
			Id:               "id0",
			Name:             "template0",
			Path:             "path/zero",
			Repository:       "repo",
			TerraformVersion: string(defaultVersion),
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}

		updatedPathTemplate := client.Template{
			Id:               "id0",
			Name:             "template0",
			Path:             "path/one",
			Repository:       "repo",
			TerraformVersion: string(defaultVersion),
			Type:             string(testType),
			OpentofuVersion:  defaultOpentofuVersion,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              pathTemplate.Name,
						"path":              "/" + pathTemplate.Path,
						"repository":        pathTemplate.Repository,
						"terraform_version": pathTemplate.TerraformVersion,
						"opentofu_version":  defaultOpentofuVersion,
						"type":              pathTemplate.Type,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", pathTemplate.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", pathTemplate.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", pathTemplate.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", pathTemplate.Type),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", pathTemplate.TerraformVersion),
						resource.TestCheckResourceAttr(resourceFullName, "path", "/"+pathTemplate.Path),
						resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", defaultOpentofuVersion),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              updatedPathTemplate.Name,
						"path":              "/" + updatedPathTemplate.Path,
						"repository":        updatedPathTemplate.Repository,
						"terraform_version": updatedPathTemplate.TerraformVersion,
						"opentofu_version":  defaultOpentofuVersion,
						"type":              updatedPathTemplate.Type,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", updatedPathTemplate.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", updatedPathTemplate.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", updatedPathTemplate.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", updatedPathTemplate.Type),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", updatedPathTemplate.TerraformVersion),
						resource.TestCheckResourceAttr(resourceFullName, "path", "/"+updatedPathTemplate.Path),
						resource.TestCheckResourceAttr(resourceFullName, "opentofu_version", defaultOpentofuVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
					Path:            "/" + pathTemplate.Path,
					Name:            pathTemplate.Name,
					Type:            pathTemplate.Type,
					Repository:      pathTemplate.Repository,
					OpentofuVersion: defaultOpentofuVersion,
				}).Times(1).Return(pathTemplate, nil),
				mock.EXPECT().Template(pathTemplate.Id).Times(2).Return(pathTemplate, nil),
				mock.EXPECT().TemplateUpdate(updatedPathTemplate.Id, client.TemplateCreatePayload{
					Path:            "/" + updatedPathTemplate.Path,
					Name:            updatedPathTemplate.Name,
					Type:            updatedPathTemplate.Type,
					Repository:      updatedPathTemplate.Repository,
					OpentofuVersion: defaultOpentofuVersion,
				}).Times(1).Return(updatedPathTemplate, nil),
				mock.EXPECT().Template(pathTemplate.Id).Times(1).Return(updatedPathTemplate, nil),
				mock.EXPECT().TemplateDelete(pathTemplate.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Invalid Terraform Version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              "template0",
						"repository":        "env0/repo",
						"type":              "terraform",
						"token_id":          "abcdefg",
						"terraform_version": "v0.15.1",
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Invalid Opentofu Version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":             "template0",
						"repository":       "env0/repo",
						"type":             "opentofu",
						"token_id":         "abcdefg",
						"opentofu_version": "v0.20.1",
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Opentofu type with no Opentofu version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":       "template0",
						"repository": "env0/repo",
						"type":       "opentofu",
						"token_id":   "abcdefg",
					}),
					ExpectError: regexp.MustCompile("must supply opentofu version"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Cloudformation type with no file_name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              "template0",
						"repository":        "env0/repo",
						"type":              "cloudformation",
						"terraform_version": "0.15.1",
					}),
					ExpectError: regexp.MustCompile("file_name is required with cloudformation template type"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Non-cloudformation type with file_name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              "template0",
						"repository":        "env0/repo",
						"type":              "terraform",
						"terraform_version": "0.15.1",
						"file_name":         "bad.yaml",
					}),
					ExpectError: regexp.MustCompile("file_name cannot be set when template type is: terraform"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("terragrunt type with no terragrunt version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              "template0",
						"repository":        "env0/repo",
						"type":              "terragrunt",
						"terraform_version": "0.15.1",
					}),
					ExpectError: regexp.MustCompile("must supply terragrunt version"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("terragrunt version with non-terragrunt type", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":               "template0",
						"repository":         "env0/repo",
						"type":               "terraform",
						"terraform_version":  "0.15.1",
						"terragrunt_version": "0.31.1",
					}),
					ExpectError: regexp.MustCompile("can't define terragrunt version for non-terragrunt template"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("run all with outdated terragrunt version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  "template0",
						"repository":            "env0/repo",
						"type":                  "terragrunt",
						"terraform_version":     "0.15.1",
						"terragrunt_version":    "0.27.50",
						"terragrunt_tf_binary":  "terraform",
						"is_terragrunt_run_all": "true",
					}),
					ExpectError: regexp.MustCompile(`can't set is_terragrunt_run_all to 'true' for terragrunt versions lower than 0.28.1`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("terragrunt_tf_binary set with a non terragrunt template type", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                 "template0",
						"repository":           "env0/repo",
						"type":                 "terraform",
						"terraform_version":    "0.15.1",
						"terragrunt_tf_binary": "opentofu",
					}),
					ExpectError: regexp.MustCompile(`can't define terragrunt_tf_binary for non-terragrunt template`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("run with terragrunt without an opentofu version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  "template0",
						"repository":            "env0/repo",
						"type":                  "terragrunt",
						"terragrunt_version":    "0.56.50",
						"is_terragrunt_run_all": "true",
					}),
					ExpectError: regexp.MustCompile("must supply opentofu version"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("run terragrunt with terraform binary and no terraform version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                  "template0",
						"repository":            "env0/repo",
						"type":                  "terragrunt",
						"terragrunt_version":    "0.56.50",
						"is_terragrunt_run_all": "true",
						"terragrunt_tf_binary":  "terraform",
					}),
					ExpectError: regexp.MustCompile("must supply terraform version"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("invalid ansible version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":            "template",
						"repository":      "env0/repo",
						"type":            "ansible",
						"ansible_version": "not-valid",
					}),
					ExpectError: regexp.MustCompile("invalid ansible version 'not-valid'"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("unsupported ansible version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":            "template",
						"repository":      "env0/repo",
						"type":            "ansible",
						"ansible_version": "2.5.6",
					}),
					ExpectError: regexp.MustCompile("supported ansible versions are 3.0.0 and above"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("ansible type with no version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":       "template",
						"repository": "env0/repo",
						"type":       "ansible",
					}),
					ExpectError: regexp.MustCompile("'ansible_version' is required"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
