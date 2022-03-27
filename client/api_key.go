package client

type APIKeyCreatePayloadWith struct {
	APIKeyCreatePayload
	OrganizationId string `json:"organizationId"`
}

func (ac *ApiClient) APIKeyCreate(payload APIKeyCreatePayload) (*APIKey, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := APIKeyCreatePayloadWith{
		APIKeyCreatePayload: payload,
		OrganizationId:      organizationId,
	}

	var result APIKey
	if err := ac.http.Post("/api-keys", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) APIKeyDelete(id string) error {
	return ac.http.Delete("/api-keys/" + id)
}

func (ac *ApiClient) APIKeys() ([]APIKey, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result []APIKey
	if err := ac.http.Get("/api-keys", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
