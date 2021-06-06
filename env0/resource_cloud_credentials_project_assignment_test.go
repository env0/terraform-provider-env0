package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitResourceCloudCredentialsProjectAssignmentResource_Create(t *testing.T) {
	resourceType := "env0_cloud_credentials_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	assignment := client.CloudCredentialsProjectAssignment{
		CredentialId: "cred-it",
		ProjectId:    "proj-it",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"credential_id": assignment.CredentialId,
					"project_id":    assignment.ProjectId,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", assignment.CredentialId+"|"+assignment.ProjectId),
					resource.TestCheckResourceAttr(accessor, "credential_id", assignment.CredentialId),
					resource.TestCheckResourceAttr(accessor, "project_id", assignment.ProjectId),
				),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
		mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{assignment.CredentialId}, nil)
		mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
	})
}
func TestUnitResourceCloudCredentialsProjectAssignmentResource_CreateApiProb(t *testing.T) {
	resourceType := "env0_cloud_credentials_project_assignment"
	resourceName := "test"
	assignment := client.CloudCredentialsProjectAssignment{
		CredentialId: "cred-it",
		ProjectId:    "proj-it",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"credential_id": assignment.CredentialId,
					"project_id":    assignment.ProjectId,
				}),
				ExpectError: regexp.MustCompile(`(could not assign cloud credentials to project)`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(client.CloudCredentialsProjectAssignment{}, errors.New("err"))
	})
}
func TestUnitResourceCloudCredentialsProjectAssignmentResource_ReadApiProb(t *testing.T) {
	resourceType := "env0_cloud_credentials_project_assignment"
	resourceName := "test"
	assignment := client.CloudCredentialsProjectAssignment{
		CredentialId: "cred-it",
		ProjectId:    "proj-it",
	}

	createTestCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"credential_id": assignment.CredentialId,
					"project_id":    assignment.ProjectId,
				}),
				ExpectError: regexp.MustCompile(`(could not get cloud_credentials:)`),
			},
		},
	}

	runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().AssignCloudCredentialsToProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(assignment, nil)
		mock.EXPECT().CloudCredentialIdsInProject(assignment.ProjectId).Times(1).Return([]string{}, errors.New("err"))
		mock.EXPECT().RemoveCloudCredentialsFromProject(assignment.ProjectId, assignment.CredentialId).Times(1).Return(nil)
	})
}
