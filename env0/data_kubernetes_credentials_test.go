package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestKubernetesCredentialsDataSource(t *testing.T) {
	tests := [][]string{
		{"env0_aws_eks_credentials", string(client.AwsEksCredentialsType)},
		{"env0_azure_aks_credentials", string(client.AzureAksCredentialsType)},
		{"env0_gcp_gke_credentials", string(client.GcpGkeCredentialsType)},
		{"env0_kubeconfig_credentials", string(client.KubeconfigCredentialsType)},
	}

	for _, test := range tests {
		credentials := client.Credentials{
			Id:   "id0",
			Name: "name0",
			Type: test[1],
		}

		credentialsOther1 := client.Credentials{
			Id:   "id1",
			Name: "name1",
			Type: test[1],
		}

		credentialsOther2 := client.Credentials{
			Id:   "id2",
			Name: "name2",
			Type: test[1],
		}

		byName := map[string]interface{}{"name": credentials.Name}
		byId := map[string]interface{}{"id": credentials.Id}

		resourceType := test[0]
		resourceName := "test"
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

		t.Run("by name - "+test[0], func(t *testing.T) {
			runUnitTest(t,
				getValidTestCase(byName),
				mockListCredentials([]client.Credentials{credentials, credentialsOther1, credentialsOther2}),
			)
		})

		t.Run("throw error when no name or id is supplied - "+test[0], func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
				func(mock *client.MockApiClientInterface) {},
			)
		})

		t.Run("throw error when by name and more than one is returned - "+test[0], func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(byName, "found multiple credentials"),
				mockListCredentials([]client.Credentials{credentials, credentialsOther1, credentialsOther2, credentials}),
			)
		})

		t.Run("Throw error when by name and not found - "+test[0], func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(byName, "not found"),
				mockListCredentials([]client.Credentials{credentialsOther1, credentialsOther2}),
			)
		})

		t.Run("Throw error when by id and not found - "+test[0], func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(byId, fmt.Sprintf("id %s not found", credentials.Id)),
				func(mock *client.MockApiClientInterface) {
					mock.EXPECT().CloudCredentials(credentials.Id).AnyTimes().Return(client.Credentials{}, http.NewMockFailedResponseError(404))
				},
			)
		})
	}
}
