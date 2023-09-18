package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestProjectDataSource(t *testing.T) {
	project := client.Project{
		Id:              "id0",
		Name:            "my-project-1",
		CreatedBy:       "env0",
		Role:            "role0",
		Description:     "A project's description",
		ParentProjectId: "parent_id",
	}

	archivedProject := client.Project{
		Id:          "otherId",
		Name:        project.Name,
		CreatedBy:   project.CreatedBy,
		Role:        project.Role,
		Description: project.Description,
		IsArchived:  true,
	}

	parentProject := client.Project{
		Id:          "id_parent1",
		Name:        "parent_project_name",
		CreatedBy:   "env0",
		Role:        "role0",
		Description: "A project's description",
	}

	otherParentProject := client.Project{
		Id:          "id_parent2",
		Name:        "other_parent_project_name",
		CreatedBy:   "env0",
		Role:        "role0",
		Description: "A project's description",
	}

	projectWithParent := client.Project{
		Id:              "id123",
		Name:            "same_name",
		CreatedBy:       "env0",
		Role:            "role0",
		Description:     "A project's description",
		ParentProjectId: parentProject.Id,
	}

	otherProjectWithParent := client.Project{
		Id:              "id234",
		Name:            "same_name",
		CreatedBy:       "env0",
		Role:            "role0",
		Description:     "A project's description",
		ParentProjectId: otherParentProject.Id,
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
						resource.TestCheckResourceAttr(accessor, "parent_project_id", project.ParentProjectId),
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

	t.Run("By Name with Parent Name", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": projectWithParent.Name, "parent_project_name": parentProject.Name}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", projectWithParent.Id),
							resource.TestCheckResourceAttr(accessor, "name", projectWithParent.Name),
							resource.TestCheckResourceAttr(accessor, "created_by", projectWithParent.CreatedBy),
							resource.TestCheckResourceAttr(accessor, "role", projectWithParent.Role),
							resource.TestCheckResourceAttr(accessor, "description", projectWithParent.Description),
							resource.TestCheckResourceAttr(accessor, "parent_project_id", projectWithParent.ParentProjectId),
						),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Projects().AnyTimes().Return([]client.Project{projectWithParent, otherProjectWithParent}, nil)
				mock.EXPECT().Project(projectWithParent.ParentProjectId).AnyTimes().Return(parentProject, nil)
				mock.EXPECT().Project(otherProjectWithParent.ParentProjectId).AnyTimes().Return(otherParentProject, nil)
			},
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
			getErrorTestCase(projectDataByName, "found multiple projects for name"),
			mockListProjectsCall([]client.Project{project, project}),
		)
	})

	t.Run("Throw error when by name and more than one project and parent project name exist", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{"name": projectWithParent.Name, "parent_project_name": parentProject.Name}, "found multiple projects for name"),
			func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Projects().Times(1).Return([]client.Project{projectWithParent, otherProjectWithParent, projectWithParent}, nil),
					mock.EXPECT().Project(projectWithParent.ParentProjectId).Times(1).Return(parentProject, nil),
					mock.EXPECT().Project(otherProjectWithParent.ParentProjectId).Times(1).Return(otherParentProject, nil),
					mock.EXPECT().Project(projectWithParent.ParentProjectId).Times(1).Return(parentProject, nil),
				)
			},
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

	t.Run("Throw error when by name and no projects found with that name and parent name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{"name": projectWithParent.Name, "parent_project_name": "noooo"}, "could not find a project with name"),
			func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Projects().Times(1).Return([]client.Project{projectWithParent}, nil),
					mock.EXPECT().Project(projectWithParent.ParentProjectId).Times(1).Return(parentProject, nil),
				)
			},
		)
	})

	t.Run("Throw error when by id not found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(projectDataById, "could not find a project with id: id0"),
			mockGetProjectCallFailed(404),
		)
	})
}
