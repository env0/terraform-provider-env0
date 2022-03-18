package client

import "encoding/json"

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
}

type DeployRequest struct {
	BlueprintId          string                `json:"blueprintId,omitempty"`
	BlueprintRevision    string                `json:"blueprintRevision,omitempty"`
	BlueprintRepository  string                `json:"blueprintRepository,omitempty"`
	ConfigurationChanges *ConfigurationChanges `json:"configurationChanges,omitempty"`
	TTL                  *TTL                  `json:"ttl,omitempty"`
	EnvName              string                `json:"envName,omitempty"`
	UserRequiresApproval *bool                 `json:"userRequiresApproval,omitempty"`
}

type GitUserData struct {
	GitUser      string `json:"gitUser,omitempty"`
	GitAvatarUrl string `json:"gitAvatarUrl,omitempty"`
	PrNumber     string `json:"prNumber,omitempty"`
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

type EnvironmentUpdate struct {
	Name                        string `json:"name,omitempty"`
	RequiresApproval            *bool  `json:"requiresApproval,omitempty"`
	IsArchived                  *bool  `json:"isArchived,omitempty"`
	ContinuousDeployment        *bool  `json:"continuousDeployment,omitempty"`
	PullRequestPlanDeployments  *bool  `json:"pullRequestPlanDeployments,omitempty"`
	AutoDeployOnPathChangesOnly *bool  `json:"autoDeployOnPathChangesOnly,omitempty"`
	AutoDeployByCustomGlob      string `json:"autoDeployByCustomGlob,omitempty"`
	TerragruntWorkingDirectory  string `json:"terragruntWorkingDirectory,omitempty"`
}

type EnvironmentDeployResponse struct {
	Id string `json:"id"`
}

type ConfigurationVariableSchema struct {
	Type   string   `json:"type"`
	Enum   []string `json:"enum"`
	Format Format   `json:"format,omitempty"`
}

type ConfigurationVariable struct {
	ScopeId        string                       `json:"scopeId,omitempty"`
	Value          string                       `json:"value"`
	OrganizationId string                       `json:"organizationId,omitempty"`
	UserId         string                       `json:"userId,omitempty"`
	IsSensitive    *bool                        `json:"isSensitive,omitempty"`
	Scope          Scope                        `json:"scope,omitempty"`
	Id             string                       `json:"id,omitempty"`
	Name           string                       `json:"name"`
	Description    string                       `json:"description,omitempty"`
	Type           *ConfigurationVariableType   `json:"type,omitempty"`
	Schema         *ConfigurationVariableSchema `json:"schema,omitempty"`
	ToDelete       *bool                        `json:"toDelete,omitempty"`
	IsReadonly     *bool                        `json:"isReadonly,omitempty"`
	IsRequired     *bool                        `json:"isRequired,omitempty"`
}

type ConfigurationVariableCreateParams struct {
	Name        string
	Value       string
	IsSensitive bool
	Scope       Scope
	ScopeId     string
	Type        ConfigurationVariableType
	EnumValues  []string
	Description string
	Format      Format
	IsReadonly  bool
	IsRequired  bool
}

type ConfigurationVariableUpdateParams struct {
	CommonParams ConfigurationVariableCreateParams
	Id           string
}

type Format string

const (
	Text Format = ""
	HCL  Format = "HCL"
	JSON Format = "JSON"
)

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

type ConfigurationChanges []ConfigurationVariable

const (
	ConfigurationVariableTypeEnvironment ConfigurationVariableType = 0
	ConfigurationVariableTypeTerraform   ConfigurationVariableType = 1
)

var VariableTypes = map[string]ConfigurationVariableType{
	"terraform":   ConfigurationVariableTypeTerraform,
	"environment": ConfigurationVariableTypeEnvironment,
}

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

type Role string

const (
	Admin    Role = "Admin"
	Deployer Role = "Deployer"
	Planner  Role = "Planner"
	Viewer   Role = "Viewer"
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
	TokenId              string           `json:"tokenId,omitempty"`
	GithubInstallationId int              `json:"githubInstallationId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	Revision             string           `json:"revision"`
	OrganizationId       string           `json:"organizationId"`
	TerraformVersion     string           `json:"terraformVersion"`
	TerragruntVersion    string           `json:"terragruntVersion,omitempty"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	BitbucketClientKey   string           `json:"bitbucketClientKey,omitempty"`
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
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	TokenId              string           `json:"tokenId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	UpdatedAt            string           `json:"updatedAt"`
	TerraformVersion     string           `json:"terraformVersion"`
	TerragruntVersion    string           `json:"terragruntVersion,omitempty"`
	IsDeleted            bool             `json:"isDeleted,omitempty"`
	BitbucketClientKey   string           `json:"bitbucketClientKey"`
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
}

type DeploymentLog struct {
	Id                  string `json:"id"`
	BlueprintId         string `json:"blueprintId"`
	BlueprintRepository string `json:"blueprintRepository"`
	BlueprintRevision   string `json:"blueprintRevision"`
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

type AwsCredentialsType string
type GcpCredentialsType string
type AzureCredentialsType string

const (
	AwsAssumedRoleCredentialsType        AwsCredentialsType   = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
	AwsAccessKeysCredentialsType         AwsCredentialsType   = "AWS_ACCESS_KEYS_FOR_DEPLOYMENT"
	GcpServiceAccountCredentialsType     GcpCredentialsType   = "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT"
	AzureServicePrincipalCredentialsType AzureCredentialsType = "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT"
)

type AwsCredentialsCreatePayload struct {
	Name           string                     `json:"name"`
	OrganizationId string                     `json:"organizationId"`
	Type           AwsCredentialsType         `json:"type"`
	Value          AwsCredentialsValuePayload `json:"value"`
}

type GcpCredentialsCreatePayload struct {
	Name           string                     `json:"name"`
	OrganizationId string                     `json:"organizationId"`
	Type           GcpCredentialsType         `json:"type"`
	Value          GcpCredentialsValuePayload `json:"value"`
}

type AzureCredentialsCreatePayload struct {
	Name           string                       `json:"name"`
	OrganizationId string                       `json:"organizationId"`
	Type           AzureCredentialsType         `json:"type"`
	Value          AzureCredentialsValuePayload `json:"value"`
}

type AwsCredentialsValuePayload struct {
	RoleArn         string `json:"roleArn"`
	ExternalId      string `json:"externalId"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type GcpCredentialsValuePayload struct {
	ProjectId         string `json:"projectId"`
	ServiceAccountKey string `json:"serviceAccountKey"`
}

type AzureCredentialsValuePayload struct {
	ClientId       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
}

type ApiKey struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}

type TeamCreatePayload struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
}

type TeamUpdatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Team struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
}

type TeamProjectAssignmentPayload struct {
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole Role   `json:"projectRole"`
}

type TeamProjectAssignment struct {
	Id          string `json:"id"`
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole Role   `json:"projectRole"`
}

type Policy struct {
	Id                          string `json:"id"`
	ProjectId                   string `json:"projectId"`
	NumberOfEnvironments        int    `json:"numberOfEnvironments"`
	NumberOfEnvironmentsTotal   int    `json:"numberOfEnvironmentsTotal"`
	RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool   `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
	UpdatedBy                   string `json:"updatedBy"`
	RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
}

type PolicyUpdatePayload struct {
	ProjectId                   string `json:"projectId"`
	NumberOfEnvironments        int    `json:"numberOfEnvironments"`
	NumberOfEnvironmentsTotal   int    `json:"numberOfEnvironmentsTotal"`
	RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool   `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
	RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
}

type EnvironmentSchedulingExpression struct {
	Cron    string `json:"cron,omitempty"`
	Enabled bool   `json:"enabled"`
}

type EnvironmentScheduling struct {
	Deploy  *EnvironmentSchedulingExpression `json:"deploy,omitempty"`
	Destroy *EnvironmentSchedulingExpression `json:"destroy,omitempty"`
}

type WorkflowTrigger struct {
	Id string `json:"id"`
}

type WorkflowTriggerUpsertPayload struct {
	DownstreamEnvironmentIds []string `json:"downstreamEnvironmentIds"`
}

type NotificationType string

const (
	NotificationTypeSlack NotificationType = "Slack"
	NotificationTypeTeams NotificationType = "Teams"
)

type Notification struct {
	Id             string           `json:"id"`
	CreatedBy      string           `json:"createdBy"`
	CreatedByUser  User             `json:"createdByUser"`
	OrganizationId string           `json:"organizationId"`
	Name           string           `json:"name"`
	Type           NotificationType `json:"type"`
	Value          string           `json:"value"`
}

type NotificationCreate struct {
	Name  string           `json:"name"`
	Type  NotificationType `json:"type"`
	Value string           `json:"value"`
}

type NotificationCreateWithOrganizationId struct {
	NotificationCreate
	OrganizationId string `json:"organizationId"`
}

type NotificationUpdate struct {
	Name  string           `json:"name,omitempty"`
	Type  NotificationType `json:"type,omitempty"`
	Value string           `json:"value,omitempty"`
}

type ModuleSshKey struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Module struct {
	ModuleName           string         `json:"moduleName"`
	ModuleProvider       string         `json:"moduleProvider"`
	Repository           string         `json:"repository"`
	Description          string         `json:"description"`
	LogoUrl              string         `json:"logoUrl"`
	TokenId              string         `json:"tokenId"`
	TokenName            string         `json:"tokenName"`
	GithubInstallationId int            `json:"githubInstallationId"`
	BitbucketClientKey   string         `json:"bitbucketClientKey"`
	IsGitLab             bool           `json:"isGitLab"`
	GitlabProjectId      int            `json:"gitlabProjectId"`
	SshKeys              []ModuleSshKey `json:"sshkeys"`
	Type                 string         `json:"type"`
	Id                   string         `json:"id"`
	OrganizationId       string         `json:"organizationId"`
	Author               User           `json:"author"`
	AuthorId             string         `json:"authorId"`
	CreatedAt            string         `json:"createdAt"`
	UpdatedAt            string         `json:"updatedAt"`
	IsDeleted            bool           `json:"isDeleted"`
}

type ModuleCreatePayload struct {
	ModuleName           string         `json:"moduleName"`
	ModuleProvider       string         `json:"moduleProvider"`
	Repository           string         `json:"repository"`
	Description          string         `json:"description,omitempty"`
	LogoUrl              string         `json:"logoUrl,omitEmpty"`
	TokenId              string         `json:"tokenId,omitempty"`
	TokenName            string         `json:"tokenName,omitempty"`
	GithubInstallationId *int           `json:"githubInstallationId,omitempty"`
	BitbucketClientKey   string         `json:"bitbucketClientKey,omitempty"`
	IsGitLab             *bool          `json:"isGitLab,omitempty"`
	GitlabProjectId      *int           `json:"gitlabProjectId,omitempty"`
	SshKeys              []ModuleSshKey `json:"sshkeys,omitempty"`
}

func (p PolicyUpdatePayload) MarshalJSON() ([]byte, error) {
	type serial struct {
		ProjectId                   string `json:"projectId"`
		NumberOfEnvironments        *int   `json:"numberOfEnvironments"`
		NumberOfEnvironmentsTotal   *int   `json:"numberOfEnvironmentsTotal"`
		RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
		IncludeCostEstimation       bool   `json:"includeCostEstimation"`
		SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
		DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
		SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
		RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
		ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
	}

	s := serial{
		ProjectId:                   p.ProjectId,
		RequiresApprovalDefault:     p.RequiresApprovalDefault,
		IncludeCostEstimation:       p.IncludeCostEstimation,
		SkipApplyWhenPlanIsEmpty:    p.SkipApplyWhenPlanIsEmpty,
		DisableDestroyEnvironments:  p.DisableDestroyEnvironments,
		SkipRedundantDeployments:    p.SkipRedundantDeployments,
		RunPullRequestPlanDefault:   p.RunPullRequestPlanDefault,
		ContinuousDeploymentDefault: p.ContinuousDeploymentDefault,
	}

	if p.NumberOfEnvironments != 0 {
		s.NumberOfEnvironments = &p.NumberOfEnvironments
	}
	if p.NumberOfEnvironmentsTotal != 0 {
		s.NumberOfEnvironmentsTotal = &p.NumberOfEnvironmentsTotal
	}

	return json.Marshal(s)
}
