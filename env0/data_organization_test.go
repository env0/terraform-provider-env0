package env0

import (
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

var resourceType = "env0_organization"
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
				Config: testEnv0OrganizationDataConfig(),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Organization().AnyTimes().Return(organization, nil)
	})
}

func testEnv0OrganizationDataConfig() string {
	return fmt.Sprintf(`	data "%s" "%s" {}`, resourceType, resourceName)
}
