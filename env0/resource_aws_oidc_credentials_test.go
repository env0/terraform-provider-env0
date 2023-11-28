package env0

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitAwsOidcCredentialsResource(t *testing.T) {
	resourceType := "env0_aws_oidc_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	duration := 3600
	updatedDuration := 2 * duration

	awsCredentialsResource := map[string]interface{}{
		"name":     "test",
		"role_arn": "11111",
		"duration": strconv.Itoa(duration),
	}

	updatedAwsCredentialsResource := map[string]interface{}{
		"name":     "test",
		"role_arn": "22222",
		"duration": strconv.Itoa(updatedDuration),
	}

	createPayload := client.AwsCredentialsCreatePayload{
		Name: awsCredentialsResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:  awsCredentialsResource["role_arn"].(string),
			Duration: duration,
		},
		Type: client.AwsOidcCredentialsType,
	}

	updatePayload := client.AwsCredentialsCreatePayload{
		Value: client.AwsCredentialsValuePayload{
			RoleArn:  updatedAwsCredentialsResource["role_arn"].(string),
			Duration: updatedDuration,
		},
		Type: client.AwsOidcCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AwsOidcCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_....",
	}

	updateReturnValues := client.Credentials{
		Id:             returnValues.Id,
		Name:           returnValues.Name,
		OrganizationId: "id",
		Type:           string(client.AwsOidcCredentialsType),
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "role_arn", awsCredentialsResource["role_arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
					resource.TestCheckResourceAttr(accessor, "duration", awsCredentialsResource["duration"].(string)),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedAwsCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedAwsCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "role_arn", updatedAwsCredentialsResource["role_arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
					resource.TestCheckResourceAttr(accessor, "duration", updatedAwsCredentialsResource["duration"].(string)),
				),
			},
		},
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, testCaseForCreateAndUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updatePayload).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, awsCredentialsResource)

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
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           awsCredentialsResource["name"].(string),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"role_arn", "duration"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           returnValues.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"role_arn", "duration"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsCredentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     otherTypeReturnValues.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("credentials not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{}),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, awsCredentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     awsCredentialsResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", awsCredentialsResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})
}
