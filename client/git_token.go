package client

type GitToken struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Value          string `json:"value"          tfschema:"-"`
	OrganizationId string `json:"organizationId"`
}

type GitTokenCreatePayload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GitTokenCreatePayloadWith struct {
	GitTokenCreatePayload
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}

func (client *ApiClient) GitTokenCreate(payload GitTokenCreatePayload) (*GitToken, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := GitTokenCreatePayloadWith{
		GitTokenCreatePayload: payload,
		OrganizationId:        organizationId,
		Type:                  "GIT",
	}

	var result GitToken
	if err := client.http.Post("/tokens", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) GitToken(id string) (*GitToken, error) {
	var result GitToken
	if err := client.http.Get("/tokens/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) GitTokenDelete(id string) error {
	return client.http.Delete("/tokens/"+id, nil)
}

func (client *ApiClient) GitTokens() ([]GitToken, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []GitToken
	if err := client.http.Get("/tokens", map[string]string{"organizationId": organizationId, "type": "GIT"}, &result); err != nil {
		return nil, err
	}

	return result, err
}
