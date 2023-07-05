package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitTeamEnvironmntAssignmentResource(t *testing.T) {
	teamId := "tid"
	environmentId := "eid"
	id := "id"
	role := "rid1"
	updatedRole := "rid2"

	createPayload := client.AssignTeamRoleToEnvironmentPayload{
		TeamId:        teamId,
		Role:          role,
		EnvironmentId: environmentId,
	}

	updatePayload := client.AssignTeamRoleToEnvironmentPayload{
		TeamId:        teamId,
		Role:          updatedRole,
		EnvironmentId: environmentId,
	}

	createResponse := client.TeamRoleEnvironmentAssignment{
		Id:     id,
		TeamId: teamId,
		Role:   role,
	}

	updateResponse := client.TeamRoleEnvironmentAssignment{
		Id:     id,
		TeamId: teamId,
		Role:   updatedRole,
	}

	otherResponse := client.TeamRoleEnvironmentAssignment{
		Id:     "id2",
		TeamId: "teamId2",
		Role:   "dasdasd",
	}

	resourceType := "env0_team_environment_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment and update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        updatedRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", updatedRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignTeamRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleEnvironmentAssignments(environmentId).Times(2).Return([]client.TeamRoleEnvironmentAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignTeamRoleToEnvironment(&updatePayload).Times(1).Return(&updateResponse, nil),
				mock.EXPECT().TeamRoleEnvironmentAssignments(environmentId).Times(1).Return([]client.TeamRoleEnvironmentAssignment{otherResponse, updateResponse}, nil),
				mock.EXPECT().RemoveTeamRoleFromEnvironment(environmentId, teamId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
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
				mock.EXPECT().AssignTeamRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleEnvironmentAssignments(environmentId).Times(1).Return([]client.TeamRoleEnvironmentAssignment{otherResponse}, nil),
			)
		})
	})

	t.Run("Create Assignment - failed to create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					ExpectError: regexp.MustCompile("could not create assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignTeamRoleToEnvironment(&createPayload).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to list assignments", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					ExpectError: regexp.MustCompile("could not get assignments: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignTeamRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleEnvironmentAssignments(environmentId).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveTeamRoleFromEnvironment(environmentId, teamId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Assignment and update role - failed to update role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", id),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "role_id", role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        teamId,
						"environment_id": environmentId,
						"role_id":        updatedRole,
					}),
					ExpectError: regexp.MustCompile("could not update assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignTeamRoleToEnvironment(&createPayload).Times(1).Return(&createResponse, nil),
				mock.EXPECT().TeamRoleEnvironmentAssignments(environmentId).Times(2).Return([]client.TeamRoleEnvironmentAssignment{otherResponse, createResponse}, nil),
				mock.EXPECT().AssignTeamRoleToEnvironment(&updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().RemoveTeamRoleFromEnvironment(environmentId, teamId).Times(1).Return(nil),
			)
		})
	})
}
