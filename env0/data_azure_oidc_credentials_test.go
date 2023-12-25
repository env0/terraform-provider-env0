package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAzureOidcCredentialDataSource(t *testing.T) {
	credentials := client.Credentials{
		Id:   "id0",
		Name: "name0",
		Type: string(client.AzureOidcCredentialsType),
	}

	credentialsOther1 := client.Credentials{
		Id:   "id1",
		Name: "name1",
		Type: string(client.AzureOidcCredentialsType),
	}

	credentialsOther2 := client.Credentials{
		Id:   "id2",
		Name: "name2",
		Type: string(client.AzureServicePrincipalCredentialsType),
	}

	byName := map[string]interface{}{"name": credentials.Name}
	byId := map[string]interface{}{"id": credentials.Id}

	resourceType := "env0_azure_oidc_credentials"
	resourceName := "test_azure_oidc_credentials"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", credentials.Id),
						resource.TestCheckResourceAttr(accessor, "name", credentials.Name),
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

	mockGetCredentials := func(returnValue client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentials(credentials.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListCredentials := func(returnValue []client.Credentials) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CloudCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("by id", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(byId),
			mockGetCredentials(credentials),
		)
	})

	t.Run("by name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(byName),
			mockListCredentials([]client.Credentials{credentials, credentialsOther1, credentialsOther2}),
		)
	})

	t.Run("throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("throw error when by name and more than one is returned", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(byName, "found multiple credentials"),
			mockListCredentials([]client.Credentials{credentials, credentialsOther1, credentialsOther2, credentials}),
		)
	})

	t.Run("Throw error when by name and not found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(byName, "not found"),
			mockListCredentials([]client.Credentials{credentialsOther1, credentialsOther2}),
		)
	})

	t.Run("Throw error when by id and not found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(byId, fmt.Sprintf("id %s not found", credentials.Id)),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().CloudCredentials(credentials.Id).AnyTimes().Return(client.Credentials{}, http.NewMockFailedResponseError(404))
			},
		)
	})
}
