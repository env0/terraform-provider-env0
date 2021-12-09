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
	ConfigurationVariables(scope Scope, scopeId string) ([]ConfigurationVariable, error)
	ConfigurationVariableCreate(name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string, description string) (ConfigurationVariable, error)
	ConfigurationVariableUpdate(id string, name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string, description string) (ConfigurationVariable, error)
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
	AwsCredentials(id string) (ApiKey, error)
	AwsCredentialsList() ([]ApiKey, error)
	AwsCredentialsCreate(request AwsCredentialsCreatePayload) (ApiKey, error)
	AwsCredentialsDelete(id string) error
	AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error)
	RemoveCloudCredentialsFromProject(projectId string, credentialId string) error
	CloudCredentialIdsInProject(projectId string) ([]string, error)
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
}

func NewApiClient(client http.HttpClientInterface) ApiClientInterface {
	return &ApiClient{
		http:                 client,
		cachedOrganizationId: "",
	}
}
