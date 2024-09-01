package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitAzureCredentialsResource(t *testing.T) {
	resourceType := "env0_azure_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
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
		Type: client.AzureServicePrincipalCredentialsType,
	}

	updateAzureCredCreatePayload := client.AzureCredentialsCreatePayload{
		Name: updateAzureCredentialResource["name"].(string),
		Value: client.AzureCredentialsValuePayload{
			ClientId:       updateAzureCredentialResource["client_id"].(string),
			ClientSecret:   updateAzureCredentialResource["client_secret"].(string),
			SubscriptionId: updateAzureCredentialResource["subscription_id"].(string),
			TenantId:       updateAzureCredentialResource["tenant_id"].(string),
		},
		Type: client.AzureServicePrincipalCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de3ff",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AzureServicePrincipalCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "GCP_....",
	}

	updateReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de3aa",
		Name:           "testUpdate",
		OrganizationId: "id",
		Type:           string(client.AzureServicePrincipalCredentialsType),
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
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().CredentialsCreate(&updateAzureCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("validate missing arguments", func(t *testing.T) {
		arguments := []string{
			"client_id",
			"client_secret",
			"subscription_id",
			"tenant_id",
			"name",
		}

		for _, argument := range arguments {
			tc := missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{}, argument)
			t.Run("validate missing arrguments "+argument, func(t *testing.T) {
				runUnitTest(t, tc, func(mock *client.MockApiClientInterface) {})
			})
		}
	})

	t.Run("Azure credentials removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, azureCredentialResource)

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
				mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     azureCredentialResource["name"].(string),
					ImportStateVerify: false,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
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
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
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
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{})
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     azureCredentialResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", azureCredentialResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})
}
