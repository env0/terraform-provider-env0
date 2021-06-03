package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitTemplateData(t *testing.T) {
	resourceType := "env0_template"
	resourceName := "test"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)
	/*retryOnDeploy := client.TemplateRetryOn{
		Times: 2,
	}
	retryOnDestroy := client.TemplateRetryOn{
		Times: 2,
	}
	templateRetry := client.TemplateRetry{
		OnDeploy:  &retryOnDeploy,
		OnDestroy: &retryOnDestroy,
	}*/
	/*
		Author               User             `json:"author"`
			AuthorId             string           `json:"authorId"`
			CreatedAt            string           `json:"createdAt"`
			Href                 string           `json:"href"`
			Id                   string           `json:"id"`
			Name                 string           `json:"name"`
			Description          string           `json:"description"`
			OrganizationId       string           `json:"organizationId"`
			Path                 string           `json:"path"`
			Revision             string           `json:"revision"`
			ProjectId            string           `json:"projectId"`
			ProjectIds           []string         `json:"projectIds"`
			Repository           string           `json:"repository"`
			Retry                TemplateRetry    `json:"retry"`
			SshKeys              []TemplateSshKey `json:"sshKeys"`
			Type                 string           `json:"type"`
			GithubInstallationId int              `json:"githubInstallationId"`
			UpdatedAt            string           `json:"updatedAt"`
			TerraformVersion     string           `json:"terraformVersion"`
	*/
	template := client.Template{
		Id:                   "id0",
		AuthorId:       "author",
		CreatedAt:      "createdAt",
		Href:           "href",
		Name:           "name0",
		Description:    "description",
		OrganizationId: "organizationId",
		Path:           "path",
		Revision:       "revision",
		ProjectId:      "projectId",
		Repository:     "repository",
		//Retry:                templateRetry,
		Type:                 "terraform",
		GithubInstallationId: 123,
		UpdatedAt:            "updatedAt",
		TerraformVersion:     "0.15.1",
	}

	templateByName := map[string]string{
		"name": template.Name,
	}

	templateById := map[string]string{
		"id": template.Id,
	}

	runScenario := func(input map[string]string, mockFunc func(mockFunc *client.MockApiClientInterface)) {
		testCase := resource.TestCase{
			ProviderFactories: testUnitProviders,
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						//resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
						//resource.TestCheckResourceAttr(resourceFullName, "author_id", template.AuthorId),
						/*resource.TestCheckResourceAttr(resourceFullName, "created_at", template.CreatedAt),
						resource.TestCheckResourceAttr(resourceFullName, "href", template.Href),
						resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
						resource.TestCheckResourceAttr(resourceFullName, "description", template.Description),
						resource.TestCheckResourceAttr(resourceFullName, "organization_id", template.OrganizationId),
						resource.TestCheckResourceAttr(resourceFullName, "path", template.Path),
						resource.TestCheckResourceAttr(resourceFullName, "revision", template.Revision),
						resource.TestCheckResourceAttr(resourceFullName, "project_id", template.ProjectId),
						resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),*/
						//resource.TestCheckResourceAttr(resourceFullName, "type", template.Type),
						resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, mockFunc)
	}

	runScenario(templateByName, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Templates().AnyTimes().Return([]client.Template{template}, nil)
	})

	runScenario(templateById, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Template(template.Id).AnyTimes().Return(template, nil)
	})

}
