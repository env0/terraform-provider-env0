package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Project Resource", func() {
	resourceName := "env0_project.test"
	project := client.Project{
		Id:          "id0",
		Name:        "name0",
		Description: "description0",
	}

	BeforeEach(func() {
		apiClientMock.EXPECT().ProjectCreate(project.Name, project.Description).Times(1).Return(project, nil)
		apiClientMock.EXPECT().Project(project.Id).Times(1)
		//apiClientMock.EXPECT().ProjectDelete(project.Id).Times(1)
	})

	It("Should validate project creation", func() {
		resource.UnitTest(utils.RecoveringGinkgoT(), resource.TestCase{
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
		})
	})
})

func testEnv0ProjectResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
	resource "env0_project" "test" {
		name = "%s"
		description = "%s"
	}
	`, name, description)
}
