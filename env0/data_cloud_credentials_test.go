package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCloudCredentialsDataSource(t *testing.T) {
	credentials1 := client.Credentials{
		Id:   "id0",
		Name: "name0",
		Type: "AWS_ASSUMED_ROLE",
	}

	credentials2 := client.Credentials{
		Id:   "id1",
		Name: "name1",
		Type: "AZURE_CREDENTIALS",
	}

	resourceType := "env0_cloud_credentials"
	resourceName := "test_cloud_credentials"
	accessor := dataSourceAccessor(resourceType, resourceName)

	mockCloudCredentials := func(returnValue []client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "names.0", credentials1.Name),
							resource.TestCheckResourceAttr(accessor, "names.1", credentials2.Name),
						),
					},
				},
			},
			mockCloudCredentials([]client.Credentials{credentials1, credentials2}),
		)
	})

	t.Run("Success With Filter", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"credential_type": credentials2.Type}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "names.0", credentials2.Name),
						),
					},
				},
			},
			mockCloudCredentials([]client.Credentials{credentials1, credentials2}),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().CloudCredentialsList().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
