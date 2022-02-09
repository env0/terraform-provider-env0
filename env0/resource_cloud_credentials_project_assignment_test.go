package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitResourceCloudCredentialsProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_cloud_credentials_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	assignment := client.CloudCredentialsProjectAssignment{
		CredentialId: "cred-it",
		ProjectId:    "proj-it",
	}
	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"credential_id": assignment.CredentialId,
		"project_id":    assignment.ProjectId,
	})
	t.Run("Create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialId+"|"+assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
			mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{assignment.CredentialId}, nil)
			mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
		})
	})
	t.Run("Create with api prob", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(could not assign cloud credentials to project)`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(client.CloudCredentialsProjectAssignment{}, errors.New("err"))
		})
	})
	t.Run("Read with api prob", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(could not get cloud_credentials:)`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
			mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{}, errors.New("err"))
			mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
		})
	})
	t.Run("Read with multi values", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialId+"|"+assignment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialId),
						resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
			mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{assignment.CredentialId, "1", "2"}, nil)
			mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
		})
	})
	t.Run("Read didnt api didnt return correct stuff", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(could not find cloud credential project assignment)`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
			mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{"1", "2"}, nil)
			mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
		})
	})
}
