package env0

import (
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
		)

		mock.EXPECT().ProjectDelete(project.Id).Times(1)
	})
}

func TestUnitProjectInvalidParams(t *testing.T) {
	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate("env0_project", "test", map[string]interface{}{"name": ""}),
				ExpectError: regexp.MustCompile("Project name cannot be empty"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}
