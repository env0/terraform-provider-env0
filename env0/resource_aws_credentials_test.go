package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitAwsCredentialsResource(t *testing.T) {
	resourceType := "env0_aws_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	awsArnCredentialResource := map[string]interface{}{
		"name":        "test",
		"arn":         "11111",
		"external_id": "22222",
	}

	updatedAwsAccessKeyCredentialResource := map[string]interface{}{
		"name":              "update",
		"access_key_id":     "33333",
		"secret_access_key": "44444",
	}

	awsArnCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: awsArnCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    awsArnCredentialResource["arn"].(string),
			ExternalId: awsArnCredentialResource["external_id"].(string),
		},
		Type: client.AwsAssumedRoleCredentialsType,
	}

	updateAwsAccessKeyCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: updatedAwsAccessKeyCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			AccessKeyId:     updatedAwsAccessKeyCredentialResource["access_key_id"].(string),
			SecretAccessKey: updatedAwsAccessKeyCredentialResource["secret_access_key"].(string),
		},
		Type: client.AwsAccessKeysCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "GCP_....",
	}

	updateReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30c",
		Name:           "update",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsArnCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsArnCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsArnCredentialResource["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsArnCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsArnCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsArnCredentialResource["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedAwsAccessKeyCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedAwsAccessKeyCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "access_key_id", updatedAwsAccessKeyCredentialResource["access_key_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret_access_key", updatedAwsAccessKeyCredentialResource["secret_access_key"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	mutuallyExclusiveErrorResource := map[string]interface{}{
		"name":          "update",
		"arn":           "11111",
		"external_id":   "22222",
		"access_key_id": "some-key",
	}
	testCaseFormMutuallyExclusiveError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, mutuallyExclusiveErrorResource),
				ExpectError: regexp.MustCompile(`"external_id": conflicts with access_key_id`),
			},
		},
	}

	missingValidInputErrorResource := map[string]interface{}{
		"name": "update",
	}
	testCaseFormMissingValidInputError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, missingValidInputErrorResource),
				ExpectError: regexp.MustCompile("one of `access_key_id,arn` must be specified"),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().AwsCredentialsCreate(updateAwsAccessKeyCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("throw error when enter mutually exclusive values", func(t *testing.T) {
		runUnitTest(t, testCaseFormMutuallyExclusiveError, func(mock *client.MockApiClientInterface) {
		})
	})

	t.Run("throw error when don't enter any valid options", func(t *testing.T) {
		runUnitTest(t, testCaseFormMissingValidInputError, func(mock *client.MockApiClientInterface) {
		})
	})

	t.Run("AWS credentials removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource)

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
				},
				{
					Config: stepConfig,
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     awsArnCredentialResource["name"].(string),
					ImportStateVerify: false,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     returnValues.Id,
					ImportStateVerify: false,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     otherTypeReturnValues.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("aws credentials resource with id %v not found", otherTypeReturnValues.Id)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{})
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsArnCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     awsArnCredentialResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", awsArnCredentialResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsArnCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})
}
