package client

type GitTokenCreatePayloadWith struct {
	GitTokenCreatePayload
	OrganizationId string `json:"organizationId"`
}

func (ac *ApiClient) GitTokenCreate(payload GitTokenCreatePayload) (*GitToken, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := GitTokenCreatePayloadWith{
		GitTokenCreatePayload: payload,
		OrganizationId:        organizationId,
	}

	var result GitToken
	if err := ac.http.Post("/tokens", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) GitToken(id string) (*GitToken, error) {
	var result GitToken
	if err := ac.http.Get("/tokens/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) GitTokenDelete(id string) error {
	return ac.http.Delete("/tokens/" + id)
}

func (ac *ApiClient) GitTokens(gtType GitTokenType) ([]GitToken, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result []GitToken
	if err := ac.http.Get("/tokens", map[string]string{"organizationId": organizationId, "type": string(gtType)}, &result); err != nil {
		return nil, err
	}

	return result, err
}
