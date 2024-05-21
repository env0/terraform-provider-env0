package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitNotificationResource(t *testing.T) {
	resourceType := "env0_notification"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	user := client.User{
		Name: "user0",
	}

	notification := client.Notification{
		Id:             "id0",
		Type:           client.NotificationTypeSlack,
		Name:           "name0",
		Value:          "https://some.url.1.com",
		OrganizationId: "org0",
		CreatedBy:      "user0",
		CreatedByUser:  user,
	}

	notificationById := client.Notification{
		Id:             uuid.NewString(),
		Type:           client.NotificationTypeSlack,
		Name:           "name0",
		Value:          "https://some.url.1.com",
		OrganizationId: "org0",
		CreatedBy:      "user0",
		CreatedByUser:  user,
	}

	updatedNotification := client.Notification{
		Id:             notification.Id,
		Type:           client.NotificationTypeWebhook,
		Name:           "name0-updated",
		Value:          "https://some.updated.url.1.com",
		OrganizationId: notification.OrganizationId,
		CreatedBy:      notification.CreatedBy,
		CreatedByUser:  notification.CreatedByUser,
	}

	webhookSecret := "my-little-secret"
	var nullString *string

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":           notification.Name,
						"type":           notification.Type,
						"value":          notification.Value,
						"webhook_secret": webhookSecret,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", notification.Id),
						resource.TestCheckResourceAttr(accessor, "name", notification.Name),
						resource.TestCheckResourceAttr(accessor, "value", notification.Value),
						resource.TestCheckResourceAttr(accessor, "type", string(notification.Type)),
						resource.TestCheckResourceAttr(accessor, "webhook_secret", webhookSecret),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  updatedNotification.Name,
						"type":  updatedNotification.Type,
						"value": updatedNotification.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedNotification.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedNotification.Name),
						resource.TestCheckResourceAttr(accessor, "value", updatedNotification.Value),
						resource.TestCheckResourceAttr(accessor, "type", string(updatedNotification.Type)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().NotificationCreate(client.NotificationCreatePayload{
				Name:          notification.Name,
				Type:          notification.Type,
				Value:         notification.Value,
				WebhookSecret: webhookSecret,
			}).Times(1).Return(&notification, nil)

			mock.EXPECT().NotificationUpdate(updatedNotification.Id, client.NotificationUpdatePayload{
				Name:          updatedNotification.Name,
				Type:          updatedNotification.Type,
				Value:         updatedNotification.Value,
				WebhookSecret: &nullString,
			}).Times(1).Return(&updatedNotification, nil)

			gomock.InOrder(
				mock.EXPECT().Notifications().Times(2).Return([]client.Notification{notification}, nil),
				mock.EXPECT().Notifications().Times(1).Return([]client.Notification{updatedNotification}, nil),
			)

			mock.EXPECT().NotificationDelete(notification.Id).Times(1)
		})
	})

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
			mock.EXPECT().NotificationCreate(client.NotificationCreatePayload{
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
			mock.EXPECT().NotificationCreate(client.NotificationCreatePayload{
				Name:  notification.Name,
				Type:  notification.Type,
				Value: notification.Value,
			}).Times(1).Return(&notification, nil)

			mock.EXPECT().NotificationUpdate(updatedNotification.Id, client.NotificationUpdatePayload{
				Name:  updatedNotification.Name,
				Type:  updatedNotification.Type,
				Value: updatedNotification.Value,
			}).Times(1).Return(nil, errors.New("error"))

			mock.EXPECT().Notifications().Times(2).Return([]client.Notification{notification}, nil)
			mock.EXPECT().NotificationDelete(notification.Id).Times(1)
		})
	})

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
			mock.EXPECT().NotificationCreate(client.NotificationCreatePayload{
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
			mock.EXPECT().NotificationCreate(client.NotificationCreatePayload{
				Name:  notificationById.Name,
				Type:  notificationById.Type,
				Value: notificationById.Value,
			}).Times(1).Return(&notificationById, nil)
			mock.EXPECT().Notifications().Times(3).Return([]client.Notification{notificationById}, nil)
			mock.EXPECT().NotificationDelete(notificationById.Id).Times(1)
		})
	})
}
