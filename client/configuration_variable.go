package client

import (
	"encoding/json"
	"errors"
)

type Scope string

const (
	ScopeGlobal         Scope = "GLOBAL"
	ScopeTemplate       Scope = "BLUEPRINT"
	ScopeSubEnvironment Scope = "SUB_ENVIRONMENT_BLUEPRINT"
	ScopeProject        Scope = "PROJECT"
	ScopeEnvironment    Scope = "ENVIRONMENT"
	ScopeDeployment     Scope = "DEPLOYMENT"
	ScopeDeploymentLog  Scope = "DEPLOYMENT_LOG"
	ScopeWorkflow       Scope = "WORKFLOW"
	ScopeSet            Scope = "SET"
)

type Format string

const (
	Text               Format = ""
	HCL                Format = "HCL"
	JSON               Format = "JSON"
	ENVIRONMENT_OUTPUT Format = "ENVIRONMENT_OUTPUT"
)

type ConfigurationVariableSchema struct {
	Type   string   `json:"type,omitempty"`
	Enum   []string `json:"enum"`
	Format Format   `json:"format,omitempty"`
}

func (c *ConfigurationVariableSchema) ResourceDataSliceStructValueWrite(values map[string]any) error {
	if len(c.Format) > 0 {
		values["format"] = c.Format
	}

	return nil
}

type ConfigurationVariableOverwrites struct {
	Value       string `json:"value"`
	Regex       string `json:"regex"`
	IsRequired  bool   `json:"isRequired"`
	IsSensitive bool   `json:"isSensitive"`
}

type ConfigurationVariable struct {
	ScopeId        string                           `json:"scopeId,omitempty"`
	Value          string                           `json:"value" tfschema:"-"`
	OrganizationId string                           `json:"organizationId,omitempty"`
	UserId         string                           `json:"userId,omitempty"`
	IsSensitive    *bool                            `json:"isSensitive,omitempty"`
	Scope          Scope                            `json:"scope,omitempty"`
	Id             string                           `json:"id,omitempty"`
	Name           string                           `json:"name"`
	Description    string                           `json:"description,omitempty"`
	Type           *ConfigurationVariableType       `json:"type,omitempty" tfschema:",omitempty"`
	Schema         *ConfigurationVariableSchema     `json:"schema,omitempty"`
	ToDelete       *bool                            `json:"toDelete,omitempty"`
	IsReadOnly     *bool                            `json:"isReadonly,omitempty"`
	IsRequired     *bool                            `json:"isRequired,omitempty"`
	Regex          string                           `json:"regex,omitempty"`
	Overwrites     *ConfigurationVariableOverwrites `json:"overwrites,omitempty"` // Is removed when marhseling to a JSON.
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

func (v ConfigurationVariable) MarshalJSON() ([]byte, error) {
	v.Overwrites = nil

	// This is done to prevent an infinite loop.
	type ConfigurationVariableDummy ConfigurationVariable

	dummy := ConfigurationVariableDummy(v)

	return json.Marshal(&dummy)
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
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []ConfigurationVariable

	params := map[string]string{"organizationId": organizationId}

	switch scope {
	case ScopeGlobal:
	case ScopeTemplate:
		params["blueprintId"] = scopeId
	case ScopeProject:
		params["projectId"] = scopeId
	case ScopeEnvironment:
		params["environmentId"] = scopeId
	case ScopeDeployment:
		return nil, errors.New("no api to fetch configuration variables by deployment")
	case ScopeDeploymentLog:
		params["deploymentLogId"] = scopeId
	case ScopeWorkflow:
		params["environmentId"] = scopeId
		params["workflowEnvironmentId"] = scopeId
	case ScopeSubEnvironment:
		params["workflowBlueprintIdSubAlias"] = scopeId
	}

	err = client.http.Get("/configuration", params, &result)
	if err != nil {
		return []ConfigurationVariable{}, err
	}

	// The API returns variables of upper scopes. Filter them out.
	var filteredVariables []ConfigurationVariable

	for _, variable := range result {
		if scopeId == variable.ScopeId && scope == variable.Scope {
			filteredVariables = append(filteredVariables, variable)
		}
	}

	return filteredVariables, nil
}

func (client *ApiClient) ConfigurationVariableCreate(params ConfigurationVariableCreateParams) (ConfigurationVariable, error) {
	if params.Scope == ScopeDeploymentLog || params.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("must not create variable on scope deployment / deploymentLog")
	}

	organizationId, err := client.OrganizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}

	var result []ConfigurationVariable

	request := map[string]any{
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

	requestInArray := []map[string]any{request}

	if err := client.http.Post("configuration", requestInArray, &result); err != nil {
		return ConfigurationVariable{}, err
	}

	return result[0], nil
}

func getSchema(params ConfigurationVariableCreateParams) map[string]any {
	schema := map[string]any{
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
	return client.http.Delete("configuration/"+id, nil)
}

func (client *ApiClient) ConfigurationVariableUpdate(updateParams ConfigurationVariableUpdateParams) (ConfigurationVariable, error) {
	commonParams := updateParams.CommonParams
	if commonParams.Scope == ScopeDeploymentLog || commonParams.Scope == ScopeDeployment {
		return ConfigurationVariable{}, errors.New("must not create variable on scope deployment / deploymentLog")
	}

	organizationId, err := client.OrganizationId()
	if err != nil {
		return ConfigurationVariable{}, err
	}

	var result []ConfigurationVariable

	request := map[string]any{
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

	requestInArray := []map[string]any{request}

	if err := client.http.Post("/configuration", requestInArray, &result); err != nil {
		return ConfigurationVariable{}, err
	}

	return result[0], nil
}
