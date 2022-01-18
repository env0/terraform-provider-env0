package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	//	"github.com/golang/mock/gomock"
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

	awsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: "test",
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    awsCredentialResoure["arn"].(string),
			ExternalId: awsCredentialResoure["external_id"].(string),
		},
	}

	returnValues := client.ApiKey{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	// updatedawsCredentialResoure := map[string]interface{}{
	// 	"name":        "update",
	// 	"arn":         "11111",
	// 	"external_id": "22222",
	// }

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

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsCreate(awsCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().AwsCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().AwsCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

}
