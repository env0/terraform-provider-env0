package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCustomRoleDataSource(t *testing.T) {
	role := client.Role{
		Id:   "id",
		Name: "name",
	}

	otherRole := client.Role{
		Id:   "id_other",
		Name: "name_other",
	}

	resourceType := "env0_custom_role"
	resourceName := "test"

	accessor := dataSourceAccessor(resourceType, resourceName)

	fieldsByName := map[string]any{"name": role.Name}
	fieldsById := map[string]any{"id": role.Id}

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", role.Id),
						resource.TestCheckResourceAttr(accessor, "name", role.Name),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockListCustomRolesCall := func(returnValue []client.Role) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Roles().AnyTimes().Return(returnValue, nil)
		}
	}

	mockCustomRoleCall := func(returnValue *client.Role) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Role(role.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(fieldsById),
			mockCustomRoleCall(&role),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(fieldsByName),
			mockListCustomRolesCall([]client.Role{role, otherRole}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one role exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsByName, "found multiple custom roles"),
			mockListCustomRolesCall([]client.Role{role, otherRole, role}),
		)
	})

	t.Run("Throw error when by name and no role found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsByName, "not found"),
			mockListCustomRolesCall([]client.Role{otherRole}),
		)
	})

	t.Run("Throw error when by id and no role found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsById, fmt.Sprintf("id %s not found", role.Id)),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Role(role.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404))
			},
		)
	})
}
