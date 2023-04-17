package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitEnvironmentStateAccessResource(t *testing.T) {
	resourceType := "env0_environment_state_access"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	remoteState := client.RemoteStateAccessConfiguration{
		EnvironmentId: "env",
		AllowedProjectIds: []string{
			"pr1",
		},
	}

	updatedRemoteState := client.RemoteStateAccessConfiguration{
		EnvironmentId:                    remoteState.EnvironmentId,
		AccessibleFromEntireOrganization: true,
	}

	createPayload := client.RemoteStateAccessConfigurationCreate{
		AllowedProjectIds: remoteState.AllowedProjectIds,
	}

	updatePayload := client.RemoteStateAccessConfigurationCreate{
		AccessibleFromEntireOrganization: true,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":      remoteState.EnvironmentId,
						"allowed_project_ids": remoteState.AllowedProjectIds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", remoteState.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "environment_id", remoteState.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "allowed_project_ids.0", remoteState.AllowedProjectIds[0]),
						resource.TestCheckResourceAttr(accessor, "accessible_from_entire_organization", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":                      remoteState.EnvironmentId,
						"accessible_from_entire_organization": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedRemoteState.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "environment_id", updatedRemoteState.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "accessible_from_entire_organization", "true"),
						resource.TestCheckNoResourceAttr(accessor, "allowed_project_ids"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().RemoteStateAccessConfigurationCreate(remoteState.EnvironmentId, createPayload).Times(1).Return(&remoteState, nil),
				mock.EXPECT().RemoteStateAccessConfiguration(remoteState.EnvironmentId).Times(2).Return(&remoteState, nil),
				mock.EXPECT().RemoteStateAccessConfigurationDelete(remoteState.EnvironmentId).Times(1).Return(nil),
				mock.EXPECT().RemoteStateAccessConfigurationCreate(remoteState.EnvironmentId, updatePayload).Times(1).Return(&updatedRemoteState, nil),
				mock.EXPECT().RemoteStateAccessConfiguration(remoteState.EnvironmentId).Times(1).Return(&updatedRemoteState, nil),
				mock.EXPECT().RemoteStateAccessConfigurationDelete(remoteState.EnvironmentId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":      remoteState.EnvironmentId,
						"allowed_project_ids": remoteState.AllowedProjectIds,
					}),
					ExpectError: regexp.MustCompile("could not create a remote state access configation: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().RemoteStateAccessConfigurationCreate(remoteState.EnvironmentId, createPayload).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Create Failure - conflict", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":                      remoteState.EnvironmentId,
						"allowed_project_ids":                 remoteState.AllowedProjectIds,
						"accessible_from_entire_organization": "true",
					}),
					ExpectError: regexp.MustCompile("'allowed_project_ids' should not be set when 'accessible_from_entire_organization' is set to 'true'"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
