package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

	returnValues := client.ApiKey{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.GoogleCostCredentiassType),
	}

	updateReturnValues := client.ApiKey{
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

	t.Run("validate missing arguments", func(t *testing.T) {
		missingArgumentsTestCases := []resource.TestCase{
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{
				"name": "update",
			}, "secret"),
			missingArgumentTestCase(resourceType, resourceName, map[string]interface{}{
				"secret": "update",
			}, "name"),
		}
		for _, testCase := range missingArgumentsTestCases {
			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			})
		}
	})

}
