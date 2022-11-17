package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func getOrgUser(userId, role string) client.OrganizationUser {
	return client.OrganizationUser{
		Role: role,
		User: client.User{
			UserId: userId,
		},
	}
}

func TestUnitUserOrganizationAssignmentResource(t *testing.T) {
	userId := "uid"
	updatedUserId := "uid2"
	role := "Admin"
	updatedRole := "User"
	customRole := "id1"
	updatedCustomRole := "id2"

	resourceType := "env0_user_organization_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "role", role),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    updatedRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "role", updatedRole),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, role).Times(1).Return(nil),
				mock.EXPECT().Users().Times(2).Return([]client.OrganizationUser{getOrgUser("dummy", "dummy"), getOrgUser(userId, role)}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(userId, updatedRole).Times(1).Return(nil),
				mock.EXPECT().Users().Times(1).Return([]client.OrganizationUser{getOrgUser("dummy", "dummy"), getOrgUser(userId, updatedRole)}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(userId, "User").Times(1).Return(nil),
			)
		})
	})

	t.Run("create assignment and update custom role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id":        userId,
						"custom_role_id": customRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", customRole),
						resource.TestCheckNoResourceAttr(accessor, "role"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id":        userId,
						"custom_role_id": updatedCustomRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", updatedCustomRole),
						resource.TestCheckNoResourceAttr(accessor, "role"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, customRole).Times(1).Return(nil),
				mock.EXPECT().Users().Times(2).Return([]client.OrganizationUser{getOrgUser("dummy", "dummy"), getOrgUser(userId, customRole)}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(userId, updatedCustomRole).Times(1).Return(nil),
				mock.EXPECT().Users().Times(1).Return([]client.OrganizationUser{getOrgUser("dummy", "dummy"), getOrgUser(userId, updatedCustomRole)}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(userId, "User").Times(1).Return(nil),
			)
		})
	})

	t.Run("create assignment and update user id (force new)", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "role", role),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": updatedUserId,
						"role":    role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedUserId),
						resource.TestCheckResourceAttr(accessor, "user_id", updatedUserId),
						resource.TestCheckResourceAttr(accessor, "role", role),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, role).Times(1).Return(nil),
				mock.EXPECT().Users().Times(2).Return([]client.OrganizationUser{getOrgUser(updatedUserId, role), getOrgUser(userId, role)}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(userId, "User").Times(1).Return(nil),
				mock.EXPECT().OrganizationUserUpdateRole(updatedUserId, role).Times(1).Return(nil),
				mock.EXPECT().Users().Times(1).Return([]client.OrganizationUser{getOrgUser(updatedUserId, role), getOrgUser(userId, "User")}, nil),
				mock.EXPECT().OrganizationUserUpdateRole(updatedUserId, "User").Times(1).Return(nil),
			)
		})
	})

	t.Run("create assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "role", role),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "role", role),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, role).Times(1).Return(nil),
				mock.EXPECT().Users().Times(1).Return([]client.OrganizationUser{getOrgUser("dummy", "dummy")}, nil),
			)
		})
	})

	t.Run("create assignment - failed to assign", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					ExpectError: regexp.MustCompile("failed to update user role organization: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, role).Times(1).Return(errors.New("error")),
			)
		})
	})

	t.Run("create Assignment - failed to list users", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"role":    role,
					}),
					ExpectError: regexp.MustCompile("could not get list of users: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationUserUpdateRole(userId, role).Times(1).Return(nil),
				mock.EXPECT().Users().Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().OrganizationUserUpdateRole(userId, "User").Times(1).Return(nil),
			)
		})
	})

}
