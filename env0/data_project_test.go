package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProjectDataSource(t *testing.T) {
	project := client.Project{
		Id:          "id0",
		Name:        "my-project-1",
		CreatedBy:   "env0",
		Role:        "role0",
		Description: "A project's description",
	}

	archivedProject := client.Project{
		Id:          "otherId",
		Name:        project.Name,
		CreatedBy:   project.CreatedBy,
		Role:        project.Role,
		Description: project.Description,
		IsArchived:  true,
	}

	projectDataByName := map[string]interface{}{"name": project.Name}
	projectDataById := map[string]interface{}{"id": project.Id}

	resourceType := "env0_project"
	resourceName := "test_project"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "created_by", project.CreatedBy),
						resource.TestCheckResourceAttr(accessor, "role", project.Role),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]interface{}, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockGetProjectCall := func(returnValue client.Project) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Project(project.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockGetProjectCallFailed := func(statusCode int) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Project("id0").AnyTimes().Return(client.Project{}, http.NewMockFailedResponseError(statusCode))
		}
	}

	mockListProjectsCall := func(returnValue []client.Project) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Projects().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(projectDataById),
			mockGetProjectCall(project),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(projectDataByName),
			mockListProjectsCall([]client.Project{project, archivedProject}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one project exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(projectDataByName, "found multiple Projects for name"),
			mockListProjectsCall([]client.Project{project, project}),
		)
	})

	t.Run("Throw error when by name and no projects found at all", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(projectDataByName, "could not find a project with name"),
			mockListProjectsCall([]client.Project{}),
		)
	})

	t.Run("Throw error when by name and no projects found with that name", func(t *testing.T) {
		projectWithOtherName := map[string]interface{}{"name": "other-name"}
		runUnitTest(t,
			getErrorTestCase(projectWithOtherName, "could not find a project with name"),
			mockListProjectsCall([]client.Project{project, project}),
		)
	})

	t.Run("Throw error when by id not found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(projectDataById, "could not find a project with id: id0"),
			mockGetProjectCallFailed(404),
		)
	})
}
