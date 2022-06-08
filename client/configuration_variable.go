package client

import (
	"errors"
)

type Scope string

const (
	ScopeGlobal        Scope = "GLOBAL"
	ScopeTemplate      Scope = "BLUEPRINT"
	ScopeProject       Scope = "PROJECT"
	ScopeEnvironment   Scope = "ENVIRONMENT"
	ScopeDeployment    Scope = "DEPLOYMENT"
	ScopeDeploymentLog Scope = "DEPLOYMENT_LOG"
)

type Format string

const (
	Text Format = ""
	HCL  Format = "HCL"
	JSON Format = "JSON"
)

type ConfigurationVariableSchema struct {
	Type   string   `json:"type"`
	Enum   []string `json:"enum"`
	Format Format   `json:"format,omitempty"`
}

func (c *ConfigurationVariableSchema) ResourceDataSliceStructValueWrite(values map[string]interface{}) error {
	if len(c.Format) > 0 {
		values["format"] = c.Format
	}
	return nil
}

type ConfigurationVariable struct {
	ScopeId        string                       `json:"scopeId,omitempty"`
	Value          string                       `json:"value"`
	OrganizationId string                       `json:"organizationId,omitempty"`
	UserId         string                       `json:"userId,omitempty"`
	IsSensitive    *bool                        `json:"isSensitive,omitempty"`
	Scope          Scope                        `json:"scope,omitempty"`
	Id             string                       `json:"id,omitempty"`
	Name           string                       `json:"name"`
	Description    string                       `json:"description,omitempty"`
	Type           *ConfigurationVariableType   `json:"type,omitempty"`
	Schema         *ConfigurationVariableSchema `json:"schema,omitempty"`
	ToDelete       *bool                        `json:"toDelete,omitempty"`
	IsReadOnly     *bool                        `json:"isReadonly,omitempty"`
	IsRequired     *bool                        `json:"isRequired,omitempty"`
	Regex          string                       `json:"regex,omitempty"`
}

type ConfigurationVariableCreateParams struct {
	Name        string
	Value       string
	IsSensitive bool
	Scope       Scope
	ScopeId     string
	Type        ConfigurationVariableType
	EnumValues  []string
	Description string
	Format      Format
	IsReadOnly  bool
	IsRequired  bool
	Regex       string
}

type ConfigurationVariableUpdateParams struct {
	CommonParams ConfigurationVariableCreateParams
	Id           string
}

func (client *ApiClient) ConfigurationVariablesById(id string) (ConfigurationVariable, error) {
	var result ConfigurationVariable

	err := client.http.Get("/configuration/"+id, nil, &result)

	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result, nil
}

func (client *ApiClient) ConfigurationVariablesByScope(scope Scope, scopeId string) ([]ConfigurationVariable, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}
	var result []ConfigurationVariable
	params := map[string]string{"organizationId": organizationId}
	switch {
	case scope == ScopeGlobal:
	case scope == ScopeTemplate:
		params["blueprintId"] = scopeId
	case scope == ScopeProject:
		params["projectId"] = scopeId
	case scope == ScopeEnvironment:
		params["environmentId"] = scopeId
	case scope == ScopeDeployment:
		return nil, errors.New("no api to fetch configuration variables by deployment")
	case scope == ScopeDeploymentLog:
		params["deploymentLogId"] = scopeId
	}
	err = client.http.Get("/configuration", params, &result)
	if err != nil {
		return []ConfigurationVariable{}, err
	}

	if scope != ScopeGlobal {
		// Filter out global scopes. If a non global scope is requested.
		filteredResult := []ConfigurationVariable{}
		for _, variable := range result {
			if variable.Scope != ScopeGlobal {
				filteredResult = append(filteredResult, variable)
			}
		}
		return filteredResult, nil
	}

	return result, nil
}

func (client *ApiClient) ConfigurationVariableCreate(params ConfigurationVariableCreateParams) (ConfigurationVariable, error) {
	if params.Scope == ScopeDeploymentLog || params.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := client.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}
	var result []ConfigurationVariable
	request := map[string]interface{}{
		"name":           params.Name,
		"description":    params.Description,
		"value":          params.Value,
		"isSensitive":    params.IsSensitive,
		"scope":          params.Scope,
		"type":           params.Type,
		"organizationId": organizationId,
		"isRequired":     params.IsRequired,
		"isReadonly":     params.IsReadOnly,
		"regex":          params.Regex,
	}
	if params.Scope != ScopeGlobal {
		request["scopeId"] = params.ScopeId
	}

	request["schema"] = getSchema(params)

	requestInArray := []map[string]interface{}{request}
	err = client.http.Post("configuration", requestInArray, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result[0], nil
}

func getSchema(params ConfigurationVariableCreateParams) map[string]interface{} {
	schema := map[string]interface{}{
		"type": "string",
	}
	if params.EnumValues != nil {
		schema["enum"] = params.EnumValues
	}
	if params.Format != Text {
		schema["format"] = params.Format
	}
	return schema
}

func (client *ApiClient) ConfigurationVariableDelete(id string) error {
	return client.http.Delete("configuration/" + id)
}

func (client *ApiClient) ConfigurationVariableUpdate(updateParams ConfigurationVariableUpdateParams) (ConfigurationVariable, error) {
	commonParams := updateParams.CommonParams
	if commonParams.Scope == ScopeDeploymentLog || commonParams.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := client.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}
	var result []ConfigurationVariable
	request := map[string]interface{}{
		"id":             updateParams.Id,
		"name":           commonParams.Name,
		"description":    commonParams.Description,
		"value":          commonParams.Value,
		"isSensitive":    commonParams.IsSensitive,
		"scope":          commonParams.Scope,
		"type":           commonParams.Type,
		"organizationId": organizationId,
		"isRequired":     commonParams.IsRequired,
		"isReadonly":     commonParams.IsReadOnly,
		"regex":          commonParams.Regex,
	}
	if commonParams.Scope != ScopeGlobal {
		request["scopeId"] = commonParams.ScopeId
	}

	request["schema"] = getSchema(updateParams.CommonParams)

	requestInArray := []map[string]interface{}{request}
	err = client.http.Post("/configuration", requestInArray, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result[0], nil
}
