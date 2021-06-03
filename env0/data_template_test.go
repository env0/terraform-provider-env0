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
	retryOnDeploy := client.TemplateRetryOn{
		Times:      2,
		ErrorRegex: "error retry on deploy",
	}
	retryOnDestroy := client.TemplateRetryOn{
		Times:      3,
		ErrorRegex: "error retry on destroy",
	}
	templateRetry := client.TemplateRetry{
		OnDeploy:  &retryOnDeploy,
		OnDestroy: &retryOnDestroy,
	}
	/*
	   d.SetId(template.Id)
	   	d.Set("name", template.Name)
	   	d.Set("repository", template.Repository)
	   	d.Set("path", template.Path)
	   	d.Set("revision", template.Revision)
	   	d.Set("type", template.Type)
	   	d.Set("project_ids", template.ProjectIds)
	   	d.Set("terraform_version", template.TerraformVersion)
	   	d.Set("ssh_keys", template.SshKeys)
	   	if template.Retry.OnDeploy != nil {
	   		d.Set("retries_on_deploy", template.Retry.OnDeploy.Times)
	   		d.Set("retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex)
	   	} else {
	   		d.Set("retries_on_deploy", 0)
	   		d.Set("retry_on_deploy_only_when_matches_regex", "")
	   	}
	   	if template.Retry.OnDestroy != nil {
	   		d.Set("retries_on_destroy", template.Retry.OnDestroy.Times)
	   		d.Set("retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex)
	   	} else {
	   		d.Set("retries_on_destroy", 0)
	   		d.Set("retry_on_destroy_only_when_matches_regex", "")
	*/
	template := client.Template{
		Id:               "id0",
		Name:             "name0",
		Repository:       "repository",
		Path:             "path",
		Revision:         "revision",
		Type:             "terraform",
		TerraformVersion: "0.15.1",
		//sshkeys
		Retry: templateRetry,
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
						resource.TestCheckResourceAttr(resourceFullName, "id", template.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", template.Name),
						resource.TestCheckResourceAttr(resourceFullName, "repository", template.Repository),
						resource.TestCheckResourceAttr(resourceFullName, "path", template.Path),
						resource.TestCheckResourceAttr(resourceFullName, "revision", template.Revision),
						resource.TestCheckResourceAttr(resourceFullName, "type", template.Type),
						resource.TestCheckResourceAttr(resourceFullName, "retries_on_deploy", strconv.Itoa(template.Retry.OnDeploy.Times)),
						resource.TestCheckResourceAttr(resourceFullName, "retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex),
						resource.TestCheckResourceAttr(resourceFullName, "retries_on_destroy", strconv.Itoa(template.Retry.OnDestroy.Times)),
						resource.TestCheckResourceAttr(resourceFullName, "retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex),
						/*

							d.Set("retries_on_destroy", template.Retry.OnDestroy.Times)
							d.Set("retries_on_deploy", template.Retry.OnDeploy.Times)
											resource.TestCheckResourceAttr(resourceFullName, "created_at", template.CreatedAt),
											resource.TestCheckResourceAttr(resourceFullName, "href", template.Href),
											resource.TestCheckResourceAttr(resourceFullName, "description", template.Description),
											resource.TestCheckResourceAttr(resourceFullName, "organization_id", template.OrganizationId),
											resource.TestCheckResourceAttr(resourceFullName, "project_id", template.ProjectId),
											resource.TestCheckResourceAttr(resourceFullName, "terraform_version", template.TerraformVersion),

												template := client.Template{
							Id:             "id0",
							Name:           "name0",
							Repository:     "repository",
							Path:           "path",
							Revision:       "revision",
							Type:                 "terraform",
							TerraformVersion:     "0.15.1",
							//sshkeys
							//Retry:                templateRetry,

						}*/
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
