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

func TestUnitAzureCloudConfigurationResource(t *testing.T) {
	resourceType := "env0_azure_cloud_configuration"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	azureConfig := client.AzureCloudAccountConfiguration{
		TenantId:                "tenant123",
		ClientId:                "client123",
		LogAnalyticsWorkspaceId: "workspace123",
	}

	updatedAzureConfig := client.AzureCloudAccountConfiguration{
		TenantId:                "tenant456",
		ClientId:                "client456",
		LogAnalyticsWorkspaceId: "workspace456",
	}

	cloudConfig := client.CloudAccount{
		Id:            uuid.NewString(),
		Provider:      "AzureLAW",
		Name:          "name1",
		Health:        false,
		Configuration: &azureConfig,
	}

	updatedCloudConfig := cloudConfig
	updatedCloudConfig.Name = "name2"
	updatedCloudConfig.Configuration = &updatedAzureConfig
	updatedCloudConfig.Health = true

	createPayload := client.CloudAccountCreatePayload{
		Name:          cloudConfig.Name,
		Provider:      "AzureLAW",
		Configuration: &azureConfig,
	}

	updatePayload := client.CloudAccountUpdatePayload{
		Name:          updatedCloudConfig.Name,
		Configuration: &updatedAzureConfig,
	}

	otherCloudConfig := client.CloudAccount{
		Id:       uuid.NewString(),
		Provider: "AzureLAW",
		Name:     "other_name",
	}

	getFields := func(cloudConfig *client.CloudAccount) map[string]interface{} {
		azureConfig := cloudConfig.Configuration.(*client.AzureCloudAccountConfiguration)

		return map[string]interface{}{
			"name":                       cloudConfig.Name,
			"tenant_id":                  azureConfig.TenantId,
			"client_id":                  azureConfig.ClientId,
			"log_analytics_workspace_id": azureConfig.LogAnalyticsWorkspaceId,
		}
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", cloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "tenant_id", azureConfig.TenantId),
						resource.TestCheckResourceAttr(accessor, "client_id", azureConfig.ClientId),
						resource.TestCheckResourceAttr(accessor, "log_analytics_workspace_id", azureConfig.LogAnalyticsWorkspaceId),
						resource.TestCheckResourceAttr(accessor, "health", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&updatedCloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", updatedCloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "tenant_id", updatedAzureConfig.TenantId),
						resource.TestCheckResourceAttr(accessor, "client_id", updatedAzureConfig.ClientId),
						resource.TestCheckResourceAttr(accessor, "log_analytics_workspace_id", updatedAzureConfig.LogAnalyticsWorkspaceId),
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
					ImportStateId:     cloudConfig.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(3).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccountDelete(cloudConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
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
