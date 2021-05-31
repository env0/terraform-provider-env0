package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
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

	testCase := resource.TestCase{
		ProviderFactories: testUnitProviders,
		Steps: []resource.TestStep{
			{
				Config: testEnv0ProjectResourceConfig(project.Name, project.Description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", project.Id),
					resource.TestCheckResourceAttr(resourceName, "name", project.Name),
					resource.TestCheckResourceAttr(resourceName, "description", project.Description),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().ProjectCreate(project.Name, project.Description).Times(1).Return(project, nil)
		mock.EXPECT().Project(project.Id).Times(1).Return(project, nil)
		//mock.EXPECT().ProjectDelete(project.Id).Times(1).Return(nil)
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
