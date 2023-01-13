package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitResourceCustomFlowAssignmentResource(t *testing.T) {
	resourceType := "env0_custom_flow_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	anotherAssignment := client.CustomFlowAssignment{
		Scope:       "PROJECT",
		ScopeId:     "scope_id",
		BlueprintId: "other_blueprint_id",
	}

	assignment := client.CustomFlowAssignment{
		Scope:       "PROJECT",
		ScopeId:     "scope_id",
		BlueprintId: "blueprint_id",
	}

	assignmentNoScope := client.CustomFlowAssignment{
		ScopeId:     assignment.ScopeId,
		BlueprintId: assignment.BlueprintId,
	}

	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"scope_id":    assignment.ScopeId,
		"template_id": assignment.BlueprintId,
	})

	t.Run("Create assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.ScopeId+"|"+assignment.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "scope_id", assignment.ScopeId),
						resource.TestCheckResourceAttr(accessor, "scope", "PROJECT"),
						resource.TestCheckResourceAttr(accessor, "template_id", assignment.BlueprintId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowAssign([]client.CustomFlowAssignment{assignment}).Times(1).Return(nil),
				mock.EXPECT().CustomFlowGetAssignments([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return([]client.CustomFlowAssignment{anotherAssignment, assignment}, nil),
				mock.EXPECT().CustomFlowUnassign([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create assignment failed", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not assign custom flow to project: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlowAssign([]client.CustomFlowAssignment{assignment}).Times(1).Return(errors.New("error"))
		})
	})

	t.Run("Create assignment read failed", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(fmt.Sprintf("could not get custom flow assignments for id %s: error", assignment.ScopeId)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowAssign([]client.CustomFlowAssignment{assignment}).Times(1).Return(nil),
				mock.EXPECT().CustomFlowGetAssignments([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().CustomFlowUnassign([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Detect drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.ScopeId+"|"+assignment.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "scope_id", assignment.ScopeId),
						resource.TestCheckResourceAttr(accessor, "scope", "PROJECT"),
						resource.TestCheckResourceAttr(accessor, "template_id", assignment.BlueprintId),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowAssign([]client.CustomFlowAssignment{assignment}).Times(1).Return(nil),
				mock.EXPECT().CustomFlowGetAssignments([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return([]client.CustomFlowAssignment{anotherAssignment}, nil),
				mock.EXPECT().CustomFlowUnassign([]client.CustomFlowAssignment{assignmentNoScope}).Times(1).Return(nil),
			)
		})
	})

}
