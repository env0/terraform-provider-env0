package client

type VcsConnection struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	VcsAgentKey string `json:"vcsAgentKey"`
}

type VcsConnectionCreatePayload struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	VcsAgentKey string `json:"vcsAgentKey,omitempty"`
}

type VcsConnectionUpdatePayload struct {
	Name        string `json:"name"`
	VcsAgentKey string `json:"vcsAgentKey,omitempty"`
}

func (client *ApiClient) VcsConnection(id string) (*VcsConnection, error) {
	var result VcsConnection

	if err := client.http.Get("/vcs/connections/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) VcsConnectionCreate(payload VcsConnectionCreatePayload) (*VcsConnection, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrg := struct {
		VcsConnectionCreatePayload
		OrganizationId string `json:"organizationId"`
	}{
		VcsConnectionCreatePayload: payload,
		OrganizationId:             organizationId,
	}

	var result VcsConnection
	if err := client.http.Post("/vcs/connections", payloadWithOrg, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) VcsConnectionUpdate(id string, payload VcsConnectionUpdatePayload) (*VcsConnection, error) {
	var result VcsConnection

	if err := client.http.Put("/vcs/connections/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) VcsConnectionDelete(id string) error {
	return client.http.Delete("/vcs/connections/"+id, nil)
}

func (client *ApiClient) VcsConnections() ([]VcsConnection, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []VcsConnection
	if err := client.http.Get("/vcs/connections", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}
