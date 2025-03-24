package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type TestCredentialsCase struct {
	credentialsType string
	cloudType       string
}

func TestCredentialsDataSource(t *testing.T) {
	testCases := []TestCredentialsCase{
		{
			credentialsType: "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
			cloudType:       "aws",
		},
		{
			credentialsType: "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT",
			cloudType:       "gcp",
		},
		{
			credentialsType: "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT",
			cloudType:       "azure",
		},
		{
			credentialsType: "GCP_CREDENTIALS",
			cloudType:       "google_cost",
		},
		{
			credentialsType: "AZURE_CREDENTIALS",
			cloudType:       "azure_cost",
		},
		{
			credentialsType: "AWS_ASSUMED_ROLE",
			cloudType:       "aws_cost",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		name := "name" + tc.cloudType
		byName := map[string]any{"name": name}
		id := "id" + tc.cloudType
		byId := map[string]any{"id": id}

		resourceType := "env0_" + tc.cloudType + "_credentials"
		resourceName := "testdata"
		accessor := dataSourceAccessor(resourceType, resourceName)

		credentials := client.Credentials{
			Id:             id,
			Name:           name,
			OrganizationId: "id",
			Type:           tc.credentialsType,
		}

		otherCredentials := client.Credentials{
			Id:             "22222",
			Name:           "notname",
			OrganizationId: "id",
			Type:           tc.credentialsType,
		}

		invalidCredentials := client.Credentials{
			Id:             credentials.Id,
			Name:           credentials.Name,
			OrganizationId: credentials.OrganizationId,
			Type:           "Invalid-type",
		}

		getValidTestCase := func(input map[string]any) resource.TestCase {
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

		getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
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

		t.Run(tc.cloudType+" by ID", func(t *testing.T) {
			runUnitTest(t,
				getValidTestCase(byId),
				mockGetCredentials(credentials),
			)
		})

		t.Run(tc.cloudType+" by Name", func(t *testing.T) {
			runUnitTest(t,
				getValidTestCase(byName),
				mockListCredentials([]client.Credentials{credentials, otherCredentials, invalidCredentials}),
			)
		})

		t.Run(tc.cloudType+" throw error when no name or id is supplied", func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(map[string]any{}, "one of `id,name` must be specified"),
				func(mock *client.MockApiClientInterface) {},
			)
		})

		t.Run(tc.cloudType+" throw error when by name and more than one aws-credential exists with the relevant name", func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(byName, "found multiple credentials with name: name"+tc.cloudType),
				mockListCredentials([]client.Credentials{credentials, credentials}),
			)
		})

		t.Run(tc.cloudType+" throw error when by name and no aws-credential found with that name", func(t *testing.T) {
			runUnitTest(t,
				getErrorTestCase(byName, "credentials with name name"+tc.cloudType+" not found"),
				mockListCredentials([]client.Credentials{otherCredentials, invalidCredentials}),
			)
		})
	}
}
