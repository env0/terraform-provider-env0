package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCustomRolesDataSource(t *testing.T) {
	role1 := client.Role{
		Id:   "id0",
		Name: "name0",
	}

	role2 := client.Role{
		Id:   "id1",
		Name: "name1",
	}

	role3 := client.Role{
		Name:          "name1",
		IsDefaultRole: true,
	}

	resourceType := "env0_custom_roles"
	resourceName := "test"

	accessor := dataSourceAccessor(resourceType, resourceName)

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "names.0", role1.Name),
						resource.TestCheckResourceAttr(accessor, "names.1", role2.Name),
					),
				},
			},
		}
	}

	mockRoles := func(returnValue []client.Role) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Roles().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockRoles([]client.Role{role3, role1, role2}),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Roles().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
