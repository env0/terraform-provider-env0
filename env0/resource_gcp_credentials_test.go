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

func TestUnitGcpCredentialsResource(t *testing.T) {
	resourceType := "env0_gcp_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	gcpCredentialResource := map[string]any{
		"name":                "test",
		"project_id":          "11111",
		"service_account_key": "22222",
	}

	updateGcpCredentialResource := map[string]any{
		"name":                "update",
		"project_id":          "333_update",
		"service_account_key": "444_update",
	}

	gcpCredCreatePayload := client.GcpCredentialsCreatePayload{
		Name: gcpCredentialResource["name"].(string),
		Value: client.GcpCredentialsValuePayload{
			ProjectId:         gcpCredentialResource["project_id"].(string),
			ServiceAccountKey: gcpCredentialResource["service_account_key"].(string),
		},
		Type: client.GcpServiceAccountCredentialsType,
	}

	updateGcpCredCreatePayload := client.GcpCredentialsCreatePayload{
		Name: updateGcpCredentialResource["name"].(string),
		Value: client.GcpCredentialsValuePayload{
			ProjectId:         updateGcpCredentialResource["project_id"].(string),
			ServiceAccountKey: updateGcpCredentialResource["service_account_key"].(string),
		},
		Type: client.GcpServiceAccountCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de300",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.GcpServiceAccountCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AZURE_....",
	}

	updateReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de311",
		Name:           "update",
		OrganizationId: "id",
		Type:           string(client.GcpServiceAccountCredentialsType),
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", gcpCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "project_id", gcpCredentialResource["project_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "service_account_key", gcpCredentialResource["service_account_key"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", gcpCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "project_id", gcpCredentialResource["project_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "service_account_key", gcpCredentialResource["service_account_key"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updateGcpCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updateGcpCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "project_id", updateGcpCredentialResource["project_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "service_account_key", updateGcpCredentialResource["service_account_key"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().CredentialsCreate(&updateGcpCredCreatePayload).Times(1).Return(updateReturnValues, nil),
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
			missingArgumentTestCase(resourceType, resourceName, map[string]any{
				"name": "update",
			}, "service_account_key"),
			missingArgumentTestCase(resourceType, resourceName, map[string]any{
				"service_account_key": "update",
			}, "name"),
		}

		for _, testCase := range missingArgumentsTestCases {
			tc := testCase

			t.Run("validate missing arguments", func(t *testing.T) {
				runUnitTest(t, tc, func(mock *client.MockApiClientInterface) {})
			})
		}
	})

	t.Run("GCP credentials removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, gcpCredentialResource)

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
				mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gcpCredentialResource["name"].(string),
					ImportStateVerify: false,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     returnValues.Id,
					ImportStateVerify: false,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
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
			mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{})
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, gcpCredentialResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gcpCredentialResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", gcpCredentialResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&gcpCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})
}
