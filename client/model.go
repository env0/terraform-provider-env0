package client

type Organization struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CreatedBy    string `json:"createdBy"`
	Role         string `json:"role"`
	IsSelfHosted bool   `json:"isSelfHosted"`
}

type User struct {
	CreatedAt   string                 `json:"created_at"`
	Email       string                 `json:"email"`
	FamilyName  string                 `json:"family_name"`
	GivenName   string                 `json:"given_name"`
	Name        string                 `json:"name"`
	Picture     string                 `json:"picture"`
	UserId      string                 `json:"user_id"`
	LastLogin   string                 `json:"last_login"`
	AppMetadata map[string]interface{} `json:"app_metadata"`
}

type Project struct {
	IsArchived     bool   `json:"isArchived"`
	OrganizationId string `json:"organizationId"`
	UpdatedAt      string `json:"updatedAt"`
	CreatedAt      string `json:"createdAt"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	CreatedBy      string `json:"createdBy"`
	Role           string `json:"role"`
	CreatedByUser  User   `json:"createdByUser"`
	Description    string `json:"description"`
}

type ProjectCreatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ConfigurationVariableSchema struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}

type ConfigurationVariable struct {
	ScopeId        string                      `json:"scopeId"`
	Value          string                      `json:"value"`
	OrganizationId string                      `json:"organizationId"`
	UserId         string                      `json:"userId"`
	IsSensitive    bool                        `json:"isSensitive"`
	Scope          Scope                       `json:"scope"`
	Id             string                      `json:"id"`
	Name           string                      `json:"name"`
	Type           ConfigurationVariableType   `json:"type"`
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

type TemplateRetryOn struct {
	Times      int    `json:"times,omitempty"`
	ErrorRegex string `json:"errorRegex"`
}

type TemplateRetry struct {
	OnDeploy  *TemplateRetryOn `json:"onDeploy"`
	OnDestroy *TemplateRetryOn `json:"onDestroy"`
}

type TemplateType string

const (
	TemplateTypeTerraform  TemplateType = "terraform"
	TemplateTypeTerragrunt TemplateType = "terragrunt"
)

type TemplateSshKey struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TemplateCreatePayload struct {
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys,omitempty"`
	Type                 TemplateType     `json:"type"`
	Description          string           `json:"description"`
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path"`
	IsGitLab             bool             `json:"isGitLab"`
	TokenName            string           `json:"tokenName"`
	TokenId              string           `json:"tokenId"`
	GithubInstallationId int              `json:"githubInstallationId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	Revision             string           `json:"revision"`
	OrganizationId       string           `json:"organizationId"`
	TerraformVersion     string           `json:"terraformVersion"`
}

type TemplateAssignmentToProjectPayload struct {
	ProjectId string `json:"projectId"`
}

type TemplateAssignmentToProject struct {
	Id         string `json:"id"`
	TemplateId string `json:"templateId"`
	ProjectId  string `json:"projectId"`
}

type CloudCredentialIdsInProjectResponse struct {
	CredentialIds []string `json:"credentialIds"`
}

type CloudCredentialsProjectAssignmentPatchPayload struct {
	CredentialIds []string `json:"credentialIds"`
}

type CloudCredentialsProjectAssignment struct {
	Id           string `json:"id"`
	CredentialId string `json:"credentialId"`
	ProjectId    string `json:"projectId"`
}

type Template struct {
	Author               User             `json:"author"`
	AuthorId             string           `json:"authorId"`
	CreatedAt            string           `json:"createdAt"`
	Href                 string           `json:"href"`
	Id                   string           `json:"id"`
	Name                 string           `json:"name"`
	Description          string           `json:"description"`
	OrganizationId       string           `json:"organizationId"`
	Path                 string           `json:"path"`
	Revision             string           `json:"revision"`
	ProjectId            string           `json:"projectId"`
	ProjectIds           []string         `json:"projectIds"`
	Repository           string           `json:"repository"`
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys"`
	Type                 string           `json:"type"`
	GithubInstallationId int              `json:"githubInstallationId"`
	TokenId              string           `json:"tokenId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	UpdatedAt            string           `json:"updatedAt"`
	TerraformVersion     string           `json:"terraformVersion"`
}

type SshKey struct {
	User           User   `json:"user"`
	UserId         string `json:"userId"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
	Value          string `json:"value"`
}

type SshKeyCreatePayload struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
	Value          string `json:"value"`
}

type AwsCredentialsCreatePayload struct {
	Name           string                     `json:"name"`
	OrganizationId string                     `json:"organizationId"`
	Type           string                     `json:"type"`
	Value          AwsCredentialsValuePayload `json:"value"`
}

type AwsCredentialsValuePayload struct {
	RoleArn    string `json:"roleArn"`
	ExternalId string `json:"externalId"`
}

type ApiKey struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}
