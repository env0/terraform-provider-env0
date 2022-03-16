package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitGoogleCostCredentialsDataSource(t *testing.T) {
	gcpCred := client.ApiKey{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "GCP_CREDENTIALS",
	}

	credWithInvalidType := client.ApiKey{
		Id:             gcpCred.Id,
		Name:           gcpCred.Name,
		OrganizationId: gcpCred.OrganizationId,
		Type:           "Invalid-type",
	}

	GcpCredFieldsByName := map[string]interface{}{"name": gcpCred.Name}
	GcpCredFieldsById := map[string]interface{}{"id": gcpCred.Id}

	resourceType := "env0_google_cost_credentials"
	resourceName := "testdata"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", gcpCred.Id),
						resource.TestCheckResourceAttr(accessor, "name", gcpCred.Name),
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

	mockGetGcpCredCall := func(returnValue client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(gcpCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListGcpCredCall := func(returnValue []client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(GcpCredFieldsById),
			mockGetGcpCredCall(gcpCred),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(GcpCredFieldsByName),
			mockListGcpCredCall([]client.ApiKey{gcpCred, credWithInvalidType}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one gcp-credential exists with the relevant name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(GcpCredFieldsByName, "Found multiple Google cost Credentials for name: testdata"),
			mockListGcpCredCall([]client.ApiKey{gcpCred, gcpCred, gcpCred}),
		)
	})

}
