package client

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ENVIRONMENT = "environment"

type ConfigurationVariableType int

func (c *ConfigurationVariableType) ReadResourceData(fieldName string, d *schema.ResourceData) error {
	val := d.Get(fieldName).(string)

	intVal, err := GetConfigurationVariableType(val)
	if err != nil {
		return err
	}

	*c = intVal

	return nil
}

func (c *ConfigurationVariableType) WriteResourceData(fieldName string, d *schema.ResourceData) error {
	var valStr string

	switch val := *c; val {
	case 0:
		valStr = ENVIRONMENT
	case 1:
		valStr = TERRAFORM
	default:
		return fmt.Errorf("unknown configuration variable type %d", val)
	}

	return d.Set(fieldName, valStr)
}

const (
	ConfigurationVariableTypeEnvironment ConfigurationVariableType = 0
	ConfigurationVariableTypeTerraform   ConfigurationVariableType = 1
)

type ConfigurationChanges []ConfigurationVariable

type TTL struct {
	Type  TTLType `json:"type"`
	Value string  `json:"value,omitempty"`
}

type TTLType string

const (
	TTLTypeDate     TTLType = "DATE"
	TTlTypeInfinite TTLType = "INFINITE"
)

type SubEnvironment struct {
	Revision             string               `json:"revision,omitempty"`
	Workspace            string               `json:"workspaceName,omitempty"`
	ConfigurationChanges ConfigurationChanges `json:"configurationChanges"`
	UserRequiresApproval bool                 `json:"userRequiresApproval"`
}

type DeployRequest struct {
	BlueprintId             string                    `json:"blueprintId,omitempty"`
	BlueprintRevision       string                    `json:"blueprintRevision,omitempty"`
	BlueprintRepository     string                    `json:"blueprintRepository,omitempty"`
	ConfigurationChanges    *ConfigurationChanges     `json:"configurationChanges,omitempty"`
	TTL                     *TTL                      `json:"ttl,omitempty"`
	EnvName                 string                    `json:"envName,omitempty"`
	UserRequiresApproval    *bool                     `json:"userRequiresApproval,omitempty"`
	SubEnvironments         map[string]SubEnvironment `json:"subEnvironments,omitempty"`
	ConfigurationSetChanges *ConfigurationSetChanges  `json:"configurationSetChanges,omitempty" tfschema:"-"`
}

type WorkflowSubEnvironment struct {
	EnvironmentId string `json:"environmentId"`
}

type WorkflowFile struct {
	Environments map[string]WorkflowSubEnvironment `json:"environments"`
}

type DeploymentLog struct {
	Id                  string          `json:"id"`
	BlueprintId         string          `json:"blueprintId"`
	BlueprintRepository string          `json:"blueprintRepository"`
	BlueprintRevision   string          `json:"blueprintRevision"`
	Output              json.RawMessage `json:"output,omitempty"`
	Error               json.RawMessage `json:"error,omitempty"`
	Type                string          `json:"type"`
	Status              string          `json:"status"`
	WorkflowFile        *WorkflowFile   `json:"workflowFile,omitempty" tfschema:"-"`
}

type DriftDetectionRequest struct {
	Enabled bool   `json:"enabled"`
	Cron    string `json:"cron"`
}

type ConfigurationSetChanges struct {
	Assign   []string `json:"assign,omitempty"`
	Unassign []string `json:"unassign,omitempty"`
}

type Environment struct {
	Id                          string        `json:"id"`
	Name                        string        `json:"name"`
	ProjectId                   string        `json:"projectId"`
	WorkspaceName               string        `json:"workspaceName,omitempty"               tfschema:"workspace"`
	RequiresApproval            *bool         `json:"requiresApproval,omitempty"            tfschema:"-"`
	ContinuousDeployment        *bool         `json:"continuousDeployment,omitempty"        tfschema:"deploy_on_push,omitempty"`
	PullRequestPlanDeployments  *bool         `json:"pullRequestPlanDeployments,omitempty"  tfschema:"run_plan_on_pull_requests,omitempty"`
	AutoDeployOnPathChangesOnly *bool         `json:"autoDeployOnPathChangesOnly,omitempty" tfschema:",omitempty"`
	AutoDeployByCustomGlob      string        `json:"autoDeployByCustomGlob,omitempty"`
	Status                      string        `json:"status"`
	LifespanEndAt               string        `json:"lifespanEndAt"                         tfschema:"ttl,omitempty"`
	LatestDeploymentLogId       string        `json:"latestDeploymentLogId"                 tfschema:"deployment_id"`
	LatestDeploymentLog         DeploymentLog `json:"latestDeploymentLog"`
	TerragruntWorkingDirectory  string        `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string        `json:"vcsCommandsAlias"`
	VcsPrCommentsEnabled        bool          `json:"vcsPrCommentsEnabled"                  tfschema:"-"`
	BlueprintId                 string        `json:"blueprintId"                           tfschema:"-"`
	IsRemoteBackend             *bool         `json:"isRemoteBackend"                       tfschema:"-"`
	IsArchived                  *bool         `json:"isArchived"                            tfschema:"-"`
	IsRemoteApplyEnabled        bool          `json:"isRemoteApplyEnabled"`
	K8sNamespace                string        `json:"k8sNamespace"`
	IsSingleUseBlueprint        bool          `json:"isSingleUseBlueprint"                  tfschema:"-"`
}

type EnvironmentCreate struct {
	Name                        string                   `json:"name"`
	ProjectId                   string                   `json:"projectId"`
	DeployRequest               *DeployRequest           `json:"deployRequest"                         tfschema:"-"`
	WorkspaceName               string                   `json:"workspaceName,omitempty"               tfschema:"workspace"`
	RequiresApproval            *bool                    `json:"requiresApproval,omitempty"            tfschema:"-"`
	ContinuousDeployment        *bool                    `json:"continuousDeployment,omitempty"        tfschema:"-"`
	PullRequestPlanDeployments  *bool                    `json:"pullRequestPlanDeployments,omitempty"  tfschema:"-"`
	AutoDeployOnPathChangesOnly *bool                    `json:"autoDeployOnPathChangesOnly,omitempty" tfchema:"-"`
	AutoDeployByCustomGlob      string                   `json:"autoDeployByCustomGlob"`
	ConfigurationChanges        *ConfigurationChanges    `json:"configurationChanges,omitempty"        tfschema:"-"`
	TTL                         *TTL                     `json:"ttl,omitempty"                         tfschema:"-"`
	TerragruntWorkingDirectory  string                   `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string                   `json:"vcsCommandsAlias"`
	VcsPrCommentsEnabled        *bool                    `json:"vcsPrCommentsEnabled,omitempty"        tfschema:"-"`
	IsRemoteBackend             *bool                    `json:"isRemoteBackend,omitempty"             tfschema:"-"`
	Type                        string                   `json:"type,omitempty"`
	DriftDetectionRequest       *DriftDetectionRequest   `json:"driftDetectionRequest,omitempty"       tfschema:"-"`
	PreventAutoDeploy           *bool                    `json:"preventAutoDeploy,omitempty"           tfschema:"-"`
	K8sNamespace                string                   `json:"k8sNamespace,omitempty"`
	ConfigurationSetChanges     *ConfigurationSetChanges `json:"configurationSetChanges,omitempty"     tfschema:"-"`
	IsRemoteApplyEnabled        bool                     `json:"isRemoteApplyEnabled"`
}

// When converted to JSON needs to be flattened. See custom MarshalJSON below.
type EnvironmentCreateWithoutTemplate struct {
	EnvironmentCreate EnvironmentCreate
	TemplateCreate    TemplateCreatePayload
}

type EnvironmentMoveRequest struct {
	ProjectId string `json:"projectId"`
}

type EnvironmentDestroyResponse struct {
	Id string `json:"id"`
}

func GetConfigurationVariableType(variableType string) (ConfigurationVariableType, error) {
	switch variableType {
	case "terraform":
		return ConfigurationVariableTypeTerraform, nil
	case "environment":
		return ConfigurationVariableTypeEnvironment, nil
	default:
		return 0, fmt.Errorf("unknown configuration variable type %s", variableType)
	}
}

// The custom marshalJSON is used to return a flat JSON.
func (create EnvironmentCreateWithoutTemplate) MarshalJSON() ([]byte, error) {
	// 1. Marshal to JSON both structs.
	ecb, err := json.Marshal(&create.EnvironmentCreate)
	if err != nil {
		return nil, err
	}

	tcb, err := json.Marshal(&create.TemplateCreate)
	if err != nil {
		return nil, err
	}

	// 2. Unmarshal both JSON byte arrays to two maps.
	var ecm, tcm map[string]any

	if err := json.Unmarshal(ecb, &ecm); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tcb, &tcm); err != nil {
		return nil, err
	}

	// 3. Merged the maps.
	maps.Copy(tcm, ecm)

	// 4. Marshal the merged map back to JSON.
	return json.Marshal(tcm)
}

type EnvironmentUpdate struct {
	Name                        string `json:"name,omitempty"`
	AutoDeployByCustomGlob      string `json:"autoDeployByCustomGlob"`
	TerragruntWorkingDirectory  string `json:"terragruntWorkingDirectory,omitempty"`
	VcsCommandsAlias            string `json:"vcsCommandsAlias,omitempty"`
	VcsPrCommentsEnabled        *bool  `json:"vcsPrCommentsEnabled,omitempty"        tfschema:"-"`
	RequiresApproval            *bool  `json:"requiresApproval,omitempty"            tfschema:"-"`
	ContinuousDeployment        *bool  `json:"continuousDeployment,omitempty"        tfschema:"-"`
	PullRequestPlanDeployments  *bool  `json:"pullRequestPlanDeployments,omitempty"  tfschema:"-"`
	AutoDeployOnPathChangesOnly *bool  `json:"autoDeployOnPathChangesOnly,omitempty" tfschema:"-"`
	IsRemoteBackend             *bool  `json:"isRemoteBackend,omitempty"             tfschema:"-"`
	IsArchived                  *bool  `json:"isArchived,omitempty"                  tfschema:"-"`
	IsRemoteApplyEnabled        bool   `json:"isRemoteApplyEnabled"`
}

type EnvironmentDeployResponse struct {
	Id string `json:"id"`
}

func (Environment) getEndpoint() string {
	return "/environments"
}

func (client *ApiClient) EnvironmentsByName(name string) ([]Environment, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	return getAll(client, map[string]string{
		"organizationId": organizationId,
		"name":           name,
	})
}

func (client *ApiClient) ProjectEnvironments(projectId string) ([]Environment, error) {
	return getAll(client, map[string]string{
		"projectId": projectId,
	})
}

func (client *ApiClient) OrganizationEnvironments(organizationId string) ([]Environment, error) {
	return getAll(client, map[string]string{
		"organizationId": organizationId,
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

func (client *ApiClient) EnvironmentDeploymentLog(id string) (*DeploymentLog, error) {
	var result DeploymentLog

	err := client.http.Get("/environments/deployments/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) EnvironmentCreate(payload EnvironmentCreate) (Environment, error) {
	var result Environment

	err := client.http.Post("/environments", payload, &result)
	if err != nil {
		return Environment{}, err
	}

	return result, nil
}

func (client *ApiClient) EnvironmentCreateWithoutTemplate(payload EnvironmentCreateWithoutTemplate) (Environment, error) {
	var result Environment

	organizationId, err := client.OrganizationId()
	if err != nil {
		return result, err
	}

	payload.TemplateCreate.OrganizationId = organizationId

	if err := client.http.Post("/environments/without-template", payload, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (client *ApiClient) EnvironmentDestroy(id string) (*EnvironmentDestroyResponse, error) {
	var result EnvironmentDestroyResponse

	err := client.http.Post("/environments/"+id+"/destroy", nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) EnvironmentUpdate(id string, payload EnvironmentUpdate) (Environment, error) {
	var result Environment

	err := client.http.Put("/environments/"+id, payload, &result)
	if err != nil {
		return Environment{}, err
	}

	return result, nil
}

func (client *ApiClient) EnvironmentMarkAsArchived(id string) error {
	payload := struct {
		IsArchived bool `json:"isArchived"`
	}{true}

	return client.http.Put("/environments/"+id, &payload, nil)
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

func (client *ApiClient) EnvironmentMove(id string, projectId string) error {
	payload := &EnvironmentMoveRequest{
		ProjectId: projectId,
	}

	return client.http.Post("/environments/"+id+"/move", payload, nil)
}
