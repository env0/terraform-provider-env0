package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataNotificationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the notification",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the notification",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "the type of the notification",
				Computed:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "the target url of the notification",
				Computed:    true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "textual description of the entity who created the notification",
				Computed:    true,
			},
		},
	}
}

func dataNotificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var notification *client.Notification
	var err error

	id, ok := d.GetOk("id")
	if ok {
		notification, err = getNotificationById(id.(string), meta)
	} else {
		name := d.Get("name").(string)
		notification, err = getNotificationByName(name, meta)
	}

	if err != nil {
		return diag.Errorf("could not read notification: %v", err)
	}

	if err := writeResourceData(notification, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
