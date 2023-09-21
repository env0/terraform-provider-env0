package client

type RemoteStateAccessConfiguration struct {
	EnvironmentId                    string   `json:"environmentId"`
	AccessibleFromEntireOrganization bool     `json:"accessibleFromEntireOrganization"`
	AllowedProjectIds                []string `json:"allowedProjectIds" tfschema:",omitempty"`
}

type RemoteStateAccessConfigurationCreate struct {
	AccessibleFromEntireOrganization bool     `json:"accessibleFromEntireOrganization"`
	AllowedProjectIds                []string `json:"allowedProjectIds"`
}

func (client *ApiClient) RemoteStateAccessConfiguration(environmentId string) (*RemoteStateAccessConfiguration, error) {
	var result RemoteStateAccessConfiguration

	if err := client.http.Get("/remote-backend/states/"+environmentId+"/access-control", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoteStateAccessConfigurationCreate(environmentId string, payload RemoteStateAccessConfigurationCreate) (*RemoteStateAccessConfiguration, error) {
	// The API doesn't accept "null". If the array isn't set, pass "[]" instead.
	if payload.AllowedProjectIds == nil {
		payload.AllowedProjectIds = []string{}
	}

	var result RemoteStateAccessConfiguration
	if err := client.http.Put("/remote-backend/states/"+environmentId+"/access-control", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoteStateAccessConfigurationDelete(environmentId string) error {
	return client.http.Delete("/remote-backend/states/"+environmentId+"/access-control", nil)
}
