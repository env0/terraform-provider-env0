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

func TestUnitResourceApprovalPolicyAssignmentResource(t *testing.T) {
	resourceType := "env0_approval_policy_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	assignment := client.ApprovalPolicyAssignment{
		Scope:       "PROJECT",
		ScopeId:     "scope_id",
		BlueprintId: "blueprint_id",
	}

	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"scope_id":     assignment.ScopeId,
		"blueprint_id": assignment.BlueprintId,
	})

	validTemplate := client.Template{
		Id:   assignment.BlueprintId,
		Type: "approval-policy",
		Name: "approval-policy-" + string(assignment.Scope) + "-" + assignment.ScopeId,
	}

	approvalPolicyByScope := client.ApprovalPolicyByScope{
		Scope:   string(assignment.Scope),
		ScopeId: assignment.ScopeId,
		ApprovalPolicy: &client.ApprovalPolicy{
			Id:   validTemplate.Id,
			Name: validTemplate.Name,
		},
	}

	t.Run("Create assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", fmt.Sprintf("%s|%s|%s", assignment.BlueprintId, assignment.Scope, assignment.ScopeId)),
						resource.TestCheckResourceAttr(accessor, "scope_id", assignment.ScopeId),
						resource.TestCheckResourceAttr(accessor, "scope", "PROJECT"),
						resource.TestCheckResourceAttr(accessor, "blueprint_id", assignment.BlueprintId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(validTemplate, nil),
				mock.EXPECT().ApprovalPolicyAssign(&assignment).Times(1).Return(&assignment, nil),
				mock.EXPECT().ApprovalPolicyByScope(string(assignment.Scope), assignment.ScopeId).Times(1).Return([]client.ApprovalPolicyByScope{approvalPolicyByScope}, nil),
				mock.EXPECT().ApprovalPolicyUnassign(string(assignment.Scope), assignment.ScopeId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create assignment - type mismatch", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("template with id"),
				},
			},
		}

		invalidTemplate := validTemplate
		invalidTemplate.Type = "terraform"

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(invalidTemplate, nil),
			)
		})
	})

	t.Run("Create assignment - name mismatch", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("template name is"),
				},
			},
		}

		invalidTemplate := validTemplate
		invalidTemplate.Name = "bad name"

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(invalidTemplate, nil),
			)
		})
	})

	t.Run("Create assignment - error when assigning", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not assign approval policy"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(validTemplate, nil),
				mock.EXPECT().ApprovalPolicyAssign(&assignment).Times(1).Return(nil, errors.New("error")),
			)
		})
	})

	t.Run("Create assignment - error when requesting template", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("unable to get template with"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(client.Template{}, errors.New("error")),
			)
		})
	})

	t.Run("Detect drift - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", fmt.Sprintf("%s|%s|%s", assignment.BlueprintId, assignment.Scope, assignment.ScopeId)),
						resource.TestCheckResourceAttr(accessor, "scope_id", assignment.ScopeId),
						resource.TestCheckResourceAttr(accessor, "scope", "PROJECT"),
						resource.TestCheckResourceAttr(accessor, "blueprint_id", assignment.BlueprintId),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(validTemplate, nil),
				mock.EXPECT().ApprovalPolicyAssign(&assignment).Times(1).Return(&assignment, nil),
				mock.EXPECT().ApprovalPolicyByScope(string(assignment.Scope), assignment.ScopeId).Times(1).Return(nil, &client.NotFoundError{}),
				mock.EXPECT().ApprovalPolicyUnassign(string(assignment.Scope), assignment.ScopeId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Detect drift - mismatch", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", fmt.Sprintf("%s|%s|%s", assignment.BlueprintId, assignment.Scope, assignment.ScopeId)),
						resource.TestCheckResourceAttr(accessor, "scope_id", assignment.ScopeId),
						resource.TestCheckResourceAttr(accessor, "scope", "PROJECT"),
						resource.TestCheckResourceAttr(accessor, "blueprint_id", assignment.BlueprintId),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		approvalPolicyByScopeMismatch := client.ApprovalPolicyByScope{
			Scope:   string(assignment.Scope),
			ScopeId: assignment.ScopeId,
			ApprovalPolicy: &client.ApprovalPolicy{
				Id:   "other_id",
				Name: validTemplate.Name,
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(assignment.BlueprintId).Times(1).Return(validTemplate, nil),
				mock.EXPECT().ApprovalPolicyAssign(&assignment).Times(1).Return(&assignment, nil),
				mock.EXPECT().ApprovalPolicyByScope(string(assignment.Scope), assignment.ScopeId).Times(1).Return([]client.ApprovalPolicyByScope{approvalPolicyByScopeMismatch}, nil),
				mock.EXPECT().ApprovalPolicyUnassign(string(assignment.Scope), assignment.ScopeId).Times(1).Return(nil),
			)
		})
	})
}
