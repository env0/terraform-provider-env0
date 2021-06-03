package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/go2hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"testing"
)

func TestUnitOrganizationData(t *testing.T) {
	resourceType := "env0_organization"
	resourceName := "test"
	accessor := go2hcl.DataSourceAccessor(resourceType, resourceName)
	organization := client.Organization{
		Id:           "id0",
		Name:         "name0",
		CreatedBy:    "env0",
		Role:         "role0",
		IsSelfHosted: false,
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: go2hcl.DataSourceConfigCreate(resourceType, resourceName, make(map[string]interface{})),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", organization.Id),
					resource.TestCheckResourceAttr(accessor, "name", organization.Name),
					resource.TestCheckResourceAttr(accessor, "created_by", organization.CreatedBy),
					resource.TestCheckResourceAttr(accessor, "role", organization.Role),
					resource.TestCheckResourceAttr(accessor, "is_self_hosted", strconv.FormatBool(organization.IsSelfHosted)),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Organization().AnyTimes().Return(organization, nil)
	})
}
