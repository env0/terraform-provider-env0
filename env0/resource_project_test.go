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
	t.Parallel()

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

	newParentProject := client.Project{
		Id:        updatedSubproject.ParentProjectId,
		Name:      "other parent name",
		Hierarchy: updatedSubproject.ParentProjectId,
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

	projectWithTags := client.Project{
		Id:          "id0",
		Name:        "name0",
		Description: "description0",
		Tags:        []string{"tag1", "tag2"},
	}

	updatedProjectWithTags := client.Project{
		Id:          projectWithTags.Id,
		Name:        "new name",
		Description: "new description",
		Tags:        []string{"tag3"},
	}

	t.Run("Test project with tags", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        projectWithTags.Name,
						"description": projectWithTags.Description,
						"tags":        projectWithTags.Tags,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", projectWithTags.Id),
						resource.TestCheckResourceAttr(accessor, "name", projectWithTags.Name),
						resource.TestCheckResourceAttr(accessor, "description", projectWithTags.Description),
						resource.TestCheckResourceAttr(accessor, "tags.0", "tag1"),
						resource.TestCheckResourceAttr(accessor, "tags.1", "tag2"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        updatedProjectWithTags.Name,
						"description": updatedProjectWithTags.Description,
						"tags":        updatedProjectWithTags.Tags,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedProjectWithTags.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedProjectWithTags.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedProjectWithTags.Description),
						resource.TestCheckResourceAttr(accessor, "tags.0", "tag3"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
				Name:        projectWithTags.Name,
				Description: projectWithTags.Description,
				Tags:        projectWithTags.Tags,
			}).Times(1).Return(projectWithTags, nil)
			mock.EXPECT().ProjectUpdate(updatedProjectWithTags.Id, client.ProjectUpdatePayload{
				Name:        updatedProjectWithTags.Name,
				Description: updatedProjectWithTags.Description,
				Tags:        updatedProjectWithTags.Tags,
			}).Times(1).Return(updatedProjectWithTags, nil)

			gomock.InOrder(
				mock.EXPECT().Project(gomock.Any()).Times(2).Return(projectWithTags, nil),        // 1 after create, 1 before update
				mock.EXPECT().Project(gomock.Any()).Times(1).Return(updatedProjectWithTags, nil), // 1 after update
				mock.EXPECT().ProjectEnvironments(projectWithTags.Id).Times(1).Return([]client.Environment{}, nil),
			)

			mock.EXPECT().ProjectDelete(projectWithTags.Id).Times(1)
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
			// The test harness plans the move step three times; the validator looks the new parent up on each.
			mock.EXPECT().Project(newParentProject.Id).Times(3).Return(newParentProject, nil)

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

func TestUnitProjectMoveValidation(t *testing.T) {
	t.Parallel()

	resourceType := "env0_project"
	resourceName := "test"

	project := client.Project{
		Id:        "id0",
		Name:      "name0",
		Hierarchy: "id0",
	}

	subProject := client.Project{
		Id:              "subProjectId",
		Name:            "sub project name",
		ParentProjectId: project.Id,
		Hierarchy:       project.Id + "|subProjectId",
	}

	createConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
		"name": project.Name,
	})

	moveConfig := func(parentProjectId string) string {
		return resourceConfigCreate(resourceType, resourceName, map[string]any{
			"name":              project.Name,
			"parent_project_id": parentProjectId,
		})
	}

	expectCreateAndDestroy := func(mock *client.MockApiClientInterface, projectReads int) {
		mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{Name: project.Name}).Times(1).Return(project, nil)
		mock.EXPECT().Project(project.Id).Times(projectReads).Return(project, nil)
		mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{}, nil)
		mock.EXPECT().ProjectDelete(project.Id).Times(1)
	}

	t.Run("Parent is the project itself", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createConfig,
				},
				{
					Config:      moveConfig(project.Id),
					ExpectError: regexp.MustCompile("a project cannot be its own parent"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			// No lookup of the new parent: the id is rejected without calling the API.
			expectCreateAndDestroy(mock, 2)
		})
	})

	t.Run("Parent is a sub-project of the project", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createConfig,
				},
				{
					Config:      moveConfig(subProject.Id),
					ExpectError: regexp.MustCompile("cannot move project under '" + subProject.Id + "': it is one of its own sub-projects"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			expectCreateAndDestroy(mock, 2)
			mock.EXPECT().Project(subProject.Id).Times(1).Return(subProject, nil)
		})
	})

	t.Run("Parent removed", func(t *testing.T) {
		rootedSubProject := client.Project{
			Id:        subProject.Id,
			Name:      subProject.Name,
			Hierarchy: subProject.Id,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              subProject.Name,
						"parent_project_id": project.Id,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": subProject.Name,
					}),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectCreate(client.ProjectCreatePayload{
					Name:            subProject.Name,
					ParentProjectId: project.Id,
				}).Times(1).Return(subProject, nil),
				mock.EXPECT().Project(subProject.Id).Times(2).Return(subProject, nil),
				mock.EXPECT().ProjectMove(subProject.Id, "").Times(1).Return(nil),
				mock.EXPECT().ProjectUpdate(subProject.Id, client.ProjectUpdatePayload{
					Name: subProject.Name,
				}).Times(1).Return(rootedSubProject, nil),
				mock.EXPECT().Project(subProject.Id).Times(1).Return(rootedSubProject, nil),
				mock.EXPECT().ProjectEnvironments(subProject.Id).Times(1).Return([]client.Environment{}, nil),
				mock.EXPECT().ProjectDelete(subProject.Id).Times(1),
			)
		})
	})

	t.Run("Parent lookup fails", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createConfig,
				},
				{
					Config:      moveConfig(subProject.Id),
					ExpectError: regexp.MustCompile("could not validate parent project '" + subProject.Id + "': error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			expectCreateAndDestroy(mock, 2)
			mock.EXPECT().Project(subProject.Id).Times(1).Return(client.Project{}, errors.New("error"))
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
	t.Parallel()

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

	t.Run("Success With Inactive Unarchived Environment", func(t *testing.T) {
		inactiveEnvironment := client.Environment{
			Name:   "destroyed-by-schedule",
			Status: "INACTIVE",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        project.Name,
						"description": project.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "force_destroy", "false"),
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
			mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{inactiveEnvironment}, nil)
			mock.EXPECT().ProjectDelete(project.Id).Times(1)
		})
	})

	t.Run("Failure With Failed Destroy Environment", func(t *testing.T) {
		failedEnvironment := client.Environment{
			Name:   "failed-destroy",
			Status: "FAILED",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        project.Name,
						"description": project.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "force_destroy", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": project.Name,
					}),
					Destroy:     true,
					ExpectError: regexp.MustCompile("could not delete project: found an active environment " + failedEnvironment.Name + " - its infrastructure may still exist"),
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
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{failedEnvironment}, nil),
				mock.EXPECT().ProjectEnvironments(project.Id).Times(1).Return([]client.Environment{}, nil),
			)

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
