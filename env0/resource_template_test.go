package env0

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitTemplateResource(t *testing.T) {
	const resourceType = "env0_template"
	const resourceName = "test"
	const defaultVersion = "0.15.1"
	const defaultType = client.TemplateTypeTerraform

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
		Type:               "terraform",
		IsGitlabEnterprise: true,
		TerraformVersion:   "0.12.24",
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
		IsGitlabEnterprise: true,
		TerraformVersion:   "0.15.1",
	}
	gitlabTemplate := client.Template{
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
		Type:             "terraform",
		TokenId:          "1",
		GitlabProjectId:  10,
		TerraformVersion: "0.12.24",
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
		Type:             "terragrunt",
		TokenId:          "2",
		GitlabProjectId:  2,
		TerraformVersion: "0.15.1",
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
		GithubInstallationId: 2,
		TerraformVersion:     "0.15.1",
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
		TerraformVersion:   "0.15.1",
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
		IsGitHubEnterprise: true,
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
		IsGitHubEnterprise: true,
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
	fullTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
		templateAsDictionary := map[string]interface{}{
			"name":       template.Name,
			"repository": template.Repository,
		}

		if template.Type != "" {
			templateAsDictionary["type"] = template.Type
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
		if template.TokenId != "" {
			templateAsDictionary["token_id"] = template.TokenId
		}
		if template.GitlabProjectId != 0 {
			templateAsDictionary["gitlab_project_id"] = template.GitlabProjectId
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
		if template.IsGitHubEnterprise != false {
			templateAsDictionary["is_github_enterprise"] = template.IsGitHubEnterprise
		}
		if template.IsBitbucketServer != false {
			templateAsDictionary["is_bitbucket_server"] = template.IsBitbucketServer
		}

		return resourceConfigCreate(resourceType, resourceName, templateAsDictionary)
	}
	fullTemplateResourceCheck := func(resourceFullName string, template client.Template) resource.TestCheckFunc {
		gitlabProjectIdAssertion := resource.TestCheckResourceAttr(resourceFullName, "gitlab_project_id", strconv.Itoa(template.GitlabProjectId))
		if template.GitlabProjectId == 0 {
			gitlabProjectIdAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "gitlab_project_id")
		}

		tokenIdAssertion := resource.TestCheckResourceAttr(resourceFullName, "token_id", template.TokenId)
		if template.TokenId == "" {
			tokenIdAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "token_id")
		}

		githubInstallationIdAssertion := resource.TestCheckResourceAttr(resourceFullName, "github_installation_id", strconv.Itoa(template.GithubInstallationId))
		if template.GithubInstallationId == 0 {
			githubInstallationIdAssertion = resource.TestCheckNoResourceAttr(resourceFullName, "github_installation_id")
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
			resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
			resource.TestCheckResourceAttr(resourceFullName, "description", template.Description),
			resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
			resource.TestCheckResourceAttr(resourceFullName, "path", template.Path),
			resource.TestCheckResourceAttr(resourceFullName, "type", template.Type),
			resource.TestCheckResourceAttr(resourceFullName, "retries_on_deploy", strconv.Itoa(template.Retry.OnDeploy.Times)),
			resource.TestCheckResourceAttr(resourceFullName, "retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex),
			resource.TestCheckResourceAttr(resourceFullName, "retries_on_destroy", strconv.Itoa(template.Retry.OnDestroy.Times)),
			resource.TestCheckResourceAttr(resourceFullName, "retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex),
			resource.TestCheckResourceAttr(resourceFullName, "is_gitlab_enterprise", strconv.FormatBool(template.IsGitlabEnterprise)),
			tokenIdAssertion,
			gitlabProjectIdAssertion,
			githubInstallationIdAssertion,
			resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),
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
	}
	for _, templateUseCase := range templateUseCases {
		t.Run("Full "+templateUseCase.vcs+" template (without SSH keys)", func(t *testing.T) {
			templateCreatePayload := client.TemplateCreatePayload{
				Name:                 templateUseCase.template.Name,
				Repository:           templateUseCase.template.Repository,
				Description:          templateUseCase.template.Description,
				GithubInstallationId: templateUseCase.template.GithubInstallationId,
				IsGitlabEnterprise:   templateUseCase.template.IsGitlabEnterprise,
				IsGitLab:             templateUseCase.template.TokenId != "",
				GitlabProjectId:      templateUseCase.template.GitlabProjectId,
				TokenId:              templateUseCase.template.TokenId,
				Path:                 templateUseCase.template.Path,
				Revision:             templateUseCase.template.Revision,
				Type:                 client.TemplateTypeTerraform,
				Retry:                templateUseCase.template.Retry,
				TerraformVersion:     templateUseCase.template.TerraformVersion,
				BitbucketClientKey:   templateUseCase.template.BitbucketClientKey,
				IsGitHubEnterprise:   templateUseCase.template.IsGitHubEnterprise,
				IsBitbucketServer:    templateUseCase.template.IsBitbucketServer,
			}
			updateTemplateCreateTemplate := client.TemplateCreatePayload{
				Name:                 templateUseCase.updatedTemplate.Name,
				Repository:           templateUseCase.updatedTemplate.Repository,
				Description:          templateUseCase.updatedTemplate.Description,
				GithubInstallationId: templateUseCase.updatedTemplate.GithubInstallationId,
				IsGitlabEnterprise:   templateUseCase.updatedTemplate.IsGitlabEnterprise,
				IsGitLab:             templateUseCase.updatedTemplate.TokenId != "",
				GitlabProjectId:      templateUseCase.updatedTemplate.GitlabProjectId,
				TokenId:              templateUseCase.updatedTemplate.TokenId,
				Path:                 templateUseCase.updatedTemplate.Path,
				Revision:             templateUseCase.updatedTemplate.Revision,
				Type:                 client.TemplateType(templateUseCase.updatedTemplate.Type),
				Retry:                templateUseCase.updatedTemplate.Retry,
				TerraformVersion:     templateUseCase.updatedTemplate.TerraformVersion,
				BitbucketClientKey:   templateUseCase.updatedTemplate.BitbucketClientKey,
				IsGitHubEnterprise:   templateUseCase.updatedTemplate.IsGitHubEnterprise,
				IsBitbucketServer:    templateUseCase.updatedTemplate.IsBitbucketServer,
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
			Type:             string(defaultType),
		}

		basicTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":       template.Name,
				"repository": template.Repository,
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
						resource.TestCheckResourceAttr(resourceFullName, "type", string(defaultType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(template.Id).AnyTimes().Return(templateWithDefaults, nil)
			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:             template.Name,
				Repository:       template.Repository,
				Type:             defaultType,
				TerraformVersion: defaultVersion,
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
					ExpectError: regexp.MustCompile(`'type' can either be 'terraform' or 'terragrunt'`),
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
			Type:             string(defaultType),
			SshKeys:          []client.TemplateSshKey{initialSshKey1, initialSshKey2},
		}

		updatedTemplate := client.Template{
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(defaultType),
			SshKeys:          []client.TemplateSshKey{updatedSshKey1, updatedSshKey2},
		}

		sshKeyTemplateResourceConfig := func(name string, repository string, sshKey1 client.TemplateSshKey, sshKey2 client.TemplateSshKey) string {
			return fmt.Sprintf(`
	resource "env0_template" "test" {
		name = "%s"
		repository = "%s"
		terraform_version = "%s"
		type = "%s"
		ssh_keys = [{
			id   = "%s"
			name = "%s"
			}, {
			id   = "%s"
			name = "%s"
		}]
	}`, name, repository, defaultVersion, string(defaultType), sshKey1.Id, sshKey1.Name, sshKey2.Id, sshKey2.Name)
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
				Name:             template.Name,
				Repository:       template.Repository,
				Type:             defaultType,
				TerraformVersion: defaultVersion,
				SshKeys:          template.SshKeys,
			}).Times(1).Return(template, nil)
			mock.EXPECT().TemplateUpdate(updatedTemplate.Id, client.TemplateCreatePayload{
				Name:             updatedTemplate.Name,
				Repository:       updatedTemplate.Repository,
				Type:             defaultType,
				TerraformVersion: defaultVersion,
				SshKeys:          updatedTemplate.SshKeys,
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
							Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "test", "repository": "env0/test", attribute: amount}),
							ExpectError: regexp.MustCompile("Retries amount must be between 1 and 3"),
						},
					},
				})
			}
		}

		for _, testCase := range testCases {
			runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
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
						Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "test", "repository": "env0/test", regexAttribute: "bla"}),
						ExpectError: regexp.MustCompile(fmt.Sprintf("`%s,%s`\\s+must\\s+be\\s+specified", timesAttribute, regexAttribute)),
					},
				},
			})
		}

		for _, testCase := range testCases {
			runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
		}
	})

	var mixedUsecases = []struct {
		firstVcs  string
		secondVcs string
		tfObject  map[string]interface{}
		exception string
	}{
		{"GitLab", "GitHub", map[string]interface{}{"name": "test", "repository": "env0/test", "github_installation_id": 1, "token_id": "2"}, "\"github_installation_id\": conflicts with token_id"},
		{"GitLab", "GitLab EE", map[string]interface{}{"name": "test", "repository": "env0/test", "token_id": "2", "is_gitlab_enterprise": "true"}, "\"is_gitlab_enterprise\": conflicts with token_id"},
		{"GitHub", "GitLab EE", map[string]interface{}{"name": "test", "repository": "env0/test", "github_installation_id": 1, "is_gitlab_enterprise": "true"}, "\"github_installation_id\": conflicts with is_gitlab_enterprise"},
		{"GitHub", "Bitbucket", map[string]interface{}{"name": "test", "repository": "env0/test", "github_installation_id": 1, "bitbucket_client_key": "3"}, "\"bitbucket_client_key\": conflicts with github_installation_id"},
		{"GitLab", "Bitbucket", map[string]interface{}{"name": "test", "repository": "env0/test", "token_id": "2", "bitbucket_client_key": "3"}, "\"bitbucket_client_key\": conflicts with token_id"},
		{"GitLab EE", "GitHub EE", map[string]interface{}{"name": "test", "repository": "env0/test", "is_gitlab_enterprise": "true", "is_github_enterprise": "true"}, "\"is_github_enterprise\": conflicts with is_gitlab_enterprise"},
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

	t.Run("Should not trigger terraform changes when gitlab_project_id is provided", func(t *testing.T) {
		template := client.Template{
			Id:               "id0",
			Name:             "template0",
			Repository:       "env0/repo",
			Type:             "terraform",
			GitlabProjectId:  123456,
			TokenId:          "abcdefg",
			TerraformVersion: defaultVersion,
		}

		tfConfig := fullTemplateResourceConfig(resourceType, resourceName, template)
		var testCase = resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: tfConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "gitlab_project_id", strconv.Itoa(template.GitlabProjectId)),
					),
				},
				{
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
					Config:             tfConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "gitlab_project_id", strconv.Itoa(template.GitlabProjectId)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(template.Id).Times(3).Return(template, nil) // 1 after create, 1 before update, 1 after update
			mock.EXPECT().TemplateCreate(gomock.Any()).Times(1).Return(template, nil)
			mock.EXPECT().TemplateDelete(template.Id).Times(1).Return(nil)
		})
	})

	t.Run("detect drift", func(t *testing.T) {
		template := client.Template{
			Id:         "id0",
			Name:       "template0",
			Repository: "env0/repo",
		}

		updateTemplate := client.Template{
			Id:         "id0-update",
			Name:       "template0-update",
			Repository: "env0/repo-update",
		}

		templateWithDefaults := client.Template{
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(defaultType),
		}
		templateWithDefaultsUpdate := client.Template{
			Id:               updateTemplate.Id,
			Name:             updateTemplate.Name,
			Repository:       updateTemplate.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(defaultType),
		}

		templateWithDrift := client.Template{
			IsDeleted:        true,
			Id:               template.Id,
			Name:             template.Name,
			Repository:       template.Repository,
			TerraformVersion: defaultVersion,
			Type:             string(defaultType),
		}

		basicTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":       template.Name,
				"repository": template.Repository,
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
						resource.TestCheckResourceAttr(resourceFullName, "type", string(defaultType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
					),
				},
				{
					Config: basicTemplateResourceConfig(resourceType, resourceName, updateTemplate),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", updateTemplate.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", updateTemplate.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", updateTemplate.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "type", string(defaultType)),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", defaultVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {

			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:             template.Name,
				Repository:       template.Repository,
				Type:             defaultType,
				TerraformVersion: defaultVersion,
			}).Times(1).Return(template, nil)

			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:             updateTemplate.Name,
				Repository:       updateTemplate.Repository,
				Type:             defaultType,
				TerraformVersion: defaultVersion,
			}).Times(1).Return(updateTemplate, nil)

			gomock.InOrder(
				mock.EXPECT().Template(template.Id).Times(1).Return(templateWithDefaults, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(templateWithDrift, nil),
				mock.EXPECT().Template(updateTemplate.Id).Times(1).Return(templateWithDefaultsUpdate, nil),
			)
			mock.EXPECT().TemplateDelete(updateTemplate.Id).Times(1).Return(nil)
		})
	})

	t.Run("Invalid Terraform Version", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"id":                "id0",
						"name":              "template0",
						"repository":        "env0/repo",
						"type":              "terraform",
						"gitlab_project_id": 123456,
						"token_id":          "abcdefg",
						"terraform_version": "v0.15.1",
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
