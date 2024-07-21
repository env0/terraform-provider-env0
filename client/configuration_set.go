package client

import "fmt"

type CreateConfigurationSetPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// if Scope is "organization", scopeId will be calculated in the functions.
	Scope                   string                  `json:"scope"`   // "project" or "organization".
	ScopeId                 string                  `json:"scopeId"` // project id or organization id.
	ConfigurationProperties []ConfigurationVariable `json:"configurationProperties" tfschema:"-"`
}

type UpdateConfigurationSetPayload struct {
	Name                           string                  `json:"name"`
	Description                    string                  `json:"description"`
	ConfigurationPropertiesChanges []ConfigurationVariable `json:"configurationPropertiesChanges" tfschema:"-"` // delta changes.
}

type ConfigurationSet struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	AssignmentScope string `json:"assignmentScope"`
	CreationScopeId string `json:"creationScopeId"`
}

func (client *ApiClient) ConfigurationSetCreate(payload *CreateConfigurationSetPayload) (*ConfigurationSet, error) {
	var result ConfigurationSet
	var err error

	if payload.Scope == "organization" && payload.ScopeId == "" {
		payload.ScopeId, err = client.OrganizationId()
		if err != nil {
			return nil, err
		}
	}

	if err := client.http.Post("/configuration-sets", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ConfigurationSetUpdate(id string, payload *UpdateConfigurationSetPayload) (*ConfigurationSet, error) {
	var result ConfigurationSet

	if err := client.http.Put("/configuration-sets/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ConfigurationSet(id string) (*ConfigurationSet, error) {
	var result ConfigurationSet

	if err := client.http.Get("/configuration-sets/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ConfigurationSets(scope string, scopeId string) ([]ConfigurationSet, error) {
	var result []ConfigurationSet

	params := map[string]string{
		"scope":   scope,
		"scopeId": scopeId,
	}

	if err := client.http.Get("/configuration-sets", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *ApiClient) ConfigurationSetDelete(id string) error {
	return client.http.Delete("/configuration-sets/"+id, nil)
}

func (client *ApiClient) ConfigurationVariablesBySetId(setId string) ([]ConfigurationVariable, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization id: %w", err)
	}

	var result []ConfigurationVariable

	if err := client.http.Get("/configuration", map[string]string{
		"setId":          setId,
		"organizationId": organizationId,
	}, &result); err != nil {
		return nil, err
	}
	return result, nil
}
