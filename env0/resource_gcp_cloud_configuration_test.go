package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitGcpCloudConfigurationResource(t *testing.T) {
	resourceType := "env0_gcp_cloud_configuration"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	gcpConfig := client.GCPCloudAccountConfiguration{
		GcpProjectId:                       "test-project-123",
		CredentialConfigurationFileContent: "test-credentials-json",
	}

	updatedGcpConfig := client.GCPCloudAccountConfiguration{
		GcpProjectId:                       "test-project-456",
		CredentialConfigurationFileContent: "updated-test-credentials-json",
	}

	cloudConfig := client.CloudAccount{
		Id:            uuid.NewString(),
		Provider:      "GCP",
		Name:          "name1",
		Health:        false,
		Configuration: &gcpConfig,
	}

	updatedCloudConfig := cloudConfig
	updatedCloudConfig.Name = "name2"
	updatedCloudConfig.Configuration = &updatedGcpConfig
	updatedCloudConfig.Health = true

	createPayload := client.CloudAccountCreatePayload{
		Name:          cloudConfig.Name,
		Provider:      "GCP",
		Configuration: &gcpConfig,
	}

	updatePayload := client.CloudAccountUpdatePayload{
		Name:          updatedCloudConfig.Name,
		Configuration: &updatedGcpConfig,
	}

	getFields := func(cloudConfig *client.CloudAccount) map[string]any {
		gcpConfig := cloudConfig.Configuration.(*client.GCPCloudAccountConfiguration)
		return map[string]any{
			"name":                            cloudConfig.Name,
			"project_id":                      gcpConfig.GcpProjectId,
			"json_configuration_file_content": gcpConfig.CredentialConfigurationFileContent,
		}
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", cloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", gcpConfig.GcpProjectId),
						resource.TestCheckResourceAttr(accessor, "json_configuration_file_content", gcpConfig.CredentialConfigurationFileContent),
						resource.TestCheckResourceAttr(accessor, "health", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&updatedCloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", updatedCloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedGcpConfig.GcpProjectId),
						resource.TestCheckResourceAttr(accessor, "json_configuration_file_content", updatedGcpConfig.CredentialConfigurationFileContent),
						resource.TestCheckResourceAttr(accessor, "health", "true"),
					),
				},
			},
		}, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(2).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccountUpdate(cloudConfig.Id, &updatePayload).Times(1).Return(&updatedCloudConfig, nil),
				mock.EXPECT().CloudAccount(updatedCloudConfig.Id).Times(1).Return(&updatedCloudConfig, nil),
				mock.EXPECT().CloudAccountDelete(updatedCloudConfig.Id).Times(1).Return(nil),
			)
		})
	})
}
