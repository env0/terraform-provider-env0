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
	IsReadonly     *bool                        `json:"isReadonly,omitempty"`
	IsRequired     *bool                        `json:"isRequired,omitempty"`
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
	IsReadonly  bool
	IsRequired  bool
}

type ConfigurationVariableUpdateParams struct {
	CommonParams ConfigurationVariableCreateParams
	Id           string
}

func (self *ApiClient) ConfigurationVariablesById(id string) (ConfigurationVariable, error) {
	var result ConfigurationVariable

	err := self.http.Get("/configuration/"+id, nil, &result)

	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result, nil
}

func (self *ApiClient) ConfigurationVariablesByScope(scope Scope, scopeId string) ([]ConfigurationVariable, error) {
	organizationId, err := self.organizationId()
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
		return nil, errors.New("No api to fetch configuration variables by deployment")
	case scope == ScopeDeploymentLog:
		params["deploymentLogId"] = scopeId
	}
	err = self.http.Get("/configuration", params, &result)
	if err != nil {
		return []ConfigurationVariable{}, err
	}
	return result, nil
}

func (self *ApiClient) ConfigurationVariableCreate(params ConfigurationVariableCreateParams) (ConfigurationVariable, error) {
	if params.Scope == ScopeDeploymentLog || params.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
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
		"isReadonly":     params.IsReadonly,
	}
	if params.Scope != ScopeGlobal {
		request["scopeId"] = params.ScopeId
	}

	request["schema"] = getSchema(params)

	requestInArray := []map[string]interface{}{request}
	err = self.http.Post("configuration", requestInArray, &result)
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

func (self *ApiClient) ConfigurationVariableDelete(id string) error {
	return self.http.Delete("configuration/" + id)
}

func (self *ApiClient) ConfigurationVariableUpdate(updateParams ConfigurationVariableUpdateParams) (ConfigurationVariable, error) {
	commonParams := updateParams.CommonParams
	if commonParams.Scope == ScopeDeploymentLog || commonParams.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
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
		"isReadonly":     commonParams.IsReadonly,
	}
	if commonParams.Scope != ScopeGlobal {
		request["scopeId"] = commonParams.ScopeId
	}

	request["schema"] = getSchema(updateParams.CommonParams)

	requestInArray := []map[string]interface{}{request}
	err = self.http.Post("/configuration", requestInArray, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result[0], nil
}
