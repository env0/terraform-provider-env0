package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestUnitOrganizationData(t *testing.T) {
	resourceType := "env0_organization"
	resourceName := "test"
	resourceFullName := DataSourceAccessor(resourceType, resourceName)
	organization := client.Organization{
		Id:           "id0",
		Name:         "name0",
		CreatedBy:    "env0",
		Role:         "role0",
		IsSelfHosted: false,
	}

	testCase := resource.TestCase{
		ProviderFactories: testUnitProviders,
		Steps: []resource.TestStep{
			{
				Config: DataSourceConfigCreate(resourceType, resourceName, make(map[string]string)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", organization.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", organization.Name),
					resource.TestCheckResourceAttr(resourceFullName, "created_by", organization.CreatedBy),
					resource.TestCheckResourceAttr(resourceFullName, "role", organization.Role),
					resource.TestCheckResourceAttr(resourceFullName, "is_self_hosted", strconv.FormatBool(organization.IsSelfHosted)),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Organization().AnyTimes().Return(organization, nil)
	})
}
