package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitUserProjectAssignmentResource(t *testing.T) {
	userId := "uid"
	projectId := "pid"
	id := "id"
	role := client.Deployer
	updatedRole := client.Viewer

	createPayload := client.AssignUserToProjectPayload{
		UserId: userId,
		Role:   role,
	}

	updatePayload := client.UpdateUserProjectAssignmentPayload{
		Role: updatedRole,
	}

	createResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   role,
	}

	updateResponse := client.UserProjectAssignment{
		Id:     id,
		UserId: userId,
		Role:   updatedRole,
	}

	otherResponse := client.UserProjectAssignment{
		Id:     "id2",
		UserId: "userId2",
		Role:   role,
	}

	resourceType := "env0_user_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id":    userId,
						"project_id": projectId,
						"role":       string(updatedRole),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(updatedRole)),
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

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
