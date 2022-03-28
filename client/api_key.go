package client

type ApiKeyCreatePayloadWith struct {
	ApiKeyCreatePayload
	OrganizationId string `json:"organizationId"`
}

func (ac *ApiClient) ApiKeyCreate(payload ApiKeyCreatePayload) (*ApiKey, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := ApiKeyCreatePayloadWith{
		ApiKeyCreatePayload: payload,
		OrganizationId:      organizationId,
	}

	var result ApiKey
	if err := ac.http.Post("/api-keys", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) ApiKeyDelete(id string) error {
	return ac.http.Delete("/api-keys/" + id)
}

func (ac *ApiClient) ApiKeys() ([]ApiKey, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result []ApiKey
	if err := ac.http.Get("/api-keys", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
