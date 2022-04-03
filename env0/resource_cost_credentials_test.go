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

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	updateReturnValues := client.Credentials{
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
				ExpectError: regexp.MustCompile("Error: ExactlyOne"),
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

func TestUnitAzureCostCredentialsResource(t *testing.T) {

	resourceType := "env0_azure_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	azureCredentialResource := map[string]interface{}{
		"name":            "test",
		"client_id":       "11111",
		"client_secret":   "client-secret1",
		"subscription_id": "subscription-id1",
		"tenant_id":       "tenant-id1",
	}

	updateAzureCredentialResource := map[string]interface{}{
		"name":            "testUpdate",
		"client_id":       "22222",
		"client_secret":   "client-secret2",
		"subscription_id": "subscription-id2",
		"tenant_id":       "tenant-id2",
	}

	azureCredCreatePayload := client.AzureCredentialsCreatePayload{
		Name: azureCredentialResource["name"].(string),
		Value: client.AzureCredentialsValuePayload{
			ClientId:       azureCredentialResource["client_id"].(string),
			ClientSecret:   azureCredentialResource["client_secret"].(string),
			SubscriptionId: azureCredentialResource["subscription_id"].(string),
			TenantId:       azureCredentialResource["tenant_id"].(string),
		},
		Type: client.AzureCostCredentialsType,
	}

	updateAzureCredCreatePayload := client.AzureCredentialsCreatePayload{
		Name: updateAzureCredentialResource["name"].(string),
		Value: client.AzureCredentialsValuePayload{
			ClientId:       updateAzureCredentialResource["client_id"].(string),
			ClientSecret:   updateAzureCredentialResource["client_secret"].(string),
			SubscriptionId: updateAzureCredentialResource["subscription_id"].(string),
			TenantId:       updateAzureCredentialResource["tenant_id"].(string),
		},
		Type: client.AzureCostCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AzureCostCredentialsType),
	}

	updateReturnValues := client.Credentials{
		Id:             "id2",
		Name:           "update",
		OrganizationId: "id",
		Type:           string(client.AzureCostCredentialsType),
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", azureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", azureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", azureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", azureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", azureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", azureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", azureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", azureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", azureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", azureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updateAzureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updateAzureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", updateAzureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", updateAzureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", updateAzureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", updateAzureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AzureCredentialsCreate(azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AzureCredentialsCreate(azureCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().AzureCredentialsCreate(updateAzureCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("validate missing arguments", func(t *testing.T) {
		missingArgumentsTestCases := []resource.TestCase{
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]interface{}{}, "client_id"),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]interface{}{}, "client_secret"),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]interface{}{}, "subscription_id"),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]interface{}{}, "tenant_id"),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]interface{}{}, "name"),
		}
		for _, testCase := range missingArgumentsTestCases {
			tc := testCase
			t.Run("validate specific argument", func(t *testing.T) {
				runUnitTest(t, tc, func(mock *client.MockApiClientInterface) {})
			})

		}
	})

}

func TestUnitGoogleCostCredentialsResource(t *testing.T) {

	resourceType := "env0_google_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	googleCostCredentialResource := map[string]interface{}{
		"name":     "test",
		"table_id": "11111",
		"secret":   "22222",
	}

	updateGoogleCostCredentialResource := map[string]interface{}{
		"name":     "testUpdate",
		"table_id": "333_update",
		"secret":   "444_update",
	}

	googleCostCredCreatePayload := client.GoogleCostCredentialsCreatePayload{
		Name: googleCostCredentialResource["name"].(string),
		Value: client.GoogleCostCredentialsValeuPayload{
			TableId: googleCostCredentialResource["table_id"].(string),
			Secret:  googleCostCredentialResource["secret"].(string),
		},
		Type: client.GoogleCostCredentiassType,
	}

	updateGoogleCostCredCreatePayload := client.GoogleCostCredentialsCreatePayload{
		Name: updateGoogleCostCredentialResource["name"].(string),
		Value: client.GoogleCostCredentialsValeuPayload{
			TableId: updateGoogleCostCredentialResource["table_id"].(string),
			Secret:  updateGoogleCostCredentialResource["secret"].(string),
		},
		Type: client.GoogleCostCredentiassType,
	}

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.GoogleCostCredentiassType),
	}

	updateReturnValues := client.Credentials{
		Id:             "id2",
		Name:           "update",
		OrganizationId: "id",
		Type:           string(client.GoogleCostCredentiassType),
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updateGoogleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updateGoogleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", updateGoogleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", updateGoogleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GoogleCostCredentialsCreate(googleCostCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().GoogleCostCredentialsCreate(googleCostCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().GoogleCostCredentialsCreate(updateGoogleCostCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})
}
