package env0

import (
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitTemplateData(t *testing.T) {
	resourceType := "env0_template"
	resourceName := "test"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)
	onDeploy := client.TemplateRetryOn{
		Times:      2,
		ErrorRegex: "error retry on deploy",
	}
	onDestroy := client.TemplateRetryOn{
		Times:      3,
		ErrorRegex: "error retry on destroy",
	}
	templateRetry := client.TemplateRetry{
		OnDeploy:  &onDeploy,
		OnDestroy: &onDestroy,
	}

	template := client.Template{
		Id:                   "id0",
		Name:                 "name0",
		Repository:           "repository",
		Path:                 "path",
		Revision:             "revision",
		Type:                 "terraform",
		TerraformVersion:     "0.15.1",
		Retry:                templateRetry,
		ProjectIds:           []string{"pId1", "pId2"},
		GithubInstallationId: 123,
		BitbucketClientKey:   "12345",
		SshKeys: []client.TemplateSshKey{
			{Id: "id", Name: "name"},
		},
		IsBitbucketServer: true,
		IsAzureDevOps:     true,
		TokenName:         "tokenname",
	}

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "path", template.Path),
						resource.TestCheckResourceAttr(resourceFullName, "revision", template.Revision),
						resource.TestCheckResourceAttr(resourceFullName, "type", template.Type),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),
						resource.TestCheckResourceAttr(resourceFullName, "retries_on_deploy", strconv.Itoa(template.Retry.OnDeploy.Times)),
						resource.TestCheckResourceAttr(resourceFullName, "retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex),
						resource.TestCheckResourceAttr(resourceFullName, "retries_on_destroy", strconv.Itoa(template.Retry.OnDestroy.Times)),
						resource.TestCheckResourceAttr(resourceFullName, "retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex),
						resource.TestCheckResourceAttr(resourceFullName, "github_installation_id", strconv.Itoa(template.GithubInstallationId)),
						resource.TestCheckResourceAttr(resourceFullName, "project_ids.0", template.ProjectIds[0]),
						resource.TestCheckResourceAttr(resourceFullName, "project_ids.1", template.ProjectIds[1]),
						resource.TestCheckResourceAttr(resourceFullName, "bitbucket_client_key", template.BitbucketClientKey),
						resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.0.id", template.SshKeys[0].Id),
						resource.TestCheckResourceAttr(resourceFullName, "ssh_keys.0.name", template.SshKeys[0].Name),
						resource.TestCheckResourceAttr(resourceFullName, "is_bitbucket_server", "true"),
						resource.TestCheckResourceAttr(resourceFullName, "is_azure_devops", "true"),
						resource.TestCheckResourceAttr(resourceFullName, "token_name", template.TokenName),
					),
				},
			},
		}
	}

	t.Run("Template By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(map[string]any{"id": template.Id}),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(template.Id).AnyTimes().Return(template, nil)
			})
	})

	t.Run("Template By Name", func(t *testing.T) {
		deletedTemplate := template
		deletedTemplate.IsDeleted = true

		runUnitTest(t,
			getValidTestCase(map[string]any{"name": template.Name}),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Templates().AnyTimes().Return([]client.Template{template, deletedTemplate}, nil)
			})
	})

	t.Run("Template with VCS Connection ID", func(t *testing.T) {
		vcsTemplate := client.Template{
			Id:               "id1",
			Name:             "vcs-template",
			Repository:       "repository",
			Path:             "path",
			Revision:         "revision",
			Type:             "terraform",
			TerraformVersion: "0.15.1",
			Retry:            templateRetry,
			ProjectIds:       []string{"pId1"},
			VcsConnectionId:  "vcs-conn-123",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"id": vcsTemplate.Id}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", vcsTemplate.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", vcsTemplate.Name),
						resource.TestCheckResourceAttr(resourceFullName, "vcs_connection_id", vcsTemplate.VcsConnectionId),
						resource.TestCheckNoResourceAttr(resourceFullName, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(vcsTemplate.Id).AnyTimes().Return(vcsTemplate, nil)
		})
	})
}
