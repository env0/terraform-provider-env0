package api

type ConfigurationVariableSchema struct {
	Type string   `json:"string"`
	Enum []string `json:"enum"`
}

type ConfigurationVariable struct {
	Value          string                      `json:"value"`
	OrganizationId string                      `json:"organizationId"`
	UserId         string                      `json:"userId"`
	IsSensitive    bool                        `json:"isSensitive"`
	Scope          string                      `json:"scope"`
	Id             string                      `json:"id"`
	Name           string                      `json:"name"`
	Type           int64                       `json:"type"`
	Schema         ConfigurationVariableSchema `json:"schema"`
}

type Scope string

const (
	ScopeGlobal        Scope = "GLOBAL"
	ScopeTemplate      Scope = "BLUEPRINT"
	ScopeProject       Scope = "PROJECT"
	ScopeEnvironment   Scope = "ENVIRONMENT"
	ScopeDeployment    Scope = "DEPLOYMENT"
	ScopeDeploymentLog Scope = "DEPLOYMENT_LOG"
)

type ConfigurationVariableType int

const (
	ConfigurationVariableTypeEnvironment ConfigurationVariableType = 0
	ConfigurationVariableTypeTerraform   ConfigurationVariableType = 1
)
