package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitProjectData(t *testing.T) {
	resourceType := "env0_project"
	resourceName := "test_project"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)
	project := client.Project{
		Id:          "id0",
		Name:        "my-project-1",
		CreatedBy:   "env0",
		Role:        "role0",
		Description: "A project's description",
	}

	projectByName := map[string]string{
		"name": project.Name,
	}

	projectById := map[string]string{
		"id": project.Id,
	}

	runScenario := func(input map[string]string, mockFunc func(mockFunc *client.MockApiClientInterface)) {
		testCase := resource.TestCase{
			ProviderFactories: testUnitProviders,
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(resourceFullName, "id", project.Id),
						resource.TestCheckResourceAttr(resourceFullName, "name", project.Name),
						resource.TestCheckResourceAttr(resourceFullName, "created_by", project.CreatedBy),
						resource.TestCheckResourceAttr(resourceFullName, "role", project.Role),
						resource.TestCheckResourceAttr(resourceFullName, "description", project.Description),
					),
				},
			},
		}

		runUnitTest(t, testCase, mockFunc)
	}

	runScenario(projectByName, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Projects().AnyTimes().Return([]client.Project{project}, nil)
	})

	runScenario(projectById, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Project(project.Id).AnyTimes().Return(project, nil)
	})
}
