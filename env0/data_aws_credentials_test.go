package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAwsCredDataSource(t *testing.T) {
	awsCred := client.ApiKey{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	credWithInvalidType := client.ApiKey{
		Id:             awsCred.Id,
		Name:           awsCred.Name,
		OrganizationId: awsCred.OrganizationId,
		Type:           "Invalid-type",
	}

	otherAwsCred := client.ApiKey{
		Id:             "22222",
		Name:           "notTestdata",
		OrganizationId: "OtherId",
		Type:           "AWS_ACCESS_KEYS_FOR_DEPLOYMENT",
	}

	AwsCredFieldsByName := map[string]interface{}{"name": awsCred.Name}
	AwsCredFieldsById := map[string]interface{}{"id": awsCred.Id}

	resourceType := "env0_aws_credentials"
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

	mockGetAwsCredCall := func(returnValue client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(awsCred.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListAwsCredCall := func(returnValue []client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
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
			mockListAwsCredCall([]client.ApiKey{awsCred, otherAwsCred, credWithInvalidType}),
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
			getErrorTestCase(AwsCredFieldsByName, "Found multiple AWS Credentials for name: testdata"),
			mockListAwsCredCall([]client.ApiKey{awsCred, awsCred, awsCred}),
		)
	})

	t.Run("Throw error when by name and no aws-credential found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(AwsCredFieldsByName, "Could not find AWS Credentials with name: testdata"),
			mockListAwsCredCall([]client.ApiKey{otherAwsCred, credWithInvalidType}),
		)
	})

}
