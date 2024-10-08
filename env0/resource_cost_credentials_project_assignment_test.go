package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitResourceCostCredentialsProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_cost_credentials_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	anotherAssignment := client.CostCredentialProjectAssignment{
		CredentialsId:   "cred-it-another",
		ProjectId:       "proj-it",
		CredentialsType: "AWS_ASSUMED_ROLE",
	}
	assignment := client.CostCredentialProjectAssignment{
		CredentialsId:   "cred-it",
		ProjectId:       "proj-it",
		CredentialsType: "AWS_ASSUMED_ROLE",
	}

	assignmentForDrift := client.CostCredentialProjectAssignment{
		CredentialsId:   "cred-it",
		ProjectId:       "proj-it-update",
		CredentialsType: "AWS_ASSUMED_ROLE",
	}
	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"credential_id": assignment.CredentialsId,
		"project_id":    assignment.ProjectId,
	})

	t.Run("Create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialsId+"|"+assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialsId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCostCredentialsToProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(assignment, nil)
			mock.EXPECT().CostCredentialIdsInProject(assignment.ProjectId).Times(1).
				Return([]client.CostCredentialProjectAssignment{assignment}, nil)
			mock.EXPECT().RemoveCostCredentialsFromProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(nil)
		})
	})
	t.Run("Create with api prob", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`could not assign cost credentials to project: err`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCostCredentialsToProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(client.CostCredentialProjectAssignment{}, errors.New("err"))
		})
	})
	t.Run("Read with api prob", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`could not get cost credentials: err`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCostCredentialsToProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(assignment, nil)
			mock.EXPECT().CostCredentialIdsInProject(assignment.ProjectId).Times(1).
				Return([]client.CostCredentialProjectAssignment{}, errors.New("err"))
			mock.EXPECT().RemoveCostCredentialsFromProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(nil)
		})
	})
	t.Run("Read with multi values", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialsId+"|"+assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialsId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCostCredentialsToProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(assignment, nil)
			mock.EXPECT().CostCredentialIdsInProject(assignment.ProjectId).Times(1).
				Return([]client.CostCredentialProjectAssignment{assignment, anotherAssignment}, nil)
			mock.EXPECT().RemoveCostCredentialsFromProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(nil)
		})
	})
	t.Run("detect drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialsId+"|"+assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialsId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"credential_id": assignmentForDrift.CredentialsId,
						"project_id":    assignmentForDrift.ProjectId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignmentForDrift.CredentialsId+"|"+assignmentForDrift.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignmentForDrift.CredentialsId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignmentForDrift.ProjectId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCostCredentialsToProject(assignment.ProjectId, assignment.CredentialsId).Times(1).Return(assignment, nil)
			mock.EXPECT().AssignCostCredentialsToProject(assignmentForDrift.ProjectId, assignmentForDrift.CredentialsId).Times(1).Return(assignmentForDrift, nil)
			mock.EXPECT().RemoveCostCredentialsFromProject(assignmentForDrift.ProjectId, assignmentForDrift.CredentialsId).Times(1).Return(nil)
			gomock.InOrder(
				mock.EXPECT().CostCredentialIdsInProject(assignment.ProjectId).Times(1).
					Return([]client.CostCredentialProjectAssignment{assignment, anotherAssignment}, nil),
				mock.EXPECT().CostCredentialIdsInProject(assignment.ProjectId).Times(1).
					Return([]client.CostCredentialProjectAssignment{anotherAssignment}, nil),
				mock.EXPECT().CostCredentialIdsInProject(assignmentForDrift.ProjectId).Times(1).
					Return([]client.CostCredentialProjectAssignment{assignmentForDrift, anotherAssignment}, nil),
			)
		})
	})
}
