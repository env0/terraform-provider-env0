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
}

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

type TemplateRetryOn struct {
	Times      int    `json:"times,omitempty"`
	ErrorRegex string `json:"errorRegex,omitempty"`
}

type TemplateRetry struct {
	OnDeploy  *TemplateRetryOn `json:"onDeploy,omitempty"`
	OnDestroy *TemplateRetryOn `json:"onDestroy,omitempty"`
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
	Retry                *TemplateRetry   `json:"retry,omitempty"`
	SshKeys              []TemplateSshKey `json:"sshKeys,omitempty"`
	Type                 TemplateType     `json:"type"`
	Description          string           `json:"description,omitempty"`
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path,omitempty"`
	IsGitLab             bool             `json:"isGitLab"`
	TokenName            string           `json:"tokenName"`
	TokenId              string           `json:"tokenId"`
	GithubInstallationId int              `json:"githubInstallationId"`
	Revision             string           `json:"revision"`
	ProjectIds           []string         `json:"projectIds,omitempty"`
	OrganizationId       string           `json:"organizationId"`
	TerraformVersion     string           `json:"terraformVersion"`
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
