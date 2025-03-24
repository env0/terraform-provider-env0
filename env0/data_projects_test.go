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
		Id:              "id0",
		Name:            "my-project-1",
		ParentProjectId: "p1",
		Hierarchy:       "adsas|fdsfsd",
	}

	project2 := client.Project{
		Id:   "id1",
		Name: "my-project-2",
	}

	project3 := client.Project{
		Id:         "id1",
		Name:       "my-project-2",
		IsArchived: true,
	}

	projects := []client.Project{project1, project2, project3}

	resourceType := "env0_projects"
	resourceName := "test_projects"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "projects.0.id", project1.Id),
						resource.TestCheckResourceAttr(accessor, "projects.0.name", project1.Name),
						resource.TestCheckResourceAttr(accessor, "projects.0.parent_project_id", project1.ParentProjectId),
						resource.TestCheckResourceAttr(accessor, "projects.0.hierarchy", project1.Hierarchy),
						resource.TestCheckResourceAttr(accessor, "projects.0.is_archived", "false"),
						resource.TestCheckResourceAttr(accessor, "projects.1.id", project2.Id),
						resource.TestCheckResourceAttr(accessor, "projects.1.name", project2.Name),
						resource.TestCheckResourceAttr(accessor, "projects.1.is_archived", "false"),
						resource.TestCheckResourceAttr(accessor, "projects.#", "2"),
					),
				},
			},
		}
	}

	getValidTestCaseWithArchived := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{
						"include_archived_projects": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "projects.0.id", project1.Id),
						resource.TestCheckResourceAttr(accessor, "projects.0.name", project1.Name),
						resource.TestCheckResourceAttr(accessor, "projects.0.parent_project_id", project1.ParentProjectId),
						resource.TestCheckResourceAttr(accessor, "projects.0.hierarchy", project1.Hierarchy),
						resource.TestCheckResourceAttr(accessor, "projects.0.is_archived", "false"),
						resource.TestCheckResourceAttr(accessor, "projects.1.id", project2.Id),
						resource.TestCheckResourceAttr(accessor, "projects.1.name", project2.Name),
						resource.TestCheckResourceAttr(accessor, "projects.1.is_archived", "false"),
						resource.TestCheckResourceAttr(accessor, "projects.2.id", project2.Id),
						resource.TestCheckResourceAttr(accessor, "projects.2.name", project2.Name),
						resource.TestCheckResourceAttr(accessor, "projects.2.is_archived", "true"),
						resource.TestCheckResourceAttr(accessor, "projects.#", "3"),
					),
				},
			},
		}
	}

	getErrorTestCase := func(expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
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

	t.Run("get all projects including archived", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCaseWithArchived(),
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
