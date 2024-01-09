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

func TestUnitGcpOidcCredentialsResource(t *testing.T) {
	resourceType := "env0_gcp_oidc_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	gcpCredentialsResource := map[string]interface{}{
		"name":                                  "test",
		"credential_configuration_file_content": "content1",
	}

	updatedGcpCredentialsResource := map[string]interface{}{
		"name":                                  "test",
		"credential_configuration_file_content": "content2",
	}

	createPayload := client.GcpCredentialsCreatePayload{
		Name: gcpCredentialsResource["name"].(string),
		Value: client.GcpCredentialsValuePayload{
			CredentialConfigurationFileContent: gcpCredentialsResource["credential_configuration_file_content"].(string),
		},
		Type: client.GcpOidcCredentialsType,
	}

	updatePayload := client.GcpCredentialsCreatePayload{
		Value: client.GcpCredentialsValuePayload{
			CredentialConfigurationFileContent: updatedGcpCredentialsResource["credential_configuration_file_content"].(string),
		},
		Type: client.GcpOidcCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.GcpOidcCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_....",
	}

	updateReturnValues := client.Credentials{
		Id:             returnValues.Id,
		Name:           returnValues.Name,
		OrganizationId: "id",
		Type:           string(client.GcpOidcCredentialsType),
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", gcpCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "credential_configuration_file_content", gcpCredentialsResource["credential_configuration_file_content"].(string)),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedGcpCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", updatedGcpCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "credential_configuration_file_content", updatedGcpCredentialsResource["credential_configuration_file_content"].(string)),
				),
			},
		},
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, testCaseForCreateAndUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updatePayload).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource)

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
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           gcpCredentialsResource["name"].(string),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"credential_configuration_file_content"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           returnValues.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"credential_configuration_file_content"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource),
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
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{}),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gcpCredentialsResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", gcpCredentialsResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})
}
