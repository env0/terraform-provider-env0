package client

type GpgKey struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	KeyId   string `json:"keyId"`
	Content string `json:"content"`
}

type GpgKeyCreatePayload struct {
	Name    string `json:"name"`
	KeyId   string `json:"keyId"`
	Content string `json:"content"`
}

func (client *ApiClient) GpgKeyCreate(payload *GpgKeyCreatePayload) (*GpgKey, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganzationId := struct {
		OrganizationId string `json:"organizationId"`
		GpgKeyCreatePayload
	}{
		organizationId,
		*payload,
	}

	var result GpgKey
	if err := client.http.Post("/gpg-keys", payloadWithOrganzationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) GpgKeyDelete(id string) error {
	return client.http.Delete("/gpg-keys/" + id)
}

func (client *ApiClient) GpgKeys() ([]GpgKey, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []GpgKey
	if err := client.http.Get("/gpg-keys", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
