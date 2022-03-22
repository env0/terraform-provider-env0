package env0

import (
	"context"
	"fmt"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		ReadContext:   resourceNotificationRead,
		UpdateContext: resourceNotificationUpdate,
		DeleteContext: resourceNotificationDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceNotificationImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "the name of the notification",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "'Slack' or 'Teams'",
				Required:    true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					notificationType := client.NotificationType(i.(string))
					if notificationType != client.NotificationTypeSlack && notificationType != client.NotificationTypeTeams {
						return diag.Errorf("Invalid notification type")
					}
					return nil
				},
			},
			"value": {
				Type:             schema.TypeString,
				Description:      "the target url of the notification",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
		},
	}
}

func setNotificationSchema(d *schema.ResourceData, notification *client.Notification) {
	d.Set("name", notification.Name)
	d.Set("type", notification.Type)
	d.Set("value", notification.Value)
}

func getNotificationById(id string, meta interface{}) (*client.Notification, error) {
	apiClient := meta.(client.ApiClientInterface)

	notifications, err := apiClient.Notifications()
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		if notification.Id == id {
			return &notification, nil
		}
	}
	return nil, nil
}

func getNotificationByName(name string, meta interface{}) (*client.Notification, error) {
	apiClient := meta.(client.ApiClientInterface)

	notifications, err := apiClient.Notifications()
	if err != nil {
		return nil, err
	}

	var foundNotifications []client.Notification
	for _, notification := range notifications {
		if notification.Name == name {
			foundNotifications = append(foundNotifications, notification)
		}
	}

	if len(foundNotifications) == 0 {
		return nil, fmt.Errorf("notification with name %v not found", name)
	}

	if len(foundNotifications) > 1 {
		return nil, fmt.Errorf("found multiple notifications with name: %s. Use id instead or make sure notification names are unique %v", name, foundNotifications)
	}

	return &foundNotifications[0], nil
}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	payload := client.NotificationCreate{
		Name:  d.Get("name").(string),
		Type:  client.NotificationType(d.Get("type").(string)),
		Value: d.Get("value").(string),
	}

	notification, err := apiClient.NotificationCreate(payload)
	if err != nil {
		return diag.Errorf("could not create notification: %v", err)
	}

	d.SetId(notification.Id)

	return nil
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	notification, err := getNotificationById(d.Id(), meta)
	if err != nil {
		return diag.Errorf("could not get notification: %v", err)
	}
	if notification == nil {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
		return nil
	}

	setNotificationSchema(d, notification)

	return nil
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.NotificationUpdate{}

	if name, ok := d.GetOk("name"); ok {
		payload.Name = name.(string)
	}

	if notificationType, ok := d.GetOk("type"); ok {
		payload.Type = client.NotificationType(notificationType.(string))
	}

	if value, ok := d.GetOk("value"); ok {
		payload.Value = value.(string)
	}

	_, err := apiClient.NotificationUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update notification: %v", err)
	}

	return nil
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.NotificationDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete notification: %v", err)
	}

	return nil
}

func getNotification(id string, meta interface{}) (*client.Notification, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		log.Println("[INFO] Resolving notification by id: ", id)
		return getNotificationById(id, meta)
	} else {
		log.Println("[INFO] Resolving notification by name: ", id)
		return getNotificationByName(id, meta)
	}
}

func resourceNotificationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	notification, err := getNotification(d.Id(), meta)
	if err != nil {
		return nil, err
	}

	d.SetId(notification.Id)
	setNotificationSchema(d, notification)

	return []*schema.ResourceData{d}, nil
}
