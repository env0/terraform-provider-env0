package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitModulenResource(t *testing.T) {
	resourceType := "env0_module"
	resourceName := "test"
	//resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	author := client.User{
		Name: "user0",
	}

	module := client.Module{
		ModuleName:     "name1",
		ModuleProvider: "provider1",
		Repository:     "repository1",
		Description:    "description1",
		TokenId:        "tokenid",
		TokenName:      "tokenname",
		IsGitlab:       true,
		OrganizationId: "org1",
		Author:         author,
		AuthorId:       "author1",
		Id:             "id1",
	}

	updatedModule := client.Module{
		ModuleName:         "name2",
		ModuleProvider:     "provider1",
		Repository:         "repository1",
		Description:        "description1",
		BitbucketClientKey: stringPtr("1234"),
		OrganizationId:     "org1",
		Author:             author,
		AuthorId:           "author1",
		Id:                 "id1",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"module_name":     module.ModuleName,
						"module_provider": module.ModuleProvider,
						"repository":      module.Repository,
						"description":     module.Description,
						"token_id":        module.TokenId,
						"token_name":      module.TokenName,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", module.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", module.ModuleName),
						resource.TestCheckResourceAttr(accessor, "module_provider", module.ModuleProvider),
						resource.TestCheckResourceAttr(accessor, "repository", module.Repository),
						resource.TestCheckResourceAttr(accessor, "description", module.Description),
						resource.TestCheckResourceAttr(accessor, "token_id", module.TokenId),
						resource.TestCheckResourceAttr(accessor, "token_name", module.TokenName),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"module_name":          updatedModule.ModuleName,
						"module_provider":      updatedModule.ModuleProvider,
						"repository":           updatedModule.Repository,
						"description":          updatedModule.Description,
						"bitbucket_client_key": *updatedModule.BitbucketClientKey,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedModule.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", updatedModule.ModuleName),
						resource.TestCheckResourceAttr(accessor, "module_provider", updatedModule.ModuleProvider),
						resource.TestCheckResourceAttr(accessor, "repository", updatedModule.Repository),
						resource.TestCheckResourceAttr(accessor, "description", updatedModule.Description),
						resource.TestCheckResourceAttr(accessor, "token_id", ""),
						resource.TestCheckResourceAttr(accessor, "token_name", ""),
						resource.TestCheckResourceAttr(accessor, "bitbucket_client_key", *updatedModule.BitbucketClientKey),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:     module.ModuleName,
				ModuleProvider: module.ModuleProvider,
				Repository:     module.Repository,
				Description:    module.Description,
				TokenId:        module.TokenId,
				TokenName:      module.TokenName,
				IsGitlab:       boolPtr(true),
			}).Times(1).Return(&module, nil)

			mock.EXPECT().ModuleUpdate(updatedModule.Id, client.ModuleUpdatePayload{
				ModuleName:           updatedModule.ModuleName,
				ModuleProvider:       updatedModule.ModuleProvider,
				Repository:           updatedModule.Repository,
				Description:          updatedModule.Description,
				TokenId:              "",
				TokenName:            "",
				IsGitlab:             false,
				GithubInstallationId: nil,
				BitbucketClientKey:   *updatedModule.BitbucketClientKey,
			}).Times(1).Return(&updatedModule, nil)

			gomock.InOrder(
				mock.EXPECT().Module(module.Id).Times(2).Return(&module, nil),
				mock.EXPECT().Module(module.Id).Times(1).Return(&updatedModule, nil),
			)

			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	/*
		t.Run("Create Failure - Invalid Type", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  notification.Name,
							"value": notification.Value,
							"type":  "bad-type",
						}),
						ExpectError: regexp.MustCompile("Invalid notification type"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
		})

		t.Run("Create Failure - Name Empty", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  "",
							"value": notification.Value,
							"type":  notification.Type,
						}),
						ExpectError: regexp.MustCompile("may not be empty"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
		})

		t.Run("Create Failure - Value Empty", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  notification.Name,
							"value": "",
							"type":  notification.Type,
						}),
						ExpectError: regexp.MustCompile("may not be empty"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
		})

		t.Run("Create Failure", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  notification.Name,
							"value": notification.Value,
							"type":  notification.Type,
						}),
						ExpectError: regexp.MustCompile("could not create notification: error"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().NotificationCreate(client.NotificationCreate{
					Name:  notification.Name,
					Type:  notification.Type,
					Value: notification.Value,
				}).Times(1).Return(nil, errors.New("error"))
			})
		})

		t.Run("Update Failure", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  notification.Name,
							"type":  notification.Type,
							"value": notification.Value,
						}),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":  updatedNotification.Name,
							"type":  updatedNotification.Type,
							"value": updatedNotification.Value,
						}),
						ExpectError: regexp.MustCompile("could not update notification: error"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().NotificationCreate(client.NotificationCreate{
					Name:  notification.Name,
					Type:  notification.Type,
					Value: notification.Value,
				}).Times(1).Return(&notification, nil)

				mock.EXPECT().NotificationUpdate(updatedNotification.Id, client.NotificationUpdate{
					Name:  updatedNotification.Name,
					Type:  updatedNotification.Type,
					Value: updatedNotification.Value,
				}).Times(1).Return(nil, errors.New("error"))

				mock.EXPECT().Notifications().Times(2).Return([]client.Notification{notification}, nil)
				mock.EXPECT().NotificationDelete(notification.Id).Times(1)
			})
		})
		/*
			t.Run("Import By Name", func(t *testing.T) {
				testCase := resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
								"name":  notification.Name,
								"type":  notification.Type,
								"value": notification.Value,
							}),
						},
						{
							ResourceName:      resourceNameImport,
							ImportState:       true,
							ImportStateId:     notification.Name,
							ImportStateVerify: true,
						},
					},
				}

				runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
					mock.EXPECT().NotificationCreate(client.NotificationCreate{
						Name:  notification.Name,
						Type:  notification.Type,
						Value: notification.Value,
					}).Times(1).Return(&notification, nil)
					mock.EXPECT().Notifications().Times(3).Return([]client.Notification{notification}, nil)
					mock.EXPECT().NotificationDelete(notification.Id).Times(1)
				})
			})

			t.Run("Import By Id", func(t *testing.T) {
				testCase := resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
								"name":  notificationById.Name,
								"type":  notificationById.Type,
								"value": notificationById.Value,
							}),
						},
						{
							ResourceName:      resourceNameImport,
							ImportState:       true,
							ImportStateId:     notificationById.Id,
							ImportStateVerify: true,
						},
					},
				}

				runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
					mock.EXPECT().NotificationCreate(client.NotificationCreate{
						Name:  notificationById.Name,
						Type:  notificationById.Type,
						Value: notificationById.Value,
					}).Times(1).Return(&notificationById, nil)
					mock.EXPECT().Notifications().Times(3).Return([]client.Notification{notificationById}, nil)
					mock.EXPECT().NotificationDelete(notificationById.Id).Times(1)
				})
			})
	*/
}
