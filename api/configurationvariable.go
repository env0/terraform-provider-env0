package api

import (
	"errors"
	"net/url"
)

func (self *ApiClient) ConfigurationVariables(scope Scope, scopeId string) ([]ConfigurationVariable, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []ConfigurationVariable
	params := url.Values{}
	params.Add("organizationId", organizationId)
	switch {
	case scope == ScopeGlobal:
	case scope == ScopeTemplate:
		params.Add("blueprintId", scopeId)
	case scope == ScopeProject:
		params.Add("projectId", scopeId)
	case scope == ScopeEnvironment:
		params.Add("environmentId", scopeId)
	case scope == ScopeDeployment:
		return nil, errors.New("No api to fetch configuration variables by deployment")
	case scope == ScopeDeploymentLog:
		params.Add("deploymentLogId", scopeId)
	}
	err = self.getJSON("/configuration", params, &result)
	if err != nil {
		return []ConfigurationVariable{}, err
	}
	return result, nil
}

func (self *ApiClient) ConfigurationVariableCreate(name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error) {
	if scope == ScopeDeploymentLog || scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}
	var result ConfigurationVariable
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
	err = self.postJSON("/configuration", request, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result, nil
}

func (self *ApiClient) ConfigurationVariableDelete(id string) error {
	return self.delete("/configuration/" + id)
}

func (self *ApiClient) ConfigurationVariableUpdate(id string, name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error) {
	if scope == ScopeDeploymentLog || scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("Must not create variable on scope deployment / deploymentLog")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}
	var result ConfigurationVariable
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
	err = self.postJSON("/configuration", request, &result)
	if err != nil {
		return ConfigurationVariable{}, err
	}
	return result, nil
}
