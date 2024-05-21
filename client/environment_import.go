package client

type EnvironmentImport struct {
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	OrganizationId string     `json:"organizationId"`
	Workspace      string     `json:"workspace"`
	Variables      []Variable `json:"variables"`
	GitConfig      GitConfig  `json:"gitConfig"`
	IacType        string     `json:"iacType"` // opentofu or terraform
	IacVersion     string     `json:"iacVersion"`
}

type Variable struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsSensitive bool   `json:"isSensitive"`
	Type        string `json:"tyoe"` // string or JSON
}

type GitConfig struct {
	Path       string `json:"path"`
	Revision   string `json:"revision"`
	Repository string `json:"repository"`
	Provider   string `json:"provider"`
}

type EnvironmentImportCreatePayload struct {
	Name       string     `json:"name"`
	Workspace  string     `json:"workspace"`
	Variables  []Variable `json:"variables"`
	GitConfig  GitConfig  `json:"gitConfig"`
	IacType    string     `json:"iacType"`
	IacVersion string     `json:"iacVersion"`
}

type EnvironmentImportUpdatePayload struct {
	Name       string     `json:"name"`
	Variables  []Variable `json:"variables"`
	GitConfig  GitConfig  `json:"gitConfig"`
	IacType    string     `json:"iacType"`
	IacVersion string     `json:"iacVersion"`
}

func (client *ApiClient) EnvironmentImportCreate(payload *EnvironmentImportCreatePayload) (*EnvironmentImport, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganzationId := struct {
		OrganizationId string `json:"organizationId"`
		EnvironmentImportCreatePayload
	}{
		organizationId,
		*payload,
	}

	var result EnvironmentImport
	if err := client.http.Post("/staging-environments", payloadWithOrganzationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) EnvironmentImportUpdate(id string, payload *EnvironmentImportUpdatePayload) (*EnvironmentImport, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganzationId := struct {
		OrganizationId string `json:"organizationId"`
		EnvironmentImportUpdatePayload
	}{
		organizationId,
		*payload,
	}

	var result EnvironmentImport
	if err := client.http.Put("/staging-environments/"+id, payloadWithOrganzationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) EnvironmentImportGet(id string) (*EnvironmentImport, error) {
	var result EnvironmentImport
	if err := client.http.Get("/staging-environments/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
