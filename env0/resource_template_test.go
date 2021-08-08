package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"strconv"
	"testing"
)

func TestUnitTemplateResource(t *testing.T) {
	const resourceType = "env0_template"
	const resourceName = "test"
	const defaultVersion = "0.15.1"
	const defaultType = client.TemplateTypeTerraform

	var resourceFullName = resourceAccessor(resourceType, resourceName)

	t.Run("Full Github template (without SSH keys)", func(t *testing.T) {
		template := client.Template{
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

		updatedTemplate := client.Template{
			Id:          template.Id,
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

		fullTemplateResourceCheck := func(resourceFullName string, template client.Template) resource.TestCheckFunc {
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
				resource.TestCheckResourceAttr(resourceFullName, "github_installation_id", strconv.Itoa(template.GithubInstallationId)),
				resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),
			)
		}

		fullTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":              template.Name,
				"description":       template.Description,
				"repository":        template.Repository,
				"path":              template.Path,
				"revision":          template.Revision,
				"type":              template.Type,
				"retries_on_deploy": template.Retry.OnDeploy.Times,
				"retry_on_deploy_only_when_matches_regex":  template.Retry.OnDeploy.ErrorRegex,
				"retries_on_destroy":                       template.Retry.OnDestroy.Times,
				"retry_on_destroy_only_when_matches_regex": template.Retry.OnDestroy.ErrorRegex,
				"github_installation_id":                   template.GithubInstallationId,
				"terraform_version":                        template.TerraformVersion,
			})
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fullTemplateResourceConfig(resourceType, resourceName, template),
					Check:  fullTemplateResourceCheck(resourceFullName, template),
				},
				{
					Config: fullTemplateResourceConfig(resourceType, resourceName, updatedTemplate),
					Check:  fullTemplateResourceCheck(resourceFullName, updatedTemplate),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(template.Id).Times(2).Return(template, nil),        // 1 after create, 1 before update
				mock.EXPECT().Template(template.Id).Times(1).Return(updatedTemplate, nil), // 1 after update
			)
			mock.EXPECT().TemplateCreate(client.TemplateCreatePayload{
				Name:                 template.Name,
				Repository:           template.Repository,
				Description:          template.Description,
				GithubInstallationId: template.GithubInstallationId,
				IsGitLab:             false,
				Path:                 template.Path,
				Revision:             template.Revision,
				Type:                 client.TemplateTypeTerraform,
				Retry:                template.Retry,
				TerraformVersion:     template.TerraformVersion,
			}).Times(1).Return(template, nil)
			mock.EXPECT().TemplateUpdate(template.Id, client.TemplateCreatePayload{
				Name:                 updatedTemplate.Name,
				Repository:           updatedTemplate.Repository,
				Description:          updatedTemplate.Description,
				GithubInstallationId: updatedTemplate.GithubInstallationId,
				IsGitLab:             false,
				Path:                 updatedTemplate.Path,
				Revision:             updatedTemplate.Revision,
				Type:                 client.TemplateTypeTerragrunt,
				Retry:                updatedTemplate.Retry,
				TerraformVersion:     updatedTemplate.TerraformVersion,
			}).Times(1).Return(updatedTemplate, nil)
			mock.EXPECT().TemplateDelete(template.Id).Times(1).Return(nil)
		})
	})

	t.Run("Full Gitlab template (without SSH keys)", func(t *testing.T) {
		template := client.Template{
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
			GitlabProjectId:  1,
			TerraformVersion: "0.12.24",
		}

		updatedTemplate := client.Template{
			Id:          template.Id,
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

		fullTemplateResourceCheck := func(resourceFullName string, template client.Template) resource.TestCheckFunc {
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
				resource.TestCheckResourceAttr(resourceFullName, "token_id", template.TokenId),
				resource.TestCheckResourceAttr(resourceFullName, "gitlab_project_id", strconv.Itoa(template.GitlabProjectId)),
				resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),
			)
		}

		fullTemplateResourceConfig := func(resourceType string, resourceName string, template client.Template) string {
			return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
				"name":              template.Name,
				"description":       template.Description,
				"repository":        template.Repository,
				"path":              template.Path,
				"revision":          template.Revision,
				"type":              template.Type,
				"retries_on_deploy": template.Retry.OnDeploy.Times,
				"retry_on_deploy_only_when_matches_regex":  template.Retry.OnDeploy.ErrorRegex,
				"retries_on_destroy":                       template.Retry.OnDestroy.Times,
				"retry_on_destroy_only_when_matches_regex": template.Retry.OnDestroy.ErrorRegex,
				"token_id":          template.TokenId,
				"gitlab_project_id": template.GitlabProjectId,
				"terraform_version": template.TerraformVersion,
			})
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fullTemplateResourceConfig(resourceType, resourceName, template),
					Check:  fullTemplateResourceCheck(resourceFullName, template),
				},
				{
					Config: fullTemplateResourceConfig(resourceType, resourceName, updatedTemplate),
					Check:  fullTemplateResourceCheck(resourceFullName, updatedTemplate),
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
				Description:      template.Description,
				TokenId:          template.TokenId,
				GitlabProjectId:  template.GitlabProjectId,
				IsGitLab:         true,
				Path:             template.Path,
				Revision:         template.Revision,
				Type:             client.TemplateTypeTerraform,
				Retry:            template.Retry,
				TerraformVersion: template.TerraformVersion,
			}).Times(1).Return(template, nil)
			mock.EXPECT().TemplateUpdate(template.Id, client.TemplateCreatePayload{
				Name:             updatedTemplate.Name,
				Repository:       updatedTemplate.Repository,
				Description:      updatedTemplate.Description,
				TokenId:          updatedTemplate.TokenId,
				GitlabProjectId:  updatedTemplate.GitlabProjectId,
				IsGitLab:         true,
				Path:             updatedTemplate.Path,
				Revision:         updatedTemplate.Revision,
				Type:             client.TemplateTypeTerragrunt,
				Retry:            updatedTemplate.Retry,
				TerraformVersion: updatedTemplate.TerraformVersion,
			}).Times(1).Return(updatedTemplate, nil)
			mock.EXPECT().TemplateDelete(template.Id).Times(1).Return(nil)
		})
	})

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
						ExpectError: regexp.MustCompile(fmt.Sprintf("`%s,%s` must be specified", timesAttribute, regexAttribute)),
					},
				},
			})
		}

	})

	t.Run("Gitlab and Github template", func(t *testing.T) {
		var testCases []resource.TestCase

		testCases = append(testCases, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "test", "repository": "env0/test", "github_installation_id": 1, "token_id": "2"}),
					ExpectError: regexp.MustCompile("Cannot set token_id and github_installation_id for the same template"),
				},
			},
		})

		for _, testCase := range testCases {
			runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
		}
	})
}
