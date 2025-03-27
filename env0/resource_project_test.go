package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
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

	updatedSubproject := client.Project{
		Id:              "subProjectId",
		Description:     "sub project des2",
		Name:            "sub project nam2",
		ParentProjectId: "other_parent_id",
	}

	t.Run("Test project", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
			mock.EXPECT().ProjectUpdate(updatedProject.Id, client.ProjectUpdatePayload{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              updatedSubproject.Name,
						"description":       updatedSubproject.Description,
						"parent_project_id": updatedSubproject.ParentProjectId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", subProject.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedSubproject.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedSubproject.Description),
						resource.TestCheckResourceAttr(accessor, "parent_project_id", updatedSubproject.ParentProjectId),
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
				mock.EXPECT().Project(subProject.Id).Times(2).Return(subProject, nil),
				mock.EXPECT().ProjectMove(subProject.Id, updatedSubproject.ParentProjectId).Times(1).Return(nil),
				mock.EXPECT().ProjectUpdate(subProject.Id, client.ProjectUpdatePayload{
					Name:        updatedSubproject.Name,
					Description: updatedSubproject.Description,
				}).Times(1).Return(updatedSubproject, nil),
				mock.EXPECT().Project(subProject.Id).Times(1).Return(updatedSubproject, nil),
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
				Config:      resourceConfigCreate("env0_project", "test", map[string]any{"name": ""}),
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

	t.Run("Success With Force Destroy", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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

	t.Run("Failure Without Force Destroy", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return(nil, errors.New("random error")),        // Second time return some random error to force the test to stop waiting.
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{environment}, nil), // Third time fail and expect the error.
				mock.EXPECT().ProjectEnvironments(project.Id).Times(2).Return([]client.Environment{}, nil),            // These calls are for destroying the project at the end of test (return no environments so it won't fail).
			)

			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})
}
