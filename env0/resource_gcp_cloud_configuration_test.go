package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitGcpCloudConfigurationResource(t *testing.T) {
	resourceType := "env0_gcp_cloud_configuration"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	gcpConfig := client.GCPCloudAccountConfiguration{
		GcpProjectId:                       "project-123",
		CredentialConfigurationFileContent: "{initialContent}",
	}

	updatedGcpConfig := client.GCPCloudAccountConfiguration{
		GcpProjectId:                       "project-456",
		CredentialConfigurationFileContent: "{updatedContent}",
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

	otherCloudConfig := client.CloudAccount{
		Id:       uuid.NewString(),
		Provider: "GCP",
		Name:     "other_name",
	}

	getFields := func(cloudConfig *client.CloudAccount) map[string]any {
		gcpConfig := cloudConfig.Configuration.(*client.GCPCloudAccountConfiguration)

		return map[string]any{
			"name":                                  cloudConfig.Name,
			"gcp_project_id":                        gcpConfig.GcpProjectId,
			"credential_configuration_file_content": gcpConfig.CredentialConfigurationFileContent,
		}
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", cloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "gcp_project_id", gcpConfig.GcpProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_configuration_file_content", gcpConfig.CredentialConfigurationFileContent),
						resource.TestCheckResourceAttr(accessor, "health", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&updatedCloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", updatedCloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "gcp_project_id", updatedGcpConfig.GcpProjectId),
						resource.TestCheckResourceAttr(accessor, "credential_configuration_file_content", updatedGcpConfig.CredentialConfigurationFileContent),
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

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     cloudConfig.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccounts().Times(1).Return([]client.CloudAccount{otherCloudConfig, cloudConfig}, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccountDelete(cloudConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     cloudConfig.Name,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("cloud configuration called '%s' was not found", cloudConfig.Name)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccounts().Times(1).Return([]client.CloudAccount{otherCloudConfig}, nil),
				mock.EXPECT().CloudAccountDelete(cloudConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
				},
			},
		}, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccountDelete(cloudConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("create failed", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
					ExpectError: regexp.MustCompile("failed to create a cloud configuration: error"),
				},
			},
		}, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("update failed", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
				},
				{
					Config:      resourceConfigCreate(resourceType, resourceName, getFields(&updatedCloudConfig)),
					ExpectError: regexp.MustCompile("failed to update cloud configuration: error"),
				},
			},
		}, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(2).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccountUpdate(cloudConfig.Id, &updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().CloudAccountDelete(updatedCloudConfig.Id).Times(1).Return(nil),
			)
		})
	})
}
