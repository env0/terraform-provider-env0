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
	t.Parallel()

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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": apiKey.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKey.Name),
						resource.TestCheckResourceAttr(accessor, "organization_role", apiKey.OrganizationRole),
						resource.TestCheckResourceAttr(accessor, "api_key_secret", apiKey.ApiKeySecret),
						resource.TestCheckResourceAttr(accessor, "api_key_id", apiKey.ApiKeyId),
						resource.TestCheckResourceAttr(accessor, "ip_allowlist.#", "0"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
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

	prefix := func(bits int) *int { return &bits }

	allowlistEntries := []string{"198.51.100.7", "198.51.100.0/28", "198.51.100.40-198.51.100.60", "2001:db8::/32"}
	allowlistRules := []client.IpAllowlistRule{
		{Type: "address", Address: "198.51.100.7", Family: "ipv4"},
		{Type: "subnet", Network: "198.51.100.0", Prefix: prefix(28), Family: "ipv4"},
		{Type: "range", Start: "198.51.100.40", End: "198.51.100.60", Family: "ipv4"},
		{Type: "subnet", Network: "2001:db8::", Prefix: prefix(32), Family: "ipv6"},
	}
	updatedAllowlistRules := []client.IpAllowlistRule{
		{Type: "subnet", Network: "192.0.2.0", Prefix: prefix(24), Family: "ipv4"},
	}

	// Reads return whatever the last PUT stored, the way the server does.
	mockStoredAllowlist := func(mock *client.MockApiClientInterface, key client.ApiKey) func(rules []client.IpAllowlistRule) *client.ApiKey {
		stored := key
		stored.ApiKeySecret = ""

		mock.EXPECT().ApiKeys().AnyTimes().DoAndReturn(func() ([]client.ApiKey, error) {
			return []client.ApiKey{stored}, nil
		})

		return func(rules []client.IpAllowlistRule) *client.ApiKey {
			stored.IpAllowlist = rules

			return &stored
		}
	}

	t.Run("ip_allowlist - create, update in place and clear", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": allowlistEntries,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "api_key_secret", apiKey.ApiKeySecret),
						resource.TestCheckResourceAttr(accessor, "ip_allowlist.#", "4"),
						resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "198.51.100.7"),
						resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "198.51.100.0/28"),
						resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "198.51.100.40-198.51.100.60"),
						resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "2001:db8::/32"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"192.0.2.0/24"},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "api_key_secret", apiKey.ApiKeySecret),
						resource.TestCheckResourceAttr(accessor, "ip_allowlist.#", "1"),
						resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "192.0.2.0/24"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": apiKey.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "ip_allowlist.#", "0"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			store := mockStoredAllowlist(mock, apiKey)
			storeRules := func(_ string, rules []client.IpAllowlistRule) (*client.ApiKey, error) { return store(rules), nil }

			// No create or delete between the PUTs: the key is updated in place, never replaced.
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(client.ApiKeyCreatePayload{
					Name:        apiKey.Name,
					Permissions: client.ApiKeyPermissions{OrganizationRole: "Admin"},
				}).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, gomock.InAnyOrder(allowlistRules)).Times(1).DoAndReturn(storeRules),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, updatedAllowlistRules).Times(1).DoAndReturn(storeRules),
				// Non-nil, so the body is "ipAllowlist": [] and not null.
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, []client.IpAllowlistRule{}).Times(1).DoAndReturn(storeRules),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
			)
		})
	})

	t.Run("ip_allowlist - a list set outside Terraform is cleared", func(t *testing.T) {
		var store func(rules []client.IpAllowlistRule) *client.ApiKey

		config := resourceConfigCreate(resourceType, resourceName, map[string]any{
			"name": apiKey.Name,
		})

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{Config: config},
				{
					// Someone sets a list in the env0 UI.
					PreConfig: func() { store(updatedAllowlistRules) },
					Config:    config,
					Check:     resource.TestCheckResourceAttr(accessor, "ip_allowlist.#", "0"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			store = mockStoredAllowlist(mock, apiKey)

			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(gomock.Any()).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, []client.IpAllowlistRule{}).Times(1).DoAndReturn(
					func(_ string, rules []client.IpAllowlistRule) (*client.ApiKey, error) { return store(rules), nil },
				),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
			)
		})
	})

	t.Run("ip_allowlist - import reads the list", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"192.0.2.0/24"},
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           apiKey.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"api_key_secret"},
					Check:                   resource.TestCheckTypeSetElemAttr(accessor, "ip_allowlist.*", "192.0.2.0/24"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			store := mockStoredAllowlist(mock, apiKey)

			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(gomock.Any()).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, updatedAllowlistRules).Times(1).DoAndReturn(
					func(_ string, rules []client.IpAllowlistRule) (*client.ApiKey, error) { return store(rules), nil },
				),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
			)
		})
	})

	t.Run("ip_allowlist - a failed PUT on create deletes the new key", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"192.0.2.0/24"},
					}),
					ExpectError: regexp.MustCompile("could not set the api key ip_allowlist: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(gomock.Any()).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, updatedAllowlistRules).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("ip_allowlist - a failed PUT and a failed delete name the key", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"192.0.2.0/24"},
					}),
					ExpectError: regexp.MustCompile("Deleting the new api key " + apiKey.Id + " also failed"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(gomock.Any()).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, updatedAllowlistRules).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1).Return(errors.New("delete error")),
			)
		})
	})

	t.Run("ip_allowlist - a failed update returns the error", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"192.0.2.0/24"},
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":         apiKey.Name,
						"ip_allowlist": []string{"198.51.100.7"},
					}),
					ExpectError: regexp.MustCompile("could not update api key ip_allowlist: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			store := mockStoredAllowlist(mock, apiKey)

			gomock.InOrder(
				mock.EXPECT().ApiKeyCreate(gomock.Any()).Times(1).Return(&apiKey, nil),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, updatedAllowlistRules).Times(1).DoAndReturn(
					func(_ string, rules []client.IpAllowlistRule) (*client.ApiKey, error) { return store(rules), nil },
				),
				mock.EXPECT().ApiKeyUpdateIpAllowlist(apiKey.Id, []client.IpAllowlistRule{
					{Type: "address", Address: "198.51.100.7", Family: "ipv4"},
				}).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().ApiKeyDelete(apiKey.Id).Times(1),
			)
		})
	})

	t.Run("ip_allowlist - invalid entries", func(t *testing.T) {
		tooMany := make([]string, 51)
		for i := range tooMany {
			tooMany[i] = fmt.Sprintf("198.51.100.%d", i)
		}

		for name, tc := range map[string]struct {
			entries []string
			err     string
		}{
			"bad IP":             {[]string{"not-an-ip"}, "invalid ip_allowlist entry 'not-an-ip'"},
			"leading-zero octet": {[]string{"198.051.100.7"}, "invalid ip_allowlist entry '198.051.100.7'"},
			"reversed range":     {[]string{"198.51.100.60-198.51.100.40"}, "the range start is after its end"},
			"mixed-family range": {[]string{"198.51.100.7-2001:db8::1"}, "both ends of a range must be in the same IP family"},
			"prefix 33":          {[]string{"198.51.100.0/33"}, "invalid ip_allowlist entry '198.51.100.0/33'"},
			"zone":               {[]string{"fe80::1%eth0"}, "IPv6 zones are not supported"},
			"whitespace":         {[]string{"198.51.100.7 "}, "invalid ip_allowlist entry '198.51.100.7 '"},
			"51 entries":         {tooMany, "50 item maximum"},
		} {
			t.Run(name, func(t *testing.T) {
				testCase := resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
								"name":         apiKey.Name,
								"ip_allowlist": tc.entries,
							}),
							ExpectError: regexp.MustCompile(regexp.QuoteMeta(tc.err)),
						},
					},
				}

				runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
			})
		}
	})
}
