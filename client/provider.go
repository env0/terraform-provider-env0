package client

type Provider struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ProviderCreatePayload struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ProviderUpdatePayload struct {
	Description string `json:"description"`
}

func (client *ApiClient) ProviderCreate(payload ProviderCreatePayload) (*Provider, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := struct {
		ProviderCreatePayload
		OrganizationId string `json:"organizationId"`
	}{
		payload,
		organizationId,
	}

	var result Provider
	if err := client.http.Post("/providers", payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Provider(providerId string) (*Provider, error) {
	var result Provider
	if err := client.http.Get("/providers/"+providerId, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ProviderDelete(providerId string) error {
	return client.http.Delete("/providers/" + providerId)
}

func (client *ApiClient) ProviderUpdate(providerId string, payload ProviderUpdatePayload) (*Provider, error) {
	var result Provider
	if err := client.http.Put("/providers/"+providerId, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Providers() ([]Provider, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Provider
	if err := client.http.Get("/providers", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
