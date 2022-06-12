package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitApiKeyResource(t *testing.T) {
	resourceType := "env0_api_key"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	apiKey := client.ApiKey{
		Id:               uuid.NewString(),
		Name:             "name",
		ApiKeyId:         "keyid",
		ApiKeySecret:     "keysecret",
		OrganizationId:   "org",
		OrganizationRole: "Admin",
	}

	apiKeyUser := client.ApiKey{
		Id:               uuid.NewString(),
		Name:             "name-user",
		ApiKeyId:         "keyid",
		ApiKeySecret:     "keysecret",
		OrganizationId:   "org",
		OrganizationRole: "User",
	}

	updatedApiKey := client.ApiKey{
		Id:               "id2",
		Name:             "name2",
		ApiKeyId:         "keyid2",
		ApiKeySecret:     "keysecret2",
		OrganizationId:   "org",
		OrganizationRole: "Admin",
	}

	t.Run("Success - Admin", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": apiKey.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKey.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", apiKey.OrganizationRole),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": updatedApiKey.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedApiKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedApiKey.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", updatedApiKey.OrganizationRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name:        apiKey.Name,
					Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
				}).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeys().Times(2).Return([]client.ApiKey{apiKey}, nil),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name:        updatedApiKey.Name,
					Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
				}).Times(1).Return(&updatedApiKey, nil),
				mock.EXPECT().ApiKeys().Times(1).Return([]client.ApiKey{updatedApiKey}, nil),
				mock.EXPECT().ApiKeyDelete(updatedApiKey.Id).Times(1),
			)
		})
	})

	t.Run("Success - User", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":              apiKeyUser.Name,
						"organization_role": apiKeyUser.OrganizationRole,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKeyUser.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKeyUser.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", apiKeyUser.OrganizationRole),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name: apiKeyUser.Name,
					Permissions: client.ApiKeyPermissions{
						OrganizationRole: apiKeyUser.OrganizationRole,
					},
				}).Times(1).Return(&apiKeyUser, nil),
				mock.EXPECT().ApiKeys().Times(1).Return([]client.ApiKey{apiKeyUser}, nil),
				mock.EXPECT().ApiKeyDelete(apiKeyUser.Id).Times(1),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": apiKey.Name,
					}),
					ExpectError: regexp.MustCompile("could not create api key: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
				Name:        apiKey.Name,
				Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": apiKey.Name,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     apiKey.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
				Name:        apiKey.Name,
				Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
			}).Times(1).Return(&apiKey, nil)
			mock.EXPECT().ApiKeys().Times(3).Return([]client.ApiKey{apiKey}, nil)
			mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": apiKey.Name,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     apiKey.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
				Name:        apiKey.Name,
				Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
			}).Times(1).Return(&apiKey, nil)
			mock.EXPECT().ApiKeys().Times(3).Return([]client.ApiKey{apiKey}, nil)
			mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1)
		})
	})

	t.Run("Import By Id Not Found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name": updatedApiKey.Name,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     apiKey.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("api key with id %s not found", apiKey.Id)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
				Name:        updatedApiKey.Name,
				Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
			}).Times(1).Return(&updatedApiKey, nil)
			mock.EXPECT().ApiKeys().Times(2).Return([]client.ApiKey{updatedApiKey}, nil)
			mock.EXPECT().ApiKeyDelete(updatedApiKey.Id).Times(1)
		})
	})

}
