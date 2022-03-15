package client

//go:generate mockgen -destination=api_client_mock.go -package=client . ApiClientInterface

import (
	"github.com/env0/terraform-provider-env0/client/http"
)

type ApiClient struct {
	http                 http.HttpClientInterface
	cachedOrganizationId string
}

type ApiClientInterface interface {
	ConfigurationVariablesByScope(scope Scope, scopeId string) ([]ConfigurationVariable, error)
	ConfigurationVariablesById(id string) (ConfigurationVariable, error)
	ConfigurationVariableCreate(params ConfigurationVariableCreateParams) (ConfigurationVariable, error)
	ConfigurationVariableUpdate(params ConfigurationVariableUpdateParams) (ConfigurationVariable, error)
	ConfigurationVariableDelete(id string) error
	Organization() (Organization, error)
	organizationId() (string, error)
	Policy(projectId string) (Policy, error)
	PolicyUpdate(payload PolicyUpdatePayload) (Policy, error)
	Projects() ([]Project, error)
	Project(id string) (Project, error)
	ProjectCreate(payload ProjectCreatePayload) (Project, error)
	ProjectUpdate(id string, payload ProjectCreatePayload) (Project, error)
	ProjectDelete(id string) error
	Template(id string) (Template, error)
	Templates() ([]Template, error)
	TemplateCreate(payload TemplateCreatePayload) (Template, error)
	TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error)
	TemplateDelete(id string) error
	AssignTemplateToProject(id string, payload TemplateAssignmentToProjectPayload) (Template, error)
	RemoveTemplateFromProject(templateId string, projectId string) error
	SshKeys() ([]SshKey, error)
	SshKeyCreate(payload SshKeyCreatePayload) (SshKey, error)
	SshKeyDelete(id string) error
	CloudCredentials(id string) (ApiKey, error)
	CloudCredentialsList() ([]ApiKey, error)
	AwsCredentialsCreate(request AwsCredentialsCreatePayload) (ApiKey, error)
	AzureCredentialsCreate(request AzureCredentialsCreatePayload) (ApiKey, error)
	CloudCredentialsDelete(id string) error
	GoogleCostCredentialsCreate(request GoogleCostCredentialsCreatePayload) (ApiKey, error)
	GcpCredentialsCreate(request GcpCredentialsCreatePayload) (ApiKey, error)
	AzureCredentialsCreate(request AzureCredentialsCreatePayload) (ApiKey, error)
	AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error)
	RemoveCloudCredentialsFromProject(projectId string, credentialId string) error
	CloudCredentialIdsInProject(projectId string) ([]string, error)
	AssignCostCredentialsToProject(projectId string, credentialId string) (CostCredentialProjectAssignment, error)
	CostCredentialIdsInProject(projectId string) ([]CostCredentialProjectAssignment, error)
	RemoveCostCredentialsFromProject(projectId string, credentialId string) error
	Team(id string) (Team, error)
	Teams() ([]Team, error)
	TeamCreate(payload TeamCreatePayload) (Team, error)
	TeamUpdate(id string, payload TeamUpdatePayload) (Team, error)
	TeamDelete(id string) error
	TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error)
	TeamProjectAssignmentDelete(assignmentId string) error
	TeamProjectAssignments(projectId string) ([]TeamProjectAssignment, error)
	Environments() ([]Environment, error)
	Environment(id string) (Environment, error)
	EnvironmentCreate(payload EnvironmentCreate) (Environment, error)
	EnvironmentDestroy(id string) (Environment, error)
	EnvironmentUpdate(id string, payload EnvironmentUpdate) (Environment, error)
	EnvironmentDeploy(id string, payload DeployRequest) (EnvironmentDeployResponse, error)
	EnvironmentUpdateTTL(id string, payload TTL) (Environment, error)
	EnvironmentScheduling(environmentId string) (EnvironmentScheduling, error)
	EnvironmentSchedulingUpdate(environmentId string, payload EnvironmentScheduling) (EnvironmentScheduling, error)
	EnvironmentSchedulingDelete(environmentId string) error
	WorkflowTrigger(environmentId string) ([]WorkflowTrigger, error)
	WorkflowTriggerUpsert(environmentId string, request WorkflowTriggerUpsertPayload) ([]WorkflowTrigger, error)
	EnvironmentDriftDetection(environmentId string) (EnvironmentSchedulingExpression, error)
	EnvironmentUpdateDriftDetection(environmentId string, payload EnvironmentSchedulingExpression) (EnvironmentSchedulingExpression, error)
	EnvironmentStopDriftDetection(environmentId string) error
}

func NewApiClient(client http.HttpClientInterface) ApiClientInterface {
	return &ApiClient{
		http:                 client,
		cachedOrganizationId: "",
	}
}
