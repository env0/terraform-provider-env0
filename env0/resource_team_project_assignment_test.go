package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
		mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignment{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
		//mock.EXPECT().TeamProjectAssignmentCreateOrUpdate().Times(1) //.Return([]client.SshKey{sshKey}, nil)
		mock.EXPECT().TeamProjectAssignmentDelete(assignment.Id).Times(1).Return(nil)
	})
}
