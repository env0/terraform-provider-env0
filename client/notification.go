package client

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

func (ac *ApiClient) NotificationUpdate(id string, payload NotificationUpdate) (*Notification, error) {
	var result Notification

	if err := ac.http.Patch("/notifications/endpoints/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
