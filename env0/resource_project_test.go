package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitProjectResource(t *testing.T) {
	resourceType := "env0_project"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	project := client.Project{
		Id:          "id0",
		Name:        "name0",
		Description: "description0",
	}

	updatedProject := client.Project{
		Id:          project.Id,
		Name:        "new name",
		Description: "new description",
	}

	subProject := client.Project{
		Id:              "subProjectId",
		Description:     "sub project des",
		Name:            "sub project nam",
		ParentProjectId: project.Id,
	}

	t.Run("Test project", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        project.Name,
						"description": project.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedProject.Name,
						"description": updatedProject.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedProject.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedProject.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedProject.Description),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
				Name:        project.Name,
				Description: project.Description,
			}).Times(1).Return(project, nil)
			mock.EXPECT().ProjectUpdate(updatedProject.Id, client.ProjectCreatePayload{
				Name:        updatedProject.Name,
				Description: updatedProject.Description,
			}).Times(1).Return(updatedProject, nil)

			gomock.InOrder(
				mock.EXPECT().Project(gomock.Any()).Times(2).Return(project, nil),        // 1 after create, 1 before update
				mock.EXPECT().Project(gomock.Any()).Times(1).Return(updatedProject, nil), // 1 after update
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{}, nil),
			)

			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})

	t.Run("Test sub-project", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              subProject.Name,
						"description":       subProject.Description,
						"parent_project_id": project.Id,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", subProject.Id),
						resource.TestCheckResourceAttr(accessor, "name", subProject.Name),
						resource.TestCheckResourceAttr(accessor, "description", subProject.Description),
						resource.TestCheckResourceAttr(accessor, "parent_project_id", project.Id),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
					Name:            subProject.Name,
					Description:     subProject.Description,
					ParentProjectId: project.Id,
				}).Times(1).Return(subProject, nil),
				mock.EXPECT().Project(subProject.Id).Times(1).Return(subProject, nil),
				mock.EXPECT().ProjectEnvironments(subProject.Id).Times(1).Return([]client.Environment{}, nil),
				mock.EXPECT().ProjectDelete(subProject.Id).Times(1),
			)
		})
	})
}

func TestUnitProjectInvalidParams(t *testing.T) {
	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate("env0_project", "test", map[string]interface{}{"name": ""}),
				ExpectError: regexp.MustCompile("may not be empty"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}

func TestUnitProjectResourceDestroyWithEnvironments(t *testing.T) {
	resourceType := "env0_project"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	project := client.Project{
		Id:          "id0",
		Name:        "name0",
		Description: "description0",
	}

	environment := client.Environment{
		Name: "name1",
	}

	t.Run("Success With Force Destory", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":          project.Name,
						"description":   project.Description,
						"force_destroy": true,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
						resource.TestCheckResourceAttr(accessor, "force_destroy", "true"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
				Name:        project.Name,
				Description: project.Description,
			}).Times(1).Return(project, nil)
			mock.EXPECT().Project(gomock.Any()).Times(1).Return(project, nil)
			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})

	t.Run("Failure Without Force Destory", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        project.Name,
						"description": project.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
						resource.TestCheckResourceAttr(accessor, "force_destroy", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": project.Name,
					}),
					Destroy:     true,
					ExpectError: regexp.MustCompile("could not delete project: found an active environment " + environment.Name),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
				Name:        project.Name,
				Description: project.Description,
			}).Times(1).Return(project, nil)

			gomock.InOrder(
				mock.EXPECT().Project(gomock.Any()).Times(2).Return(project, nil),
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{environment}, nil),
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{}, nil),
			)

			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})

	t.Run("Test wait", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        project.Name,
						"description": project.Description,
						"wait":        "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
						resource.TestCheckResourceAttr(accessor, "force_destroy", "false"),
						resource.TestCheckResourceAttr(accessor, "wait", "true"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": project.Name,
					}),
					Destroy:     true,
					ExpectError: regexp.MustCompile("could not delete project: found an active environment"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
				Name:        project.Name,
				Description: project.Description,
			}).Times(1).Return(project, nil)

			gomock.InOrder(
				mock.EXPECT().Project(gomock.Any()).Times(2).Return(project, nil),
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{environment}, nil), // First time wait - an environment is still active.
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return(nil, errors.New("random error")),        // Second time return some random error will stop waiting.
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{environment}, nil), // Third time will fail.
				mock.EXPECT().ProjectEnvironments(project.Id).Times(2).Return([]client.Environment{}, nil),            // This will allow the project to get destroyed (no environments).
			)

			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})
}
