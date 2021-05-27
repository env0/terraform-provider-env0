package env0

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

type testEnv0ProjectConfig struct {
	name        string
	description string
}

func (c testEnv0ProjectConfig) hcl() string {
	return fmt.Sprintf(`
	resource "env0_project" "test" {
		name = "%s"
		description = "%s"
	}
	`, c.name, c.description)
}

func TestUnitEnv0Project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testUnitProviders,
		IsUnitTest:        true,
		Steps:             []resource.TestStep{},
	})
}
