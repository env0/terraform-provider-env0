package client

//go:generate mockgen -destination=api_client_mock.go -package=client . ApiClientInterface

import (
	"github.com/env0/terraform-provider-env0/client/http"
)

type ApiClient struct {
	http                  http.HttpClientInterface
	cachedOrganizationId  string
	defaultOrganizationId string
}

type ApiClientInterface interface {
	ConfigurationVariablesByScope(scope Scope, scopeId string) ([]ConfigurationVariable, error)
	ConfigurationVariablesById(id string) (ConfigurationVariable, error)
	ConfigurationVariableCreate(params ConfigurationVariableCreateParams) (ConfigurationVariable, error)
	ConfigurationVariableUpdate(params ConfigurationVariableUpdateParams) (ConfigurationVariable, error)
	ConfigurationVariableDelete(id string) error
	Organization() (Organization, error)
	OrganizationId() (string, error)
	OrganizationPolicyUpdate(OrganizationPolicyUpdatePayload) (*Organization, error)
	OrganizationUserUpdateRole(userId string, roleId string) error
	Policy(projectId string) (Policy, error)
	PolicyUpdate(payload PolicyUpdatePayload) (Policy, error)
	Projects() ([]Project, error)
	Project(id string) (Project, error)
	ProjectCreate(payload ProjectCreatePayload) (Project, error)
	ProjectUpdate(id string, payload ProjectUpdatePayload) (Project, error)
	ProjectDelete(id string) error
	ProjectMove(id string, targetProjectId string) error
	Template(id string) (Template, error)
	Templates() ([]Template, error)
	TemplateCreate(payload TemplateCreatePayload) (Template, error)
	TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error)
	TemplateDelete(id string) error
	AssignTemplateToProject(id string, payload TemplateAssignmentToProjectPayload) (Template, error)
	RemoveTemplateFromProject(templateId string, projectId string) error
	VariablesFromRepository(payload *VariablesFromRepositoryPayload) ([]ConfigurationVariable, error)
	SshKeys() ([]SshKey, error)
	SshKeyCreate(payload SshKeyCreatePayload) (*SshKey, error)
	SshKeyDelete(id string) error
	CredentialsCreate(request CredentialCreatePayload) (Credentials, error)
	CloudCredentials(id string) (Credentials, error)
	CloudCredentialsList() ([]Credentials, error)
	CloudCredentialsDelete(id string) error
	AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error)
	RemoveCloudCredentialsFromProject(projectId string, credentialId string) error
	CloudCredentialIdsInProject(projectId string) ([]string, error)
	AssignCostCredentialsToProject(projectId string, credentialId string) (CostCredentialProjectAssignment, error)
	CostCredentialIdsInProject(projectId string) ([]CostCredentialProjectAssignment, error)
	RemoveCostCredentialsFromProject(projectId string, credentialId string) error
	Team(id string) (Team, error)
	Teams() ([]Team, error)
	TeamsByName(name string) ([]Team, error)
	TeamCreate(payload TeamCreatePayload) (Team, error)
	TeamUpdate(id string, payload TeamUpdatePayload) (Team, error)
	TeamDelete(id string) error
	TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error)
	TeamProjectAssignmentDelete(assignmentId string) error
	TeamProjectAssignments(projectId string) ([]TeamProjectAssignment, error)
	Environments() ([]Environment, error)
	ProjectEnvironments(projectId string) ([]Environment, error)
	Environment(id string) (Environment, error)
	EnvironmentCreate(payload EnvironmentCreate) (Environment, error)
	EnvironmentCreateWithoutTemplate(payload EnvironmentCreateWithoutTemplate) (Environment, error)
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
	Notifications() ([]Notification, error)
	NotificationCreate(payload NotificationCreatePayload) (*Notification, error)
	NotificationDelete(id string) error
	NotificationUpdate(id string, payload NotificationUpdatePayload) (*Notification, error)
	NotificationProjectAssignments(projectId string) ([]NotificationProjectAssignment, error)
	NotificationProjectAssignmentUpdate(projectId string, endpointId string, payload NotificationProjectAssignmentUpdatePayload) (*NotificationProjectAssignment, error)
	ModuleCreate(payload ModuleCreatePayload) (*Module, error)
	Module(id string) (*Module, error)
	ModuleDelete(id string) error
	ModuleUpdate(id string, payload ModuleUpdatePayload) (*Module, error)
	Modules() ([]Module, error)
	GitToken(id string) (*GitToken, error)
	GitTokens() ([]GitToken, error)
	GitTokenCreate(payload GitTokenCreatePayload) (*GitToken, error)
	GitTokenDelete(id string) error
	ApiKeyCreate(payload ApiKeyCreatePayload) (*ApiKey, error)
	ApiKeyDelete(id string) error
	ApiKeys() ([]ApiKey, error)
	AssignAgentsToProjects(payload AssignProjectsAgentsAssignmentsPayload) (*ProjectsAgentsAssignments, error)
	ProjectsAgentsAssignments() (*ProjectsAgentsAssignments, error)
	Agents() ([]Agent, error)
	AgentValues(id string) (string, error)
	Users() ([]OrganizationUser, error)
	AssignUserToProject(projectId string, payload *AssignUserToProjectPayload) (*UserProjectAssignment, error)
	RemoveUserFromProject(projectId string, userId string) error
	UserProjectAssignments(projectId string) ([]UserProjectAssignment, error)
	UpdateUserProjectAssignment(projectId string, userId string, payload *UpdateUserProjectAssignmentPayload) (*UserProjectAssignment, error)
	RoleCreate(payload RoleCreatePayload) (*Role, error)
	Role(id string) (*Role, error)
	RoleDelete(id string) error
	RoleUpdate(id string, payload RoleUpdatePayload) (*Role, error)
	Roles() ([]Role, error)
	CustomFlow(id string) (*CustomFlow, error)
	CustomFlows(name string) ([]CustomFlow, error)
	CustomFlowCreate(payload CustomFlowCreatePayload) (*CustomFlow, error)
	CustomFlowUpdate(id string, payload CustomFlowCreatePayload) (*CustomFlow, error)
	CustomFlowDelete(id string) error
	CustomFlowAssign(assignments []CustomFlowAssignment) error
	CustomFlowUnassign(assignments []CustomFlowAssignment) error
	CustomFlowGetAssignments(assignments []CustomFlowAssignment) ([]CustomFlowAssignment, error)
	SubscribeWorkflowTrigger(environmentId string, payload WorkflowTriggerEnvironments) error
	UnsubscribeWorkflowTrigger(environmentId string, payload WorkflowTriggerEnvironments) error
	RemoteStateAccessConfiguration(environmentId string) (*RemoteStateAccessConfiguration, error)
	RemoteStateAccessConfigurationCreate(environmentId string, payload RemoteStateAccessConfigurationCreate) (*RemoteStateAccessConfiguration, error)
	RemoteStateAccessConfigurationDelete(environmentId string) error
	GpgKeyCreate(payload *GpgKeyCreatePayload) (*GpgKey, error)
	GpgKeyDelete(id string) error
	GpgKeys() ([]GpgKey, error)
	ProviderCreate(payload ProviderCreatePayload) (*Provider, error)
	Provider(providerId string) (*Provider, error)
	ProviderDelete(providerId string) error
	ProviderUpdate(providerId string, payload ProviderUpdatePayload) (*Provider, error)
	Providers() ([]Provider, error)
}

func NewApiClient(client http.HttpClientInterface, defaultOrganizationId string) ApiClientInterface {
	return &ApiClient{
		http:                  client,
		cachedOrganizationId:  "",
		defaultOrganizationId: defaultOrganizationId,
	}
}
