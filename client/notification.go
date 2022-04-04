package client

type NotificationType string

const (
	NotificationTypeSlack NotificationType = "Slack"
	NotificationTypeTeams NotificationType = "Teams"
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
	Name  string           `json:"name"`
	Type  NotificationType `json:"type"`
	Value string           `json:"value"`
}

type NotificationCreatePayloadWith struct {
	NotificationCreatePayload
	OrganizationId string `json:"organizationId"`
}

type NotificationUpdatePayload struct {
	Name  string           `json:"name,omitempty"`
	Type  NotificationType `json:"type,omitempty"`
	Value string           `json:"value,omitempty"`
}

func (ac *ApiClient) Notifications() ([]Notification, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result []Notification
	if err := ac.http.Get("/notifications/endpoints", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (ac *ApiClient) NotificationCreate(payload NotificationCreatePayload) (*Notification, error) {
	var result Notification

	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := NotificationCreatePayloadWith{
		NotificationCreatePayload: payload,
		OrganizationId:            organizationId,
	}

	if err = ac.http.Post("/notifications/endpoints", payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ac *ApiClient) NotificationDelete(id string) error {
	if err := ac.http.Delete("/notifications/endpoints/" + id); err != nil {
		return err
	}
	return nil
}

func (ac *ApiClient) NotificationUpdate(id string, payload NotificationUpdatePayload) (*Notification, error) {
	var result Notification

	if err := ac.http.Patch("/notifications/endpoints/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
