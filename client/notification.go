package client

import "errors"

func (ac *ApiClient) Notifications() ([]Notification, error) {
	var result []Notification
	err := ac.http.Get("/notifications/endpoints", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ac *ApiClient) Notification(id string) (*Notification, error) {
	notifications, err := ac.Notifications()
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		if notification.Id == id {
			return &notification, nil
		}
	}
	return nil, errors.New("not found")
}

func (ac *ApiClient) NotificationCreate(payload NotificationCreate) (*Notification, error) {
	var result Notification

	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := NotificationCreateWithOrganizationId{
		NotificationCreate: payload,
		OrganizationId:     organizationId,
	}

	err = ac.http.Post("/notifications/endpoints", payloadWithOrganizationId, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ac *ApiClient) NotificationDelete(id string) error {
	err := ac.http.Delete("/notifications/endpoints/" + id)
	if err != nil {
		return err
	}
	return nil
}

func (ac *ApiClient) NotificationUpdate(id string, payload NotificationUpdate) (*Notification, error) {
	var result Notification

	err := ac.http.Patch("/notifications/endpoints/"+id, payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
