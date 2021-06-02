package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitProjectResource(t *testing.T) {
	resourceName := "env0_project.test"

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
				Config: testEnv0ProjectResourceConfig(project.Name, project.Description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", project.Id),
					resource.TestCheckResourceAttr(resourceName, "name", project.Name),
					resource.TestCheckResourceAttr(resourceName, "description", project.Description),
				),
			},
			{
				Config: testEnv0ProjectResourceConfig(updatedProject.Name, updatedProject.Description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", updatedProject.Id),
					resource.TestCheckResourceAttr(resourceName, "name", updatedProject.Name),
					resource.TestCheckResourceAttr(resourceName, "description", updatedProject.Description),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ProjectCreate(project.Name, project.Description).Times(1).Return(project, nil)
		mock.EXPECT().ProjectUpdate(updatedProject.Id, client.UpdateProjectPayload{
			Name:        updatedProject.Name,
			Description: updatedProject.Description,
		}).Times(1).Return(updatedProject, nil)

		gomock.InOrder(
			mock.EXPECT().Project(gomock.Any()).Times(2).Return(project, nil), // 1 after create, 1 before update
			mock.EXPECT().Project(gomock.Any()).Return(updatedProject, nil),   // 1 after update
		)

		mock.EXPECT().ProjectDelete(project.Id).Times(1)
	})
}

func testEnv0ProjectResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
	resource "env0_project" "test" {
		name = "%s"
		description = "%s"
	}
	`, name, description)
}
