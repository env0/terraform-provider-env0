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
	Value       string `json:"value,omitempty"`
	IsSensitive bool   `json:"isSensitive"`
	Type        string `json:"type"` // string or JSON
}

type GitConfig struct {
	Path       string `json:"path,omitempty"`
	Revision   string `json:"revision,omitempty"`
	Repository string `json:"repository,omitempty"`
	Provider   string `json:"provider,omitempty"`
}

type EnvironmentImportCreatePayload struct {
	Name       string     `json:"name,omitempty"`
	Workspace  string     `json:"workspace,omitempty"`
	Variables  []Variable `json:"variables,omitempty"`
	GitConfig  GitConfig  `json:"gitConfig,omitempty"`
	IacType    string     `json:"iacType,omitempty"`
	IacVersion string     `json:"iacVersion,omitempty"`
}

type EnvironmentImportUpdatePayload struct {
	Name       string     `json:"name,omitempty"`
	Variables  []Variable `json:"variables,omitempty"`
	GitConfig  GitConfig  `json:"gitConfig,omitempty"`
	IacType    string     `json:"iacType,omitempty"`
	IacVersion string     `json:"iacVersion,omitempty"`
}

func (client *ApiClient) EnvironmentImportCreate(payload *EnvironmentImportCreatePayload) (*EnvironmentImport, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := struct {
		OrganizationId string `json:"organizationId"`
		EnvironmentImportCreatePayload
	}{
		organizationId,
		*payload,
	}

	var result EnvironmentImport
	if err := client.http.Post("/environment-imports", payloadWithOrganizationId, &result); err != nil {
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
	if err := client.http.Put("/environment-imports/"+id, payloadWithOrganzationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) EnvironmentImportGet(id string) (*EnvironmentImport, error) {
	var result EnvironmentImport
	if err := client.http.Get("/environment-imports/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
