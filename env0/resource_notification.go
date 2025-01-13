package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Description: "'Slack', 'Teams', 'Email' or 'Webhook'",
				Required:    true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					notificationType := client.NotificationType(i.(string))
					if notificationType != client.NotificationTypeSlack && notificationType != client.NotificationTypeTeams && notificationType != client.NotificationTypeEmail && notificationType != client.NotificationTypeWebhook {
						return diag.Errorf("Invalid notification type")
					}

					return nil
				},
			},
			"value": {
				Type:             schema.TypeString,
				Description:      "URL for Slack, Teams or Webhooks endpoint. Coma separated list of email addresses for email endpoint, you can use `$ENVIRONMENT_CREATOR$`, and `$DEPLOYER$` to resolve emails dynamically.",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
			"webhook_secret": {
				Type:        schema.TypeString,
				Description: "the webhook secret to use for signing the webhook payload",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
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

	return nil, ErrNotFound
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

	var payload client.NotificationCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
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
		if errors.Is(err, ErrNotFound) {
			if notification == nil {
				tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
				d.SetId("")

				return nil
			}
		}

		return diag.Errorf("could not get notification: %v", err)
	}

	if err := writeResourceData(notification, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.NotificationUpdatePayload{}

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if d.HasChanges("webhook_secret") {
		webhookSecret := d.Get("webhook_secret").(string)

		var strPtr *string

		if webhookSecret == "" {
			// webhook secret was removed - pass 'null pointer' (will pass 'null' in json).
			// see https://docs.env0.com/reference/notifications-update-notification-endpoint
			payload.WebhookSecret = &strPtr
		} else {
			strPtr = &webhookSecret
			payload.WebhookSecret = &strPtr
		}
	}

	if _, err := apiClient.NotificationUpdate(d.Id(), payload); err != nil {
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

func getNotification(ctx context.Context, id string, meta interface{}) (*client.Notification, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving notification by id", map[string]interface{}{"id": id})

		return getNotificationById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving notification by name", map[string]interface{}{"name	": id})

		return getNotificationByName(id, meta)
	}
}

func resourceNotificationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	notification, err := getNotification(ctx, d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("notification with id %v not found", d.Id())
		}

		return nil, err
	}

	if err := writeResourceData(notification, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
