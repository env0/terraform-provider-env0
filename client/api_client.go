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
	ConfigurationVariableCreate(name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error)
	ConfigurationVariableUpdate(id string, name string, value string, isSensitive bool, scope Scope, scopeId string, type_ ConfigurationVariableType, enumValues []string) (ConfigurationVariable, error)
	ConfigurationVariableDelete(id string) error
	Organization() (Organization, error)
	organizationId() (string, error)
	Projects() ([]Project, error)
	Project(id string) (Project, error)
	ProjectCreate(name string, description string) (Project, error)
	ProjectUpdate(id string, payload UpdateProjectPayload) (Project, error)
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
	RemoveCloudCredentialsFromProject(credentialId string, projectId string) error
	CloudCredentialProjectAssignments(projectId string) ([]CloudCredentialsProjectAssignment, error)
}

func NewApiClient(client http.HttpClientInterface) ApiClientInterface {
	return &ApiClient{
		http:                 client,
		cachedOrganizationId: "",
	}
}
