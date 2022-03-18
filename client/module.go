package client

type ModuleCreatePayloadWith struct {
	ModuleCreatePayload
	Type           string `json:"type"`
	OrganizationId string `json:"organizationId"`
}

func (ac *ApiClient) ModuleCreate(payload ModuleCreatePayload) (*Module, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := ModuleCreatePayloadWith{
		ModuleCreatePayload: payload,
		OrganizationId:      organizationId,
		Type:                "module",
	}

	var result Module
	err = ac.http.Post("/modules", payloadWith, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
