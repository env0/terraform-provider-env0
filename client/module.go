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
	if err := ac.http.Post("/modules", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) Module(id string) (*Module, error) {
	var result Module
	if err := ac.http.Get("/modules/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) ModuleDelete(id string) error {
	return ac.http.Delete("/modules/" + id)
}

func (ac *ApiClient) ModuleUpdate(id string, payload ModuleUpdatePayload) (*Module, error) {
	var result Module
	if err := ac.http.Patch("/modules/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ac *ApiClient) Modules() ([]Module, error) {
	organizationId, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result []Module
	if err := ac.http.Get("/modules", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
