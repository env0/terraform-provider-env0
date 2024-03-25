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

	id := "id"
	role := "rid1"
	updatedRole := "rid2"

	createPayload := client.AssignOrganizationRoleToTeamPayload{
		TeamId: teamId,
		Role:   role,
	}

	updatePayload := client.AssignOrganizationRoleToTeamPayload{
		TeamId: teamId,
		Role:   updatedRole,
	}

	createResponse := client.OrganizationRoleTeamAssignment{
		Id:     id,
		TeamId: teamId,
		Role:   role,
	}

	updateResponse := client.OrganizationRoleTeamAssignment{
		Id:     id,
		TeamId: teamId,
		Role:   updatedRole,
	}

	otherResponse := client.OrganizationRoleTeamAssignment{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
			gomock.InOrder(
				mock.EXPECT().AssignOrganizationRoleToTeam(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().OrganizationRoleTeamAssignments().Times(2).Return([]client.OrganizationRoleTeamAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignOrganizationRoleToTeam(&updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().OrganizationRoleTeamAssignments().Times(1).Return([]client.OrganizationRoleTeamAssignment{otherResponse, updateResponse}, nil),
				mock.EXPECT().RemoveOrganizationRoleFromTeam(teamId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
			gomock.InOrder(
				mock.EXPECT().AssignOrganizationRoleToTeam(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().OrganizationRoleTeamAssignments().Times(1).Return([]client.OrganizationRoleTeamAssignment{otherResponse}, nil),
			)
		})
	})

	t.Run("Create Assignment - failed to create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id": teamId,
						"role_id": role,
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignOrganizationRoleToTeam(&createPayload).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to list assignments", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id": teamId,
						"role_id": role,
					}),
					ExpectError: regexp.MustCompile("could not get assignments: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignOrganizationRoleToTeam(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().OrganizationRoleTeamAssignments().Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveOrganizationRoleFromTeam(teamId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment and update role - failed to update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id": teamId,
						"role_id": updatedRole,
					}),
					ExpectError: regexp.MustCompile("could not update assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignOrganizationRoleToTeam(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().OrganizationRoleTeamAssignments().Times(2).Return([]client.OrganizationRoleTeamAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignOrganizationRoleToTeam(&updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveOrganizationRoleFromTeam(teamId).Times(1).Return(nil),
			)
		})
	})
}
