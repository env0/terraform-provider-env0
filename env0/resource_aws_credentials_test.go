package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitAwsCredentialsResource(t *testing.T) {

	resourceType := "env0_aws_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	awsCredentialResoure := map[string]interface{}{
		"name":        "test",
		"arn":         "11111",
		"external_id": "22222",
	}

	updatedawsCredentialResoure := map[string]interface{}{
		"name":        "update",
		"arn":         "11111",
		"external_id": "22222",
	}

	awsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: awsCredentialResoure["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    awsCredentialResoure["arn"].(string),
			ExternalId: awsCredentialResoure["external_id"].(string),
		},
	}

	updateAwsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: updatedawsCredentialResoure["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    updatedawsCredentialResoure["arn"].(string),
			ExternalId: updatedawsCredentialResoure["external_id"].(string),
		},
	}

	returnValues := client.ApiKey{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	updateReturnValues := client.ApiKey{
		Id:             "id2",
		Name:           "update",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResoure),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResoure["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResoure["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsCredentialResoure["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", "id"),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResoure),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResoure["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResoure["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", awsCredentialResoure["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedawsCredentialResoure),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedawsCredentialResoure["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", updatedawsCredentialResoure["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "external_id", updatedawsCredentialResoure["external_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	testCaseForError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().AwsCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().AwsCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("destroy when one of the values changed", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().AwsCredentialsDelete(returnValues.Id).Times(1).Return(nil)
			mock.EXPECT().AwsCredentialsCreate(updateAwsCredCreatePayload).Times(1).Return(updateReturnValues, nil)

			gomock.InOrder(
				mock.EXPECT().AwsCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().AwsCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().AwsCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("throw error when one of the values is missing", func(t *testing.T) {
		runUnitTest(t, testCaseForError, func(mock *client.MockApiClientInterface) {
			//mock.EXPECT().AwsCredentialsCreate(gomock.Any()).Return(client.ApiKey{}, "error")
		})
	})

}
