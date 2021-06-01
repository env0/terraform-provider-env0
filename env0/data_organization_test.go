package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

var resourceType = "data_organization"
var resourceName = "test"
var resourceFullName = fmt.Sprintf("%s.%s", resourceType, resourceName)
var organization = client.Organization{
	Id:   "id0",
	Name: "name0",
}

func TestUnitOrganizationDataById(t *testing.T) {
	testCase := resource.TestCase{
		ProviderFactories: testUnitProviders,
		Steps: []resource.TestStep{
			{
				Config: testEnv0OrganizationDataConfig(organization),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", organization.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", organization.Name),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Organization().Times(1).Return(organization, nil)
	})
}

func testEnv0OrganizationDataConfig(organization client.Organization) string {
	return fmt.Sprintf(`
	data "%s" "%s" {
		id = "%s"
	}
	`, resourceType, resourceName, organization.Id)
}
