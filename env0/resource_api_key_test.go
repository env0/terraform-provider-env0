package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
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

	apiKeyWithProjectPermissions := client.ApiKey{
		Id:               uuid.NewString(),
		Name:             "name-with-permissions",
		ApiKeyId:         "keyid",
		ApiKeySecret:     "keysecret",
		OrganizationId:   "org",
		OrganizationRole: "User",
	}

	customRoleId := uuid.NewString()
	apiKeyWithCustomRole := client.ApiKey{
		Id:               uuid.NewString(),
		Name:             "name-custom-role",
		ApiKeyId:         "keyid",
		ApiKeySecret:     "keysecret",
		OrganizationId:   "org",
		OrganizationRole: customRoleId,
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
						resource.TestCheckResourceAttr(accessor, "api_key_secret", apiKey.ApiKeySecret),
						resource.TestCheckResourceAttr(accessor, "api_key_id", apiKey.ApiKeyId),
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
						resource.TestCheckResourceAttr(accessor, "api_key_secret", updatedApiKey.ApiKeySecret),
						resource.TestCheckResourceAttr(accessor, "api_key_id", updatedApiKey.ApiKeyId),
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

	t.Run("Omit API Key Secret", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                apiKey.Name,
						"omit_api_key_secret": true,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKey.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", apiKey.OrganizationRole),
						resource.TestCheckResourceAttr(accessor, "api_key_secret", "omitted"),
						resource.TestCheckResourceAttr(accessor, "api_key_id", apiKey.ApiKeyId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name: apiKey.Name,
					Permissions: client.ApiKeyPermissions{
						OrganizationRole: apiKey.OrganizationRole,
					},
				}).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeys().Times(1).Return([]client.ApiKey{apiKey}, nil),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
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

	t.Run("Success - User with Project Permissions", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
                        resource "%s" "%s" {
                            name              = "%s"
                            organization_role = "%s"
                            
                            project_permissions {
                                project_id   = "proj1"
                                project_role = "Viewer"
                            }
                        }
                    `, resourceType, resourceName, apiKeyWithProjectPermissions.Name, apiKeyWithProjectPermissions.OrganizationRole),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKeyWithProjectPermissions.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKeyWithProjectPermissions.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", apiKeyWithProjectPermissions.OrganizationRole),
						resource.TestCheckResourceAttr(accessor, "project_permissions.#", "1"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name: apiKeyWithProjectPermissions.Name,
					Permissions: client.ApiKeyPermissions{
						OrganizationRole: apiKeyWithProjectPermissions.OrganizationRole,
						ProjectPermissions: []client.ProjectPermission{
							{ProjectId: "proj1", ProjectRole: "Viewer"},
						},
					},
				}).Times(1).Return(&apiKeyWithProjectPermissions, nil),
				mock.EXPECT().ApiKeys().Times(1).Return([]client.ApiKey{apiKeyWithProjectPermissions}, nil),
				mock.EXPECT().ApiKeyDelete(apiKeyWithProjectPermissions.Id).Times(1),
			)
		})
	})

	t.Run("Success - Custom Role with Project Permissions", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
                        resource "%s" "%s" {
                            name              = "%s"
                            organization_role = "%s"
                            
                            project_permissions {
                                project_id   = "proj1"
                                project_role = "Planner"
                            }
                        }
                    `, resourceType, resourceName, apiKeyWithCustomRole.Name, customRoleId),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKeyWithCustomRole.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKeyWithCustomRole.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", customRoleId),
						resource.TestCheckResourceAttr(accessor, "project_permissions.#", "1"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name: apiKeyWithCustomRole.Name,
					Permissions: client.ApiKeyPermissions{
						OrganizationRole: customRoleId,
						ProjectPermissions: []client.ProjectPermission{
							{ProjectId: "proj1", ProjectRole: "Planner"},
						},
					},
				}).Times(1).Return(&apiKeyWithCustomRole, nil),
				mock.EXPECT().ApiKeys().Times(1).Return([]client.ApiKey{apiKeyWithCustomRole}, nil),
				mock.EXPECT().ApiKeyDelete(apiKeyWithCustomRole.Id).Times(1),
			)
		})
	})

	t.Run("Failure - Admin with Project Permissions", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
                        resource "%s" "%s" {
                            name = "%s"
                            
                            project_permissions {
                                project_id   = "proj1"
                                project_role = "Viewer"
                            }
                        }
                    `, resourceType, resourceName, apiKey.Name),
					ExpectError: regexp.MustCompile("project_permissions cannot be set when organization_role is Admin"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			// No API calls expected since validation fails
		})
	})
}
