package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitAwsCostCredentialsResource(t *testing.T) {

	resourceType := "env0_aws_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	awsCredentialResource := map[string]interface{}{
		"name":        "test",
		"arn":         "11111",
		"external_id": "22222",
	}

	updatedAwsCredentialResource := map[string]interface{}{
		"name":        "update",
		"arn":         "33333",
		"external_id": "44444",
	}

	awsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: awsCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    awsCredentialResource["arn"].(string),
			ExternalId: awsCredentialResource["external_id"].(string),
		},
		Type: client.AwsCostCredentialsType,
	}

	updateAwsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: updatedAwsCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    updatedAwsCredentialResource["arn"].(string),
			ExternalId: updatedAwsCredentialResource["external_id"].(string),
		},
		Type: client.AwsCostCredentialsType,
	}

	returnValues := client.ApiKey{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	updateReturnValues := client.ApiKey{
		Id:             "id2",
		Name:           "update",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsCredentialResource["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", "id"),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsCredentialResource["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedAwsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedAwsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", updatedAwsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", updatedAwsCredentialResource["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
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
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AwsCredentialsCreate(awsCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().AwsCredentialsCreate(updateAwsCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("throw error when don't enter any valid options", func(t *testing.T) {
		runUnitTest(t, testCaseFormMissingValidInputError, func(mock *client.MockApiClientInterface) {
		})
	})

}
