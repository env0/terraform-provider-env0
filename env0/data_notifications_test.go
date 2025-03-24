package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestNotificationsDataSource(t *testing.T) {
	notification1 := client.Notification{
		Id:    "id0",
		Name:  "name0",
		Type:  client.NotificationTypeSlack,
		Value: "http://my-notification-0.com",
	}

	notification2 := client.Notification{
		Id:    "id1",
		Name:  "name10",
		Type:  client.NotificationTypeTeams,
		Value: "http://my-notification-0.com",
	}

	resourceType := "env0_notifications"
	resourceName := "test_notifications"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "names.0", notification1.Name),
						resource.TestCheckResourceAttr(accessor, "names.1", notification2.Name),
					),
				},
			},
		}
	}

	mockNotifications := func(returnValue []client.Notification) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Notifications().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockNotifications([]client.Notification{notification1, notification2}),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Notifications().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
