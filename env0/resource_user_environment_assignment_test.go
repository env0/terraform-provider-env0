package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitUserEnvironmntAssignmentResource(t *testing.T) {
	userId := "uid"
	environmentId := "eid"
	id := "id"
	role := "rid1"
	updatedRole := "rid2"

	createPayload := client.AssignUserRoleToEnvironmentPayload{
		UserId:        userId,
		Role:          role,
		EnvironmentId: environmentId,
	}

	updatePayload := client.AssignUserRoleToEnvironmentPayload{
		UserId:        userId,
		Role:          updatedRole,
		EnvironmentId: environmentId,
	}

	createResponse := client.UserRoleEnvironmentAssignment{
		Id:     id,
		UserId: userId,
		Role:   role,
	}

	updateResponse := client.UserRoleEnvironmentAssignment{
		Id:     id,
		UserId: userId,
		Role:   updatedRole,
	}

	otherResponse := client.UserRoleEnvironmentAssignment{
		Id:     "id2",
		UserId: "userId2",
		Role:   "dasdasd",
	}

	resourceType := "env0_user_environment_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        updatedRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", updatedRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserRoleEnvironmentAssignments(environmentId).Times(2).Return([]client.UserRoleEnvironmentAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignUserRoleToEnvironment(&updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().UserRoleEnvironmentAssignments(environmentId).Times(1).Return([]client.UserRoleEnvironmentAssignment{otherResponse, updateResponse}, nil),
				mock.EXPECT().RemoveUserRoleFromEnvironment(environmentId, userId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserRoleEnvironmentAssignments(environmentId).Times(1).Return([]client.UserRoleEnvironmentAssignment{otherResponse}, nil),
			)
		})
	})

	t.Run("Create Assignment - failed to create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserRoleToEnvironment(&createPayload).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to list assignments", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					ExpectError: regexp.MustCompile("could not get assignments: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserRoleEnvironmentAssignments(environmentId).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveUserRoleFromEnvironment(environmentId, userId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment and update role - failed to update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"environment_id": environmentId,
						"role_id":        updatedRole,
					}),
					ExpectError: regexp.MustCompile("could not update assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserRoleEnvironmentAssignments(environmentId).Times(2).Return([]client.UserRoleEnvironmentAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignUserRoleToEnvironment(&updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveUserRoleFromEnvironment(environmentId, userId).Times(1).Return(nil),
			)
		})
	})
}
