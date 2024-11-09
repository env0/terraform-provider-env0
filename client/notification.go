package client

type NotificationType string

const (
	NotificationTypeSlack   NotificationType = "Slack"
	NotificationTypeTeams   NotificationType = "Teams"
	NotificationTypeEmail   NotificationType = "Email"
	NotificationTypeWebhook NotificationType = "Webhook"
)

type Notification struct {
	Id             string           `json:"id"`
	CreatedBy      string           `json:"createdBy"`
	CreatedByUser  User             `json:"createdByUser"`
	OrganizationId string           `json:"organizationId"`
	Name           string           `json:"name"`
	Type           NotificationType `json:"type"`
	Value          string           `json:"value"`
}

type NotificationCreatePayload struct {
	Name          string           `json:"name"`
	Type          NotificationType `json:"type"`
	Value         string           `json:"value"`
	WebhookSecret string           `json:"webhookSecret,omitempty"`
}

type NotificationCreatePayloadWith struct {
	NotificationCreatePayload
	OrganizationId string `json:"organizationId"`
}

type NotificationUpdatePayload struct {
	Name          string           `json:"name,omitempty"`
	Type          NotificationType `json:"type,omitempty"`
	Value         string           `json:"value,omitempty"`
	WebhookSecret **string         `json:"webhookSecret,omitempty" tfschema:"-"`
}

func (client *ApiClient) Notifications() ([]Notification, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Notification

	if err := client.http.Get("/notifications/endpoints", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *ApiClient) NotificationCreate(payload NotificationCreatePayload) (*Notification, error) {
	var result Notification

	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := NotificationCreatePayloadWith{
		NotificationCreatePayload: payload,
		OrganizationId:            organizationId,
	}

	if err = client.http.Post("/notifications/endpoints", payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) NotificationDelete(id string) error {
	if err := client.http.Delete("/notifications/endpoints/"+id, nil); err != nil {
		return err
	}

	return nil
}

func (client *ApiClient) NotificationUpdate(id string, payload NotificationUpdatePayload) (*Notification, error) {
	var result Notification

	if err := client.http.Patch("/notifications/endpoints/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
