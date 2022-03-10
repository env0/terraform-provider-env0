package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAzureCredDataSource(t *testing.T) {
	azureCred := client.ApiKey{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT",
	}

	credWithInvalidType := client.ApiKey{
		Id:             azureCred.Id,
		Name:           azureCred.Name,
		OrganizationId: azureCred.OrganizationId,
		Type:           "Invalid-type",
	}

	AzureCredFieldsByName := map[string]interface{}{"name": azureCred.Name}
	AzureCredFieldsById := map[string]interface{}{"id": azureCred.Id}

	resourceType := "env0_azure_credentials"
	resourceName := "testdata"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", azureCred.Id),
						resource.TestCheckResourceAttr(accessor, "name", azureCred.Name),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]interface{}, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockGetAzureCredCall := func(returnValue client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(azureCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListAzureCredCall := func(returnValue []client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(AzureCredFieldsById),
			mockGetAzureCredCall(azureCred),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(AzureCredFieldsByName),
			mockListAzureCredCall([]client.ApiKey{azureCred, credWithInvalidType}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one azure-credential exists with the relevant name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(AzureCredFieldsByName, "Found multiple Azure Credentials for name: testdata"),
			mockListAzureCredCall([]client.ApiKey{azureCred, azureCred, azureCred}),
		)
	})

}
