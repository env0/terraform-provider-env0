package client

type EnvironmentImport struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationId string    `json:"organizationId"`
	Workspace      string    `json:"workspace"`
	GitConfig      GitConfig `json:"gitConfig"`
	IacType        string    `json:"iacType"` // opentofu or terraform
	IacVersion     string    `json:"iacVersion"`
}

type GitConfig struct {
	Path       string `json:"path,omitempty"`
	Revision   string `json:"revision,omitempty"`
	Repository string `json:"repository,omitempty"`
	Provider   string `json:"provider,omitempty" tfschema:"git_provider"`
}

type EnvironmentImportCreatePayload struct {
	Name       string    `json:"name,omitempty"`
	Workspace  string    `json:"workspace,omitempty"`
	GitConfig  GitConfig `json:"gitConfig,omitempty"`
	IacType    string    `json:"iacType,omitempty"`
	IacVersion string    `json:"iacVersion,omitempty"`
}

type EnvironmentImportUpdatePayload struct {
	Name       string    `json:"name,omitempty"`
	GitConfig  GitConfig `json:"gitConfig,omitempty"`
	IacType    string    `json:"iacType,omitempty"`
	IacVersion string    `json:"iacVersion,omitempty"`
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

func (client *ApiClient) EnvironmentImportDelete(id string) error {
	return client.http.Delete("/environment-imports/"+id, nil)
}
