package client

import (
	"errors"
)

func (self *ApiClient) ConfigurationVariables(scope Scope, scopeId string) ([]ConfigurationVariable, error) {
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
		"type":   "string",
		"format": params.Format,
	}
	if params.EnumValues != nil {
		schema["enum"] = params.EnumValues
	}
	return schema
}

func (self *ApiClient) ConfigurationVariableDelete(id string) error {
	return self.http.Delete("configuration/" + id)
}

func (self *ApiClient) ConfigurationVariableUpdate(params ConfigurationVariableUpdateParams) (ConfigurationVariable, error) {
	basicParams := params.BasicParams
	if basicParams.Scope == ScopeDeploymentLog || basicParams.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}
	var result []ConfigurationVariable
	request := map[string]interface{}{
		"id":             params.Id,
		"name":           basicParams.Name,
		"description":    basicParams.Description,
		"value":          basicParams.Value,
		"isSensitive":    basicParams.IsSensitive,
		"scope":          basicParams.Scope,
		"type":           basicParams.Type,
		"organizationId": organizationId,
	}
	if basicParams.Scope != ScopeGlobal {
		request["scopeId"] = basicParams.ScopeId
	}

	request["schema"] = getSchema(params.BasicParams)

	requestInArray := []map[string]interface{}{request}
	err = self.http.Post("/configuration", requestInArray, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result[0], nil
}
