package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitTeamProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_team_project_assignment"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName

	accessor := resourceAccessor(resourceType, resourceName)

	customRole := "id1"
	updatedCustomRole := "id2"

	assignment := client.TeamProjectAssignment{
		Id:          "assignmentId",
		TeamId:      "teamId0",
		ProjectId:   "projectId0",
		ProjectRole: string(client.Admin),
	}

	updateAssignment := client.TeamProjectAssignment{
		Id:          "assignmentIdupdate",
		TeamId:      "teamIdUupdate",
		ProjectId:   "projectId0",
		ProjectRole: string(client.Admin),
	}

	assignmentCustom := client.TeamProjectAssignment{
		Id:          "assignmentId",
		TeamId:      "teamId0",
		ProjectId:   "projectId0",
		ProjectRole: customRole,
	}

	updateAssignmentCustom := client.TeamProjectAssignment{
		Id:          "assignmentIdupdate",
		TeamId:      "teamIdUupdate",
		ProjectId:   "projectId0",
		ProjectRole: updatedCustomRole,
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
						resource.TestCheckResourceAttr(accessor, "role", assignment.ProjectRole),
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
						resource.TestCheckResourceAttr(accessor, "role", assignment.ProjectRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    updateAssignment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "role", string(assignment.ProjectRole)),
					),
				},
			},
		}

		runUnitTest(t, driftTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: updateAssignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(updateAssignment, nil)
			gomock.InOrder(
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{assignment}, nil),
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{updateAssignment}, nil),
				mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{updateAssignment}, nil),
			)
			mock.EXPECT().TeamProjectAssignmentDelete(updateAssignment.Id).Times(1).Return(nil)
		})
	})

	t.Run("create and update custom assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        assignmentCustom.TeamId,
						"project_id":     assignmentCustom.ProjectId,
						"custom_role_id": customRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignmentCustom.ProjectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", customRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        updateAssignmentCustom.TeamId,
						"project_id":     updateAssignmentCustom.ProjectId,
						"custom_role_id": updatedCustomRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", updateAssignmentCustom.ProjectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", updatedCustomRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    updateAssignmentCustom.TeamId,
						"project_id": updateAssignmentCustom.ProjectId,
						"role":       updateAssignment.ProjectRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", updateAssignmentCustom.ProjectId),
						resource.TestCheckResourceAttr(accessor, "role", updateAssignment.ProjectRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignmentCustom.TeamId, ProjectId: assignmentCustom.ProjectId, ProjectRole: assignmentCustom.ProjectRole}).Times(1).Return(assignmentCustom, nil),
				mock.EXPECT().TeamProjectAssignments(assignmentCustom.ProjectId).Times(2).Return([]client.TeamProjectAssignment{assignmentCustom}, nil),
				mock.EXPECT().TeamProjectAssignmentDelete(assignmentCustom.Id).Times(1).Return(nil),
				mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: updateAssignmentCustom.TeamId, ProjectId: updateAssignmentCustom.ProjectId, ProjectRole: updateAssignmentCustom.ProjectRole}).Times(1).Return(updateAssignmentCustom, nil),
				mock.EXPECT().TeamProjectAssignments(assignmentCustom.ProjectId).Times(2).Return([]client.TeamProjectAssignment{updateAssignmentCustom}, nil),
				mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: updateAssignmentCustom.TeamId, ProjectId: updateAssignmentCustom.ProjectId, ProjectRole: updateAssignment.ProjectRole}).Times(1).Return(updateAssignment, nil),
				mock.EXPECT().TeamProjectAssignments(assignmentCustom.ProjectId).Times(1).Return([]client.TeamProjectAssignment{updateAssignment}, nil),
				mock.EXPECT().TeamProjectAssignmentDelete(updateAssignment.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import - built-in role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignment.TeamId + "_" + assignment.ProjectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
			mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(3).Return([]client.TeamProjectAssignment{assignment}, nil)
			mock.EXPECT().TeamProjectAssignmentDelete(assignment.Id).Times(1).Return(nil)
		})
	})

	t.Run("import - custom role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        assignmentCustom.TeamId,
						"project_id":     assignmentCustom.ProjectId,
						"custom_role_id": assignmentCustom.ProjectRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignmentCustom.TeamId + "_" + assignmentCustom.ProjectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignmentCustom.TeamId, ProjectId: assignmentCustom.ProjectId, ProjectRole: assignmentCustom.ProjectRole}).Times(1).Return(assignmentCustom, nil)
			mock.EXPECT().TeamProjectAssignments(assignmentCustom.ProjectId).Times(3).Return([]client.TeamProjectAssignment{assignmentCustom}, nil)
			mock.EXPECT().TeamProjectAssignmentDelete(assignmentCustom.Id).Times(1).Return(nil)
		})
	})

	t.Run("Import role - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": assignment.ProjectId,
						"role":       assignment.ProjectRole,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignment.TeamId + "_" + assignment.ProjectId,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamProjectAssignmentCreateOrUpdate(client.TeamProjectAssignmentPayload{TeamId: assignment.TeamId, ProjectId: assignment.ProjectId, ProjectRole: assignment.ProjectRole}).Times(1).Return(assignment, nil)
			mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{assignment}, nil)
			mock.EXPECT().TeamProjectAssignments(assignment.ProjectId).Times(1).Return([]client.TeamProjectAssignment{}, nil)
			mock.EXPECT().TeamProjectAssignmentDelete(assignment.Id).Times(1)
		})
	})
}
