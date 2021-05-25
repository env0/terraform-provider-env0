package client

import (
	"errors"
)

func (self *ApiClient) ConfigurationVariables(scope Scope, scopeId string) ([]ConfigurationVariable, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}

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

	result, err := self.http.Get("/configuration", params)
	if err != nil {
		return []ConfigurationVariable{}, err
	}
	return result.([]ConfigurationVariable), nil
}

func (self *ApiClient) ConfigurationVariableCreate(name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error) {
	if scope == ScopeDeploymentLog || scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}

	request := map[string]interface{}{
		"name":           name,
		"value":          value,
		"isSensitive":    isSensitive,
		"scope":          scope,
		"type":           type_,
		"organizationId": organizationId,
	}
	if scope != ScopeGlobal {
		request["scopeId"] = scopeId
	}
	if enumValues != nil {
		request["schema"] = map[string]interface{}{
			"type": "string",
			"enum": enumValues,
		}
	}
	requestInArray := []map[string]interface{}{request}
	result, err := self.http.Post("configuration", requestInArray)
	if err != nil {
		return ConfigurationVariable{}, err
	}

	configurations := result.([]ConfigurationVariable)
	return configurations[0], nil
}

func (self *ApiClient) ConfigurationVariableDelete(id string) error {
	return self.http.Delete("configuration/" + id)
}

func (self *ApiClient) ConfigurationVariableUpdate(id string, name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error) {
	if scope == ScopeDeploymentLog || scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}

	request := map[string]interface{}{
		"id":             id,
		"name":           name,
		"value":          value,
		"isSensitive":    isSensitive,
		"scope":          scope,
		"type":           type_,
		"organizationId": organizationId,
	}
	if scope != ScopeGlobal {
		request["scopeId"] = scopeId
	}
	if enumValues != nil {
		request["schema"] = map[string]interface{}{
			"type": "string",
			"enum": enumValues,
		}
	}
	requestInArray := []map[string]interface{}{request}
	result, err := self.http.Post("/configuration", requestInArray)
	if err != nil {
		return ConfigurationVariable{}, err
	}

	configurations := result.([]ConfigurationVariable)
	return configurations[0], nil
}
