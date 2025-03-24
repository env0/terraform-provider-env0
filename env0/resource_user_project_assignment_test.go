package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitUserProjectAssignmentResource(t *testing.T) {
	userId := "uid"
	projectId := "pid"
	id := "id"
	role := client.DeployerRole
	updatedRole := client.ViewerRole
	customRole := "id1"
	updatedCustomRole := "id2"

	createPayload := client.AssignUserToProjectPayload{
		UserId: userId,
		Role:   string(role),
	}

	updatePayload := client.UpdateUserProjectAssignmentPayload{
		Role: string(updatedRole),
	}

	createResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   string(role),
	}

	updateResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   string(updatedRole),
	}

	createCustomPayload := client.AssignUserToProjectPayload{
		UserId: userId,
		Role:   customRole,
	}

	updateCustomPayload := client.UpdateUserProjectAssignmentPayload{
		Role: updatedCustomRole,
	}

	createCustomResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   customRole,
	}

	updateCustomResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   updatedCustomRole,
	}

	otherResponse := client.UserProjectAssignment{
		Id:     "id2",
		UserId: "userId2",
		Role:   string(role),
	}

	resourceType := "env0_user_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(role)),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(updatedRole),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(updatedRole)),
						resource.TestCheckNoResourceAttr(accessor, "custom_role_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(2).Return([]client.UserProjectAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().UpdateUserProjectAssignment(projectId, userId, &updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(1).Return([]client.UserProjectAssignment{otherResponse, updateResponse}, nil),
				mock.EXPECT().RemoveUserFromProject(projectId, userId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create assignment and update custom role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"project_id":     projectId,
						"custom_role_id": customRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", customRole),
						resource.TestCheckNoResourceAttr(accessor, "role"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":        userId,
						"project_id":     projectId,
						"custom_role_id": updatedCustomRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", updatedCustomRole),
						resource.TestCheckNoResourceAttr(accessor, "role"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(updatedRole),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", ""),
						resource.TestCheckResourceAttr(accessor, "role", string(updatedRole)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createCustomPayload).Times(1).Return(&createCustomResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(2).Return([]client.UserProjectAssignment{otherResponse, createCustomResponse}, nil),
				mock.EXPECT().UpdateUserProjectAssignment(projectId, userId, &updateCustomPayload).Times(1).Return(&updateCustomResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(2).Return([]client.UserProjectAssignment{otherResponse, updateCustomResponse}, nil),
				mock.EXPECT().UpdateUserProjectAssignment(projectId, userId, &updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(1).Return([]client.UserProjectAssignment{otherResponse, updateResponse}, nil),
				mock.EXPECT().RemoveUserFromProject(projectId, userId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(role)),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(role)),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(1).Return([]client.UserProjectAssignment{otherResponse}, nil),
			)
		})
	})

	t.Run("Create Assignment - failed to create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createPayload).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to list assignments", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					ExpectError: regexp.MustCompile("could not get UserProjectAssignments: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveUserFromProject(projectId, userId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment and update role - failed to update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(role),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(role)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(updatedRole),
					}),
					ExpectError: regexp.MustCompile("could not update role for UserProjectAssignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignUserToProject(projectId, &createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().UserProjectAssignments(projectId).Times(2).Return([]client.UserProjectAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().UpdateUserProjectAssignment(projectId, userId, &updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveUserFromProject(projectId, userId).Times(1).Return(nil),
			)
		})
	})
}
