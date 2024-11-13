package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
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
							resource.TestCheckResourceAttr(accessor, "parent_project_name", parentProject.Name),
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

	createProject := func(name string, ancestors []client.Project) *client.Project {
		p := client.Project{
			Id:   uuid.NewString(),
			Name: name,
		}

		for _, ancestor := range ancestors {
			p.Hierarchy += ancestor.Id + "|"
		}

		p.Hierarchy += p.Id

		return &p
	}

	t.Run("By name with parent path", func(t *testing.T) {
		p1 := createProject("p1", nil)
		p2 := createProject("p2", []client.Project{*p1})
		p3 := createProject("p3", []client.Project{*p1, *p2})
		p4 := createProject("p4", []client.Project{*p1, *p2, *p3})

		p3other := createProject("p3", []client.Project{*p1})
		p4other := createProject("p4", []client.Project{*p1})

		pother1 := createProject("pother1", nil)
		pother2 := createProject("p2", []client.Project{*pother1})
		pother3 := createProject("p3", []client.Project{*pother1, *pother2})

		t.Run("exact match", func(t *testing.T) {
			runUnitTest(t,
				resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "p3", "parent_project_path": "p1|p2"}),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr(accessor, "id", p3.Id),
								resource.TestCheckResourceAttr(accessor, "name", p3.Name),
								resource.TestCheckResourceAttr(accessor, "hierarchy", p1.Id+"|"+p2.Id+"|"+p3.Id),
							),
						},
					},
				},
				func(mock *client.MockApiClientInterface) {
					mock.EXPECT().Projects().AnyTimes().Return([]client.Project{*p3, *p3other, *pother3}, nil)
					mock.EXPECT().Project(p1.Id).AnyTimes().Return(*p1, nil)
					mock.EXPECT().Project(p2.Id).AnyTimes().Return(*p2, nil)
					mock.EXPECT().Project(p3.Id).AnyTimes().Return(*p3, nil)
					mock.EXPECT().Project(p3other.Id).AnyTimes().Return(*p3other, nil)
					mock.EXPECT().Project(pother1.Id).AnyTimes().Return(*pother1, nil)
					mock.EXPECT().Project(pother2.Id).AnyTimes().Return(*pother2, nil)
					mock.EXPECT().Project(pother3.Id).AnyTimes().Return(*pother3, nil)
				},
			)
		})

		t.Run("prefix match", func(t *testing.T) {
			runUnitTest(t,
				resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": "p4", "parent_project_path": "p1|p2"}),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr(accessor, "id", p4.Id),
								resource.TestCheckResourceAttr(accessor, "name", p4.Name),
								resource.TestCheckResourceAttr(accessor, "hierarchy", p1.Id+"|"+p2.Id+"|"+p3.Id+"|"+p4.Id),
							),
						},
					},
				},
				func(mock *client.MockApiClientInterface) {
					mock.EXPECT().Projects().AnyTimes().Return([]client.Project{*p4, *p4other}, nil)
					mock.EXPECT().Project(p1.Id).AnyTimes().Return(*p1, nil)
					mock.EXPECT().Project(p2.Id).AnyTimes().Return(*p2, nil)
					mock.EXPECT().Project(p3.Id).AnyTimes().Return(*p3, nil)
					mock.EXPECT().Project(p4.Id).AnyTimes().Return(*p3, nil)
					mock.EXPECT().Project(p4other.Id).AnyTimes().Return(*p4other, nil)
				},
			)
		})
	})

	t.Run("By Name with Parent Id", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"name": projectWithParent.Name, "parent_project_id": parentProject.Id}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", projectWithParent.Id),
							resource.TestCheckResourceAttr(accessor, "name", projectWithParent.Name),
							resource.TestCheckResourceAttr(accessor, "created_by", projectWithParent.CreatedBy),
							resource.TestCheckResourceAttr(accessor, "role", projectWithParent.Role),
							resource.TestCheckResourceAttr(accessor, "description", projectWithParent.Description),
							resource.TestCheckResourceAttr(accessor, "parent_project_id", projectWithParent.ParentProjectId),
							resource.TestCheckNoResourceAttr(accessor, "parent_project_name"),
						),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Projects().AnyTimes().Return([]client.Project{projectWithParent, otherProjectWithParent}, nil)
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
