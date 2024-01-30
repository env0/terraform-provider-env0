package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProjectsDataSource(t *testing.T) {
	project1 := client.Project{
		Id:   "id0",
		Name: "my-project-1",
	}

	project2 := client.Project{
		Id:   "id1",
		Name: "my-project-2",
	}

	projects := []client.Project{project1, project2}

	resourceType := "env0_projects"
	resourceName := "test_projects"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "projects.0.id", project1.Id),
						resource.TestCheckResourceAttr(accessor, "projects.0.name", project1.Name),
						resource.TestCheckResourceAttr(accessor, "projects.1.id", project2.Id),
						resource.TestCheckResourceAttr(accessor, "projects.1.name", project2.Name),
					),
				},
			},
		}
	}

	getErrorTestCase := func(expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockGetProjectsCall := func(returnValue []client.Project) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Projects().AnyTimes().Return(returnValue, nil)
		}
	}

	mockGetProjectsCallFailed := func(err string) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Projects().AnyTimes().Return([]client.Project{}, errors.New(err))
		}
	}

	t.Run("get all projects", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(),
			mockGetProjectsCall(projects),
		)
	})

	t.Run("Error when API call fails", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase("failed to get list of projects: error"),
			mockGetProjectsCallFailed("error"),
		)
	})
}
