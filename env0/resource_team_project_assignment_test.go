package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitTeamProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_team_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	assignment := client.TeamProjectAssignment{
		Id:          "assignmentId",
		TeamId:      "teamId0",
		ProjectId:   "projectId0",
		ProjectRole: client.Admin,
	}

	updateAsigment := client.TeamProjectAssignment{
		Id:          "assignmentIdupdate",
		TeamId:      "teamIdUupdate",
		ProjectId:   "projectId0",
		ProjectRole: client.Admin,
	}
	t.Run("create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "role", string(assignment.ProjectRole)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
			mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{assignment}, nil)
			mock.EXPECT().TeamProjectAssignmentDelete(assignment.Id).Times(1).Return(nil)
		})
	})

	t.Run("detect drift", func(t *testing.T) {
		driftTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "role", string(assignment.ProjectRole)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    updateAsigment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAsigment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "role", string(assignment.ProjectRole)),
					),
				},
			},
		}

		runUnitTest(t, driftTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: updateAsigment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(updateAsigment, nil)
			gomock.InOrder(
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{assignment}, nil),
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{updateAsigment}, nil),
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{updateAsigment}, nil),
			)
			mock.EXPECT().TeamProjectAssignmentDelete(updateAsigment.Id).Times(1).Return(nil)
		})
	})

}
