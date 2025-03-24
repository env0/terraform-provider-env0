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

func TestUnitVaultOidcCredentialsResource(t *testing.T) {
	resourceType := "env0_vault_oidc_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	vaultCredentialsResource := map[string]any{
		"name":                  "test",
		"address":               "http://fake1.com:80",
		"version":               "version1",
		"role_name":             "rolename1",
		"jwt_auth_backend_path": "path1",
		"namespace":             "namespace1",
	}

	updatedVaultCredentialsResource := map[string]any{
		"name":                  "test",
		"address":               "http://fake2.com:80",
		"version":               "version2",
		"role_name":             "rolename2",
		"jwt_auth_backend_path": "path2",
	}

	createPayload := client.VaultCredentialsCreatePayload{
		Name: vaultCredentialsResource["name"].(string),
		Value: client.VaultCredentialsValuePayload{
			Address:            vaultCredentialsResource["address"].(string),
			Version:            vaultCredentialsResource["version"].(string),
			RoleName:           vaultCredentialsResource["role_name"].(string),
			JwtAuthBackendPath: vaultCredentialsResource["jwt_auth_backend_path"].(string),
			Namespace:          vaultCredentialsResource["namespace"].(string),
		},
		Type: client.VaultOidcCredentialsType,
	}

	updatePayload := client.VaultCredentialsCreatePayload{
		Value: client.VaultCredentialsValuePayload{
			Address:            updatedVaultCredentialsResource["address"].(string),
			Version:            updatedVaultCredentialsResource["version"].(string),
			RoleName:           updatedVaultCredentialsResource["role_name"].(string),
			JwtAuthBackendPath: updatedVaultCredentialsResource["jwt_auth_backend_path"].(string),
		},
		Type: client.VaultOidcCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.VaultOidcCredentialsType),
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
		Type:           string(client.VaultOidcCredentialsType),
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", vaultCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "address", vaultCredentialsResource["address"].(string)),
					resource.TestCheckResourceAttr(accessor, "version", vaultCredentialsResource["version"].(string)),
					resource.TestCheckResourceAttr(accessor, "role_name", vaultCredentialsResource["role_name"].(string)),
					resource.TestCheckResourceAttr(accessor, "jwt_auth_backend_path", vaultCredentialsResource["jwt_auth_backend_path"].(string)),
					resource.TestCheckResourceAttr(accessor, "namespace", vaultCredentialsResource["namespace"].(string)),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedVaultCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", updatedVaultCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "address", updatedVaultCredentialsResource["address"].(string)),
					resource.TestCheckResourceAttr(accessor, "version", updatedVaultCredentialsResource["version"].(string)),
					resource.TestCheckResourceAttr(accessor, "role_name", updatedVaultCredentialsResource["role_name"].(string)),
					resource.TestCheckResourceAttr(accessor, "jwt_auth_backend_path", updatedVaultCredentialsResource["jwt_auth_backend_path"].(string)),
					resource.TestCheckResourceAttr(accessor, "namespace", ""),
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
		stepConfig := resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource)

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
					Config: resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           vaultCredentialsResource["name"].(string),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"address", "version", "role_name", "jwt_auth_backend_path", "namespace"},
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
					Config: resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           returnValues.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"address", "version", "role_name", "jwt_auth_backend_path", "namespace"},
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
					Config: resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource),
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
					Config: resourceConfigCreate(resourceType, resourceName, vaultCredentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     vaultCredentialsResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", vaultCredentialsResource["name"].(string))),
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
