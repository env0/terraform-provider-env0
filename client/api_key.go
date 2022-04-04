package client

type ApiKey struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ApiKeyId       string `json:"apiKeyId"`
	ApiKeySecret   string `json:"apiKeySecret"`
	LastUsedAt     string `json:"lastUsedAt"`
	OrganizationId string `json:"organizationId"`
	CreatedAt      string `json:"createdAt"`
	CreatedBy      string `json:"createdBy"`
	CreatedByUser  User   `json:"createdByUser"`
}

type ApiKeyCreatePayload struct {
	Name string `json:"name"`
}

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
