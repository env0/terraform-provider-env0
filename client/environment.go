package client

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigurationVariableType int

func (c *ConfigurationVariableType) ReadResourceData(fieldName string, d *schema.ResourceData) error {
	val := d.Get(fieldName).(string)
	intVal, ok := VariableTypes[val]
	if !ok {
		return fmt.Errorf("unknown configuration variable type %s", val)
	}
	*c = intVal

	return nil
}

func (c *ConfigurationVariableType) WriteResourceData(fieldName string, d *schema.ResourceData) error {
	val := *c
	valStr := ""
	if val == 0 {
		valStr = "environment"
	} else if val == 1 {
		valStr = "terraform"
	} else {
		return fmt.Errorf("unknown configuration variable type %d", val)
	}

	return d.Set(fieldName, valStr)
}

const (
	ConfigurationVariableTypeEnvironment ConfigurationVariableType = 0
	ConfigurationVariableTypeTerraform   ConfigurationVariableType = 1
)

type ConfigurationChanges []ConfigurationVariable

var VariableTypes = map[string]ConfigurationVariableType{
	"terraform":   ConfigurationVariableTypeTerraform,
	"environment": ConfigurationVariableTypeEnvironment,
}

type TTL struct {
	Type  TTLType `json:"type"`
	Value string  `json:"value,omitempty"`
}

type TTLType string

const (
	TTLTypeDate     TTLType = "DATE"
	TTlTypeInfinite TTLType = "INFINITE"
)

type DeployRequest struct {
	BlueprintId          string                `json:"blueprintId,omitempty"`
	BlueprintRevision    string                `json:"blueprintRevision,omitempty"`
	BlueprintRepository  string                `json:"blueprintRepository,omitempty"`
	ConfigurationChanges *ConfigurationChanges `json:"configurationChanges,omitempty"`
	TTL                  *TTL                  `json:"ttl,omitempty"`
	EnvName              string                `json:"envName,omitempty"`
	UserRequiresApproval *bool                 `json:"userRequiresApproval,omitempty"`
}

type DeploymentLog struct {
	Id                  string `json:"id"`
	BlueprintId         string `json:"blueprintId"`
	BlueprintRepository string `json:"blueprintRepository"`
	BlueprintRevision   string `json:"blueprintRevision"`
	Status              string `json:"status"`
}

type Environment struct {
	Id                          string        `json:"id"`
	Name                        string        `json:"name"`
	ProjectId                   string        `json:"projectId"`
	WorkspaceName               string        `json:"workspaceName,omitempty"`
	RequiresApproval            *bool         `json:"requiresApproval,omitempty"`
	ContinuousDeployment        *bool         `json:"continuousDeployment,omitempty"`
	PullRequestPlanDeployments  *bool         `json:"pullRequestPlanDeployments,omitempty"`
	AutoDeployOnPathChangesOnly *bool         `json:"autoDeployOnPathChangesOnly,omitempty"`
	AutoDeployByCustomGlob      string        `json:"autoDeployByCustomGlob,omitempty"`
	Status                      string        `json:"status"`
	LifespanEndAt               string        `json:"lifespanEndAt"`
	LatestDeploymentLogId       string        `json:"latestDeploymentLogId"`
	LatestDeploymentLog         DeploymentLog `json:"latestDeploymentLog"`
	IsArchived                  bool          `json:"isArchived"`
	TerragruntWorkingDirectory  string        `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string        `json:"vcsCommandsAlias"`
}

type EnvironmentCreate struct {
	Name                        string                `json:"name"`
	ProjectId                   string                `json:"projectId"`
	DeployRequest               *DeployRequest        `json:"deployRequest"`
	WorkspaceName               string                `json:"workspaceName,omitempty"`
	RequiresApproval            *bool                 `json:"requiresApproval,omitempty"`
	ContinuousDeployment        *bool                 `json:"continuousDeployment,omitempty"`
	PullRequestPlanDeployments  *bool                 `json:"pullRequestPlanDeployments,omitempty"`
	AutoDeployOnPathChangesOnly *bool                 `json:"autoDeployOnPathChangesOnly,omitempty"`
	AutoDeployByCustomGlob      string                `json:"autoDeployByCustomGlob,omitempty"`
	ConfigurationChanges        *ConfigurationChanges `json:"configurationChanges,omitempty"`
	TTL                         *TTL                  `json:"ttl,omitempty"`
	TerragruntWorkingDirectory  string                `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string                `json:"vcsCommandsAlias"`
}

type EnvironmentUpdate struct {
	Name                        string `json:"name,omitempty"`
	RequiresApproval            *bool  `json:"requiresApproval,omitempty"`
	IsArchived                  *bool  `json:"isArchived,omitempty"`
	ContinuousDeployment        *bool  `json:"continuousDeployment,omitempty"`
	PullRequestPlanDeployments  *bool  `json:"pullRequestPlanDeployments,omitempty"`
	AutoDeployOnPathChangesOnly *bool  `json:"autoDeployOnPathChangesOnly,omitempty"`
	AutoDeployByCustomGlob      string `json:"autoDeployByCustomGlob,omitempty"`
	TerragruntWorkingDirectory  string `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string `json:"vcsCommandsAlias,omitempty"`
}

type EnvironmentDeployResponse struct {
	Id string `json:"id"`
}

func (Environment) getEndpoint() string {
	return "/environments"
}

func (client *ApiClient) Environments() ([]Environment, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	return getAll(client, map[string]string{
		"organizationId": organizationId,
	})
}

func (client *ApiClient) ProjectEnvironments(projectId string) ([]Environment, error) {
	return getAll(client, map[string]string{
		"projectId": projectId,
	})
}

func (client *ApiClient) Environment(id string) (Environment, error) {
	var result Environment
	err := client.http.Get("/environments/"+id, nil, &result)
	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentCreate(payload EnvironmentCreate) (Environment, error) {
	var result Environment

	err := client.http.Post("/environments", payload, &result)
	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentDestroy(id string) (EnvironmentDeployResponse, error) {
	var result EnvironmentDeployResponse
	err := client.http.Post("/environments/"+id+"/destroy", nil, &result)
	if err != nil {
		return EnvironmentDeployResponse{}, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentUpdate(id string, payload EnvironmentUpdate) (Environment, error) {
	var result Environment
	err := client.http.Put("/environments/"+id, payload, &result)

	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentUpdateTTL(id string, payload TTL) (Environment, error) {
	var result Environment
	err := client.http.Put("/environments/"+id+"/ttl", payload, &result)

	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (client *ApiClient) EnvironmentDeploy(id string, payload DeployRequest) (EnvironmentDeployResponse, error) {
	var result EnvironmentDeployResponse
	err := client.http.Post("/environments/"+id+"/deployments", payload, &result)

	if err != nil {
		return EnvironmentDeployResponse{}, err
	}
	return result, nil
}

func (client *ApiClient) Deployment(id string) (DeploymentLog, error) {
	var result DeploymentLog
	err := client.http.Get("/environments/deployments/"+id, nil, &result)

	if err != nil {
		return DeploymentLog{}, err
	}
	return result, nil
}
