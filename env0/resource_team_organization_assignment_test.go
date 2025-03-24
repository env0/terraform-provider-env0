package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitTeamOrganizationAssignmentResource(t *testing.T) {
	teamId := "tid"
	organizationId := "oid"

	id := "id"
	role := "rid1"
	updatedRole := "rid2"

	createPayload := client.TeamRoleAssignmentCreateOrUpdatePayload{
		TeamId:         teamId,
		Role:           role,
		OrganizationId: organizationId,
	}

	updatePayload := client.TeamRoleAssignmentCreateOrUpdatePayload{
		TeamId:         teamId,
		Role:           updatedRole,
		OrganizationId: organizationId,
	}

	createResponse := client.TeamRoleAssignmentPayload{
		Id:     id,
		TeamId: teamId,
		Role:   role,
	}

	updateResponse := client.TeamRoleAssignmentPayload{
		Id:     id,
		TeamId: teamId,
		Role:   updatedRole,
	}

	otherResponse := client.TeamRoleAssignmentPayload{
		Id:     "id2",
		TeamId: "teamId2",
		Role:   "dasdasd",
	}

	resourceType := "env0_team_organization_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": updatedRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "role_id", updatedRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{OrganizationId: organizationId}).Times(2).Return([]client.TeamRoleAssignmentPayload{otherResponse, createResponse}, nil),
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{OrganizationId: organizationId}).Times(1).Return([]client.TeamRoleAssignmentPayload{otherResponse, updateResponse}, nil),
				mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: teamId, OrganizationId: organizationId}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{OrganizationId: organizationId}).Times(1).Return([]client.TeamRoleAssignmentPayload{otherResponse}, nil),
			)
		})
	})

	t.Run("Create Assignment - failed to create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&createPayload).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to list assignments", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					ExpectError: regexp.MustCompile("could not get assignments: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{OrganizationId: organizationId}).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: teamId, OrganizationId: organizationId}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment and update role - failed to update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": updatedRole,
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{OrganizationId: organizationId}).Times(2).Return([]client.TeamRoleAssignmentPayload{otherResponse, createResponse}, nil),
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: teamId, OrganizationId: organizationId}).Times(1).Return(nil),
			)
		})
	})

	t.Run("unsupported built-in role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"team_id": teamId,
						"role_id": "Planner",
					}),
					ExpectError: regexp.MustCompile(`the following built-in role 'Planner' is not supported for this resource, must be one of \[User,Admin\]`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
