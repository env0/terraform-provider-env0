package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

	returnValues := client.ApiKey{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AzureCostCredentialsType),
	}

	updateReturnValues := client.ApiKey{
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
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, "client_id"),
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, "client_secret"),
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, "subscription_id"),
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, "tenant_id"),
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, "name"),
		}
		for _, testCase := range missingArgumentsTestCases {
			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			})
		}
	})

}
