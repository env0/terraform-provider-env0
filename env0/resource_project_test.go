package env0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func testEnv0ProjectResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
	resource "env0_project" "test" {
		name = "%s"
		description = "%s"
	}
	`, name, description)
}

func TestUnitEnv0ProjectCreate(t *testing.T) {
	const name = "name0"
	const description = "description0"

	ctrl, _ := mockApiClient(t)
	defer ctrl.Finish()

	runUnitTest(t, resource.TestCase{
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				Config: testEnv0ProjectResourceConfig("name0", "description0"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
