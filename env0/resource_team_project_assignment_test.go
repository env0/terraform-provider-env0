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

	projectId := "projectId0"

	assignment := client.TeamRoleAssignmentPayload{
		Id:     "assignmentId",
		TeamId: "teamId0",
		Role:   string(client.AdminRole),
	}

	updateAssignment := client.TeamRoleAssignmentPayload{
		Id:     "assignmentIdupdate",
		TeamId: "teamIdUupdate",
		Role:   string(client.AdminRole),
	}

	assignmentCustom := client.TeamRoleAssignmentPayload{
		Id:     "assignmentId",
		TeamId: "teamId0",
		Role:   customRole,
	}

	updateAssignmentCustom := client.TeamRoleAssignmentPayload{
		Id:     "assignmentIdupdate",
		TeamId: "teamIdUupdate",
		Role:   updatedCustomRole,
	}

	t.Run("create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": projectId,
						"role":       assignment.Role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", assignment.Role),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignment.TeamId, ProjectId: projectId, Role: assignment.Role}).Times(1).Return(&assignment, nil)
			mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{assignment}, nil)
			mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: assignment.TeamId, ProjectId: projectId}).Times(1).Return(nil)
		})
	})

	t.Run("detect drift", func(t *testing.T) {
		driftTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": projectId,
						"role":       assignment.Role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", assignment.Role),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    updateAssignment.TeamId,
						"project_id": projectId,
						"role":       assignment.Role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignment.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", string(assignment.Role)),
					),
				},
			},
		}

		runUnitTest(t, driftTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignment.TeamId, ProjectId: projectId, Role: assignment.Role}).Times(1).Return(&assignment, nil)
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: updateAssignment.TeamId, ProjectId: projectId, Role: assignment.Role}).Times(1).Return(&updateAssignment, nil)
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{assignment}, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{updateAssignment}, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{updateAssignment}, nil),
			)
			mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: updateAssignment.TeamId, ProjectId: projectId}).Times(1).Return(nil)
		})
	})

	t.Run("create and update custom assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        assignmentCustom.TeamId,
						"project_id":     projectId,
						"custom_role_id": customRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", assignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", customRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        updateAssignmentCustom.TeamId,
						"project_id":     projectId,
						"custom_role_id": updatedCustomRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "custom_role_id", updatedCustomRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    updateAssignmentCustom.TeamId,
						"project_id": projectId,
						"role":       updateAssignment.Role,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "team_id", updateAssignmentCustom.TeamId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "role", updateAssignment.Role),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignmentCustom.TeamId, ProjectId: projectId, Role: assignmentCustom.Role}).Times(1).Return(&assignmentCustom, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(2).Return([]client.TeamRoleAssignmentPayload{assignmentCustom}, nil),
				mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: assignmentCustom.TeamId, ProjectId: projectId}).Times(1).Return(nil),
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: updateAssignmentCustom.TeamId, ProjectId: projectId, Role: updateAssignmentCustom.Role}).Times(1).Return(&updateAssignmentCustom, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(2).Return([]client.TeamRoleAssignmentPayload{updateAssignmentCustom}, nil),
				mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: updateAssignmentCustom.TeamId, ProjectId: projectId, Role: updateAssignment.Role}).Times(1).Return(&updateAssignment, nil),
				mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{updateAssignment}, nil),
				mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: updateAssignment.TeamId, ProjectId: projectId}).Times(1).Return(nil),
			)
		})
	})

	t.Run("import - built-in role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": projectId,
						"role":       assignment.Role,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignment.TeamId + "_" + projectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignment.TeamId, ProjectId: projectId, Role: assignment.Role}).Times(1).Return(&assignment, nil)
			mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(3).Return([]client.TeamRoleAssignmentPayload{assignment}, nil)
			mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: assignment.TeamId, ProjectId: projectId}).Times(1).Return(nil)
		})
	})

	t.Run("import - custom role", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":        assignmentCustom.TeamId,
						"project_id":     projectId,
						"custom_role_id": assignmentCustom.Role,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignmentCustom.TeamId + "_" + projectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignmentCustom.TeamId, ProjectId: projectId, Role: assignmentCustom.Role}).Times(1).Return(&assignmentCustom, nil)
			mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(3).Return([]client.TeamRoleAssignmentPayload{assignmentCustom}, nil)
			mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: assignmentCustom.TeamId, ProjectId: projectId}).Times(1).Return(nil)
		})
	})

	t.Run("Import role - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"team_id":    assignment.TeamId,
						"project_id": projectId,
						"role":       assignment.Role,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignment.TeamId + "_" + projectId,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamRoleAssignmentCreateOrUpdate(&client.TeamRoleAssignmentCreateOrUpdatePayload{TeamId: assignment.TeamId, ProjectId: projectId, Role: assignment.Role}).Times(1).Return(&assignment, nil)
			mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{assignment}, nil)
			mock.EXPECT().TeamRoleAssignments(&client.TeamRoleAssignmentListPayload{ProjectId: projectId}).Times(1).Return([]client.TeamRoleAssignmentPayload{}, nil)
			mock.EXPECT().TeamRoleAssignmentDelete(&client.TeamRoleAssignmentDeletePayload{TeamId: assignment.TeamId, ProjectId: projectId}).Times(1)
		})
	})
}
