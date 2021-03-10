package env0apiclient

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
	ScopeBlueprint     Scope = "BLUEPRINT"
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

type TemplateRetry struct {
	OnDeploy  string `json:"onDeploy"`
	OnDestroy string `json:"onDestroy"`
}

type TemplateType string

const (
	TemplateTypeTerraform  TemplateType = "terraform"
	TemplateTypeTerragrunt TemplateType = "terragrunt"
)

type TemplateCreatePayload struct {
	Retry                TemplateRetry `json:"retry"`
	SshKeys              []string      `json:"sshKeys"`
	Type                 TemplateType  `json:"type"`
	Description          string        `json:"description"`
	Name                 string        `json:"name"`
	Repository           string        `json:"repository"`
	IsGitLab             bool          `json:"isGitLab"`
	TokenName            string        `json:"tokenName"`
	TokenId              string        `json:"tokenId"`
	GithubInstallationId string        `json:"githubInstallationId"`
	Revision             string        `json:"revision"`
	ProjectIds           []string      `json:"projectIds"`
	OrganizationId       string        `json:"organizationId"`
}

type Template struct {
	Author         User          `json:"author"`
	AuthorId       string        `json:"authorId"`
	CreatedAt      string        `json:"createdAt"`
	Href           string        `json:"href"`
	Id             string        `json:"id"`
	Name           string        `json:"name"`
	OrganizationId string        `json:"organizationId"`
	Path           string        `json:"path"`
	ProjectId      string        `json:"projectId"`
	ProjectIds     []string      `json:"projectIds"`
	Repository     string        `json:"repository"`
	Retry          TemplateRetry `json:"retry"`
	SshKeys        []string      `json:"sshKeys"`
	Type           string        `json:"type"`
	UpdatedAt      string        `json:"updatedAt"`
}
