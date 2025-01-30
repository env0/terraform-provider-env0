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

func TestUnitAwsCloudConfigurationResource(t *testing.T) {
	resourceType := "env0_aws_cloud_configuration"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	awsConfig := client.AWSCloudAccountConfiguration{
		AccountId:                   "acc1",
		BucketName:                  "b1",
		Regions:                     []string{"us-west-1", "us-east-1000"},
		ShouldPrefixUnderLogsFolder: true,
	}

	updatedAwsConfig := client.AWSCloudAccountConfiguration{
		AccountId:  "acc2",
		BucketName: "b2",
		Regions:    []string{"us-west-12", "us-east-1040"},
		Prefix:     "////",
	}

	cloudConfig := client.CloudAccount{
		Id:            uuid.NewString(),
		Provider:      "AWS",
		Name:          "name1",
		Health:        false,
		Configuration: &awsConfig,
	}

	updatedCloudConfig := cloudConfig
	updatedCloudConfig.Name = "name2"
	updatedCloudConfig.Configuration = &updatedAwsConfig
	updatedCloudConfig.Health = true

	createPayload := client.CloudAccountCreatePayload{
		Name:          cloudConfig.Name,
		Provider:      "AWS",
		Configuration: &awsConfig,
	}

	updatePayload := client.CloudAccountUpdatePayload{
		Name:          updatedCloudConfig.Name,
		Configuration: &updatedAwsConfig,
	}

	otherCloudConfig := client.CloudAccount{
		Id:       uuid.NewString(),
		Provider: "AWS",
		Name:     "other_name",
	}

	getFields := func(cloudConfig *client.CloudAccount) map[string]interface{} {
		awsConfig := cloudConfig.Configuration.(*client.AWSCloudAccountConfiguration)

		fields := map[string]interface{}{
			"name":                            cloudConfig.Name,
			"account_id":                      awsConfig.AccountId,
			"bucket_name":                     awsConfig.BucketName,
			"regions":                         awsConfig.Regions,
			"should_prefix_under_logs_folder": awsConfig.ShouldPrefixUnderLogsFolder,
		}

		if awsConfig.Prefix != "" {
			fields["prefix"] = awsConfig.Prefix
		}

		return fields
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&cloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", cloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "account_id", awsConfig.AccountId),
						resource.TestCheckResourceAttr(accessor, "bucket_name", awsConfig.BucketName),
						resource.TestCheckResourceAttr(accessor, "regions.0", awsConfig.Regions[0]),
						resource.TestCheckResourceAttr(accessor, "regions.1", awsConfig.Regions[1]),
						resource.TestCheckResourceAttr(accessor, "health", "false"),
						resource.TestCheckResourceAttr(accessor, "should_prefix_under_logs_folder", "true"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, getFields(&updatedCloudConfig)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", updatedCloudConfig.Name),
						resource.TestCheckResourceAttr(accessor, "account_id", updatedAwsConfig.AccountId),
						resource.TestCheckResourceAttr(accessor, "bucket_name", updatedAwsConfig.BucketName),
						resource.TestCheckResourceAttr(accessor, "regions.0", updatedAwsConfig.Regions[0]),
						resource.TestCheckResourceAttr(accessor, "regions.1", updatedAwsConfig.Regions[1]),
						resource.TestCheckResourceAttr(accessor, "prefix", updatedAwsConfig.Prefix),
						resource.TestCheckResourceAttr(accessor, "health", "true"),
						resource.TestCheckResourceAttr(accessor, "should_prefix_under_logs_folder", "false"),
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

	t.Run("import by id not found", func(t *testing.T) {
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
					ExpectError:       regexp.MustCompile("failed to get clound configuration: not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(&cloudConfig, nil),
				mock.EXPECT().CloudAccount(cloudConfig.Id).Times(1).Return(nil, &client.NotFoundError{}),
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
			gomock.InOrder(
				mock.EXPECT().CloudAccountCreate(&createPayload).Times(1).Return(nil, errors.New("error")),
			)
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
