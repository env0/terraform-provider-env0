package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitAzureCostCredentialsDataSource(t *testing.T) {
	azureCred := client.Credentials{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "AZURE_CREDENTIALS",
	}

	credWithInvalidType := client.Credentials{
		Id:             azureCred.Id,
		Name:           azureCred.Name,
		OrganizationId: azureCred.OrganizationId,
		Type:           "Invalid-type",
	}

	credWithDiffName := client.Credentials{
		Id:             "22222",
		Name:           "diff name",
		OrganizationId: azureCred.OrganizationId,
		Type:           "AZURE_CREDENTIALS",
	}

	AzureCredFieldsByName := map[string]interface{}{"name": azureCred.Name}
	AzureCredFieldsById := map[string]interface{}{"id": azureCred.Id}

	resourceType := "env0_azure_cost_credentials"
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

	mockGetAzureCredCall := func(returnValue client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(azureCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListAzureCredCall := func(returnValue []client.Credentials) func(mockFunc *client.MockApiClientInterface) {
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
			mockListAzureCredCall([]client.Credentials{azureCred, credWithInvalidType}),
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
			getErrorTestCase(AzureCredFieldsByName, "Error: Found multiple Cost Credentials for name: testdata"),
			mockListAzureCredCall([]client.Credentials{azureCred, azureCred, azureCred}),
		)
	})

	t.Run("Throw error when by name and no azure-credential exists with the relevant name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(AzureCredFieldsByName, "Error: Could not find Cost Credentials with name: testdata"),
			mockListAzureCredCall([]client.Credentials{credWithDiffName, credWithDiffName}),
		)
	})

}

func TestUnitGoogleCostCredentialsDataSource(t *testing.T) {
	gcpCred := client.Credentials{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "GCP_CREDENTIALS",
	}

	credWithInvalidType := client.Credentials{
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

	mockGetGcpCredCall := func(returnValue client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(gcpCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListGcpCredCall := func(returnValue []client.Credentials) func(mockFunc *client.MockApiClientInterface) {
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
			mockListGcpCredCall([]client.Credentials{gcpCred, credWithInvalidType}),
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
			getErrorTestCase(GcpCredFieldsByName, "Error: Found multiple Cost Credentials for name: testdata"),
			mockListGcpCredCall([]client.Credentials{gcpCred, gcpCred, gcpCred}),
		)
	})

}

func TestUnitAwsCostCredentialsData(t *testing.T) {
	awsCred := client.Credentials{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	credWithInvalidType := client.Credentials{
		Id:             awsCred.Id,
		Name:           awsCred.Name,
		OrganizationId: awsCred.OrganizationId,
		Type:           "Invalid-type",
	}

	otherAwsCred := client.Credentials{
		Id:             "22222",
		Name:           "notTestdata",
		OrganizationId: "OtherId",
		Type:           "AWS_ACCESS_KEYS_FOR_DEPLOYMENT",
	}

	AwsCredFieldsByName := map[string]interface{}{"name": awsCred.Name}
	AwsCredFieldsById := map[string]interface{}{"id": awsCred.Id}

	resourceType := "env0_aws_cost_credentials"
	resourceName := "testdata"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", awsCred.Id),
						resource.TestCheckResourceAttr(accessor, "name", awsCred.Name),
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

	mockGetAwsCredCall := func(returnValue client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(awsCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListAwsCredCall := func(returnValue []client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(AwsCredFieldsById),
			mockGetAwsCredCall(awsCred),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(AwsCredFieldsByName),
			mockListAwsCredCall([]client.Credentials{awsCred, otherAwsCred, credWithInvalidType}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one aws-credential exists with the relevant name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(AwsCredFieldsByName, "Found multiple Cost Credentials for name: testdata"),
			mockListAwsCredCall([]client.Credentials{awsCred, awsCred, awsCred}),
		)
	})

	t.Run("Throw error when by name and no aws-credential found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(AwsCredFieldsByName, "Could not find Cost Credentials with name: testdata"),
			mockListAwsCredCall([]client.Credentials{otherAwsCred, credWithInvalidType}),
		)
	})

}
