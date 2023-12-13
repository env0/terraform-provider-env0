package client

type ApiKey struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ApiKeyId         string `json:"apiKeyId"`
	ApiKeySecret     string `json:"apiKeySecret" tfschema:",omitempty"`
	LastUsedAt       string `json:"lastUsedAt"`
	OrganizationId   string `json:"organizationId"`
	OrganizationRole string `json:"organizationRole"`
	CreatedAt        string `json:"createdAt"`
	CreatedBy        string `json:"createdBy"`
	CreatedByUser    User   `json:"createdByUser"`
}

type ApiKeyPermissions struct {
	OrganizationRole string `json:"organizationRole"`
}

type ApiKeyCreatePayload struct {
	Name        string            `json:"name"`
	Permissions ApiKeyPermissions `json:"permissions"`
}

type ApiKeyCreatePayloadWith struct {
	ApiKeyCreatePayload
	OrganizationId string `json:"organizationId"`
}

func (client *ApiClient) ApiKeyCreate(payload ApiKeyCreatePayload) (*ApiKey, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := ApiKeyCreatePayloadWith{
		ApiKeyCreatePayload: payload,
		OrganizationId:      organizationId,
	}

	var result ApiKey
	if err := client.http.Post("/api-keys", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ApiKeyDelete(id string) error {
	return client.http.Delete("/api-keys/"+id, nil)
}

func (client *ApiClient) ApiKeys() ([]ApiKey, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []ApiKey
	if err := client.http.Get("/api-keys", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (client *ApiClient) OidcSub() (string, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return "", err
	}

	var result string
	if err := client.http.Get("/api-keys/oidc-sub", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return "", err
	}

	return result, nil
}
