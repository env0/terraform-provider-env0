package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestNotificationDataSource(t *testing.T) {
	notification := client.Notification{
		Id:    "id0",
		Name:  "my-notification-0",
		Type:  client.NotificationTypeSlack,
		Value: "http://my-notification-0.com",
	}

	otherNotification := client.Notification{
		Id:    "id1",
		Name:  "my-notification-1",
		Type:  client.NotificationTypeTeams,
		Value: "http://my-notification-1.com",
	}

	notificationFieldsByName := map[string]any{"name": notification.Name}
	notificationFieldsById := map[string]any{"id": notification.Id}

	resourceType := "env0_notification"
	resourceName := "test_notification"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", notification.Id),
						resource.TestCheckResourceAttr(accessor, "name", notification.Name),
						resource.TestCheckResourceAttr(accessor, "type", string(notification.Type)),
						resource.TestCheckResourceAttr(accessor, "value", notification.Value),
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

	mockListNotificationsCall := func(returnValue []client.Notification) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Notifications().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(notificationFieldsById),
			mockListNotificationsCall([]client.Notification{notification, otherNotification}),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(notificationFieldsByName),
			mockListNotificationsCall([]client.Notification{notification, otherNotification}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one notification exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(notificationFieldsByName, "found multiple notifications"),
			mockListNotificationsCall([]client.Notification{notification, notification, otherNotification}),
		)
	})

	t.Run("Throw error when by name and no notifications found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(notificationFieldsByName, "not found"),
			mockListNotificationsCall([]client.Notification{otherNotification}),
		)
	})
}
