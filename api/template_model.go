package api

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
}

type Template struct {
	Author         User             `json:"author"`
	AuthorId       string           `json:"authorId"`
	CreatedAt      string           `json:"createdAt"`
	Href           string           `json:"href"`
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	OrganizationId string           `json:"organizationId"`
	Path           string           `json:"path"`
	Revision       string           `json:"revision"`
	ProjectId      string           `json:"projectId"`
	ProjectIds     []string         `json:"projectIds"`
	Repository     string           `json:"repository"`
	Retry          TemplateRetry    `json:"retry"`
	SshKeys        []TemplateSshKey `json:"sshKeys"`
	Type           string           `json:"type"`
	UpdatedAt      string           `json:"updatedAt"`
}
