package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataNotificationsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Description: "list of all notifications (by name)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the notification name",
				},
			},
		},
	}
}

func dataNotificationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	notifications, err := apiClient.Notifications()
	if err != nil {
		return diag.Errorf("could not get notifications: %v", err)
	}

	names := []string{}

	for _, notification := range notifications {
		names = append(names, notification.Name)
	}

	if err := d.Set("names", names); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("all_notification_names")

	return nil
}
