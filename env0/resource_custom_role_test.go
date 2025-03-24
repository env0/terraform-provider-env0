package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitCustomRoleResource(t *testing.T) {
	resourceType := "env0_custom_role"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	role := client.Role{
		Id:   uuid.NewString(),
		Name: "name",
		Permissions: []string{
			"MANAGE_PROJECT_TEMPLATES",
			"CREATE_CUSTOM_ROLES",
		},
		IsDefaultRole:  true,
		OrganizationId: "orgid",
	}

	updatedRole := client.Role{
		Id:             role.Id,
		Name:           "name2",
		Permissions:    []string{"OVERRIDE_MAX_TTL"},
		IsDefaultRole:  false,
		OrganizationId: "orgid",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", role.Id),
						resource.TestCheckResourceAttr(accessor, "name", role.Name),
						resource.TestCheckResourceAttr(accessor, "permissions.#", "2"),
						resource.TestCheckResourceAttr(accessor, "permissions.0", role.Permissions[0]),
						resource.TestCheckResourceAttr(accessor, "permissions.1", role.Permissions[1]),
						resource.TestCheckResourceAttr(accessor, "is_default_role", strconv.FormatBool(role.IsDefaultRole)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedRole.Name,
						"permissions":     updatedRole.Permissions,
						"is_default_role": updatedRole.IsDefaultRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedRole.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedRole.Name),
						resource.TestCheckResourceAttr(accessor, "permissions.#", "1"),
						resource.TestCheckResourceAttr(accessor, "permissions.0", updatedRole.Permissions[0]),
						resource.TestCheckResourceAttr(accessor, "is_default_role", strconv.FormatBool(updatedRole.IsDefaultRole)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().RoleUpdate(role.Id, client.RoleUpdatePayload{
				Name:          updatedRole.Name,
				Permissions:   updatedRole.Permissions,
				IsDefaultRole: updatedRole.IsDefaultRole,
			}).Times(1).Return(&updatedRole, nil)

			gomock.InOrder(
				mock.EXPECT().Role(gomock.Any()).Times(2).Return(&role, nil),        // 1 after create, 1 before update
				mock.EXPECT().Role(gomock.Any()).Times(1).Return(&updatedRole, nil), // 1 after update
			)

			mock.EXPECT().RoleDelete(role.Id).Times(1)
		})
	})

	t.Run("Failure in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
					ExpectError: regexp.MustCompile("could not create a custom role: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Failure in update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedRole.Name,
						"permissions":     updatedRole.Permissions,
						"is_default_role": updatedRole.IsDefaultRole,
					}),
					ExpectError: regexp.MustCompile("could not update custom role: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().RoleUpdate(role.Id, client.RoleUpdatePayload{
				Name:          updatedRole.Name,
				Permissions:   updatedRole.Permissions,
				IsDefaultRole: updatedRole.IsDefaultRole,
			}).Times(1).Return(nil, errors.New("error"))

			mock.EXPECT().Role(gomock.Any()).Times(2).Return(&role, nil) // 1 after create, 1 before update
			mock.EXPECT().RoleDelete(role.Id).Times(1)
		})
	})

	t.Run("Drift", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
			"name":            role.Name,
			"permissions":     role.Permissions,
			"is_default_role": role.IsDefaultRole,
		})

		createRoleParams := client.RoleCreatePayload{
			Name:          role.Name,
			Permissions:   role.Permissions,
			IsDefaultRole: role.IsDefaultRole,
		}

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
				},
				{
					Config: stepConfig,
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().RoleCreate(createRoleParams).Times(1).Return(&role, nil),
				mock.EXPECT().Role(role.Id).Times(1).Return(&role, nil),
				mock.EXPECT().Role(role.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404)),
				mock.EXPECT().RoleCreate(createRoleParams).Times(1).Return(&role, nil),
				mock.EXPECT().Role(role.Id).Times(1).Return(&role, nil),
				mock.EXPECT().RoleDelete(role.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     role.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().Roles().Times(1).Return([]client.Role{role}, nil)
			mock.EXPECT().Role(role.Id).Times(2).Return(&role, nil)
			mock.EXPECT().RoleDelete(role.Id).Times(1)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     role.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().Role(role.Id).Times(3).Return(&role, nil)
			mock.EXPECT().RoleDelete(role.Id).Times(1)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		otherUuid := uuid.New().String()

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     otherUuid,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().Role(role.Id).Times(1).Return(&role, nil)
			mock.EXPECT().Role(otherUuid).Times(1).Return(nil, &client.NotFoundError{})
			mock.EXPECT().RoleDelete(role.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		otherName := "otherName"

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            role.Name,
						"permissions":     role.Permissions,
						"is_default_role": role.IsDefaultRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     otherName,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RoleCreate(client.RoleCreatePayload{
				Name:          role.Name,
				Permissions:   role.Permissions,
				IsDefaultRole: role.IsDefaultRole,
			}).Times(1).Return(&role, nil)
			mock.EXPECT().Role(role.Id).Times(1).Return(&role, nil)
			mock.EXPECT().Roles().Times(1).Return([]client.Role{role}, nil)
			mock.EXPECT().RoleDelete(role.Id).Times(1).Return(nil)
		})
	})
}
